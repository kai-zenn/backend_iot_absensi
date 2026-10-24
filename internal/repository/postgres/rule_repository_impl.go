package postgres

import (
	"context"
	"time"

	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ruleRepositoryImpl struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) repository.IRuleRepository {
  return &ruleRepositoryImpl{db: db}
}

func (rs *ruleRepositoryImpl) FindEffectiveForDate(ctx context.Context, date time.Time) (*entity.AttendanceRule, error) {
  var attendanceRule entity.AttendanceRule
  err := rs.db.WithContext(ctx).Where("effective_from <= ?", date).Order("effective_from DESC").First(&attendanceRule).Error
  if err != nil {
    return nil, err
  }

  return &attendanceRule, nil
}

func (rs *ruleRepositoryImpl)	Upsert(ctx context.Context, rule *entity.AttendanceRule) error {
  return rs.db.WithContext(ctx).Clauses(clause.OnConflict{
    Columns: []clause.Column{{Name: "effective_from"}},
    DoUpdates: clause.AssignmentColumns([]string{"on_time_until"}),
  }).Create(rule).Error
}

func (rs *ruleRepositoryImpl) List(ctx context.Context) ([]entity.AttendanceRule, error) {
  var attendanceRule []entity.AttendanceRule

  err := rs.db.WithContext(ctx).Order("effective_from DESC").Find(&attendanceRule).Error
	if err != nil {
		return nil, err
	}
  
	return attendanceRule, nil
}
