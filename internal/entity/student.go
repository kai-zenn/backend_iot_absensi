package entity

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
    ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
    NIS       string    `gorm:"column:nis"`
    Name      string    `gorm:"column:name"`
    ClassID   uuid.UUID `gorm:"column:class_id"`
    Status    string    `gorm:"column:status"` // ACTIVE | INACTIVE
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

    Class *Class `gorm:"foreignKey:ClassID"`
}

func (Student) TableName() string { return "students" }
