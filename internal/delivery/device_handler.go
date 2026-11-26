package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/usecase"
)

type deviceHandler struct {
	useCase usecase.IDeviceUsecase
}

func NewDeviceHandler(useCase usecase.IDeviceUsecase) *deviceHandler {
	return &deviceHandler{
		useCase: useCase,
	}
}

func (h *deviceHandler) Register(ctx *gin.Context) {
	var req dto.RegisterDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Lengkapi data perangkat yang kurang",
		})
		return
	}

	res, err := h.useCase.RegisterDevice(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"message": "Gagal mendaftarkan perangkat, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Berhasil mendaftarkan perangkat",
		"data":    res,
	})
}

func (h *deviceHandler) RotateAPIKey(ctx *gin.Context) {
	deviceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "ID perangkat tidak valid",
		})
		return
	}

	res, err := h.useCase.RotateAPIKey(ctx.Request.Context(), deviceID)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"message": "Gagal mengganti API key perangkat, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil mengganti API key perangkat",
		"data":    res,
	})
}

func (h *deviceHandler) UpdateStatus(ctx *gin.Context) {
	deviceID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "ID perangkat tidak valid",
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Status perangkat wajib diisi",
		})
		return
	}

	if err := h.useCase.UpdateStatus(ctx.Request.Context(), deviceID, req.Status); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Status perangkat tidak dapat diperbarui, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil memperbarui status perangkat",
	})
}

func (h *deviceHandler) ListDevices(ctx *gin.Context) {
	res, err := h.useCase.ListDevices(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengambil data perangkat, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil menemukan data perangkat",
		"data":    res,
	})
}

func (h *deviceHandler) GetCardSnapshot(ctx *gin.Context) {
	res, err := h.useCase.GetCardSnapshot(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengambil snapshot kartu, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil mengambil snapshot kartu",
		"data":    res,
	})
}

func (h *deviceHandler) Heartbeat(ctx *gin.Context) {
	var req dto.HeartbeatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data heartbeat perangkat tidak valid",
		})
		return
	}

	res, err := h.useCase.Heartbeat(ctx.Request.Context(), req.DeviceID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Heartbeat perangkat gagal diproses, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Heartbeat perangkat berhasil diproses",
		"data":    res,
	})
}
