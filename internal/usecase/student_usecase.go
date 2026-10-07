package usecase

import (
  "context"
  "github.com/kai-zenn/backend_iot_absensi/internal/dto"
)

type IStudentUsecase interface {
  Search(ctx context.Context, req dto.SearchStudentRequest) (dto.SearchStudentResponse, error)
  Create(ctx context.Context, req dto.CreateStudentRequest) (dto.StudentDetailResponse, error)
  Update(ctx context.Context, id uint64, req dto.UpdateStudentRequest) error
  ImportCSV(ctx context.Context, rows [][]string) (dto.ImportResultResponse, error)
  ListClasses(ctx context.Context) (dto.ListClassResponse, error)
}
