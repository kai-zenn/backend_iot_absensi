package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type studentRepoImpl struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) repository.IStudentRepository {
  return &studentRepoImpl{db: db}
}

func (sr *studentRepoImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Student, error) {
  var student entity.Student
  err := sr.db.WithContext(ctx).Where("id = ?", id).First(&student).Error
  if err != nil {
    return nil, err
  }

  return &student, nil
}

func (sr *studentRepoImpl) FindByNIS(ctx context.Context, nis string) (*entity.Student, error) {
  var student entity.Student
  err := sr.db.WithContext(ctx).Where("nis = ?", nis).Error
  if err != nil {
    return nil, err
  }

  return &student, nil
}

func (sr *studentRepoImpl) Search(ctx context.Context, query string, classID *uuid.UUID, status *string, limit, offset int) ([]entity.Student, int64, error) {
  var student []entity.Student
  var total int64

  db := sr.db.WithContext(ctx).Model(&entity.Student{})

  if query != "" {
    searchParam := "%" + query + "%"
    db = db.Where("name ILIKE ? OR nis ILIKE ?", searchParam, searchParam)
  }

  if classID != nil {
    db = db.Where("class_id = ?", *classID)
  }

  if status != nil {
    db = db.Where("status = ?", *status)
  }

  if err := db.Count(&total).Error; err != nil {
    return nil, 0, err
  }

  err := db.Limit(limit).Offset(offset).Find(&student).Error
  if err != nil {
    return nil, 0, err
  }

  return student, total, nil
}

func (sr *studentRepoImpl) CountActive(ctx context.Context, classID *uuid.UUID) (int, error) {
  var activeCount int64

  db := sr.db.WithContext(ctx).Model(&entity.Student{}).Where("status = ?", "ACTIVE")

  if classID != nil {
    db = db.Where("class_id = ?", *classID)
  }

  err := db.Count(&activeCount).Error
  if err != nil {
    return 0, err
  }

  return int(activeCount), err
}

func (sr *studentRepoImpl) Create(ctx context.Context, s *entity.Student) (id uuid.UUID, err error) {
  if err := sr.db.WithContext(ctx).Create(s).Error; err != nil {
    return uuid.Nil, err
  }

  return s.ID, nil
}

func (sr *studentRepoImpl) Update(ctx context.Context, s *entity.Student) error {
  var updates map[string]interface{}
  err := sr.db.WithContext(ctx).Model(&entity.Student{}).Where("id = ?", s.ID).Updates(updates).Error 
  if err != nil {
    return err
  }

  return nil
}

func (sr *studentRepoImpl) BulkUpsertFromImport(ctx context.Context, rows []entity.Student) (successCount int, failures []dto.ImportRowError, err error) {
  if len(rows) == 0 {
    return 0, nil, nil
  }

  if err := sr.db.WithContext(ctx).Clauses(clause.OnConflict{
    Columns: []clause.Column{{Name: "nis"}},
    DoUpdates: clause.AssignmentColumns([]string{"name", "class_id", "updated_at"}),
  }).CreateInBatches(&rows, 100).Error; err != nil {
    return 0, nil, err
  }

  return len(rows), nil, nil
}

type classRepositoryImpl struct {
  db *gorm.DB
}

func NewClassRepository(db *gorm.DB) repository.IClassRepository {
  return &classRepositoryImpl{db: db}
}

func (cr *classRepositoryImpl) List(ctx context.Context) ([]entity.Class, error) {
  var class []entity.Class

  err := cr.db.WithContext(ctx).Order("name DESC").Find(&class).Error
  if err != nil {
    return nil, err
  }

  return class, nil
}

func (cr *classRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Class, error) {
  var class entity.Class
  
	err := cr.db.WithContext(ctx).Where("name = ?", name).First(&class).Error
	if err != nil {
		return nil, err
	}
  
	return &class, nil
}


func (cr *classRepositoryImpl) Create(ctx context.Context, name string) (id uuid.UUID, err error) {
  newClass := entity.Class{
		Name: name,
	}
  
	if err := cr.db.WithContext(ctx).Create(&newClass).Error; err != nil {
		return uuid.Nil, err
	}
  
	return newClass.ID, nil
}
