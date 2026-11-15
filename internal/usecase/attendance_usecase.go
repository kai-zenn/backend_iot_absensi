package usecase

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
)

type IAttendanceUsecase interface {
	ProcessScanEvent(ctx context.Context, deviceID uuid.UUID, req dto.ScanEventRequest) (*dto.ScanEventResponse, error)
	ProcessSyncBatch(ctx context.Context, deviceID uuid.UUID, req dto.SyncBatchRequest) (*dto.SyncBatchResponse, error)
	ManualCheckIn(ctx context.Context, actorUserID uuid.UUID, req dto.ManualCheckInRequest) (*dto.AttendanceDetailResponse, error)
	GetSummary(ctx context.Context, date time.Time, classID *uuid.UUID) (*dto.AttendanceSummaryResponse, error)
	ListAttendance(ctx context.Context, req dto.ListAttendanceRequest) (*dto.ListAttendanceResponse, error)
	GetStudentHistory(ctx context.Context, studentID uuid.UUID, from, to time.Time) (*dto.StudentHistoryResponse, error)
	ListNeedsReview(ctx context.Context, page, limit int) (*dto.ListEventsResponse, error)
}

type attendanceUsecase struct {
	attendanceRepo repository.IAttendanceRepository
	eventRepo      repository.IAttendanceEventRepository
	cardRepo       repository.ICardRepository
	studentRepo    repository.IStudentRepository
	ruleRepo       repository.IRuleRepository
	deviceRepo     repository.IDeviceRepository
	txManager      repository.TxManager
	schoolTZ       *time.Location
}

func NewAttendanceUsecase(
	ar repository.IAttendanceRepository,
	aer repository.IAttendanceEventRepository,
	cr repository.ICardRepository,
	sr repository.IStudentRepository,
	rr repository.IRuleRepository,
	dr repository.IDeviceRepository,
	tx repository.TxManager,
) IAttendanceUsecase {
	return &attendanceUsecase{
		attendanceRepo: ar,
		eventRepo:      aer,
		cardRepo:       cr,
		studentRepo:    sr,
		ruleRepo:       rr,
		deviceRepo:     dr,
		txManager:      tx,
		schoolTZ:       time.FixedZone("WIB", 7*3600),
	}
}

func normalizeUID(uid string) string {
	r := strings.NewReplacer(" ", "", ":", "", "-", "")
	return strings.ToUpper(r.Replace(uid))
}

func (au *attendanceUsecase) toSchoolLocal(t time.Time) (localDate, localTime time.Time) {
	lt := t.In(au.schoolTZ)
	ld := time.Date(lt.Year(), lt.Month(), lt.Day(), 0, 0, 0, 0, au.schoolTZ)
	return ld, lt
}

func computeCheckInStatus(localTap time.Time, rule *entity.AttendanceRule) string {
	if rule == nil {
		return "PRESENT"
	}
	deadline := time.Date(
		localTap.Year(), localTap.Month(), localTap.Day(),
		rule.OnTimeUntil.Hour(), rule.OnTimeUntil.Minute(), rule.OnTimeUntil.Second(),
		0, localTap.Location(),
	)
	if localTap.After(deadline) {
		return "LATE"
	}
	return "PRESENT"
}

func mapStudentBrief(s *entity.Student) dto.StudentBrief {
	if s == nil {
		return dto.StudentBrief{}
	}
	className := ""
	if s.Class != nil {
		className = s.Class.Name
	}
	return dto.StudentBrief{
		NIS:   s.NIS,
		Name:  s.Name,
		Class: className,
	}
}

func (au *attendanceUsecase) ProcessScanEvent(ctx context.Context, deviceID uuid.UUID, req dto.ScanEventRequest) (*dto.ScanEventResponse, error) {
	return au.processSingleEvent(ctx, deviceID, req, "LIVE")
}

func (au *attendanceUsecase) processSingleEvent(ctx context.Context, deviceID uuid.UUID, req dto.ScanEventRequest, source string) (*dto.ScanEventResponse, error) {
	now := time.Now().In(au.schoolTZ)
	normalizedUID := normalizeUID(req.CardUID)

	var resp *dto.ScanEventResponse

	err := au.txManager.WithinTx(ctx, func(repos repository.RepoSet) error {
		event := &entity.AttendanceEvent{
			EventID:    req.EventID,
			DeviceID:   deviceID,
			CardUID:    normalizedUID,
			OccurredAt: req.OccurredAt,
			Source:     source,
		}

		existingRow, alreadyExisted, err := repos.AttendanceEvent.TryInsert(ctx, event)
		if err != nil {
			return err
		}

		if alreadyExisted && existingRow != nil {
			if existingRow.CardUID != normalizedUID || !existingRow.OccurredAt.Equal(req.OccurredAt) {
				resp = &dto.ScanEventResponse{
					EventID:    req.EventID,
					Result:     "EVENT_ID_CONFLICT",
					Replayed:   false,
					Message:    "Event ID bentrok dengan data berbeda",
					OccurredAt: req.OccurredAt,
					ServerTime: now,
				}
				return nil
			}

			resp = &dto.ScanEventResponse{
				EventID:    existingRow.EventID,
				Result:     existingRow.Result,
				Replayed:   true,
				OccurredAt: existingRow.OccurredAt,
				ServerTime: now,
			}
			return nil
		}

		var clockSkewPtr *int
		if source == "LIVE" {
			skew := int(now.Sub(req.OccurredAt).Seconds())
			clockSkewPtr = &skew
		}

		if req.OccurredAt.After(now.Add(5*time.Minute)) || (source == "LIVE" && req.OccurredAt.Before(now.Add(-5*time.Minute))) {
			resp = &dto.ScanEventResponse{
				EventID:    req.EventID,
				Result:     "REJECTED_CLOCK",
				Message:    "Waktu perangkat tidak wajar",
				OccurredAt: req.OccurredAt,
				ServerTime: now,
			}
			return repos.AttendanceEvent.UpdateResult(ctx, event.ID, "REJECTED_CLOCK", nil, nil, clockSkewPtr)
		}

		card, err := repos.Card.FindByUID(ctx, normalizedUID)
		if err != nil || card == nil {
			resp = &dto.ScanEventResponse{
				EventID:    req.EventID,
				Result:     "UNKNOWN_CARD",
				Message:    "Kartu belum terdaftar",
				OccurredAt: req.OccurredAt,
				ServerTime: now,
			}
			return repos.AttendanceEvent.UpdateResult(ctx, event.ID, "UNKNOWN_CARD", nil, nil, clockSkewPtr)
		}
		if card.Status != "ACTIVE" {
			resp = &dto.ScanEventResponse{
				EventID:    req.EventID,
				Result:     "CARD_BLOCKED",
				Message:    "Kartu diblokir atau nonaktif",
				OccurredAt: req.OccurredAt,
				ServerTime: now,
			}
			return repos.AttendanceEvent.UpdateResult(ctx, event.ID, "CARD_BLOCKED", &card.StudentID, nil, clockSkewPtr)
		}

		student, err := repos.Student.FindByID(ctx, card.StudentID)
		if err != nil || student == nil || student.Status != "ACTIVE" {
			resp = &dto.ScanEventResponse{
				EventID:    req.EventID,
				Result:     "STUDENT_INACTIVE",
				Message:    "Siswa tidak aktif",
				OccurredAt: req.OccurredAt,
				ServerTime: now,
			}
			return repos.AttendanceEvent.UpdateResult(ctx, event.ID, "STUDENT_INACTIVE", &card.StudentID, nil, clockSkewPtr)
		}

		localDate, localTime := au.toSchoolLocal(req.OccurredAt)
		rule, _ := repos.Rule.FindEffectiveForDate(ctx, localDate)
		status := computeCheckInStatus(localTime, rule)

		attID, inserted, err := repos.Attendance.CreateCheckIn(ctx, &entity.Attendance{
			StudentID:      student.ID,
			AttendanceDate: localDate,
			CheckInAt:      req.OccurredAt,
			CheckInStatus:  status,
			CheckInSource:  "CARD",
			CheckInEventID: &event.ID,
		})
		if err != nil {
			return err
		}

		resultStr := "CHECKED_IN"
		finalStatus := status
		finalAttID := attID

		if !inserted {
			existingAtt, err := repos.Attendance.FindByStudentAndDateForUpdate(ctx, student.ID, localDate)
			if err != nil {
				return err
			}
			finalAttID = existingAtt.ID
			resultStr = "ALREADY_CHECKED_IN"

			if req.OccurredAt.Before(existingAtt.CheckInAt) {
				if err := repos.Attendance.CorrectCheckIn(ctx, existingAtt.ID, req.OccurredAt, status, event.ID); err != nil {
					return err
				}
				finalStatus = status
			} else {
				finalStatus = existingAtt.CheckInStatus
			}
		}

		if err := repos.AttendanceEvent.UpdateResult(ctx, event.ID, resultStr, &student.ID, &finalAttID, clockSkewPtr); err != nil {
			return err
		}

		brief := mapStudentBrief(student)
		resp = &dto.ScanEventResponse{
			EventID:    req.EventID,
			Result:     resultStr,
			Status:     finalStatus,
			Replayed:   false,
			Student:    &brief,
			OccurredAt: req.OccurredAt,
			ServerTime: now,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (au *attendanceUsecase) ProcessSyncBatch(ctx context.Context, deviceID uuid.UUID, req dto.SyncBatchRequest) (*dto.SyncBatchResponse, error) {
	sorted := make([]dto.SyncEventItem, len(req.Events))
	copy(sorted, req.Events)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].OccurredAt.Before(sorted[j].OccurredAt)
	})

	resultMap := make(map[string]dto.SyncEventResult, len(sorted))
	for _, item := range sorted {
		res, err := au.processSingleEvent(ctx, deviceID, dto.ScanEventRequest{
			EventID:    item.EventID,
			CardUID:    item.CardUID,
			DeviceID:   deviceID,
			OccurredAt: item.OccurredAt,
		}, "SYNC")
		if err != nil {
			return nil, err
		}
		resultMap[item.EventID] = dto.SyncEventResult{
			EventID:  res.EventID,
			Result:   res.Result,
			Status:   res.Status,
			Replayed: res.Replayed,
			Message:  res.Message,
		}
	}

	// Susun ulang sesuai urutan request asli (§4b langkah 5)
	orderedResults := make([]dto.SyncEventResult, 0, len(req.Events))
	for _, orig := range req.Events {
		orderedResults = append(orderedResults, resultMap[orig.EventID])
	}

	return &dto.SyncBatchResponse{
		Results:    orderedResults,
		ServerTime: time.Now().In(au.schoolTZ),
	}, nil
}

func (au *attendanceUsecase) ManualCheckIn(ctx context.Context, actorUserID uuid.UUID, req dto.ManualCheckInRequest) (*dto.AttendanceDetailResponse, error) {
	if strings.TrimSpace(req.Note) == "" {
		return nil, errors.New("catatan wajib diisi")
	}

	student, err := au.studentRepo.FindByID(ctx, req.StudentID)
	if err != nil || student == nil {
		return nil, errors.New("siswa tidak ditemukan")
	}
	if student.Status != "ACTIVE" {
		return nil, errors.New("siswa tidak aktif")
	}

	var resp *dto.AttendanceDetailResponse
	err = au.txManager.WithinTx(ctx, func(repos repository.RepoSet) error {
		now := time.Now().In(au.schoolTZ)
		localDate, localTime := au.toSchoolLocal(now)
		rule, _ := repos.Rule.FindEffectiveForDate(ctx, localDate)
		status := computeCheckInStatus(localTime, rule)

		note := req.Note
		id, inserted, err := repos.Attendance.CreateManual(ctx, &entity.Attendance{
			StudentID:       student.ID,
			AttendanceDate:  localDate,
			CheckInAt:       now,
			CheckInStatus:   status,
			CheckInSource:   "MANUAL",
			CheckInByUserID: &actorUserID,
			Note:            &note,
		})
		if err != nil {
			return err
		}
		if !inserted {
			return errors.New("CONFLICT: siswa sudah tercatat hadir hari ini")
		}

		resp = &dto.AttendanceDetailResponse{
			ID:            id,
			StudentID:     student.ID,
			Student:       mapStudentBrief(student),
			CheckInAt:     now,
			CheckInStatus: status,
			CheckInSource: "MANUAL",
			Note:          req.Note,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (au *attendanceUsecase) GetSummary(ctx context.Context, date time.Time, classID *uuid.UUID) (*dto.AttendanceSummaryResponse, error) {
	present, late, activeStudents, err := au.attendanceRepo.Summary(ctx, date, classID)
	if err != nil {
		return nil, err
	}

	notYet := activeStudents - (present + late)
	if notYet < 0 {
		notYet = 0
	}

	_, needsReview, err := au.eventRepo.ListNeedsReview(ctx, 1, 0)
	if err != nil {
		return nil, err
	}

	devices, err := au.deviceRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	pendingSync := 0
	for _, d := range devices {
		pendingSync += d.PendingQueue
	}

	return &dto.AttendanceSummaryResponse{
		Date:           date.Format("2006-01-02"),
		ActiveStudents: activeStudents,
		Present:        present,
		Late:           late,
		NotYetArrived:  notYet,
		PendingSync:    pendingSync,
		NeedsReview:    int(needsReview),
	}, nil
}

func (au *attendanceUsecase) ListAttendance(ctx context.Context, req dto.ListAttendanceRequest) (*dto.ListAttendanceResponse, error) {
	targetDate, _ := au.toSchoolLocal(time.Now())
	if req.Date != "" {
		parsed, err := time.ParseInLocation("2006-01-02", req.Date, au.schoolTZ)
		if err != nil {
			return nil, errors.New("format tanggal harus YYYY-MM-DD")
		}
		targetDate = parsed
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit

	rows, totalData, err := au.attendanceRepo.ListByFilter(ctx, targetDate, req.ClassID, req.Status, req.Query, req.Limit, offset)
	if err != nil {
		return nil, err
	}

	items := make([]dto.AttendanceDetailResponse, 0, len(rows))
	for _, r := range rows {
		noteStr := ""
		if r.Note != nil {
			noteStr = *r.Note
		}
		items = append(items, dto.AttendanceDetailResponse{
			ID:            r.ID,
			StudentID:     r.StudentID,
			Student:       mapStudentBrief(r.Student),
			CheckInAt:     r.CheckInAt,
			CheckInStatus: r.CheckInStatus,
			CheckInSource: r.CheckInSource,
			Note:          noteStr,
		})
	}

	totalPage := int(totalData) / req.Limit
	if int(totalData)%req.Limit != 0 {
		totalPage++
	}

	return &dto.ListAttendanceResponse{
		Data: items,
		Meta: dto.Pagination{
			Page:      req.Page,
			Limit:     req.Limit,
			TotalData: totalData,
			TotalPage: totalPage,
		},
	}, nil
}

func (au *attendanceUsecase) GetStudentHistory(ctx context.Context, studentID uuid.UUID, from, to time.Time) (*dto.StudentHistoryResponse, error) {
	rows, err := au.attendanceRepo.ListHistoryByStudent(ctx, studentID, from, to)
	if err != nil {
		return nil, err
	}

	items := make([]dto.AttendanceDetailResponse, 0, len(rows))
	for _, r := range rows {
		noteStr := ""
		if r.Note != nil {
			noteStr = *r.Note
		}
		items = append(items, dto.AttendanceDetailResponse{
			ID:            r.ID,
			StudentID:     r.StudentID,
			Student:       mapStudentBrief(r.Student),
			CheckInAt:     r.CheckInAt,
			CheckInStatus: r.CheckInStatus,
			CheckInSource: r.CheckInSource,
			Note:          noteStr,
		})
	}

	return &dto.StudentHistoryResponse{
		StudentID: studentID,
		Data:      items,
	}, nil
}

func (au *attendanceUsecase) ListNeedsReview(ctx context.Context, page, limit int) (*dto.ListEventsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	events, totalData, err := au.eventRepo.ListNeedsReview(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	items := make([]dto.EventLogItem, 0, len(events))
	for _, ev := range events {
		items = append(items, dto.EventLogItem{
			ID:         ev.ID,
			EventID:    ev.EventID,
			DeviceID:   ev.DeviceID,
			CardUID:    ev.CardUID,
			Result:     ev.Result,
			OccurredAt: ev.OccurredAt,
			ReceivedAt: ev.ReceivedAt,
		})
	}

	totalPage := int(totalData) / limit
	if int(totalData)%limit != 0 {
		totalPage++
	}

	return &dto.ListEventsResponse{
		Data: items,
		Meta: dto.Pagination{
			Page:      page,
			Limit:     limit,
			TotalData: totalData,
			TotalPage: totalPage,
		},
	}, nil
}
