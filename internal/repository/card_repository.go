package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type ICardRepository interface {
  FindByUID(ctx context.Context, uid string) (*entity.RFIDCard, error)
  FindActiveByStudentID(ctx context.Context, studentID uuid.UUID) (*entity.RFIDCard, error)
  Create(ctx context.Context, card *entity.RFIDCard) error
  UpdateStatus(ctx context.Context, cardID uuid.UUID, status string) error

  ListUnregisteredUIDs(ctx context.Context, sinceDays int) ([]dto.UnregisteredCardSummary, error)
  ListActiveUIDs(ctx context.Context) ([]string, error)
}
