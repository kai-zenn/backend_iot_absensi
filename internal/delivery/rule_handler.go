package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/usecase"
)

type ruleHandler struct {
	useCase usecase.IRuleUsecase
}

func NewRuleHandler(useCase usecase.IRuleUsecase) *ruleHandler {
	return &ruleHandler{useCase: useCase}
}

func (h *ruleHandler) List(ctx *gin.Context) {
	res, err := h.useCase.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengambil aturan kehadiran, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil menemukan aturan kehadiran",
		"data":    res,
	})
}

func (h *ruleHandler) Upsert(ctx *gin.Context) {
	var req dto.UpsertRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data aturan kehadiran tidak valid",
		})
		return
	}

	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Sesi anda tidak valid, silahkan login ulang",
		})
		return
	}
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Sesi anda tidak valid, silahkan login ulang",
		})
		return
	}

	if err := h.useCase.Upsert(ctx.Request.Context(), userID, req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"message": "Gagal menyimpan aturan kehadiran, " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil menyimpan aturan kehadiran",
	})
}
