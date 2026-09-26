package handler

import (
	"context"
	"net/http"

	"github.com/artvault/artvault/internal/dto"
	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// SessionHandler 场次接口。
type SessionHandler struct {
	svc *service.SessionService
}

func NewSessionHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{svc: svc}
}

// List 公开查询某展览的场次列表。
func (h *SessionHandler) List(c *gin.Context) {
	list, err := h.svc.ListByExhibition(context.Background(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

// Create 策展人为自己的展览新增场次。
func (h *SessionHandler) Create(c *gin.Context) {
	var req dto.SessionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	userID, _ := c.Get("userID")
	created, err := h.svc.Create(context.Background(), c.Param("id"), userID.(string), req.Date, req.StartTime, req.EndTime, req.Capacity)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, created)
}
