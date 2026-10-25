package service

import (
	"context"
	"database/sql"
	"fmt"
	"mosanda-vacation-attendace/models"
	"mosanda-vacation-attendace/pkg/osioutil"
	"mosanda-vacation-attendace/store"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

type Service struct {
	store  *store.Store
	dbConn *sql.DB
}

func NewService(parStore *store.Store, parDBConn *sql.DB) *Service {
	return &Service{
		store:  parStore,
		dbConn: parDBConn,
	}
}

func (objService *Service) AuthenticateUser(parEmail, parPassword string) (int, error) {
	var pvHashPassowd string
	var pvUserID int

	psRow := objService.dbConn.QueryRow("SELECT id , hash_password FROM users WHERE email = $1", parEmail)

	pvError := psRow.Scan(&pvUserID, &pvHashPassowd)
	if pvError != nil {
		return -1, fmt.Errorf("SVC0001 %v", pvError)
	}

	// 4. Compare Password With Hash
	pfMatch, pvError := argon2id.ComparePasswordAndHash(parPassword, pvHashPassowd)
	if pvError != nil || !pfMatch {
		return -1, fmt.Errorf("SVC0002 Failed to login with username and password %v", pvError)
	}

	return pvUserID, nil
}

func (objService *Service) StartSession(parUserID int) (*models.Session, error) {
	pvSesssionID := uuid.New()
	pvExpiresAt := time.Now().Add(30 * time.Minute)

	_, pvError := objService.dbConn.Exec("INSERT INTO sessions (id,user_id,expires_at,created_at) VALUES ($1,$2,$3,$4)",
		pvSesssionID,
		parUserID,
		pvExpiresAt,
		time.Now(),
	)
	if pvError != nil {
		return nil, osioutil.GetRefError("SVC0201", pvError)
	}

	return &models.Session{
		ID:        pvSesssionID,
		ExpiresAt: pvExpiresAt,
	}, nil
}

func (objService *Service) GetUserFromSession(parSessionID uuid.UUID) (*models.AuthenticatedUser, error) {
	psUser := &models.AuthenticatedUser{}
	pvError := objService.dbConn.QueryRow(
		`SELECT users.id, COALESCE(users.email, '')
		FROM sessions
		INNER JOIN users ON users.id = sessions.user_id
		WHERE sessions.id = $1 AND sessions.expires_at > CURRENT_TIMESTAMP`,
		parSessionID,
	).Scan(&psUser.ID, &psUser.Email)
	if pvError != nil {
		return nil, fmt.Errorf("SVC0202: %w", pvError)
	}
	return psUser, nil
}

func (objService *Service) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	psResult, pvError := objService.dbConn.ExecContext(
		ctx,
		"DELETE FROM sessions WHERE expires_at <= CURRENT_TIMESTAMP",
	)
	if pvError != nil {
		return 0, fmt.Errorf("SVC0203: %w", pvError)
	}
	pvDeletedCount, pvError := psResult.RowsAffected()
	if pvError != nil {
		return 0, fmt.Errorf("SVC0204: %w", pvError)
	}
	return pvDeletedCount, nil
}

func (objService *Service) CreateVacationBatch(parBatch *models.VacationBatch) (*models.VacationBatch, error) {
	psBatch := &models.VacationBatch{}
	pvError := objService.dbConn.QueryRow(
		`INSERT INTO vacation_batches (name, start_date, end_date, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, start_date, end_date, created_by, created_at, updated_at`,
		parBatch.Name,
		parBatch.StartDate,
		parBatch.EndDate,
		parBatch.CreatedBy,
	).Scan(
		&psBatch.ID,
		&psBatch.Name,
		&psBatch.StartDate,
		&psBatch.EndDate,
		&psBatch.CreatedBy,
		&psBatch.CreatedAt,
		&psBatch.UpdatedAt,
	)
	if pvError != nil {
		return nil, fmt.Errorf("SVC0301: %w", pvError)
	}
	return psBatch, nil
}

func (objService *Service) GetVacationBatches() ([]models.VacationBatch, error) {
	psRows, pvError := objService.dbConn.Query(
		`SELECT id, name, start_date, end_date, created_by, created_at, updated_at
		FROM vacation_batches
		ORDER BY id`,
	)
	if pvError != nil {
		return nil, fmt.Errorf("SVC0302: %w", pvError)
	}
	defer psRows.Close()

	pvBatches := make([]models.VacationBatch, 0)
	for psRows.Next() {
		psBatch := models.VacationBatch{}
		pvError = psRows.Scan(
			&psBatch.ID,
			&psBatch.Name,
			&psBatch.StartDate,
			&psBatch.EndDate,
			&psBatch.CreatedBy,
			&psBatch.CreatedAt,
			&psBatch.UpdatedAt,
		)
		if pvError != nil {
			return nil, fmt.Errorf("SVC0303: %w", pvError)
		}
		pvBatches = append(pvBatches, psBatch)
	}
	if pvError = psRows.Err(); pvError != nil {
		return nil, fmt.Errorf("SVC0304: %w", pvError)
	}
	return pvBatches, nil
}

func (objService *Service) GetVacationBatch(parBatchID int) (*models.VacationBatch, error) {
	psBatch := &models.VacationBatch{}
	pvError := objService.dbConn.QueryRow(
		`SELECT id, name, start_date, end_date, created_by, created_at, updated_at
		FROM vacation_batches
		WHERE id = $1`,
		parBatchID,
	).Scan(
		&psBatch.ID,
		&psBatch.Name,
		&psBatch.StartDate,
		&psBatch.EndDate,
		&psBatch.CreatedBy,
		&psBatch.CreatedAt,
		&psBatch.UpdatedAt,
	)
	if pvError != nil {
		return nil, fmt.Errorf("SVC0305: %w", pvError)
	}
	return psBatch, nil
}

func (objService *Service) UpdateVacationBatch(parBatchID int, parBatch *models.VacationBatch) (*models.VacationBatch, error) {
	psBatch := &models.VacationBatch{}
	pvError := objService.dbConn.QueryRow(
		`UPDATE vacation_batches
		SET name = $1, start_date = $2, end_date = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, name, start_date, end_date, created_by, created_at, updated_at`,
		parBatch.Name,
		parBatch.StartDate,
		parBatch.EndDate,
		parBatchID,
	).Scan(
		&psBatch.ID,
		&psBatch.Name,
		&psBatch.StartDate,
		&psBatch.EndDate,
		&psBatch.CreatedBy,
		&psBatch.CreatedAt,
		&psBatch.UpdatedAt,
	)
	if pvError != nil {
		return nil, fmt.Errorf("SVC0306: %w", pvError)
	}
	return psBatch, nil
}

func (objService *Service) DeleteVacationBatch(parBatchID int) error {
	psResult, pvError := objService.dbConn.Exec(
		"DELETE FROM vacation_batches WHERE id = $1",
		parBatchID,
	)
	if pvError != nil {
		return fmt.Errorf("SVC0307: %w", pvError)
	}
	pvRowsAffected, pvError := psResult.RowsAffected()
	if pvError != nil {
		return fmt.Errorf("SVC0308: %w", pvError)
	}
	if pvRowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

type attendanceScanner interface {
	Scan(...any) error
}

func scanAttendance(parScanner attendanceScanner) (*models.Attendance, error) {
	psAttendance := &models.Attendance{}
	var pvCheckedInUserID sql.NullInt64
	var pvCheckedInUserEmail sql.NullString
	pvError := parScanner.Scan(
		&psAttendance.ID,
		&psAttendance.BatchID,
		&psAttendance.Name,
		&psAttendance.Location,
		&psAttendance.PhoneNumber,
		&psAttendance.MilitaryNumber,
		&psAttendance.Status,
		&psAttendance.Notes,
		&psAttendance.StartDate,
		&psAttendance.EndDate,
		&psAttendance.CheckedInAt,
		&psAttendance.CheckedInByID,
		&psAttendance.CreatedBy,
		&psAttendance.CreatedAt,
		&psAttendance.UpdatedAt,
		&pvCheckedInUserID,
		&pvCheckedInUserEmail,
	)
	if pvError != nil {
		return nil, pvError
	}
	if pvCheckedInUserID.Valid {
		psAttendance.CheckedInBy = &models.AuthenticatedUser{
			ID:    int(pvCheckedInUserID.Int64),
			Email: pvCheckedInUserEmail.String,
		}
	}
	return psAttendance, nil
}

const attendanceSelectColumns = `attendance_row.id, attendance_row.batch_id, attendance_row.name,
	attendance_row.location, attendance_row.phone_number, attendance_row.military_number,
	COALESCE(attendance_row.status::text, 'absent'), attendance_row.notes, attendance_row.start_date,
	attendance_row.end_date, attendance_row.checked_in_at, attendance_row.checked_in_by,
	attendance_row.created_by, attendance_row.created_at, attendance_row.updated_at,
	check_in_user.id, COALESCE(check_in_user.email, '')`

func (objService *Service) CreateAttendance(parAttendance *models.Attendance) (*models.Attendance, error) {
	psRow := objService.dbConn.QueryRow(
		`WITH created_attendance AS (
			INSERT INTO attendance
			(batch_id, name, location, phone_number, military_number, status, notes, start_date, end_date,
			checked_in_at, checked_in_by, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			RETURNING *
		)
		SELECT `+attendanceSelectColumns+`
		FROM created_attendance AS attendance_row
		LEFT JOIN users AS check_in_user ON check_in_user.id = attendance_row.checked_in_by`,
		parAttendance.BatchID,
		parAttendance.Name,
		parAttendance.Location,
		parAttendance.PhoneNumber,
		parAttendance.MilitaryNumber,
		parAttendance.Status,
		parAttendance.Notes,
		parAttendance.StartDate,
		parAttendance.EndDate,
		parAttendance.CheckedInAt,
		parAttendance.CheckedInByID,
		parAttendance.CreatedBy,
	)
	psAttendance, pvError := scanAttendance(psRow)
	if pvError != nil {
		return nil, fmt.Errorf("SVC0401: %w", pvError)
	}
	return psAttendance, nil
}

func (objService *Service) GetAttendances() ([]models.Attendance, error) {
	psRows, pvError := objService.dbConn.Query(
		`SELECT ` + attendanceSelectColumns + `
		FROM attendance AS attendance_row
		LEFT JOIN users AS check_in_user ON check_in_user.id = attendance_row.checked_in_by
		ORDER BY attendance_row.id`,
	)
	if pvError != nil {
		return nil, fmt.Errorf("SVC0402: %w", pvError)
	}
	defer psRows.Close()

	pvAttendances := make([]models.Attendance, 0)
	for psRows.Next() {
		psAttendance, pvError := scanAttendance(psRows)
		if pvError != nil {
			return nil, fmt.Errorf("SVC0403: %w", pvError)
		}
		pvAttendances = append(pvAttendances, *psAttendance)
	}
	if pvError = psRows.Err(); pvError != nil {
		return nil, fmt.Errorf("SVC0404: %w", pvError)
	}
	return pvAttendances, nil
}

func (objService *Service) GetAttendance(parAttendanceID int) (*models.Attendance, error) {
	psRow := objService.dbConn.QueryRow(
		`SELECT `+attendanceSelectColumns+`
		FROM attendance AS attendance_row
		LEFT JOIN users AS check_in_user ON check_in_user.id = attendance_row.checked_in_by
		WHERE attendance_row.id = $1`,
		parAttendanceID,
	)
	psAttendance, pvError := scanAttendance(psRow)
	if pvError != nil {
		return nil, fmt.Errorf("SVC0405: %w", pvError)
	}
	return psAttendance, nil
}

func (objService *Service) UpdateAttendance(parAttendanceID int, parAttendance *models.Attendance) (*models.Attendance, error) {
	psRow := objService.dbConn.QueryRow(
		`WITH updated_attendance AS (
			UPDATE attendance
			SET batch_id = $1, name = $2, phone_number = $3, military_number = $4, status = $5,
			notes = $6, start_date = $7, end_date = $8, checked_in_at = $9, checked_in_by = $10,
			location = $11, updated_at = CURRENT_TIMESTAMP
			WHERE id = $12
			RETURNING *
		)
		SELECT `+attendanceSelectColumns+`
		FROM updated_attendance AS attendance_row
		LEFT JOIN users AS check_in_user ON check_in_user.id = attendance_row.checked_in_by`,
		parAttendance.BatchID,
		parAttendance.Name,
		parAttendance.PhoneNumber,
		parAttendance.MilitaryNumber,
		parAttendance.Status,
		parAttendance.Notes,
		parAttendance.StartDate,
		parAttendance.EndDate,
		parAttendance.CheckedInAt,
		parAttendance.CheckedInByID,
		parAttendance.Location,
		parAttendanceID,
	)
	psAttendance, pvError := scanAttendance(psRow)
	if pvError != nil {
		return nil, fmt.Errorf("SVC0406: %w", pvError)
	}
	return psAttendance, nil
}

func (objService *Service) DeleteAttendance(parAttendanceID int) error {
	psResult, pvError := objService.dbConn.Exec(
		"DELETE FROM attendance WHERE id = $1",
		parAttendanceID,
	)
	if pvError != nil {
		return fmt.Errorf("SVC0407: %w", pvError)
	}
	pvRowsAffected, pvError := psResult.RowsAffected()
	if pvError != nil {
		return fmt.Errorf("SVC0408: %w", pvError)
	}
	if pvRowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
