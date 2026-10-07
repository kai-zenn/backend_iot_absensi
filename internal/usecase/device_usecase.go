package usecase

import (
    "context"
    "github.com/kai-zenn/backend_iot_absensi/internal/dto"
)

type IDeviceUsecase interface {
  RegisterDevice(ctx context.Context, req dto.RegisterDeviceRequest) (dto.DeviceCreatedResponse, error)
  RotateAPIKey(ctx context.Context, deviceID uint64) (dto.DeviceCreatedResponse, error)
  UpdateStatus(ctx context.Context, deviceID uint64, status string) error
  ListDevices(ctx context.Context) (dto.ListDevicesResponse, error)
  Heartbeat(ctx context.Context, deviceID uint64, req dto.HeartbeatRequest) (dto.HeartbeatResponse, error)
  GetCardSnapshot(ctx context.Context) (dto.CardSnapshotResponse, error)
}
