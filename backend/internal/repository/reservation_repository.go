package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrDuplicateReservation 同一手机号在同一场次已有有效预约。
var ErrDuplicateReservation = errors.New("duplicate active reservation")

// ReservationRepository 预约数据访问。
type ReservationRepository struct {
	coll *mongo.Collection
}

func NewReservationRepository(db *mongo.Database) *ReservationRepository {
	return &ReservationRepository{coll: db.Collection("reservations")}
}

func (r *ReservationRepository) Create(ctx context.Context, rv *model.Reservation) error {
	_, err := r.coll.InsertOne(ctx, rv)
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicateReservation
	}
	if err != nil {
		return fmt.Errorf("insert reservation: %w", err)
	}
	return nil
}

func (r *ReservationRepository) FindByID(ctx context.Context, id string) (*model.Reservation, error) {
	var rv model.Reservation
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&rv)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find reservation: %w", err)
	}
	return &rv, nil
}

// ListByExhibition 策展人查看某展览的全部预约（最新在前）。
func (r *ReservationRepository) ListByExhibition(ctx context.Context, exhibitionID string) ([]model.Reservation, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"exhibition_id": exhibitionID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("list reservations: %w", err)
	}
	defer cursor.Close(ctx)
	list := []model.Reservation{}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, fmt.Errorf("decode reservations: %w", err)
	}
	return list, nil
}

// CountBySession 统计场次有效预约人数。
func (r *ReservationRepository) CountBySession(ctx context.Context, sessionID, status string) (int64, error) {
	n, err := r.coll.CountDocuments(ctx, bson.M{"session_id": sessionID, "status": status})
	if err != nil {
		return 0, fmt.Errorf("count reservations: %w", err)
	}
	return n, nil
}

// Cancel 将有效预约置为取消状态。
func (r *ReservationRepository) Cancel(ctx context.Context, id string) (bool, error) {
	canceledAt := time.Now()
	res, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": id, "status": model.ReservationBooked},
		bson.M{"$set": bson.M{"status": model.ReservationCanceled, "canceled_at": &canceledAt}},
	)
	if err != nil {
		return false, fmt.Errorf("cancel reservation: %w", err)
	}
	return res.MatchedCount > 0, nil
}
