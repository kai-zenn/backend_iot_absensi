package postgres

import (
	"context"
	"gorm.io/gorm"
	
	"github.com/kai-zenn/backend_iot_absensi/internal/repository" // Sesuaikan dengan nama module lu
)

type txManagerImpl struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) repository.TxManager {
	return &txManagerImpl{db: db}
}

func (t *txManagerImpl) WithinTx(ctx context.Context, fn func(repos repository.RepoSet) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		
		repos := repository.RepoSet{
			AttendanceEvent: NewAttendanceEventRepository(tx),
			Attendance:      NewAttendanceRepository(tx),
			Card:            NewCardRepository(tx),
			Student:         NewStudentRepository(tx),
			Class:           NewClassRepository(tx),
			Rule:            NewRuleRepository(tx),
			Device:          NewDeviceRepository(tx),
			User:            NewUserRepository(tx),
			Session:         NewSessionRepository(tx),
		}

		return fn(repos)
	})
}
