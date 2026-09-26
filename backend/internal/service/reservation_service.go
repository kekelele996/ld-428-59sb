package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/model"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/util"
)

// ReservationService 预约业务逻辑。
type ReservationService struct {
	rvRepo  *repository.ReservationRepository
	sesRepo *repository.SessionRepository
	exhRepo *repository.ExhibitionRepository
	logger  *slog.Logger
}

func NewReservationService(rvRepo *repository.ReservationRepository, sesRepo *repository.SessionRepository, exhRepo *repository.ExhibitionRepository, logger *slog.Logger) *ReservationService {
	return &ReservationService{rvRepo: rvRepo, sesRepo: sesRepo, exhRepo: exhRepo, logger: logger}
}

// Create 观众提交预约：先原子占座，再写入记录（重复手机号/超员均拒绝）。
func (s *ReservationService) Create(ctx context.Context, sessionID, phone string, partySize int) (*model.Reservation, error) {
	ses, err := s.sesRepo.FindByID(ctx, sessionID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, "场次不存在", nil)
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if err := s.ensureBookable(ctx, ses); err != nil {
		return nil, err
	}

	// 原子预占名额：超出上限直接拒绝。
	ok, err := s.sesRepo.TryBookCapacity(ctx, sessionID, partySize)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if !ok {
		return nil, util.NewAppError(constants.CodeSessionFull, "该场次剩余名额不足，预约失败", nil)
	}

	rv := &model.Reservation{
		ID:           util.NewID("rsv"),
		ExhibitionID: ses.ExhibitionID,
		SessionID:    sessionID,
		Phone:        phone,
		PartySize:    partySize,
		Status:       model.ReservationBooked,
		CreatedAt:    time.Now(),
	}
	if err := s.rvRepo.Create(ctx, rv); err != nil {
		// 同手机号同场次重复提交：回滚刚占的名额。
		if errors.Is(err, repository.ErrDuplicateReservation) {
			if rerr := s.sesRepo.ReleaseCapacity(ctx, sessionID, partySize); rerr != nil {
				s.logger.Error("release capacity after duplicate reservation failed", "error", rerr, "sessionID", sessionID)
			}
			return nil, util.NewAppError(constants.CodeDuplicateReservation, "该手机号已预约本场次，请勿重复提交", nil)
		}
		// 写入失败同样回滚，避免名额被幽灵占用。
		if rerr := s.sesRepo.ReleaseCapacity(ctx, sessionID, partySize); rerr != nil {
			s.logger.Error("release capacity after reservation insert failed", "error", rerr, "sessionID", sessionID)
		}
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}

	s.logger.Info(constants.LogReservationMade, "reservationID", rv.ID, "sessionID", sessionID, "partySize", partySize)
	return rv, nil
}

// ListByExhibition 策展人查看某展览的预约记录。
func (s *ReservationService) ListByExhibition(ctx context.Context, exhibitionID, curatorID string) ([]model.Reservation, error) {
	if err := s.ensureOwnedExhibition(ctx, exhibitionID, curatorID); err != nil {
		return nil, err
	}
	list, err := s.rvRepo.ListByExhibition(ctx, exhibitionID)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return list, nil
}

// Cancel 策展人取消预约，名额立即释放。
func (s *ReservationService) Cancel(ctx context.Context, reservationID, curatorID string) (*model.Reservation, error) {
	rv, err := s.rvRepo.FindByID(ctx, reservationID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, "预约记录不存在", nil)
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if err := s.ensureOwnedExhibition(ctx, rv.ExhibitionID, curatorID); err != nil {
		return nil, err
	}
	if rv.Status != model.ReservationBooked {
		return nil, util.NewAppError(constants.CodeConflict, "该预约已取消，无需重复操作", nil)
	}

	matched, err := s.rvRepo.Cancel(ctx, reservationID)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if !matched {
		return nil, util.NewAppError(constants.CodeConflict, "该预约已被取消", nil)
	}
	// 名额立即释放。
	if err := s.sesRepo.ReleaseCapacity(ctx, rv.SessionID, rv.PartySize); err != nil {
		s.logger.Error("release capacity after cancel failed", "error", err, "reservationID", reservationID)
	}

	rv.Status = model.ReservationCanceled
	canceledAt := time.Now()
	rv.CanceledAt = &canceledAt
	s.logger.Info(constants.LogReservationCanceled, "reservationID", reservationID, "sessionID", rv.SessionID)
	return rv, nil
}

// ensureOwnedExhibition 校验展览存在且属于当前策展人。
func (s *ReservationService) ensureOwnedExhibition(ctx context.Context, exhibitionID, curatorID string) error {
	e, err := s.exhRepo.FindByID(ctx, exhibitionID)
	if errors.Is(err, repository.ErrNotFound) {
		return util.NewAppError(constants.CodeNotFound, "展览不存在", nil)
	}
	if err != nil {
		return util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if e.CuratorID != curatorID {
		return util.NewAppError(constants.CodeForbidden, constants.MsgForbidden, nil)
	}
	return nil
}

// ensureBookable 校验预约开放条件：
// 展览已审核发布、未结束/归档、场次日期未过。
func (s *ReservationService) ensureBookable(ctx context.Context, ses *model.Session) error {
	e, err := s.exhRepo.FindByID(ctx, ses.ExhibitionID)
	if errors.Is(err, repository.ErrNotFound) {
		return util.NewAppError(constants.CodeNotFound, "展览不存在", nil)
	}
	if err != nil {
		return util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if e.ReviewStatus != "Approved" {
		return util.NewAppError(constants.CodeReservationClosed, "展览尚未发布，暂不接受预约", nil)
	}
	if e.Status == constants.ExhibitionEnded || e.Status == constants.ExhibitionArchived {
		return util.NewAppError(constants.CodeReservationClosed, "展览已结束，预约通道已关闭", nil)
	}
	if e.Status != constants.ExhibitionActive {
		return util.NewAppError(constants.CodeReservationClosed, "展览尚未发布，暂不接受预约", nil)
	}
	if e.EndDate != "" {
		if end, perr := time.Parse("2006-01-02", e.EndDate); perr == nil && time.Now().After(end.Add(24*time.Hour)) {
			return util.NewAppError(constants.CodeReservationClosed, "展览展期已过，预约通道已关闭", nil)
		}
	}
	if day, perr := time.Parse("2006-01-02", ses.Date); perr == nil {
		today, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		if day.Before(today) {
			return util.NewAppError(constants.CodeReservationClosed, "该场次已过期，无法预约", nil)
		}
	}
	return nil
}
