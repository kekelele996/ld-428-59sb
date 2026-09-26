// artvault MongoDB 初始化脚本（参考用）
// 实际集合与种子数据由后端启动时自动写入（database/seeds/seed.go）。
db = db.getSiblingDB('artvault');
db.createCollection('users');
db.createCollection('artists');
db.createCollection('artworks');
db.createCollection('exhibitions');
db.createCollection('interactions');
db.createCollection('review_logs');
db.createCollection('audit_logs');
db.createCollection('sessions');
db.createCollection('reservations');

// 同一手机号在同一场次仅允许一条有效预约（取消后可重新预约）。
db.reservations.createIndex(
  { session_id: 1, phone: 1 },
  { name: 'uniq_active_session_phone', unique: true, partialFilterExpression: { status: 'Booked' } },
);
db.sessions.createIndex({ exhibition_id: 1, date: 1 }, { name: 'exhibition_date' });
db.reservations.createIndex({ exhibition_id: 1, created_at: -1 }, { name: 'exhibition_created' });
