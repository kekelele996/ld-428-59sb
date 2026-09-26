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

// SessionService 场次业务逻辑。
type SessionService struct {
	sessionRepo *repository.SessionRepository
	exhRepo     *repository.ExhibitionRepository
	logger      *slog.Logger
}

func NewSessionService(sessionRepo *repository.SessionRepository, exhRepo *repository.ExhibitionRepository, logger *slog.Logger) *SessionService {
	return &SessionService{sessionRepo: sessionRepo, exhRepo: exhRepo, logger: logger}
}

func (s *SessionService) Create(ctx context.Context, exhibitionID, curatorID, date, startTime, endTime string, capacity int) (*model.Session, error) {
	e, err := s.exhRepo.FindByID(ctx, exhibitionID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, "展览不存在", nil)
	}
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if e.CuratorID != curatorID {
		return nil, util.NewAppError(constants.CodeForbidden, constants.MsgForbidden, nil)
	}
	if err := validateSessionWindow(e, date, startTime, endTime); err != nil {
		return nil, err
	}
	session := &model.Session{
		ID:           util.NewID("ses"),
		ExhibitionID: exhibitionID,
		Date:         date,
		StartTime:    startTime,
		EndTime:      endTime,
		Capacity:     capacity,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogSessionCreated, "sessionID", session.ID, "exhibitionID", exhibitionID)
	return session, nil
}

// ListByExhibition 公开查询展览场次。
func (s *SessionService) ListByExhibition(ctx context.Context, exhibitionID string) ([]model.Session, error) {
	list, err := s.sessionRepo.ListByExhibition(ctx, exhibitionID)
	if err != nil {
		return nil, util.Wrap(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return list, nil
}

// validateSessionWindow 校验时间格式、时段顺序与日期落在展期内。
func validateSessionWindow(e *model.Exhibition, date, startTime, endTime string) error {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return util.NewAppError(constants.CodeValidationFailed, "日期格式应为 YYYY-MM-DD", nil)
	}
	st, err1 := time.Parse("15:04", startTime)
	et, err2 := time.Parse("15:04", endTime)
	if err1 != nil || err2 != nil {
		return util.NewAppError(constants.CodeValidationFailed, "场次时间格式应为 HH:MM", nil)
	}
	if !et.After(st) {
		return util.NewAppError(constants.CodeValidationFailed, "结束时间必须晚于开始时间", nil)
	}
	if e.StartDate != "" {
		if start, perr := time.Parse("2006-01-02", e.StartDate); perr == nil && day.Before(start) {
			return util.NewAppError(constants.CodeValidationFailed, "场次日期不能早于展览开始日期", nil)
		}
	}
	if e.EndDate != "" {
		if end, perr := time.Parse("2006-01-02", e.EndDate); perr == nil && day.After(end) {
			return util.NewAppError(constants.CodeValidationFailed, "场次日期不能晚于展览结束日期", nil)
		}
	}
	return nil
}
