package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kai-zenn/backend_iot_absensi/internal/pkgs/helper"
	"github.com/kai-zenn/backend_iot_absensi/internal/pkgs/jwt"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
)

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exist := ctx.Get("role")
		if !exist {
			ctx.JSON(http.StatusAccepted, gin.H{
				"status":  false,
				"message": "Sesi tidak teridentifikasi",
			})
			ctx.Abort()
			return
		}

		role := fmt.Sprint(userRole)
		isAllowed := slices.Contains(allowedRoles, role)

		if !isAllowed {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Role tidak ditemukan",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func AuthenticateDevice(deviceRepo repository.IDeviceRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")
		if apiKey == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "API key perangkat tidak ditemukan",
			})
			ctx.Abort()
			return
		}

		device, err := deviceRepo.FindByAPIKeyHash(ctx.Request.Context(), helper.HashAPIKey(apiKey))
		if err != nil || device == nil || device.Status != "ACTIVE" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "API key perangkat tidak valid",
			})
			ctx.Abort()
			return
		}

		ctx.Set("device_id", device.ID)
		ctx.Next()
	}
}

func Authenticate(jwtService jwt.JwtInterface) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    authHeader := ctx.GetHeader("Authorization")
    
    if authHeader == "" {
      ctx.JSON(http.StatusUnauthorized, gin.H{
        "status": false,
        "message": "Token tidak ditemukan",
      })
      ctx.Abort()
      return
    }
    
    tokenParts := strings.Split(authHeader, " ")
    
    if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Format token salah, wajib Bearer Token)",
			})
			ctx.Abort()
			return
		}

		token := tokenParts[1]
		
		claims, err := jwtService.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Sesi login kadaluarsa",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserId)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
