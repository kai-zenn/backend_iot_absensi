package usecase

import (
	"context"
	"errors"
	"sort"
	"time"
	"strings"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/pkgs/helper"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
)

type IDeviceUsecase interface {
  RegisterDevice(ctx context.Context, req dto.RegisterDeviceRequest) (*dto.DeviceCreatedResponse, error)
  RotateAPIKey(ctx context.Context, deviceID uuid.UUID) (*dto.DeviceCreatedResponse, error)
  UpdateStatus(ctx context.Context, deviceID uuid.UUID, status string) error
  ListDevices(ctx context.Context) (*dto.ListDevicesResponse, error)
  GetCardSnapshot(ctx context.Context) (*dto.CardSnapshotResponse, error)
  Heartbeat(ctx context.Context, deviceID uuid.UUID, req dto.HeartbeatRequest) (*dto.HeartbeatResponse, error)
}

type deviceUsecase struct {
  deviceRepo repository.IDeviceRepository
  cardRepo   repository.ICardRepository
	ruleRepo   repository.IRuleRepository
}

func NewDeviceUsecase(dr repository.IDeviceRepository) IDeviceUsecase {
  return &deviceUsecase{
    deviceRepo: dr,
  }
}

func (du *deviceUsecase) RegisterDevice(ctx context.Context, req dto.RegisterDeviceRequest) (*dto.DeviceCreatedResponse, error) {
  existing, err := du.deviceRepo.FindByCode(ctx, req.DeviceCode)
  if existing != nil && err == nil {
    return nil, errors.New("Device sudah terdaftar")
  }

  rawApiKey, err := helper.GenerateRandomAPIKey(32)
  if err != nil {
    return nil, err
  }
  hashedApiKey := helper.HashAPIKey(rawApiKey)

  var loc *string
  if req.Location != "" {
    loc = &req.Location
  }
  
  newDevice := &entity.Device{
    DeviceCode: req.DeviceCode,
    Name: req.Name,
    Location: loc,
    Role: "CHECK-IN",
    Status: "ACTIVE",
    APIKeyHash: hashedApiKey,
  }

  deviceId, err := du.deviceRepo.Create(ctx, newDevice)
  if err != nil {
    return nil, err
  }

  return &dto.DeviceCreatedResponse{
    ID: deviceId,
    DeviceCode: newDevice.DeviceCode,
    APIKey: hashedApiKey,
  }, nil
}

func (du *deviceUsecase) RotateAPIKey(ctx context.Context, deviceID uuid.UUID) (*dto.DeviceCreatedResponse, error) {
  device, err := du.deviceRepo.FindByID(ctx, deviceID)
  if err != nil {
    return nil, err
  }

  rawApiKey, err := helper.GenerateRandomAPIKey(32)
  if err != nil {
    return nil, err
  }
  
  hashedApiKey := helper.HashAPIKey(rawApiKey)

  if err := du.deviceRepo.UpdateAPIKeyHash(ctx, deviceID, hashedApiKey); err != nil {
    return nil, err
  }

 	return &dto.DeviceCreatedResponse{
		ID: device.ID,
		DeviceCode: device.DeviceCode,
		APIKey: rawApiKey,
	}, nil
}

func (du *deviceUsecase) UpdateStatus(ctx context.Context, deviceID uuid.UUID, status string) error {
 	if status != "ACTIVE" && status != "DISABLED" {
		return errors.New("invalid status: must be ACTIVE or DISABLED")
	}
 
	err := du.deviceRepo.UpdateStatus(ctx, deviceID, status)
	if err != nil {
		return err
	}
 
	return nil
}

func (du *deviceUsecase) ListDevices(ctx context.Context) (*dto.ListDevicesResponse, error) {
  devices, err := du.deviceRepo.List(ctx)
  if err != nil {
    return nil, err
  }
  
  deviceItems := make([]dto.DeviceStatusItem, 0, len(devices))

  for _, d := range devices {
    connState := "DISCONNECTED"
    if d.LastSeenAt != nil && time.Since(*d.LastSeenAt) <= 2*time.Minute {
      connState = "ONLINE"
    }
    item := dto.DeviceStatusItem{
      ID: d.ID,
      DeviceCode: d.DeviceCode,
      Name: d.Name,
      Status: d.Status,
      ConnectionState: connState,
      LastSeenAt: d.LastSeenAt,
      PendingQueue: d.PendingQueue,
    }
    deviceItems = append(deviceItems, item)
  }

  return &dto.ListDevicesResponse{
    Data: deviceItems,
  }, nil
}


func (du *deviceUsecase) GetCardSnapshot(ctx context.Context) (*dto.CardSnapshotResponse, error) {
  uids, err := du.cardRepo.ListActiveUIDs(ctx)
  if err != nil {
	return nil, err
  }
  if uids == nil {
	uids = []string{}
  }

  sort.Strings(uids)

  joined := strings.Join(uids, ",")
  fullHash := helper.HashAPIKey(joined)
  version := fullHash[:8]

  return &dto.CardSnapshotResponse{
	Version: version,
	UIDs:    uids,
  }, nil
}


func (du *deviceUsecase) Heartbeat(ctx context.Context, deviceID uuid.UUID, req dto.HeartbeatRequest) (*dto.HeartbeatResponse, error) {
  if req.DeviceID != uuid.Nil && req.DeviceID != deviceID {
    return nil, errors.New("Device ID tidak cocok")
  }

  now := helper.NowWIB()

  clockSkew := int(now.Sub(req.DeviceTime).Seconds())
  
  err := du.deviceRepo.UpdateHeartbeat(ctx, deviceID, req.PendingQueue, req.FirmwareVersion, now)
	if err != nil {
		return nil, err
	}

	snapShot, err := du.GetCardSnapshot(ctx)

	rule, err := du.ruleRepo.FindEffectiveForDate(ctx, now) 
	onTimeUntil := "06:30:00"
	if err == nil && rule != nil {
    onTimeUntil = rule.OnTimeUntil.Format("15:04:05")
	}

	return &dto.HeartbeatResponse{
		ServerTime: now,
		ClockSkewSeconds: clockSkew,
		CardsVersion: snapShot.Version,
		Rules: dto.RuleSnapshot{OnTimeUntil: onTimeUntil, TimezoneOffset: "+07:00"},
	}, nil
}
