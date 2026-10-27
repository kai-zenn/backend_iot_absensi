package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type IStudentRepository interface {
  FindByID(ctx context.Context, id uuid.UUID) (*entity.Student, error)
  FindByNIS(ctx context.Context, nis string) (*entity.Student, error)
  Search(ctx context.Context, query string, classID *uuid.UUID, status *string, limit, offset int) ([]entity.Student, int64, error)
  CountActive(ctx context.Context, classID *uuid.UUID) (int, error)

  Create(ctx context.Context, s *entity.Student) (uuid.UUID, error)
  Update(ctx context.Context, s *entity.Student) error

  BulkUpsertFromImport(ctx context.Context, rows []entity.Student) (int, []dto.ImportRowError, error)
}

type IClassRepository interface {
  List(ctx context.Context) ([]entity.Class, error)
  FindByName(ctx context.Context, name string) (*entity.Class, error)
  Create(ctx context.Context, name string) (uuid.UUID, error)
}
