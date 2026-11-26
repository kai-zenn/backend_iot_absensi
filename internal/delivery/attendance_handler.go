package delivery

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/usecase"
)

type attendanceHandler struct {
	useCase usecase.IAttendanceUsecase
}

func NewAttendanceHandler(useCase usecase.IAttendanceUsecase) *attendanceHandler {
	return &attendanceHandler{
		useCase: useCase,
	}
}

func (h *attendanceHandler) ProcessScanEvent(ctx *gin.Context) {
	var req dto.ScanEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data scan kehadiran tidak valid",
		})
		return
	}

	deviceID, exists := DeviceIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Identitas perangkat tidak valid",
		})
		return
	}
	req.DeviceID = deviceID

	res, err := h.useCase.ProcessScanEvent(ctx.Request.Context(), deviceID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal memproses scan kehadiran, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Scan kehadiran berhasil diproses",
		"data":    res,
	})
}

func (h *attendanceHandler) ProcessSyncBatch(ctx *gin.Context) {
	var req dto.SyncBatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data sinkronisasi kehadiran tidak valid",
		})
		return
	}

	deviceID, exists := DeviceIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Identitas perangkat tidak valid",
		})
		return
	}
	req.DeviceID = deviceID

	res, err := h.useCase.ProcessSyncBatch(ctx.Request.Context(), deviceID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal memproses sinkronisasi kehadiran, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Sinkronisasi kehadiran berhasil diproses",
		"data":    res,
	})
}

func (h *attendanceHandler) ManualCheckIn(ctx *gin.Context) {
	var req dto.ManualCheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Lengkapi data check-in manual yang kurang",
		})
		return
	}

	actorUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Sesi anda tidak valid, silahkan login ulang",
		})
		return
	}

	userID, ok := actorUserID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Sesi anda tidak valid, silahkan login ulang",
		})
		return
	}

	res, err := h.useCase.ManualCheckIn(ctx.Request.Context(), userID, req)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"message": "Check-in manual gagal diproses, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Check-in manual berhasil dicatat",
		"data":    res,
	})
}

func (h *attendanceHandler) GetSummary(ctx *gin.Context) {
	date := time.Now()
	if dateParam := ctx.Query("date"); dateParam != "" {
		parsedDate, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Format tanggal harus YYYY-MM-DD",
			})
			return
		}
		date = parsedDate
	}

	var classID *uuid.UUID
	if classIDParam := ctx.Query("classId"); classIDParam != "" {
		parsedID, err := uuid.Parse(classIDParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "ID kelas tidak valid",
			})
			return
		}
		classID = &parsedID
	}

	res, err := h.useCase.GetSummary(ctx.Request.Context(), date, classID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengambil ringkasan kehadiran, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil mengambil ringkasan kehadiran",
		"data":    res,
	})
}

func (h *attendanceHandler) ListAttendance(ctx *gin.Context) {
	var req dto.ListAttendanceRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Parameter pencarian kehadiran tidak valid",
		})
		return
	}

	res, err := h.useCase.ListAttendance(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data kehadiran tidak dapat ditemukan, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil menemukan data kehadiran",
		"data":    res,
	})
}

func (h *attendanceHandler) GetStudentHistory(ctx *gin.Context) {
	studentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "ID siswa tidak valid",
		})
		return
	}

	fromParam, toParam := ctx.Query("from"), ctx.Query("to")
	if fromParam == "" || toParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Tanggal from dan to wajib diisi dengan format YYYY-MM-DD",
		})
		return
	}

	from, err := time.Parse("2006-01-02", fromParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Format tanggal from harus YYYY-MM-DD",
		})
		return
	}
	to, err := time.Parse("2006-01-02", toParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Format tanggal to harus YYYY-MM-DD",
		})
		return
	}

	res, err := h.useCase.GetStudentHistory(ctx.Request.Context(), studentID, from, to)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Riwayat kehadiran tidak ditemukan, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil menemukan riwayat kehadiran siswa",
		"data":    res,
	})
}

func (h *attendanceHandler) ListNeedsReview(ctx *gin.Context) {
	page, err := parseQueryInt(ctx, "page", 1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Parameter page harus berupa angka",
		})
		return
	}
	limit, err := parseQueryInt(ctx, "limit", 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Parameter limit harus berupa angka",
		})
		return
	}

	res, err := h.useCase.ListNeedsReview(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengambil daftar event yang perlu ditinjau, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil mengambil daftar event yang perlu ditinjau",
		"data":    res,
	})
}

func parseQueryInt(ctx *gin.Context, key string, defaultValue int) (int, error) {
	value := ctx.Query(key)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}
