package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/artvault/artvault/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrSessionFull 场次名额已满哨兵错误。
var ErrSessionFull = fmt.Errorf("session capacity exceeded")

// ReservationRepository 预约数据访问。
type ReservationRepository struct {
	coll     *mongo.Collection
	counters *mongo.Collection
}

func NewReservationRepository(db *mongo.Database) *ReservationRepository {
	return &ReservationRepository{
		coll:     db.Collection("reservations"),
		counters: db.Collection("reservation_counters"),
	}
}

// EnsureIndexes 创建幂等索引：同一手机号在同一展览同一日期同一场次仅允许一条有效预约。
func (r *ReservationRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "exhibition_id", Value: 1},
			{Key: "visit_date", Value: 1},
			{Key: "session", Value: 1},
			{Key: "phone", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetName("uniq_active_reservation").
			SetPartialFilterExpression(bson.M{"status": "Confirmed"}),
	})
	if err != nil {
		return fmt.Errorf("create reservation index: %w", err)
	}
	return nil
}

// IsDuplicateKey 判断错误是否为唯一索引冲突。
func IsDuplicateKey(err error) bool {
	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				return true
			}
		}
	}
	var ce mongo.CommandError
	if errors.As(err, &ce) && ce.Code == 11000 {
		return true
	}
	return false
}

// counterKey 场次计数器主键。
func counterKey(exhibitionID, visitDate, session string) string {
	return exhibitionID + "|" + visitDate + "|" + session
}

// TryReserve 原子占用场次名额；名额不足时返回 ErrSessionFull。
func (r *ReservationRepository) TryReserve(ctx context.Context, exhibitionID, visitDate, session string, count, capacity int) error {
	filter := bson.M{
		"_id":      counterKey(exhibitionID, visitDate, session),
		"reserved": bson.M{"$lte": capacity - count},
	}
	update := bson.M{
		"$inc": bson.M{"reserved": count},
		"$setOnInsert": bson.M{
			"exhibition_id": exhibitionID,
			"visit_date":    visitDate,
			"session":       session,
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err := r.counters.FindOneAndUpdate(ctx, filter, update, opts).Err()
	if err == nil {
		return nil
	}
	if !IsDuplicateKey(err) {
		return fmt.Errorf("reserve session capacity: %w", err)
	}
	// upsert 冲突有两种可能：并发首次占额（可重试）或名额已满。
	// 去掉 upsert 重试一次，仍无匹配文档说明计数器已存在且超出上限。
	err = r.counters.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"reserved": count}}).Err()
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrSessionFull
	}
	return fmt.Errorf("reserve session capacity: %w", err)
}

// Release 释放场次名额（取消预约时调用）。
func (r *ReservationRepository) Release(ctx context.Context, exhibitionID, visitDate, session string, count int) error {
	_, err := r.counters.UpdateOne(ctx,
		bson.M{"_id": counterKey(exhibitionID, visitDate, session)},
		bson.M{"$inc": bson.M{"reserved": -count}},
	)
	if err != nil {
		return fmt.Errorf("release session capacity: %w", err)
	}
	return nil
}

// ReservedCount 查询某场次已确认人数。
func (r *ReservationRepository) ReservedCount(ctx context.Context, exhibitionID, visitDate, session string) (int, error) {
	var doc struct {
		Reserved int `bson:"reserved"`
	}
	err := r.counters.FindOne(ctx, bson.M{"_id": counterKey(exhibitionID, visitDate, session)}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("find session counter: %w", err)
	}
	return doc.Reserved, nil
}

func (r *ReservationRepository) Create(ctx context.Context, rv *model.Reservation) error {
	_, err := r.coll.InsertOne(ctx, rv)
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

// FindConfirmedByPhone 查询同一手机号在同一场次的有效预约（幂等防重）。
func (r *ReservationRepository) FindConfirmedByPhone(ctx context.Context, exhibitionID, visitDate, session, phone string) (*model.Reservation, error) {
	var rv model.Reservation
	err := r.coll.FindOne(ctx, bson.M{
		"exhibition_id": exhibitionID,
		"visit_date":    visitDate,
		"session":       session,
		"phone":         phone,
		"status":        "Confirmed",
	}).Decode(&rv)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find reservation by phone: %w", err)
	}
	return &rv, nil
}

// ListByExhibition 策展人查看展览的全部预约；exhibitionID 为空时返回全部。
func (r *ReservationRepository) ListByExhibition(ctx context.Context, exhibitionID string) ([]model.Reservation, error) {
	filter := bson.M{}
	if exhibitionID != "" {
		filter["exhibition_id"] = exhibitionID
	}
	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(500))
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

func (r *ReservationRepository) Update(ctx context.Context, rv *model.Reservation) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": rv.ID}, rv)
	if err != nil {
		return fmt.Errorf("update reservation: %w", err)
	}
	return nil
}
