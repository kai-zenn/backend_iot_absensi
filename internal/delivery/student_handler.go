package delivery

import (
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/usecase"
)

type studentHandler struct {
	useCase usecase.IStudentUsecase
}

func NewStudentHandler(useCase usecase.IStudentUsecase) *studentHandler {
  return &studentHandler{
    useCase: useCase,
  }
}

func (h *studentHandler) Create(ctx *gin.Context) {
  var req dto.CreateStudentRequest

  err := ctx.ShouldBindJSON(&req)
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "status": false,
      "message": "Lengkapi data yang kurang",
    })
    return
  }

  res, err := h.useCase.Create(ctx.Request.Context(), req)
  if err != nil {
    ctx.JSON(http.StatusUnprocessableEntity, gin.H{
      "status": false,
      "message": "Gagal membuat data siswa, " + err.Error(),
    })
    return 
  }

  ctx.JSON(http.StatusCreated, gin.H{
    "status": true,
    "message": "Berhasil membuat data siswa",
    "data": res,
  })
}

func (h *studentHandler) Search(ctx *gin.Context) {
  var req dto.SearchStudentRequest

  err := ctx.ShouldBind(&req)
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "status": false,
      "message": "Tidak boleh kosong",
    })
    return
  }

  res, err := h.useCase.Search(ctx.Request.Context(), req)
  if err != nil {
    ctx.JSON(http.StatusNotFound, gin.H{
      "status": false,
      "message": "Data siswa tidak ditemukan, " + err.Error(),
    })
    return 
  }

  ctx.JSON(http.StatusOK, gin.H{
    "status": true,
    "message": "Berhasil menemukan data siswa",
    "data": res,
  })
}

func (h *studentHandler) Update(ctx *gin.Context) {
  var req dto.UpdateStudentRequest

  err := ctx.ShouldBindJSON(&req)
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "status": false,
      "message": "Tidak boleh kosong",
    })
    return
  }

  studentID, err := uuid.Parse(ctx.Param("id"))
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "status": false,
      "message": "ID siswa tidak valid",
    })
    return
  }

  err = h.useCase.Update(ctx.Request.Context(), studentID, req)
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "status": false,
      "message": "Data siswa tidak dapat diperbarui, " + err.Error(),
    })
    return 
  }

  ctx.JSON(http.StatusOK, gin.H{
    "status": true,
    "message": "Berhasil memperbarui data siswa",
  })
}

func (h *studentHandler) ImportCSV(ctx *gin.Context) {
  file, _, err := ctx.Request.FormFile("file")
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "status": false,
			"message": "File CSV wajib diupload",
    })
    return
  }

  defer file.Close()

  reader := csv.NewReader(file)

  rows, err := reader.ReadAll()
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"message": "Format isi CSV rusak, " + err.Error(),
		})
		return
  }

  if len(rows) > 0 && len(rows[0]) > 0 && strings.EqualFold(strings.TrimSpace(rows[0][0]), "nis") {
		rows = rows[1:] 
	}

	res, err := h.useCase.ImportCSV(ctx.Request.Context(), rows)
	if err != nil {
	  var failures any
		if res != nil {
		  failures = res.Failures
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
		  "status": false,
			"message": err.Error(),
			"failures": failures,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"message": "Import selesai",
		"data": res,
	})
}

func (h *studentHandler) ListClasses(ctx *gin.Context) {
  res, err := h.useCase.ListClasses(ctx.Request.Context())
  if err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "status": false,
      "message": "Data tidak ditemukan",
    })
    return
  }

 	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"message": "Data berhasil ditemukan",
		"data": res,
	})
}
