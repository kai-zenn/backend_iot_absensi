package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type IAttendanceEventRepository interface {
    // TryInsert diimplementasikan dengan:
    //   db.Clauses(clause.OnConflict{
    //       Columns:   []clause.Column{{Name: "device_id"}, {Name: "event_id"}},
    //       DoNothing: true,
    //   }).Create(&event)
    // Cek result.RowsAffected: 0 berarti conflict (alreadyExisted=true) — lanjutkan
    // dengan db.Where("device_id = ? AND event_id = ?", ...).First(&existing) untuk
    // mengambil row LAMA, dipakai membandingkan cardUid/occurredAt (deteksi
    // EVENT_ID_CONFLICT, lihat 4a).
    TryInsert(ctx context.Context, event *entity.AttendanceEvent) (*entity.AttendanceEvent, bool, error)

    UpdateResult(ctx context.Context, eventID uuid.UUID, result string, studentID *uuid.UUID, attendanceID *uuid.UUID, clockSkewSeconds *int) error

    // ListNeedsReview memanfaatkan idx_attendance_events_needs_review.
    ListNeedsReview(ctx context.Context, limit, offset int) ([]entity.AttendanceEvent, int64, error)

    ListByFilter(ctx context.Context, from, to time.Time, result *string, limit, offset int) ([]entity.AttendanceEvent, int64, error)
}

type IAttendanceRepository interface {
    CreateCheckIn(ctx context.Context, a *entity.Attendance) (uuid.UUID, bool, error)
    FindByStudentAndDateForUpdate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.Attendance, error)
    CorrectCheckIn(ctx context.Context, attendanceID uuid.UUID, checkInAt time.Time, status string, eventID uuid.UUID) error

    CreateManual(ctx context.Context, a *entity.Attendance) (uuid.UUID, bool, error)
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error)
    Summary(ctx context.Context, date time.Time, classID *uuid.UUID) (present, late, activeStudents int, err error)

    ListByFilter(ctx context.Context, date time.Time, classID *uuid.UUID, status *string, searchQuery string, limit, offset int) ([]entity.Attendance, int64, error)
    ListHistoryByStudent(ctx context.Context, studentID uuid.UUID, from, to time.Time) ([]entity.Attendance, error)
}
