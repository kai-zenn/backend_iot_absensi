package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type IDeviceRepository interface {
  FindByID(ctx context.Context, id uuid.UUID) (*entity.Device, error)
  FindByAPIKeyHash(ctx context.Context, hash string) (*entity.Device, error)
  FindByCode(ctx context.Context, code string) (*entity.Device, error)
  List(ctx context.Context) ([]entity.Device, error)

  Create(ctx context.Context, d *entity.Device) (id uuid.UUID, err error)
  UpdateAPIKeyHash(ctx context.Context, deviceID uuid.UUID, newHash string) error
  UpdateStatus(ctx context.Context, deviceID uuid.UUID, status string) error

  UpdateHeartbeat(ctx context.Context, deviceID uuid.UUID, pendingQueue int, firmwareVersion string, seenAt time.Time) error
}
