package entity

import (
	"time"

	"github.com/google/uuid"
)

type RFIDCard struct {
    ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
    UID       string    `gorm:"column:uid"` // sudah dinormalisasi (hex uppercase)
    StudentID uuid.UUID `gorm:"column:student_id"`
    Status    string    `gorm:"column:status"` // ACTIVE | BLOCKED | RETIRED
    Note      *string   `gorm:"column:note"`
    IssuedAt  time.Time `gorm:"column:issued_at"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

    Student *Student `gorm:"foreignKey:StudentID"`
}

func (RFIDCard) TableName() string { return "rfid_cards" }
