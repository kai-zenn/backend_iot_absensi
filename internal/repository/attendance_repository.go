package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
)


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
