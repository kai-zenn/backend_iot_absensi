package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterDeviceRequest struct {
  DeviceCode string `json:"deviceCode" binding:"required,max=30"`
  Name       string `json:"name" binding:"required"`
  Location   string `json:"location"`
}

// DeviceCreatedResponse — APIKey HANYA muncul satu kali di response ini,
// tidak pernah bisa diambil ulang (hanya hash SHA-256 yang disimpan).
type DeviceCreatedResponse struct {
  ID         uuid.UUID `json:"id"`
  DeviceCode string `json:"deviceCode"`
  APIKey     string `json:"apiKey"`
}

type HeartbeatRequest struct {
  DeviceID        uuid.UUID `json:"deviceId" binding:"required"`
  PendingQueue    int       `json:"pendingQueue"`
  FirmwareVersion string    `json:"firmwareVersion"`
  DeviceTime      time.Time `json:"deviceTime" binding:"required"`
  WifiRSSI        *int      `json:"wifiRssi,omitempty"`
}

type RuleSnapshot struct {
  OnTimeUntil    string `json:"onTimeUntil"`    // "06:30:00"
  TimezoneOffset string `json:"timezoneOffset"` // "+07:00"
}

type HeartbeatResponse struct {
  ServerTime       time.Time    `json:"serverTime"`
  ClockSkewSeconds int          `json:"clockSkewSeconds"`
  CardsVersion     string       `json:"cardsVersion"`
  Rules            RuleSnapshot `json:"rules"`
}

type CardSnapshotResponse struct {
  Version string   `json:"version"`
  UIDs    []string `json:"uids"`
}

type DeviceStatusItem struct {
  ID              uuid.UUID  `json:"id"`
  DeviceCode      string     `json:"deviceCode"`
  Name            string     `json:"name"`
  Status          string     `json:"status"`
  ConnectionState string     `json:"connectionState"` // ONLINE | DISCONNECTED (turunan, lihat 4c)
  LastSeenAt      *time.Time `json:"lastSeenAt,omitempty"`
  PendingQueue    int        `json:"pendingQueue"`
}

type ListDevicesResponse struct {
  Data []DeviceStatusItem `json:"data"`
}
