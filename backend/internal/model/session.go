package model

import "time"

// Session 展览预约场次。
type Session struct {
	ID           string    `bson:"_id" json:"id"`
	ExhibitionID string    `bson:"exhibition_id" json:"exhibitionId"`
	Date         string    `bson:"date" json:"date"` // YYYY-MM-DD
	StartTime    string    `bson:"start_time" json:"startTime"`
	EndTime      string    `bson:"end_time" json:"endTime"`
	Capacity     int       `bson:"capacity" json:"capacity"`
	BookedCount  int       `bson:"booked_count" json:"bookedCount"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`
}
