package dto

import "github.com/google/uuid"

type SearchStudentRequest struct {
  Query   string  `form:"q"`
  ClassID *uuid.UUID `form:"classId"`
  Status  *string `form:"status"`
  Page    int     `form:"page"`
  Limit   int     `form:"limit"`
}

type StudentDetailResponse struct {
  ID     uuid.UUID `json:"id"`
  NIS    string `json:"nis"`
  Name   string `json:"name"`
  Class  string `json:"class"`
  Status string `json:"status"`
}

type SearchStudentResponse struct {
  Data []StudentDetailResponse `json:"data"`
  Meta Pagination              `json:"meta"`
}

type CreateStudentRequest struct {
  NIS     string `json:"nis" binding:"required"`
  Name    string `json:"name" binding:"required"`
  ClassID uuid.UUID `json:"classId" binding:"required"`
}

type UpdateStudentRequest struct {
  Name    *string `json:"name"`
  ClassID *uuid.UUID `json:"classId"`
  Status  *string `json:"status"`
}

type ImportRowError struct {
  Row     int    `json:"row"`
  Message string `json:"message"`
}

type ImportResultResponse struct {
  SuccessCount int              `json:"successCount"`
  Failures     []ImportRowError `json:"failures"`
}

type ListClassResponse struct {
  Data []string `json:"data"`
}
