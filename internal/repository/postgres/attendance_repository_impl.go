package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)



  
type attendanceRepositoryImpl struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) repository.IAttendanceRepository {
  return &attendanceRepositoryImpl{db: db}
}

func (ar *attendanceRepositoryImpl) CreateCheckIn(ctx context.Context, a *entity.Attendance) (uuid.UUID, bool, error) {
  result := ar.db.WithContext(ctx).Clauses(clause.OnConflict{
    Columns: []clause.Column{{Name: "student_id"}, {Name: "attendance_date"}},
    DoNothing: true,
  }).Create(a)

  if result.Error != nil {
    return uuid.Nil, false, result.Error
  }

  inserted := result.RowsAffected > 0
  return a.ID, inserted, nil
}

func (ar *attendanceRepositoryImpl)	FindByStudentAndDateForUpdate(ctx context.Context, studentID uuid.UUID, date time.Time) (*entity.Attendance, error) {
  var att entity.Attendance

  err := ar.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
    Where("student_id = ? AND attendance_date = ?", studentID, date).
    Find(&att).Error

  if err != nil {
    return nil, err
  }

  return &att, nil
}

func (ar *attendanceRepositoryImpl) CorrectCheckIn(ctx context.Context, attendanceID uuid.UUID, checkInAt time.Time, status string, eventID uuid.UUID) error {
  return ar.db.WithContext(ctx).
		Model(&entity.Attendance{}).
		Where("id = ?", attendanceID).
		Updates(map[string]any{
			"check_in_at":       checkInAt,
			"check_in_status":   status,
			"check_in_event_id": eventID,
		}).Error
}

func (ar *attendanceRepositoryImpl) CreateManual(ctx context.Context, a *entity.Attendance) (uuid.UUID, bool, error) {
  result := ar.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "student_id"}, {Name: "attendance_date"}},
		DoNothing: true,
	}).Create(a)
  
	if result.Error != nil {
		return uuid.Nil, false, result.Error
	}
  
	inserted := result.RowsAffected > 0
	return a.ID, inserted, nil
}

func (ar *attendanceRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error) {
  var att entity.Attendance
	err := ar.db.WithContext(ctx).Where("id = ?", id).First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (ar *attendanceRepositoryImpl) Summary(ctx context.Context, date time.Time, classID *uuid.UUID) (present, late, activeStudents int, err error) {
  var presentCount, lateCount, activeStudentsCount int64
  
	sq := ar.db.WithContext(ctx).Model(&entity.Student{}).Where("status = ?", "ACTIVE")
	if classID != nil {
		sq = sq.Where("class_id = ?", *classID)
	}
	if err := sq.Count(&activeStudentsCount).Error; err != nil {
		return 0, 0, 0, err
	}
  
	aq := ar.db.WithContext(ctx).Model(&entity.Attendance{}).
		Joins("JOIN students ON students.id = attendance.student_id").
		Where("attendance.attendance_date = ?", date)
	
	if classID != nil {
		aq = aq.Where("students.class_id = ?", *classID)
	}
  
	// Hitung yang PRESENT
	if err := aq.Session(&gorm.Session{}).Where("attendance.check_in_status = ?", "PRESENT").Count(&presentCount).Error; err != nil {
		return 0, 0, 0, err
	}
  
	// Hitung yang LATE
	if err := aq.Session(&gorm.Session{}).Where("attendance.check_in_status = ?", "LATE").Count(&lateCount).Error; err != nil {
		return 0, 0, 0, err
	}
  
	return int(presentCount), int(lateCount), int(activeStudentsCount), nil
}

func (ar *attendanceRepositoryImpl)	ListByFilter(ctx context.Context, date time.Time, classID *uuid.UUID, status *string, searchQuery string, limit, offset int) ([]entity.Attendance, int64, error) {
  var attendances []entity.Attendance
	var totalData int64
  
	q := ar.db.WithContext(ctx).Model(&entity.Attendance{}).
		Joins("JOIN students ON students.id = attendance.student_id").
		Preload("Student").
		Preload("Student.Class").
		Where("attendance.attendance_date = ?", date)
  
	if classID != nil {
		q = q.Where("students.class_id = ?", *classID)
	}
	if status != nil {
		q = q.Where("attendance.check_in_status = ?", *status)
	}
	if searchQuery != "" {
		search := "%" + searchQuery + "%"
		q = q.Where("students.name ILIKE ? OR students.nis ILIKE ?", search, search)
	}
  
	if err := q.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}
  
	if err := q.Limit(limit).Offset(offset).Order("attendance.check_in_at DESC").Find(&attendances).Error; err != nil {
		return nil, 0, err
	}
  
	return attendances, totalData, nil
}

func (ar *attendanceRepositoryImpl) ListHistoryByStudent(ctx context.Context, studentID uuid.UUID, from, to time.Time) ([]entity.Attendance, error) {
  var attendances []entity.Attendance
  
	err := ar.db.WithContext(ctx).
		Where("student_id = ? AND attendance_date >= ? AND attendance_date <= ?", studentID, from, to).
		Order("attendance_date DESC").
		Find(&attendances).Error
  
	if err != nil {
		return nil, err
	}
  
	return attendances, nil
}
