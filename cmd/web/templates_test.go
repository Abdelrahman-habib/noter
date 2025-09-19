package main

import (
	"testing"
	"time"

	"github.com/Abdelrahman-habib/noter/internal/assert"
)

// func humanDate(t time.Time) string {
// 	return t.Format("02 Jan 2006 at 15:04")
// }

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 9, 13, 10, 15, 0, 0, time.UTC),
			want: "13 Sep 2025 at 10:15",
		},
		{
			name: "zero time",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2025, 9, 13, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "13 Sep 2025 at 09:15",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, humanDate(tt.tm), tt.want)
		})
	}
}
