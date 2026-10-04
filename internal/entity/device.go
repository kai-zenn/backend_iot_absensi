package entity

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
    ID              uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
    DeviceCode      string     `gorm:"column:device_code"`
    Name            string     `gorm:"column:name"`
    Location        *string    `gorm:"column:location"`
    Role            string     `gorm:"column:role"`     // CHECK_IN | CHECK_OUT (reserved)
    Status          string     `gorm:"column:status"` // ACTIVE | DISABLED
    APIKeyHash      string     `gorm:"column:api_key_hash"`
    LastSeenAt      *time.Time `gorm:"column:last_seen_at"`
    LastSyncAt      *time.Time `gorm:"column:last_sync_at"`
    PendingQueue    int        `gorm:"column:pending_queue"`
    FirmwareVersion *string    `gorm:"column:firmware_version"`
    CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (Device) TableName() string { return "devices" }
