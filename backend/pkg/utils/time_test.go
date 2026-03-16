package utils

import (
	"testing"
	"time"
)

func TestNowUTC(t *testing.T) {
	before := time.Now().UTC()
	got := NowUTC()
	after := time.Now().UTC()

	if got.Before(before) || got.After(after) {
		t.Errorf("NowUTC() = %v, expected between %v and %v", got, before, after)
	}

	if got.Location() != time.UTC {
		t.Errorf("NowUTC() location = %v, want UTC", got.Location())
	}
}
