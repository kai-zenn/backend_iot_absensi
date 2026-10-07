package repository
import (
  "context"
  "time"
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
    TryInsert(ctx context.Context, event *entity.AttendanceEvent) (row *entity.AttendanceEvent, alreadyExisted bool, err error)

    UpdateResult(ctx context.Context, eventID uint64, result string, studentID *uint64, attendanceID *uint64, clockSkewSeconds *int) error

    // ListNeedsReview memanfaatkan idx_attendance_events_needs_review.
    ListNeedsReview(ctx context.Context, limit, offset int) ([]entity.AttendanceEvent, int64, error)

    ListByFilter(ctx context.Context, from, to time.Time, result *string, limit, offset int) ([]entity.AttendanceEvent, int64, error)
}

type IAttendanceRepository interface {
    // CreateCheckIn diimplementasikan dengan:
    //   db.Clauses(clause.OnConflict{
    //       Columns:   []clause.Column{{Name: "student_id"}, {Name: "attendance_date"}},
    //       DoNothing: true,
    //   }).Create(a)
    // inserted=false berarti sudah ada attendance untuk siswa+tanggal ini — caller
    // wajib lanjut ke FindByStudentAndDateForUpdate (earliest-wins, lihat 4a).
    CreateCheckIn(ctx context.Context, a *entity.Attendance) (id uint64, inserted bool, err error)

    // FindByStudentAndDateForUpdate: db.Clauses(clause.Locking{Strength: "UPDATE"}).
    //   Where("student_id = ? AND attendance_date = ?", studentID, date).First(&att)
    // WAJIB dipanggil di dalam transaksi aktif — mengunci baris ini sampai transaksi
    // selesai, mencegah dua koreksi berjalan bersamaan untuk siswa+tanggal yang sama.
    FindByStudentAndDateForUpdate(ctx context.Context, studentID uint64, date time.Time) (*entity.Attendance, error)

    // CorrectCheckIn menimpa check_in_at/check_in_status/check_in_event_id pada
    // baris yang SUDAH dikunci lewat FindByStudentAndDateForUpdate di tx yang sama.
    CorrectCheckIn(ctx context.Context, attendanceID uint64, checkInAt time.Time, status string, eventID uint64) error

    CreateManual(ctx context.Context, a *entity.Attendance) (id uint64, inserted bool, err error)

    FindByID(ctx context.Context, id uint64) (*entity.Attendance, error)

    Summary(ctx context.Context, date time.Time, classID *uint64) (present, late, activeStudents int, err error)

    // ListByFilter memakai Preload("Student.Class") supaya nama & kelas siswa
    // ikut termuat dalam satu query (hindari N+1 saat menampilkan daftar di dashboard).
    ListByFilter(ctx context.Context, date time.Time, classID *uint64, status *string, searchQuery string, limit, offset int) ([]entity.Attendance, int64, error)

    ListHistoryByStudent(ctx context.Context, studentID uint64, from, to time.Time) ([]entity.Attendance, error)
}
