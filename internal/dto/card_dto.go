package dto

import "github.com/google/uuid"

type AssignCardRequest struct {
  UID             string `json:"uid" binding:"required,rfiduid"`
  StudentID       uuid.UUID `json:"studentId" binding:"required"`
  ReplaceExisting bool   `json:"replaceExisting"`
}

type UnregisteredCardSummary struct {
  UID          string `json:"uid"`
  Occurrences  int    `json:"occurrences"`
  LastSeenAt   string `json:"lastSeenAt"`
  LastDeviceID uuid.UUID `json:"lastDeviceId"`
}

type ListUnregisteredCardsResponse struct {
  Data []UnregisteredCardSummary `json:"data"`
}
