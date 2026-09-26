package seeds

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes 创建业务索引（幂等）。
func EnsureIndexes(ctx context.Context, db *mongo.Database, logger *slog.Logger) {
	// 同一手机号在同一场次仅允许存在一条有效（Booked）预约；
	// 取消后状态变为 Canceled，可重新预约并由新记录占位。
	_, err := db.Collection("reservations").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "session_id", Value: 1}, {Key: "phone", Value: 1}},
		Options: options.Index().
			SetName("uniq_active_session_phone").
			SetUnique(true).
			SetPartialFilterExpression(bson.M{"status": "Booked"}),
	})
	if err != nil {
		logger.Error("create reservations index failed", "error", err)
	}

	// 按展览查询场次/预约的常用过滤索引。
	_, _ = db.Collection("sessions").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "exhibition_id", Value: 1}, {Key: "date", Value: 1}},
		Options: options.Index().SetName("exhibition_date"),
	})
	_, _ = db.Collection("reservations").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "exhibition_id", Value: 1}, {Key: "created_at", Value: -1}},
		Options: options.Index().SetName("exhibition_created"),
	})
}
