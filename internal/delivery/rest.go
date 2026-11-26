package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/delivery/middleware"
	bcryptpkg "github.com/kai-zenn/backend_iot_absensi/internal/pkgs/bcrypt"
	jwtpkg "github.com/kai-zenn/backend_iot_absensi/internal/pkgs/jwt"
	postgresrepo "github.com/kai-zenn/backend_iot_absensi/internal/repository/postgres"
	"github.com/kai-zenn/backend_iot_absensi/internal/usecase"
	"gorm.io/gorm"
)

type Rest struct {
	router  *gin.Engine
	db      *gorm.DB
	jwtAuth jwtpkg.JwtInterface
	bcrypt  bcryptpkg.Interface
}

func NewRest(router *gin.Engine, db *gorm.DB, jwtAuth jwtpkg.JwtInterface, bcrypt bcryptpkg.Interface) *Rest {
	return &Rest{
		router:  router,
		db:      db,
		jwtAuth: jwtAuth,
		bcrypt:  bcrypt,
	}
}

func (r *Rest) MountEndPoint() {
	// Repositories
	txManager := postgresrepo.NewTxManager(r.db)
	attendanceEventRepo := postgresrepo.NewAttendanceEventRepository(r.db)
	attendanceRepo := postgresrepo.NewAttendanceRepository(r.db)
	cardRepo := postgresrepo.NewCardRepository(r.db)
	deviceRepo := postgresrepo.NewDeviceRepository(r.db)
	ruleRepo := postgresrepo.NewRuleRepository(r.db)
	studentRepo := postgresrepo.NewStudentRepository(r.db)
	classRepo := postgresrepo.NewClassRepository(r.db)
	userRepo := postgresrepo.NewUserRepository(r.db)

	// Use case
	authUseCase := usecase.NewAuthUseCase(userRepo, r.bcrypt, r.jwtAuth)
	studentUseCase := usecase.NewStudentUsecase(studentRepo, classRepo)
	cardUseCase := usecase.NewCardUsecase(cardRepo, txManager)
	deviceUseCase := usecase.NewDeviceUsecase(deviceRepo, cardRepo, ruleRepo)
	ruleUseCase := usecase.NewRuleUsecase(ruleRepo)
	attendanceUseCase := usecase.NewAttendanceUsecase(
		attendanceRepo,
		attendanceEventRepo,
		cardRepo,
		studentRepo,
		ruleRepo,
		deviceRepo,
		txManager,
	)

	// Handler
	authHandler := NewAuthHandler(authUseCase)
	studentHandler := NewStudentHandler(studentUseCase)
	cardHandler := NewCardHandler(cardUseCase)
	deviceHandler := NewDeviceHandler(deviceUseCase)
	ruleHandler := NewRuleHandler(ruleUseCase)
	attendanceHandler := NewAttendanceHandler(attendanceUseCase)

	api := r.router.Group("/api/v1")
	api.POST("/auth/login", authHandler.Login)

	// Authenticated endpoints.
	protected := api.Group("")
	protected.Use(middleware.Authenticate(r.jwtAuth))

	students := protected.Group("/students")
	students.Use(middleware.RequireRoles("admin", "staff"))
	{
		students.POST("", studentHandler.Create)
		students.GET("", studentHandler.Search)
		students.PATCH("/:id", studentHandler.Update)
		students.POST("/import", studentHandler.ImportCSV)
		students.GET("/classes", studentHandler.ListClasses)
	}

	cards := protected.Group("/cards")
	cards.Use(middleware.RequireRoles("admin", "staff"))
	{
		cards.GET("/unregistered", cardHandler.ListUnregistered)
		cards.POST("/assign", cardHandler.AssignCard)
		cards.PATCH("/:id/status", cardHandler.UpdateCardStatus)
	}

	devices := protected.Group("/devices")
	devices.Use(middleware.RequireRoles("admin"))
	{
		devices.POST("", deviceHandler.Register)
		devices.GET("", deviceHandler.ListDevices)
		devices.POST("/:id/rotate-key", deviceHandler.RotateAPIKey)
		devices.PATCH("/:id/status", deviceHandler.UpdateStatus)
	}

	rules := protected.Group("/rules")
	rules.Use(middleware.RequireRoles("admin"))
	{
		rules.GET("", ruleHandler.List)
		rules.PUT("", ruleHandler.Upsert)
	}

	attendance := protected.Group("/attendance")
	attendance.Use(middleware.RequireRoles("admin", "staff"))
	{
		attendance.GET("", attendanceHandler.ListAttendance)
		attendance.GET("/summary", attendanceHandler.GetSummary)
		attendance.POST("/manual", attendanceHandler.ManualCheckIn)
		attendance.GET("/students/:id/history", attendanceHandler.GetStudentHistory)
		attendance.GET("/events/review", attendanceHandler.ListNeedsReview)
	}

	deviceAPI := api.Group("")
	deviceAPI.Use(middleware.AuthenticateDevice(deviceRepo))
	deviceAPI.POST("/devices/heartbeat", deviceHandler.Heartbeat)
	deviceAPI.GET("/devices/cards/snapshot", deviceHandler.GetCardSnapshot)
	deviceAPI.POST("/attendance/events", attendanceHandler.ProcessScanEvent)
	deviceAPI.POST("/attendance/sync", attendanceHandler.ProcessSyncBatch)
}

func (r *Rest) Serve(port string) error {
	return r.router.Run(":" + port)
}

func DeviceIDFromContext(ctx *gin.Context) (uuid.UUID, bool) {
	deviceID, exists := ctx.Get("device_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := deviceID.(uuid.UUID)
	return id, ok
}
