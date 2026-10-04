package entity

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceEvent struct {
    ID               uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
    EventID          string     `gorm:"column:event_id"`
    DeviceID         uuid.UUID  `gorm:"column:device_id"`
    CardUID          string     `gorm:"column:card_uid"`
    OccurredAt       time.Time  `gorm:"column:occurred_at"`
    ReceivedAt       time.Time  `gorm:"column:received_at;autoCreateTime"`
    Source           string     `gorm:"column:source"` // LIVE | SYNC
    Result           string     `gorm:"column:result"`
    StudentID        *uuid.UUID `gorm:"column:student_id"`
    AttendanceID     *uuid.UUID `gorm:"column:attendance_id"`
    ClockSkewSeconds *int       `gorm:"column:clock_skew_seconds"`
    CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (AttendanceEvent) TableName() string { return "attendance_events" }
