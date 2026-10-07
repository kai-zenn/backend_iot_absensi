package repository

import "context"

type RepoSet struct {
	AttendanceEvent IAttendanceEventRepository
	Attendance      IAttendanceRepository
	Card            ICardRepository
	Student         IStudentRepository
	Class           IClassRepository
	Rule            IRuleRepository
	Device          IDeviceRepository
	User            IUserRepository
	Session         ISessionRepository
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(repos RepoSet) error) error
}
