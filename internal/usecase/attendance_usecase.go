package usecase


import (
  "context"
  "time"
  "github.com/kai-zenn/backend_iot_absensi/internal/dto"
)

type IAttendanceUsecase interface {
    ProcessScanEvent(ctx context.Context, deviceID uint64, req dto.ScanEventRequest) (dto.ScanEventResponse, error)
    ProcessSyncBatch(ctx context.Context, deviceID uint64, req dto.SyncBatchRequest) (dto.SyncBatchResponse, error)
    ManualCheckIn(ctx context.Context, actorUserID uint64, req dto.ManualCheckInRequest) (dto.AttendanceDetailResponse, error)
    GetSummary(ctx context.Context, date time.Time, classID *uint64) (dto.AttendanceSummaryResponse, error)
    ListAttendance(ctx context.Context, req dto.ListAttendanceRequest) (dto.ListAttendanceResponse, error)
    GetStudentHistory(ctx context.Context, studentID uint64, from, to time.Time) (dto.StudentHistoryResponse, error)
    ListNeedsReview(ctx context.Context, page, limit int) (dto.ListEventsResponse, error)
}
