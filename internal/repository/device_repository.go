package repository

import (
  "context"
  "time"
  "github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type IDeviceRepository interface {
  FindByID(ctx context.Context, id uint64) (*entity.Device, error)
  FindByAPIKeyHash(ctx context.Context, hash string) (*entity.Device, error)
  FindByCode(ctx context.Context, code string) (*entity.Device, error)
  List(ctx context.Context) ([]entity.Device, error)

  Create(ctx context.Context, d *entity.Device) (id uint64, err error)
  UpdateAPIKeyHash(ctx context.Context, deviceID uint64, newHash string) error
  UpdateStatus(ctx context.Context, deviceID uint64, status string) error

    // UpdateHeartbeat: satu db.Model(&entity.Device{}).Where("id = ?", deviceID).
    //   Updates(map[string]any{...}) — single UPDATE, bukan tiga panggilan terpisah.
  UpdateHeartbeat(ctx context.Context, deviceID uint64, pendingQueue int, firmwareVersion string, seenAt time.Time) error
}
