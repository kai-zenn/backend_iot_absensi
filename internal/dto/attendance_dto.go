package dto

import (
	"time"

	"github.com/google/uuid"
)

type ScanEventRequest struct {
  EventID    string    `json:"eventId" binding:"required,max=64"`
  CardUID    string    `json:"cardUid" binding:"required,rfiduid"` // custom rule, lihat §1.3
  DeviceID   uuid.UUID `json:"deviceId" binding:"required"`
  OccurredAt time.Time `json:"occurredAt" binding:"required"`
}

type StudentBrief struct {
  NIS   string `json:"nis"`
  Name  string `json:"name"`
  Class string `json:"class"`
}

type ScanEventResponse struct {
  EventID    string        `json:"eventId"`
  Result     string        `json:"result"`
  Status     string        `json:"status,omitempty"`
  Replayed   bool          `json:"replayed"`
  Message    string        `json:"message,omitempty"`
  Student    *StudentBrief `json:"student,omitempty"`
  OccurredAt time.Time     `json:"occurredAt"`
  ServerTime time.Time     `json:"serverTime"`
}

type SyncEventItem struct {
  EventID    string    `json:"eventId" binding:"required,max=64"`
  CardUID    string    `json:"cardUid" binding:"required,rfiduid"`
  OccurredAt time.Time `json:"occurredAt" binding:"required"`
}

type SyncBatchRequest struct {
  DeviceID uuid.UUID       `json:"deviceId" binding:"required"`
  Events   []SyncEventItem `json:"events" binding:"required,max=50,dive"`
}

type SyncEventResult struct {
  EventID  string `json:"eventId"`
  Result   string `json:"result"`
  Status   string `json:"status,omitempty"`
  Replayed bool   `json:"replayed"`
  Message  string `json:"message,omitempty"`
}

type SyncBatchResponse struct {
  Results    []SyncEventResult `json:"results"`
  ServerTime time.Time         `json:"serverTime"`
}

type ManualCheckInRequest struct {
  StudentID uuid.UUID `json:"studentId" binding:"required"`
  Note      string `json:"note" binding:"required"`
}

type AttendanceDetailResponse struct {
  ID            uuid.UUID    `json:"id"`
  StudentID     uuid.UUID    `json:"studentId"`
  Student       StudentBrief `json:"student"`
  CheckInAt     time.Time    `json:"checkInAt"`
  CheckInStatus string       `json:"checkInStatus"`
  CheckInSource string       `json:"checkInSource"`
  Note          string       `json:"note,omitempty"`
}

type AttendanceSummaryResponse struct {
  Date           string `json:"date"`
  ActiveStudents int    `json:"activeStudents"`
  Present        int    `json:"present"`
  Late           int    `json:"late"`
  NotYetArrived  int    `json:"notYetArrived"`
  PendingSync    int    `json:"pendingSync"`
  NeedsReview    int    `json:"needsReview"`
}

type ListAttendanceRequest struct {
  Date    string  `form:"date"`
  ClassID *uuid.UUID `form:"classId"`
  Status  *string `form:"status"`
  Query   string  `form:"q"`
  Page    int     `form:"page"`
  Limit   int     `form:"limit"`
}

type ListAttendanceResponse struct {
  Data []AttendanceDetailResponse `json:"data"`
  Meta Pagination                 `json:"meta"`
}

type StudentHistoryResponse struct {
  StudentID uuid.UUID                  `json:"studentId"`
  Data      []AttendanceDetailResponse `json:"data"`
}

type EventLogItem struct {
  ID         uuid.UUID `json:"id"`
  EventID    string    `json:"eventId"`
  DeviceID   uuid.UUID `json:"deviceId"`
  CardUID    string    `json:"cardUid"`
  Result     string    `json:"result"`
  OccurredAt time.Time `json:"occurredAt"`
  ReceivedAt time.Time `json:"receivedAt"`
}

type ListEventsResponse struct {
  Data []EventLogItem `json:"data"`
  Meta Pagination     `json:"meta"`
}
