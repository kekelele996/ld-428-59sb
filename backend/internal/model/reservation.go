package model

import "time"

// 预约状态。
const (
	ReservationBooked   = "Booked"   // 已预约，占名额
	ReservationCanceled = "Canceled" // 已取消，名额已释放
)

// Reservation 观众预约记录。
type Reservation struct {
	ID           string     `bson:"_id" json:"id"`
	ExhibitionID string     `bson:"exhibition_id" json:"exhibitionId"`
	SessionID    string     `bson:"session_id" json:"sessionId"`
	Phone        string     `bson:"phone" json:"phone"`
	PartySize    int        `bson:"party_size" json:"partySize"`
	Status       string     `bson:"status" json:"status"`
	CreatedAt    time.Time  `bson:"created_at" json:"createdAt"`
	CanceledAt   *time.Time `bson:"canceled_at,omitempty" json:"canceledAt,omitempty"`
}
