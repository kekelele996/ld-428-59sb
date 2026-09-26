package model

import "time"

// Reservation 展览参观预约。
type Reservation struct {
	ID           string    `bson:"_id" json:"id"`
	ExhibitionID string    `bson:"exhibition_id" json:"exhibitionId"`
	VisitDate    string    `bson:"visit_date" json:"visitDate"`
	Session      string    `bson:"session" json:"session"`
	VisitorCount int       `bson:"visitor_count" json:"visitorCount"`
	Phone        string    `bson:"phone" json:"phone"`
	Status       string    `bson:"status" json:"status"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`
}
