package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/usecase"
)

type AuthHandler struct {
	useCase usecase.IAuthUsecase
}

func NewAuthHandler(useCase usecase.IAuthUsecase) *AuthHandler {
  return &AuthHandler{
    useCase: useCase,
  }
}

func (h *AuthHandler) Login(ctx *gin.Context) {
  var req dto.LoginRequest

  err := ctx.ShouldBindJSON(&req)
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "status": false,
      "message": "Username Wajib diisi, ",
    })
    return
  }

  token, err := h.useCase.Login(ctx.Request.Context(), req)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "status": false,
      "message": "Gagal login, " + err.Error(),
    })
    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "status":  true,
    "message": "Berhasil login",
    "token":   token.Token,
  })
}
