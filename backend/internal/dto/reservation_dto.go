package dto

// ReservationCreateRequest 观众提交参观预约。
type ReservationCreateRequest struct {
	VisitDate    string `json:"visitDate" binding:"required"`
	Session      string `json:"session" binding:"required,oneof=Morning Afternoon Evening"`
	VisitorCount int    `json:"visitorCount" binding:"required,min=1,max=10"`
	Phone        string `json:"phone" binding:"required"`
}

// SessionAvailability 场次余票。
type SessionAvailability struct {
	Session   string `json:"session"`
	Capacity  int    `json:"capacity"`
	Reserved  int    `json:"reserved"`
	Remaining int    `json:"remaining"`
}
