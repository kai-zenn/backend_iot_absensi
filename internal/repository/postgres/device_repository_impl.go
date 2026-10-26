package postgres

import (
 	"context"
	"time"
 
	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
	"gorm.io/gorm"
)

type deviceRepositoryImpl struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) repository.IDeviceRepository {
  return &deviceRepositoryImpl{db: db}
}

func (dr *deviceRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
  var gawai entity.Device
  err := dr.db.WithContext(ctx).Where("id = ?", id).First(&gawai).Error
  if err != nil {
    return nil, err
  }

  return &gawai, nil
}

func (dr *deviceRepositoryImpl) FindByAPIKeyHash(ctx context.Context, hash string) (*entity.Device, error) {
  var gawai entity.Device
	err := dr.db.WithContext(ctx).Where("api_key_hash = ?", hash).First(&gawai).Error
	if err != nil {
		return nil, err
	}
	return &gawai, nil
}

func (dr *deviceRepositoryImpl) FindByCode(ctx context.Context, code string) (*entity.Device, error) {
  var gawai entity.Device
	err := dr.db.WithContext(ctx).Where("device_code = ?", code).First(&gawai).Error
	if err != nil {
		return nil, err
	}
	return &gawai, nil
}

func (dr *deviceRepositoryImpl) List(ctx context.Context) ([]entity.Device, error) {
  var gawai []entity.Device

  err := dr.db.WithContext(ctx).Order("device_code DESC").Find(&gawai).Error
  if err != nil {
    return nil, err
  }

  return gawai, nil
}

func (dr *deviceRepositoryImpl) Create(ctx context.Context, d *entity.Device) (id uuid.UUID, err error) {
  if err := dr.db.WithContext(ctx).Create(d).Error; err != nil {
		return uuid.Nil, err
	}
	
	return d.ID, nil
}

func (dr *deviceRepositoryImpl) UpdateAPIKeyHash(ctx context.Context, deviceID uuid.UUID, newHash string) error {
  return dr.db.WithContext(ctx).Model(&entity.Device{}).Where("id = ?", deviceID).Update("api_key_hash", newHash).Error
}

func (dr *deviceRepositoryImpl) UpdateStatus(ctx context.Context, deviceID uuid.UUID, status string) error {
  return dr.db.WithContext(ctx).Model(&entity.Device{}).Where("id = ?", deviceID).Update("status", status).Error
}

func (dr *deviceRepositoryImpl) UpdateHeartbeat(ctx context.Context, deviceID uuid.UUID, pendingQueue int, firmwareVersion string, seenAt time.Time) error {
  return dr.db.WithContext(ctx).Model(&entity.Device{}).Where("id = ?", deviceID).
		Updates(map[string]any{
			"pending_queue":    pendingQueue,
			"firmware_version": firmwareVersion,
			"last_seen_at":     seenAt,
		}).Error
}
