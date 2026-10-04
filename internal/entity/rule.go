package entity

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceRule struct {
    ID            uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
    EffectiveFrom time.Time  `gorm:"column:effective_from"` // DATE
    OnTimeUntil   time.Time  `gorm:"column:on_time_until"`   // hanya komponen jam dipakai
    CreatedBy     *uuid.UUID `gorm:"column:created_by"`
    CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (AttendanceRule) TableName() string { return "attendance_rules" }
