package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/pkgs/helper"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
)

type ICardUsecase interface {
  ListUnregistered(ctx context.Context) (*dto.ListUnregisteredCardsResponse, error)
  AssignCard(ctx context.Context, req dto.AssignCardRequest) error
  UpdateCardStatus(ctx context.Context, cardID uuid.UUID, status string) error
}

type cardUsecase struct {
  cardRepo repository.ICardRepository
  txManager repository.TxManager
}

func NewCardUsecase(cr repository.ICardRepository, tx repository.TxManager) ICardUsecase {
  return &cardUsecase{
    cardRepo: cr,
    txManager: tx,
  }
}

func (cu *cardUsecase) ListUnregistered(ctx context.Context) (*dto.ListUnregisteredCardsResponse, error) {
  item, err := cu.cardRepo.ListUnregisteredUIDs(ctx, 7)
  if err != nil {
    return nil, errors.New("")
  }

  if item == nil {
    item = []dto.UnregisteredCardSummary{} 
  }

  return &dto.ListUnregisteredCardsResponse{
    Data: item,
  }, nil
}

func (cu *cardUsecase) AssignCard(ctx context.Context, req dto.AssignCardRequest) error {
  normalizedUID, err := helper.ParseUID(req.UID)
	if err != nil {
		return err
	}
	
  return cu.txManager.WithinTx(ctx, func(repos repository.RepoSet) error {
    existingCard, err := repos.Card.FindActiveByStudentID(ctx, req.StudentID)

    if err == nil && existingCard != nil {
      if !req.ReplaceExisting {
        return errors.New("siswa sudah memiliki kartu aktif")
      }
      if err := repos.Card.UpdateStatus(ctx, existingCard.ID, "RETIRED"); err != nil {
        return err
      }
    }

    newCard := &entity.RFIDCard{
			UID:       normalizedUID,
			StudentID: req.StudentID,
			Status:    "ACTIVE",
		}
    
		return repos.Card.Create(ctx, newCard)
  },)
}

func (cu *cardUsecase) UpdateCardStatus(ctx context.Context, cardID uuid.UUID, status string) error {
  if status != "ACTIVE" && status != "BLOCKED" && status != "RETIRED" {
		return errors.New("status kartu harus ACTIVE, BLOCKED, atau RETIRED")
	}
  
	return cu.cardRepo.UpdateStatus(ctx, cardID, status)
}
