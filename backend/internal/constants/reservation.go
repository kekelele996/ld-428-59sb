package constants

// VisitSession 参观场次。
const (
	SessionMorning   = "Morning"
	SessionAfternoon = "Afternoon"
	SessionEvening   = "Evening"
)

// VisitSessions 全部场次（用于校验与遍历）。
var VisitSessions = []string{SessionMorning, SessionAfternoon, SessionEvening}

// IsValidSession 校验场次是否合法。
func IsValidSession(s string) bool {
	for _, item := range VisitSessions {
		if item == s {
			return true
		}
	}
	return false
}

// SessionCapacity 每个场次的预约人数上限。
const SessionCapacity = 50

// MaxVisitorsPerReservation 单次预约人数上限。
const MaxVisitorsPerReservation = 10

// ReservationStatus 预约状态。
const (
	ReservationConfirmed = "Confirmed"
	ReservationCancelled = "Cancelled"
)
