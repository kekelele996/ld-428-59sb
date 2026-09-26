package handler

import (
	"context"
	"net/http"

	"github.com/artvault/artvault/internal/dto"
	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// ReservationHandler 预约接口。
type ReservationHandler struct {
	svc *service.ReservationService
}

func NewReservationHandler(svc *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{svc: svc}
}

// Create 观众提交预约（公开接口，凭手机号联系）。
func (h *ReservationHandler) Create(c *gin.Context) {
	var req dto.ReservationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	created, err := h.svc.Create(context.Background(), req.SessionID, req.Phone, req.PartySize)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, created)
}

// ListByExhibition 策展人查看某展览的预约记录。
func (h *ReservationHandler) ListByExhibition(c *gin.Context) {
	userID, _ := c.Get("userID")
	list, err := h.svc.ListByExhibition(context.Background(), c.Param("id"), userID.(string))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

// Cancel 策展人取消预约并释放名额。
func (h *ReservationHandler) Cancel(c *gin.Context) {
	userID, _ := c.Get("userID")
	rv, err := h.svc.Cancel(context.Background(), c.Param("reservationId"), userID.(string))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, rv)
}
