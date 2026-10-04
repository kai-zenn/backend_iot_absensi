package entity

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
    ID              uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
    StudentID       uuid.UUID  `gorm:"column:student_id"`
    AttendanceDate  time.Time  `gorm:"column:attendance_date"`
    CheckInAt       time.Time  `gorm:"column:check_in_at"`
    CheckInStatus   string     `gorm:"column:check_in_status"` // PRESENT | LATE
    CheckInSource   string     `gorm:"column:check_in_source"` // CARD | MANUAL
    CheckInEventID  *uuid.UUID `gorm:"column:check_in_event_id"`
    CheckInByUserID *uuid.UUID `gorm:"column:check_in_by_user_id"`
    Note            *string    `gorm:"column:note"`
    CheckOutAt      *time.Time `gorm:"column:check_out_at"`     // reserved
    CheckOutStatus  *string    `gorm:"column:check_out_status"` // reserved
    CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime"`

    Student *Student `gorm:"foreignKey:StudentID"`
}

func (Attendance) TableName() string { return "attendance" }
