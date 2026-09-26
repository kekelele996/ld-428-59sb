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

// Create 观众提交预约（公开接口，凭手机号预约）。
func (h *ReservationHandler) Create(c *gin.Context) {
	var req dto.ReservationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 40001, util.TranslateError(err))
		return
	}
	rv, duplicated, err := h.svc.Create(context.Background(), c.Param("id"), req.VisitDate, req.Session, req.VisitorCount, req.Phone)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, gin.H{"reservation": rv, "duplicated": duplicated})
}

// Availability 查询某日期各场次余票。
func (h *ReservationHandler) Availability(c *gin.Context) {
	visitDate := c.DefaultQuery("date", "")
	if visitDate == "" {
		util.Fail(c, http.StatusBadRequest, 40001, "字段 date 校验失败：不能为空")
		return
	}
	list, err := h.svc.Availability(context.Background(), c.Param("id"), visitDate)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

// List 策展人在工作台查看展览预约。
func (h *ReservationHandler) List(c *gin.Context) {
	list, err := h.svc.ListByExhibition(context.Background(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, list)
}

// Cancel 策展人取消预约，名额立即释放。
func (h *ReservationHandler) Cancel(c *gin.Context) {
	rv, err := h.svc.Cancel(context.Background(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, rv)
}
