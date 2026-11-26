package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/usecase"
)

type cardHandler struct {
	useCase usecase.ICardUsecase
}

func NewCardHandler(useCase usecase.ICardUsecase) *cardHandler {
	return &cardHandler{
		useCase: useCase,
	}
}

func (h *cardHandler) ListUnregistered(ctx *gin.Context) {
	res, err := h.useCase.ListUnregistered(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data kartu tidak ditemukan, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil menemukan data kartu",
		"data":    res,
	})
}

func (h *cardHandler) AssignCard(ctx *gin.Context) {
	var req dto.AssignCardRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Lengkapi data kartu yang kurang",
		})
		return
	}

	err = h.useCase.AssignCard(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"message": "Gagal menghubungkan kartu ke siswa, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Berhasil menghubungkan kartu ke siswa",
	})
}

func (h *cardHandler) UpdateCardStatus(ctx *gin.Context) {
	cardID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "ID kartu tidak valid",
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Status kartu wajib diisi",
		})
		return
	}

	if err := h.useCase.UpdateCardStatus(ctx.Request.Context(), cardID, req.Status); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Status kartu tidak dapat diperbarui, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil memperbarui status kartu",
	})
}
