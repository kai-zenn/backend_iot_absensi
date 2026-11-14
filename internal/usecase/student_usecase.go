package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
)

type IStudentUsecase interface {
  Create(ctx context.Context, req dto.CreateStudentRequest) (*dto.StudentDetailResponse, error)
  Search(ctx context.Context, req dto.SearchStudentRequest) (*dto.SearchStudentResponse, error)
  Update(ctx context.Context, id uuid.UUID, req dto.UpdateStudentRequest) error
  ImportCSV(ctx context.Context, rows [][]string) (*dto.ImportResultResponse, error)
  ListClasses(ctx context.Context) (*dto.ListClassResponse, error)
}

type studentUsecase struct {
	studentRepo repository.IStudentRepository
	classRepo   repository.IClassRepository
}

func NewStudentUsecase(sr repository.IStudentRepository, cr repository.IClassRepository) IStudentUsecase {
	return &studentUsecase{
		studentRepo: sr,
		classRepo:   cr,
	}
}

func (su *studentUsecase) Create(ctx context.Context, req dto.CreateStudentRequest) (*dto.StudentDetailResponse, error) {
  existing, _ := su.studentRepo.FindByNIS(ctx, req.NIS)
	if existing != nil {
		return nil, errors.New("NIS sudah terdaftar")
	}

	newStudent := &entity.Student{
	   NIS: req.NIS,
		 Name: req.Name,
		 ClassID: req.ClassID,
	}

	id, err := su.studentRepo.Create(ctx, newStudent)
	if err != nil {
	 return nil, err
	}

	return &dto.StudentDetailResponse{
		ID:     id,
		NIS:    newStudent.NIS,
		Name:   newStudent.Name,
		Status: newStudent.Status,
	}, nil
}

func (su *studentUsecase) Search(ctx context.Context, req dto.SearchStudentRequest) (*dto.SearchStudentResponse, error) {
  if req.Page <= 0 {
    req.Page = 1
  }
  if req.Limit <= 0 {
    req.Limit = 10
  }
  offset := (req.Page - 1) * req.Limit

  students, totalData, err := su.studentRepo.Search(ctx, req.Query, req.ClassID, req.Status, req.Limit, offset)
  if err != nil {
    return nil, err
  }

  var dtoData []dto.StudentDetailResponse
  for _, s := range students {
    className := ""
    if s.Class != nil {
      className = s.Class.Name 
    } 
    dtoData = append(dtoData, dto.StudentDetailResponse{
      ID: s.ID,
      NIS: s.NIS,
      Name: s.Name,
      Class: className,
      Status: s.Status,
    })
  }

  totalPage := int(totalData) / req.Limit
  if int(totalData)%req.Limit != 0 {
    totalPage++
  }

  return &dto.SearchStudentResponse{
    Data: dtoData,
    Meta: dto.Pagination{
      Page: req.Page,
      Limit: req.Limit,
      TotalData: totalData,
      TotalPage: totalPage,
    },
  }, nil
}

func (su *studentUsecase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateStudentRequest) error {
  student, err := su.studentRepo.FindByID(ctx, id)
  if err != nil {
    return errors.New("Siswa tidak ditemukan")
  }

  if req.Name != nil {
		student.Name = *req.Name
	}
	if req.ClassID != nil {
		student.ClassID = *req.ClassID
	}
	if req.Status != nil {
		student.Status = *req.Status
	}
  
	return su.studentRepo.Update(ctx, student)
}

func (su *studentUsecase) ImportCSV(ctx context.Context, rows [][]string) (*dto.ImportResultResponse, error) {
  var students []entity.Student
  var errorLists []dto.ImportRowError

  for i, row := range rows {
    if len(row) < 3 {
      errorLists = append(errorLists, dto.ImportRowError{
        Row: i + 1,
        Message: "Format baris tidak lengkap",
      })
      continue
    }

    nis, name, className := strings.TrimSpace(row[0]), strings.TrimSpace(row[1]), strings.TrimSpace(row[2])
    class, err := su.classRepo.FindByName(ctx, className)
    var classID uuid.UUID
    if err != nil {
      newClassID, errCreate := su.classRepo.Create(ctx, className)
     if errCreate != nil {
       errorLists = append(errorLists, dto.ImportRowError{
         Row: i + 1,
         Message: "Gagal membuat kelas ",
       })
       continue
     } 
     classID = newClassID
    } else {
      classID = class.ID
    }

    students = append(students, entity.Student{
      NIS: nis,
      Name: name,
      ClassID: classID,
      Status: "ACTIVE",
    })
  }

  if len(students) == 0 {
    return &dto.ImportResultResponse{
      SuccessCount: 0,
      Failures: errorLists,
    }, errors.New("Tidak ada data valid yang di import")
  }


  insertedCount, repoErr, err := su.studentRepo.BulkUpsertFromImport(ctx, students)
  if err != nil {
    return nil, err
  }

  errorLists = append(errorLists, repoErr...)
  
	return &dto.ImportResultResponse{
		SuccessCount: insertedCount, 
		Failures:     errorLists,    
	}, nil
}


func (su *studentUsecase) ListClasses(ctx context.Context) (*dto.ListClassResponse, error) {
	classes, err := su.classRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	var classNames []string

	for _, c := range classes {
	  classNames = append(classNames, c.Name)
	}

	return &dto.ListClassResponse{
	Data: classNames,
	}, nil
}
