package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SessionRepository 场次数据访问。
type SessionRepository struct {
	coll *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{coll: db.Collection("sessions")}
}

func (r *SessionRepository) Create(ctx context.Context, s *model.Session) error {
	_, err := r.coll.InsertOne(ctx, s)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

func (r *SessionRepository) FindByID(ctx context.Context, id string) (*model.Session, error) {
	var s model.Session
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}
	return &s, nil
}

func (r *SessionRepository) ListByExhibition(ctx context.Context, exhibitionID string) ([]model.Session, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"exhibition_id": exhibitionID},
		options.Find().SetSort(bson.D{{Key: "date", Value: 1}, {Key: "start_time", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer cursor.Close(ctx)
	list := []model.Session{}
	if err := cursor.All(ctx, &list); err != nil {
		return nil, fmt.Errorf("decode sessions: %w", err)
	}
	return list, nil
}

// TryBookCapacity 原子预占名额：仅当剩余名额充足时 booked_count += delta，
// 返回 false 表示名额不足，调用方应拒绝预约。
func (r *SessionRepository) TryBookCapacity(ctx context.Context, id string, delta int) (bool, error) {
	res, err := r.coll.UpdateOne(ctx,
		bson.M{
			"_id": id,
			"$expr": bson.M{
				"$lte": bson.A{bson.M{"$add": bson.A{"$booked_count", delta}}, "$capacity"},
			},
		},
		bson.M{"$inc": bson.M{"booked_count": delta}, "$set": bson.M{"updated_at": time.Now()}},
	)
	if err != nil {
		return false, fmt.Errorf("book session capacity: %w", err)
	}
	return res.MatchedCount > 0, nil
}

// ReleaseCapacity 取消预约时释放名额（booked_count 不会减成负数）。
func (r *SessionRepository) ReleaseCapacity(ctx context.Context, id string, delta int) error {
	_, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": id, "booked_count": bson.M{"$gte": delta}},
		bson.M{"$inc": bson.M{"booked_count": -delta}, "$set": bson.M{"updated_at": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("release session capacity: %w", err)
	}
	return nil
}
