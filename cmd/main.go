package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/kai-zenn/backend_iot_absensi/internal/config"
	bcryptpkg "github.com/kai-zenn/backend_iot_absensi/internal/pkgs/bcrypt"
	"github.com/kai-zenn/backend_iot_absensi/internal/pkgs/helper"
	jwtpkg "github.com/kai-zenn/backend_iot_absensi/internal/pkgs/jwt"
	postgresrepo "github.com/kai-zenn/backend_iot_absensi/internal/repository/postgres"
	"github.com/kai-zenn/backend_iot_absensi/internal/usecase"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConf()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dsn := (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   net.JoinHostPort(cfg.DBHost, cfg.DBPort),
		Path:   "/" + cfg.DBName,
		RawQuery: url.Values{
			"sslmode": []string{cfg.DBSSLMode},
		}.Encode(),
	}).String()

	db, err := gorm.Open(postgresDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open PostgreSQL connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get PostgreSQL connection pool: %w", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Error("close PostgreSQL connection", "error", err)
		}
	}()

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping PostgreSQL: %w", err)
	}

	engine := binding.Validator.Engine()
	validate, ok := engine.(*validator.Validate)
	if !ok {
		return errors.New("Gin validator engine is not go-playground/validator")
	}
	if err := validate.RegisterValidation("rfiduid", func(fl validator.FieldLevel) bool {
		return helper.IsValidUID(fl.Field().String())
	}); err != nil {
		return fmt.Errorf("register rfiduid validation: %w", err)
	}

	bcryptService := bcryptpkg.NewCryptoBcrypt()
	jwtService := jwtpkg.NewJWTToken(cfg.JWTSecret, cfg.JWTExpiry)

	txManager := postgresrepo.NewTxManager(db)
	attendanceEventRepository := postgresrepo.NewAttendanceEventRepository(db)
	attendanceRepository := postgresrepo.NewAttendanceRepository(db)
	cardRepository := postgresrepo.NewCardRepository(db)
	deviceRepository := postgresrepo.NewDeviceRepository(db)
	ruleRepository := postgresrepo.NewRuleRepository(db)
	studentRepository := postgresrepo.NewStudentRepository(db)
	classRepository := postgresrepo.NewClassRepository(db)
	userRepository := postgresrepo.NewUserRepository(db)

	attendanceUsecase := usecase.NewAttendanceUsecase(
		attendanceRepository,
		attendanceEventRepository,
		cardRepository,
		studentRepository,
		ruleRepository,
		deviceRepository,
		txManager,
	)
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepository)
	cardUsecase := usecase.NewCardUsecase(cardRepository, txManager)
	studentUsecase := usecase.NewStudentUsecase(studentRepository, classRepository)
	ruleUsecase := usecase.NewRuleUsecase(ruleRepository)
	authUsecase := usecase.NewAuthUseCase(userRepository, bcryptService, jwtService)

	// TODO: Delivery Layer Wiring (Handlers, Middlewares, & rest.go)
	// Initialize AuthHandler, StudentHandler, RuleHandler, CardHandler,
	// DeviceHandler, and AttendanceHandler, then mount their routes under /api/v1.
	// Handler files and route wiring are placeholders until the delivery layer is ready.
	_ = attendanceUsecase
	_ = deviceUsecase
	_ = cardUsecase
	_ = studentUsecase
	_ = ruleUsecase
	_ = authUsecase

	router := gin.New()
	router.Use(gin.Recovery(), requestLogger())
	router.GET("/healthz", func(c *gin.Context) {
		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("HTTP server listening", "address", srv.Addr, "environment", cfg.AppEnv)
		serverErrors <- srv.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server: %w", err)
		}
	case <-ctx.Done():
		slog.Info("shutting down HTTP server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		if err := <-serverErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server: %w", err)
		}
	}

	return nil
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()
		slog.Info("HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(startedAt),
			"client_ip", c.ClientIP(),
		)
	}
}
