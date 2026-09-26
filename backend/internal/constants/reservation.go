package constants

// 预约相关错误码。
const (
	CodeSessionFull          = 40902
	CodeDuplicateReservation = 40903
	CodeReservationClosed    = 40904
)

// 预约人数限制。
const (
	MinPartySize = 1
	MaxPartySize = 10
)
