package seed

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/pkgs/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultRuleEffectiveFrom = "2000-01-01"
	defaultOnTimeUntil       = "06:30:00"
)

func Run(db *gorm.DB, bcryptService bcrypt.Interface, adminUsername, adminPassword, adminName string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := seedDefaultRule(tx); err != nil {
			return err
		}

		var adminCount int64
		if err := tx.Model(&entity.User{}).Where("role = ?", entity.RoleAdmin).Count(&adminCount).Error; err != nil {
			return fmt.Errorf("check for existing admin: %w", err)
		}
		if adminCount > 0 {
			return nil
		}

		adminUsername = strings.TrimSpace(adminUsername)
		adminName = strings.TrimSpace(adminName)
		if adminUsername == "" || adminPassword == "" || adminName == "" {
			return errors.New("no admin account exists; configure ADMIN_USERNAME, ADMIN_PASSWORD, and ADMIN_NAME to bootstrap the first admin")
		}

		var existing entity.User
		err := tx.Where("username = ?", adminUsername).First(&existing).Error
		if err == nil {
			return fmt.Errorf("configured admin username %q already belongs to a non-admin user", adminUsername)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check configured admin username: %w", err)
		}

		passwordHash, err := bcryptService.GenerateHashPassword(adminPassword)
		if err != nil {
			return fmt.Errorf("hash initial admin password: %w", err)
		}

		admin := entity.User{
			Username:     adminUsername,
			PasswordHash: &passwordHash,
			Name:         adminName,
			Role:         entity.RoleAdmin,
			IsActive:     true,
		}
		if err := tx.Create(&admin).Error; err != nil {
			return fmt.Errorf("create initial admin: %w", err)
		}

		return nil
	})
}

func seedDefaultRule(db *gorm.DB) error {
	effectiveFrom, err := time.Parse("2006-01-02", defaultRuleEffectiveFrom)
	if err != nil {
		return fmt.Errorf("parse default rule effective date: %w", err)
	}
	onTimeUntil, err := time.Parse("15:04:05", defaultOnTimeUntil)
	if err != nil {
		return fmt.Errorf("parse default rule time: %w", err)
	}

	rule := entity.AttendanceRule{}
	result := db.Where(&entity.AttendanceRule{EffectiveFrom: effectiveFrom}).
		Attrs(entity.AttendanceRule{OnTimeUntil: onTimeUntil}).
		FirstOrCreate(&rule)
	if result.Error != nil {
		return fmt.Errorf("seed default attendance rule: %w", result.Error)
	}
	return nil
}
