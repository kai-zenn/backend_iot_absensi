package repository

import (
  "context"
  "time"
  "github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type IRuleRepository interface {
    // FindEffectiveForDate: db.Where("effective_from <= ?", date).
    //   Order("effective_from DESC").First(&rule)
    // Selalu ada hasil karena migrasi awal menyisipkan seed 1970-01-01 (lihat 2.6).
  FindEffectiveForDate(ctx context.Context, date time.Time) (*entity.AttendanceRule, error)

    // Upsert: db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name:
    //   "effective_from"}}, DoUpdates: clause.AssignmentColumns([]string{"on_time_until"})}).
    //   Create(rule)
  Upsert(ctx context.Context, rule *entity.AttendanceRule) error

  List(ctx context.Context) ([]entity.AttendanceRule, error)
}
