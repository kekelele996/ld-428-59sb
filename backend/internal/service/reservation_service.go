package service

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"time"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/dto"
	"github.com/artvault/artvault/internal/model"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/util"
)

// phonePattern 中国大陆手机号。
var phonePattern = regexp.MustCompile(`^1[3-9]\d{9}$`)

// ReservationService 预约业务逻辑。
type ReservationService struct {
	repo           *repository.ReservationRepository
	exhibitionRepo *repository.ExhibitionRepository
	logger         *slog.Logger
}

func NewReservationService(repo *repository.ReservationRepository, exhibitionRepo *repository.ExhibitionRepository, logger *slog.Logger) *ReservationService {
	return &ReservationService{repo: repo, exhibitionRepo: exhibitionRepo, logger: logger}
}

// Create 提交预约。同一手机号在同一场次重复提交时返回已有记录（duplicated=true），不重复占名额。
func (s *ReservationService) Create(ctx context.Context, exhibitionID, visitDate, session string, visitorCount int, phone string) (*model.Reservation, bool, error) {
	e, err := s.exhibitionRepo.FindByID(ctx, exhibitionID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, false, util.NewAppError(constants.CodeNotFound, "展览不存在", nil)
	}
	if err != nil {
		return nil, false, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	// 展览未发布或已结束时停止预约。
	if e.Status != constants.ExhibitionActive {
		s.logger.Info(constants.LogReservationRejected, "exhibitionID", exhibitionID, "reason", "status="+e.Status)
		return nil, false, util.NewAppError(constants.CodeReservationClosed, constants.MsgReservationClosed, nil)
	}
	if !phonePattern.MatchString(phone) {
		return nil, false, util.NewAppError(constants.CodeValidationFailed, constants.MsgPhoneInvalid, nil)
	}
	if !constants.IsValidSession(session) {
		return nil, false, util.NewAppError(constants.CodeValidationFailed, constants.MsgSessionInvalid, nil)
	}
	if visitorCount < 1 || visitorCount > constants.MaxVisitorsPerReservation {
		return nil, false, util.NewAppError(constants.CodeValidationFailed, "单次预约人数须在 1-10 之间", nil)
	}
	// 参观日期必须为 2006-01-02 格式，且在展期内、不早于今天（字典序即时间序）。
	if _, err := time.Parse("2006-01-02", visitDate); err != nil {
		return nil, false, util.NewAppError(constants.CodeValidationFailed, constants.MsgVisitDateInvalid, nil)
	}
	today := util.FormatDate(time.Now())
	if visitDate < today || visitDate < e.StartDate || visitDate > e.EndDate {
		return nil, false, util.NewAppError(constants.CodeValidationFailed, constants.MsgVisitDateInvalid, nil)
	}
	// 幂等：同一手机号在同一场次已有有效预约时直接返回，不占两次名额。
	existing, err := s.repo.FindConfirmedByPhone(ctx, exhibitionID, visitDate, session, phone)
	if err == nil {
		return existing, true, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, false, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	// 原子占用场次名额，超出上限拒绝。
	if err := s.repo.TryReserve(ctx, exhibitionID, visitDate, session, visitorCount, constants.SessionCapacity); err != nil {
		if errors.Is(err, repository.ErrSessionFull) {
			s.logger.Info(constants.LogReservationRejected, "exhibitionID", exhibitionID, "session", session, "reason", "session full")
			return nil, false, util.NewAppError(constants.CodeReservationFull, constants.MsgReservationFull, nil)
		}
		return nil, false, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	rv := &model.Reservation{
		ID:           util.NewID("rsv"),
		ExhibitionID: exhibitionID,
		VisitDate:    visitDate,
		Session:      session,
		VisitorCount: visitorCount,
		Phone:        phone,
		Status:       constants.ReservationConfirmed,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.repo.Create(ctx, rv); err != nil {
		// 写入失败必须回滚已占名额。
		if releaseErr := s.repo.Release(ctx, exhibitionID, visitDate, session, visitorCount); releaseErr != nil {
			s.logger.Error("release session capacity failed", "exhibitionID", exhibitionID, "session", session, "error", releaseErr)
		}
		if repository.IsDuplicateKey(err) {
			// 并发重复提交：返回先到的记录，不占两次名额。
			if existing, findErr := s.repo.FindConfirmedByPhone(ctx, exhibitionID, visitDate, session, phone); findErr == nil {
				return existing, true, nil
			}
		}
		return nil, false, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogReservationCreated, "reservationID", rv.ID, "exhibitionID", exhibitionID, "session", session, "visitors", visitorCount)
	return rv, false, nil
}

// ListByExhibition 策展人在工作台查看预约记录。
func (s *ReservationService) ListByExhibition(ctx context.Context, exhibitionID string) ([]model.Reservation, error) {
	list, err := s.repo.ListByExhibition(ctx, exhibitionID)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return list, nil
}

// Cancel 策展人取消预约，名额立即释放。
func (s *ReservationService) Cancel(ctx context.Context, id string) (*model.Reservation, error) {
	rv, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, "预约不存在", nil)
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if rv.Status == constants.ReservationCancelled {
		return rv, nil
	}
	rv.Status = constants.ReservationCancelled
	rv.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, rv); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if err := s.repo.Release(ctx, rv.ExhibitionID, rv.VisitDate, rv.Session, rv.VisitorCount); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogReservationCancelled, "reservationID", id, "exhibitionID", rv.ExhibitionID)
	return rv, nil
}

// Availability 查询某日期各场次的剩余名额。
func (s *ReservationService) Availability(ctx context.Context, exhibitionID, visitDate string) ([]dto.SessionAvailability, error) {
	if _, err := s.exhibitionRepo.FindByID(ctx, exhibitionID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "展览不存在", nil)
		}
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	result := make([]dto.SessionAvailability, 0, len(constants.VisitSessions))
	for _, session := range constants.VisitSessions {
		reserved, err := s.repo.ReservedCount(ctx, exhibitionID, visitDate, session)
		if err != nil {
			return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
		}
		remaining := constants.SessionCapacity - reserved
		if remaining < 0 {
			remaining = 0
		}
		result = append(result, dto.SessionAvailability{
			Session:   session,
			Capacity:  constants.SessionCapacity,
			Reserved:  reserved,
			Remaining: remaining,
		})
	}
	return result, nil
}
