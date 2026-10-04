package entity

import (
	"time"

	"github.com/google/uuid"
)

type Class struct {
  ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
  Name      string    `gorm:"column:name"`
  CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Class) TableName() string { return "classes" }
