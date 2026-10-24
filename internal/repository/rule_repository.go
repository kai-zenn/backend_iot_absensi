package repository

import (
  "context"
  "time"
  "github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type IRuleRepository interface {
  FindEffectiveForDate(ctx context.Context, date time.Time) (*entity.AttendanceRule, error)
  Upsert(ctx context.Context, rule *entity.AttendanceRule) error
  List(ctx context.Context) ([]entity.AttendanceRule, error)
}
