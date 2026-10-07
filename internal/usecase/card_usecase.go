package usecase

import (
  "context"
  "github.com/kai-zenn/backend_iot_absensi/internal/dto"
)

type ICardUsecase interface {
  AssignCard(ctx context.Context, req dto.AssignCardRequest) error
  UpdateCardStatus(ctx context.Context, cardID uint64, status string) error
  ListUnregistered(ctx context.Context) (dto.ListUnregisteredCardsResponse, error)
}
