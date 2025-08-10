package repository

import (
	"reflect"
	"testing"
	"time"

	"github.com/RoGogDBD/subscription-aggregator/internal/models"
)

func TestParseMonthYear(t *testing.T) {
	tests := []struct {
		input       string
		wantYear    int
		wantMonth   time.Month
		expectError bool
	}{
		{"01-20", 2020, time.January, false},
		{"12-1999", 1999, time.December, false},
		{"07-2025", 2025, time.July, false},
		{"00-2020", 0, 0, true},
		{"13-2020", 0, 0, true},
		{"05-1899", 0, 0, true},
		{"05-10000", 0, 0, true},
		{"05-2a20", 0, 0, true},
		{"05-2", 0, 0, true},
		{"5-2020", 0, 0, true},
		{"052020", 0, 0, true},
		{"05-20-10", 0, 0, true},
	}

	for _, tt := range tests {
		got, err := models.ParseMonthYear(tt.input)
		if (err != nil) != tt.expectError {
			t.Errorf("parseMonthYear(%q) error = %v, wantErr %v", tt.input, err, tt.expectError)
			continue
		}
		if !tt.expectError {
			if got.Year() != tt.wantYear || got.Month() != tt.wantMonth || got.Day() != 1 {
				t.Errorf("parseMonthYear(%q) = %v, want year %d month %d", tt.input, got, tt.wantYear, tt.wantMonth)
			}
		}
	}
}

func TestFormatToMMYYYY(t *testing.T) {
	tm := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	want := "07-2025"
	got := models.FormatToMMYYYY(tm)
	if got != want {
		t.Errorf("formatToMMYYYY() = %q; want %q", got, want)
	}
}

func TestSubscription_MarshalJSON(t *testing.T) {
	start := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		sub     models.Subscription
		wantStr string
	}{
		{
			"With end date",
			models.Subscription{
				ID:          "123",
				ServiceName: "Service",
				Price:       100,
				UserID:      "user-1",
				StartDate:   start,
				EndDate:     &end,
			},
			`{"id":"123","service_name":"Service","price":100,"user_id":"user-1","start_date":"03-2023","end_date":"12-2023"}`,
		},
		{
			"Without end date",
			models.Subscription{
				ID:          "456",
				ServiceName: "OtherService",
				Price:       200,
				UserID:      "user-2",
				StartDate:   start,
				EndDate:     nil,
			},
			`{"id":"456","service_name":"OtherService","price":200,"user_id":"user-2","start_date":"03-2023"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, err := tt.sub.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON error: %v", err)
			}
			got := string(gotBytes)
			if got != tt.wantStr {
				t.Errorf("MarshalJSON() = %s; want %s", got, tt.wantStr)
			}
		})
	}
}

func TestSubscription_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonStr     string
		want        models.Subscription
		expectError bool
	}{
		{
			"Valid with end_date",
			`{"id":"1","service_name":"S","price":10,"user_id":"u1","start_date":"01-2023","end_date":"12-2023"}`,
			models.Subscription{
				ID:          "1",
				ServiceName: "S",
				Price:       10,
				UserID:      "u1",
				StartDate:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     ptrTime(time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)),
			},
			false,
		},
		{
			"Valid without end_date",
			`{"id":"2","service_name":"S2","price":20,"user_id":"u2","start_date":"02-2023"}`,
			models.Subscription{
				ID:          "2",
				ServiceName: "S2",
				Price:       20,
				UserID:      "u2",
				StartDate:   time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     nil,
			},
			false,
		},
		{
			"Invalid start_date",
			`{"id":"3","service_name":"S3","price":30,"user_id":"u3","start_date":"invalid"}`,
			models.Subscription{},
			true,
		},
		{
			"Invalid end_date",
			`{"id":"4","service_name":"S4","price":40,"user_id":"u4","start_date":"01-2023","end_date":"bad"}`,
			models.Subscription{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s models.Subscription
			err := s.UnmarshalJSON([]byte(tt.jsonStr))
			if (err != nil) != tt.expectError {
				t.Fatalf("UnmarshalJSON() error = %v, expectError %v", err, tt.expectError)
			}
			if !tt.expectError && !reflect.DeepEqual(s, tt.want) {
				t.Errorf("UnmarshalJSON() = %+v, want %+v", s, tt.want)
			}
		})
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
