package dto

// SessionCreateRequest 策展人创建场次。
type SessionCreateRequest struct {
	Date      string `json:"date" binding:"required,datetime=2006-01-02"`
	StartTime string `json:"startTime" binding:"required,len=5"` // HH:MM
	EndTime   string `json:"endTime" binding:"required,len=5"`
	Capacity  int    `json:"capacity" binding:"required,min=1,max=10000"`
}

// ReservationCreateRequest 观众提交预约。
type ReservationCreateRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Phone     string `json:"phone" binding:"required,cn_mobile"`
	PartySize int    `json:"partySize" binding:"required,min=1,max=10"`
}
