package repository

import (
  "context"
  "github.com/kai-zenn/backend_iot_absensi/internal/dto"
  "github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type ICardRepository interface {
  FindByUID(ctx context.Context, uid string) (*entity.RFIDCard, error)
  FindActiveByStudentID(ctx context.Context, studentID uint64) (*entity.RFIDCard, error)

  Create(ctx context.Context, card *entity.RFIDCard) error

    // UpdateStatus: db.Model(&entity.RFIDCard{}).Where("id = ?", cardID).
    //   Update("status", status)
  UpdateStatus(ctx context.Context, cardID uint64, status string) error

    // ListUnregisteredUIDs mengambil UID distinct dari attendance_events ber-result
    // UNKNOWN_CARD dalam N hari terakhir — join sederhana, BUKAN kolom fisik di rfid_cards.
  ListUnregisteredUIDs(ctx context.Context, sinceDays int) ([]dto.UnregisteredCardSummary, error)

    // ListActiveUIDs dipakai untuk snapshot kartu ke device (lihat 4c):
    //   db.Model(&entity.RFIDCard{}).Joins("JOIN students ON students.id = rfid_cards.student_id").
    //     Where("rfid_cards.status = ? AND students.status = ?", "ACTIVE", "ACTIVE").
    //     Pluck("rfid_cards.uid", &uids)
  ListActiveUIDs(ctx context.Context) ([]string, error)
}
