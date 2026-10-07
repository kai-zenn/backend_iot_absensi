package repository

import (
  "context"
  "github.com/kai-zenn/backend_iot_absensi/internal/dto"
  "github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type IStudentRepository interface {
  FindByID(ctx context.Context, id uint64) (*entity.Student, error)
  FindByNIS(ctx context.Context, nis string) (*entity.Student, error)
  Search(ctx context.Context, query string, classID *uint64, status *string, limit, offset int) ([]entity.Student, int64, error)
  CountActive(ctx context.Context, classID *uint64) (int, error)

  Create(ctx context.Context, s *entity.Student) (id uint64, err error)
  Update(ctx context.Context, s *entity.Student) error

    // BulkUpsertFromImport: db.Clauses(clause.OnConflict{
    //   Columns: []clause.Column{{Name: "nis"}}, DoUpdates: clause.AssignmentColumns(
    //   []string{"name", "class_id", "updated_at"})}).CreateInBatches(rows, 100)
  BulkUpsertFromImport(ctx context.Context, rows []entity.Student) (successCount int, failures []dto.ImportRowError, err error)
}

type IClassRepository interface {
  List(ctx context.Context) ([]entity.Class, error)
  FindByName(ctx context.Context, name string) (*entity.Class, error)
  Create(ctx context.Context, name string) (id uint64, err error)
}
