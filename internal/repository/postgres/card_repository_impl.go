package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
	"gorm.io/gorm"
)

type cardRepositoryImpl struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) repository.ICardRepository {
  return &cardRepositoryImpl{db: db}
}

func (cr *cardRepositoryImpl) FindByUID(ctx context.Context, uid string) (*entity.RFIDCard, error) {
  var card entity.RFIDCard

  err := cr.db.WithContext(ctx).Where("uid = ?", uid).First(&card).Error
  if err != nil {
    return nil, err
  }
  return &card, err
}

func (cr *cardRepositoryImpl)	FindActiveByStudentID(ctx context.Context, studentID uuid.UUID) (*entity.RFIDCard, error) {
  var card entity.RFIDCard

  err := cr.db.WithContext(ctx).Where("student_id = ?", studentID).First(&card).Error
  if err != nil {
    return nil, err
  }
  return &card, err
}

func (cr *cardRepositoryImpl)	Create(ctx context.Context, card *entity.RFIDCard) error {
  return cr.db.WithContext(ctx).Create(card).Error
}

func (cr *cardRepositoryImpl) UpdateStatus(ctx context.Context, cardID uuid.UUID, status string) error {
  return cr.db.WithContext(ctx).Model(entity.RFIDCard{}).Where("id = ?", cardID).Update("status", status).Error
} 

func (cr *cardRepositoryImpl) ListUnregisteredUIDs(ctx context.Context, sinceDays int) ([]dto.UnregisteredCardSummary, error) {
  var results []dto.UnregisteredCardSummary
  
	timeThreshold := time.Now().AddDate(0, 0, -sinceDays)
  
	err := cr.db.WithContext(ctx).
		Model(&entity.AttendanceEvent{}).
		Select(`
			card_uid AS uid, 
			COUNT(id) AS occurrences, 
			MAX(occurred_at) AS last_seen_at, 
			(ARRAY_AGG(device_id ORDER BY occurred_at DESC))[1] AS last_device_id
		`).
		Where("result = ? AND occurred_at >= ?", "UNKNOWN_CARD", timeThreshold).
		Group("card_uid").
		Order("last_seen_at DESC").
		Scan(&results).Error 
  
	if err != nil {
		return nil, err
	}

	return results, nil
}


func (cr *cardRepositoryImpl) ListActiveUIDs(ctx context.Context) ([]string, error) {
  var uids []string
  
	err := cr.db.WithContext(ctx).Model(&entity.RFIDCard{}).
		Joins("JOIN students ON students.id = rfid_cards.student_id").
		Where("rfid_cards.status = ? AND students.status = ?", "ACTIVE", "ACTIVE").
		Pluck("rfid_cards.uid", &uids).Error
  
	if err != nil {
		return nil, err
	}
  
	return uids, nil
}
