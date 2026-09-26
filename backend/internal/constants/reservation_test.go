package constants

import "testing"

func TestVisitSessions(t *testing.T) {
	cases := []struct {
		session string
		valid   bool
	}{
		{SessionMorning, true},
		{SessionAfternoon, true},
		{SessionEvening, true},
		{"Night", false},
		{"", false},
	}
	for _, c := range cases {
		if got := IsValidSession(c.session); got != c.valid {
			t.Fatalf("IsValidSession(%q) = %v, want %v", c.session, got, c.valid)
		}
	}
	if len(VisitSessions) != 3 {
		t.Fatalf("VisitSessions len = %d, want 3", len(VisitSessions))
	}
}

func TestReservationStatuses(t *testing.T) {
	if ReservationConfirmed != "Confirmed" || ReservationCancelled != "Cancelled" {
		t.Fatal("reservation status mismatch")
	}
}

func TestSessionCapacity(t *testing.T) {
	if SessionCapacity <= 0 {
		t.Fatal("session capacity must be positive")
	}
	if MaxVisitorsPerReservation > SessionCapacity {
		t.Fatal("per-reservation limit must not exceed session capacity")
	}
}
