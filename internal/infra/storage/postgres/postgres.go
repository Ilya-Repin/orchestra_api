package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Ilya-Repin/orchestra_api/internal/config"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/storage"
	"github.com/Ilya-Repin/orchestra_api/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

type PostgresStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func InitDB(cfg *config.StorageConfig) (db *sql.DB, err error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
		cfg.Sslmode,
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
func (s *PostgresStorage) AddMember(ctx context.Context, fullName, email, phone string) (id uuid.UUID, err error) {
	const op = "infra.storage.postgres.AddMember"

	if !model.IsValidEmail(email) {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, storage.ErrInvalidEmail)
	}

	if !model.IsValidPhone(phone) {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, storage.ErrInvalidPhone)
	}

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO club_members (full_name, email, phone) VALUES ($1, $2, $3) RETURNING id;")
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, fullName, email, phone).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				if isEmailDuplicate(ctx, s.db, email) {
					return uuid.UUID{}, fmt.Errorf("%s: %w", op, storage.ErrEmailDuplicate)
				}
			}
			if isPhoneDuplicate(ctx, s.db, phone) {
				return uuid.UUID{}, fmt.Errorf("%s: %w", op, storage.ErrPhoneDuplicate)
			}
		}
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func isEmailDuplicate(ctx context.Context, db *sql.DB, email string) bool {
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM club_members WHERE email = $1", email).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func isPhoneDuplicate(ctx context.Context, db *sql.DB, phone string) bool {
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM club_members WHERE phone = $1", phone).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (s *PostgresStorage) GetMember(ctx context.Context, id uuid.UUID) (model.Member, error) {
	const op = "infra.storage.postgres.GetMember"

	stmt, err := s.db.PrepareContext(ctx, "SELECT id, full_name, email, phone, status, created_at, updated_at FROM club_members WHERE id = $1;")
	if err != nil {
		return model.Member{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var member model.Member
	err = stmt.QueryRowContext(ctx, id).Scan(
		&member.ID,
		&member.FullName,
		&member.Email,
		&member.Phone,
		&member.Status,
		&member.CreatedAt,
		&member.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Member{}, storage.ErrMemberNotFound
		}
		return model.Member{}, fmt.Errorf("%s: %w", op, err)
	}

	return member, nil
}

func (s *PostgresStorage) GetMembers(ctx context.Context) ([]model.Member, error) {
	const op = "infra.storage.postgres.GetMembers"

	query := `
		SELECT id, full_name, email, phone, status, created_at
		FROM club_members
		ORDER BY created_at DESC;
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var members []model.Member
	for rows.Next() {
		var m model.Member
		err := rows.Scan(&m.ID, &m.FullName, &m.Email, &m.Phone, &m.Status, &m.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return members, nil
}

func (s *PostgresStorage) GetMembersWithStatus(ctx context.Context, status model.MemberStatus) ([]model.Member, error) {
	const op = "infra.storage.postgres.GetMembersWithStatus"

	query := `
		SELECT id, full_name, email, phone, status, created_at
		FROM club_members
		WHERE status = $1
		ORDER BY created_at DESC;
	`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var members []model.Member
	for rows.Next() {
		var m model.Member
		if err := rows.Scan(&m.ID, &m.FullName, &m.Email, &m.Phone, &m.Status, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return members, nil
}

func (s *PostgresStorage) DeleteMember(ctx context.Context, id uuid.UUID) error {
	const op = "infra.storage.postgres.DeleteMember"

	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM club_members WHERE id = $1;")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if affected == 0 {
		return storage.ErrMemberNotFound
	}

	return nil
}

func (s *PostgresStorage) UpdateMember(ctx context.Context, id uuid.UUID, fullName, email, phone string) error {
	const op = "infra.storage.postgres.UpdateMember"

	stmt, err := s.db.PrepareContext(ctx, `
		UPDATE club_members 
		SET full_name = $1, email = $2, phone = $3 
		WHERE id = $4;
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, fullName, email, phone, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if affected == 0 {
		return storage.ErrMemberNotFound
	}

	return nil
}

func (s *PostgresStorage) UpdateMemberStatus(ctx context.Context, id uuid.UUID, status model.MemberStatus) error {
	const op = "infra.storage.postgres.UpdateMemberStatus"

	stmt, err := s.db.PrepareContext(ctx, "UPDATE club_members SET status=$1 WHERE id=$2;")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, status, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows affected: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrMemberNotFound)
	}

	return nil
}

func (s *PostgresStorage) CheckIsApproved(ctx context.Context, id uuid.UUID) (bool, error) {
	const op = "infra.storage.postgres.CheckIsApproved"

	query := "SELECT status FROM club_members WHERE id = $1"

	var status string

	err := s.db.QueryRowContext(ctx, query, id).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, storage.ErrMemberNotFound
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return status == string(model.StatusApproved), nil
}

func (s *PostgresStorage) GetEvents(
	ctx context.Context,
	eventType *int,
	begin, end *time.Time,
) ([]model.Event, error) {
	const op = "infra.storage.postgres.GetEvents"

	var (
		events []model.Event
		args   []interface{}
		conds  []string
		argNum = 1
	)

	query := `
		SELECT 
			e.id, e.title, e.description, e.event_date, e.capacity, e.created_at, e.updated_at,
			et.id, et.name, et.description,
			l.id, l.name, l.route, l.features
		FROM events e
		JOIN event_types et ON e.event_type = et.id
		JOIN locations l ON e.location = l.id
	`

	if eventType != nil {
		conds = append(conds, fmt.Sprintf("e.event_type = $%d", argNum))
		args = append(args, *eventType)
		argNum++
	}
	if begin != nil {
		conds = append(conds, fmt.Sprintf("e.event_date >= $%d", argNum))
		args = append(args, *begin)
		argNum++
	}
	if end != nil {
		conds = append(conds, fmt.Sprintf("e.event_date <= $%d", argNum))
		args = append(args, *end)
		argNum++
	}

	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	query += " ORDER BY e.event_date ASC"

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var ev model.Event
		err := rows.Scan(
			&ev.ID, &ev.Title, &ev.Description, &ev.EventDate, &ev.Capacity, &ev.CreatedAt, &ev.UpdatedAt,
			&ev.EventType.ID, &ev.EventType.Name, &ev.EventType.Description,
			&ev.Location.ID, &ev.Location.Name, &ev.Location.Route, &ev.Location.Features,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		events = append(events, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

func (s *PostgresStorage) GetUpcomingEvents(ctx context.Context) ([]model.Event, error) {
	const op = "infra.storage.postgres.GetUpcomingEvents"

	query := `
		SELECT 
		e.id, e.title, e.description, e.event_date, e.capacity, e.created_at, e.updated_at,
		et.id, et.name, et.description,
		l.id, l.name, l.route, l.features
		FROM events e
		JOIN event_types et ON e.event_type = et.id
		JOIN locations l ON e.location = l.id
		WHERE e.event_date >= CURRENT_TIMESTAMP
		ORDER BY e.event_date ASC;
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var events []model.Event

	for rows.Next() {
		var ev model.Event

		err := rows.Scan(&ev.ID, &ev.Title, &ev.Description, &ev.EventDate, &ev.Capacity, &ev.CreatedAt, &ev.UpdatedAt,
			&ev.EventType.ID, &ev.EventType.Name, &ev.EventType.Description,
			&ev.Location.ID, &ev.Location.Name, &ev.Location.Route, &ev.Location.Features,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		events = append(events, ev)
	}
	return events, nil
}

func (s *PostgresStorage) GetAvailableEvents(ctx context.Context, memberID uuid.UUID) ([]model.Event, error) {
	const op = "infra.storage.postgres.GetAvailableEvents"

	query := `
		SELECT 
		e.id, e.title, e.description, e.event_date, e.capacity, e.created_at, e.updated_at,
		et.id, et.name, et.description,
		l.id, l.name, l.route, l.features
		FROM events e
		JOIN event_types et ON e.event_type = et.id
		JOIN locations l ON e.location = l.id
		WHERE e.event_date >= CURRENT_TIMESTAMP
		AND e.id NOT IN (
			SELECT reg.event_id
			FROM registrations reg
			WHERE reg.user_id = $1 AND reg.registration_status = 'registered'
		)
		ORDER BY e.event_date ASC;
	`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var events []model.Event

	for rows.Next() {
		var ev model.Event

		err := rows.Scan(&ev.ID, &ev.Title, &ev.Description, &ev.EventDate, &ev.Capacity, &ev.CreatedAt, &ev.UpdatedAt,
			&ev.EventType.ID, &ev.EventType.Name, &ev.EventType.Description,
			&ev.Location.ID, &ev.Location.Name, &ev.Location.Route, &ev.Location.Features,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		events = append(events, ev)
	}
	return events, nil
}

func (s *PostgresStorage) GetRegisteredEvents(ctx context.Context, memberID uuid.UUID) ([]model.Event, error) {
	const op = "infra.storage.postgres.GetRegisteredEvents"

	query := `
		SELECT 
		e.id, e.title, e.description, e.event_date, e.capacity, e.created_at, e.updated_at,
		et.id, et.name, et.description,
		l.id, l.name, l.route, l.features
		FROM events e
		JOIN event_types et ON e.event_type = et.id
		JOIN locations l ON e.location = l.id
		WHERE e.event_date >= CURRENT_TIMESTAMP
		AND e.id IN (
			SELECT reg.event_id
			FROM registrations reg
			WHERE reg.user_id = $1
		)
		ORDER BY e.event_date ASC;
	`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var events []model.Event

	for rows.Next() {
		var ev model.Event

		err := rows.Scan(&ev.ID, &ev.Title, &ev.Description, &ev.EventDate, &ev.Capacity, &ev.CreatedAt, &ev.UpdatedAt,
			&ev.EventType.ID, &ev.EventType.Name, &ev.EventType.Description,
			&ev.Location.ID, &ev.Location.Name, &ev.Location.Route, &ev.Location.Features,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		events = append(events, ev)
	}
	return events, nil
}

func (s *PostgresStorage) GetEvent(ctx context.Context, id int) (model.Event, error) {
	const op = "infra.storage.postgres.GetEvent"

	query := `
		SELECT e.id, e.title, e.description, e.event_date, e.capacity, e.created_at, e.updated_at,
		       l.id, l.name, et.id, et.name
		FROM events e
		JOIN locations l ON e.location = l.id
		JOIN event_types et ON e.event_type = et.id
		WHERE e.id = $1;
	`

	var ev model.Event

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&ev.ID, &ev.Title, &ev.Description, &ev.EventDate, &ev.Capacity, &ev.CreatedAt, &ev.UpdatedAt,
		&ev.Location.ID, &ev.Location.Name,
		&ev.EventType.ID, &ev.EventType.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Event{}, storage.ErrEventNotFound
		}
		return model.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	return ev, nil
}

func (s *PostgresStorage) AddEvent(ctx context.Context, title, description string, evType int, evDate time.Time, location int, capacity int) (int, error) {
	const op = "infra.storage.postgres.AddEvent"

	query := `
		INSERT INTO events (title, description, event_type, event_date, location, capacity)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	var id int
	err := s.db.QueryRowContext(ctx, query,
		title, description, evType, evDate, location, capacity,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *PostgresStorage) DeleteEvent(ctx context.Context, id int) error {
	const op = "infra.storage.postgres.DeleteEvent"

	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM events WHERE id = $1;")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *PostgresStorage) UpdateEvent(ctx context.Context, id int, title, description string, evType int, evDate time.Time, location int, capacity int) error {
	const op = "infra.storage.postgres.UpdateEvent"

	query := `
		UPDATE events
		SET title = $1, description = $2, event_type = $3, event_date = $4, location = $5, capacity = $6
		WHERE id = $7
	`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, title, description, evType, evDate, location, capacity, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rows == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *PostgresStorage) RegisterForEvent(ctx context.Context, memberID uuid.UUID, eventID int) (string, error) {
	const op = "infra.storage.postgres.RegisterForEvent"

	query := `
	WITH event_data AS (
		SELECT 
			e.capacity,
			(SELECT COUNT(*) FROM registrations r WHERE r.event_id = e.id AND r.registration_status = 'registered') AS current_count
		FROM events e
		WHERE e.id = $1
	),
	upd AS (
		UPDATE registrations
		SET registration_status = 'registered'
		FROM event_data
		WHERE registrations.user_id = $2
		  AND registrations.event_id = $1
		  AND registrations.registration_status = 'cancelled'
		  AND event_data.current_count < event_data.capacity
		RETURNING registrations.registration_status
	),
	ins AS (
		INSERT INTO registrations(user_id, event_id, registration_status)
		SELECT $2, $1, 'registered'
		FROM event_data
		WHERE event_data.current_count < event_data.capacity
		  AND NOT EXISTS (
		      SELECT 1 FROM registrations
		      WHERE user_id = $2 AND event_id = $1
		  )
		RETURNING registration_status
	)
	SELECT registration_status FROM upd
	UNION ALL
	SELECT registration_status FROM ins;
	`

	var status string
	err := s.db.QueryRowContext(ctx, query, eventID, memberID).Scan(&status)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return "", fmt.Errorf("%s: %w", op, storage.ErrRegAlreadyExists)
			}
		}

		var exists bool
		errEvent := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)", eventID).Scan(&exists)
		if errEvent != nil {
			return "", fmt.Errorf("%s: %w", op, errEvent)
		}
		if !exists {
			return "", fmt.Errorf("%s: %w", op, storage.ErrEventNotFound)
		}

		errMember := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM club_members WHERE id = $1)", memberID).Scan(&exists)
		if errMember != nil {
			return "", fmt.Errorf("%s: %w", op, errMember)
		}
		if !exists {
			return "", fmt.Errorf("%s: %w", op, storage.ErrMemberNotFound)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrEventFull)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (s *PostgresStorage) CancelRegistration(ctx context.Context, memberID uuid.UUID, eventID int) (string, error) {
	const op = "infra.storage.postgres.CancelRegistration"

	query := `
		UPDATE registrations
		SET registration_status = 'cancelled'
		WHERE user_id = $1 AND event_id = $2
		RETURNING registration_status;
	`

	var status string
	err := s.db.QueryRowContext(ctx, query, memberID, eventID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrRegNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (s *PostgresStorage) GetRegistrationStatus(ctx context.Context, memberID uuid.UUID, eventID int) (string, error) {
	const op = "infra.storage.postgres.GetRegistrationStatus"

	query := `
		SELECT registration_status
		FROM registrations
		WHERE user_id = $1 AND event_id = $2;
	`

	var status string
	err := s.db.QueryRowContext(ctx, query, memberID, eventID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrRegNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (s *PostgresStorage) GetEventTypes(ctx context.Context) ([]model.EventType, error) {
	const op = "infra.storage.postgres.GetEventTypes"

	rows, err := s.db.QueryContext(ctx, "SELECT id, name, description FROM event_types")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var types []model.EventType
	for rows.Next() {
		var et model.EventType
		if err := rows.Scan(&et.ID, &et.Name, &et.Description); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		types = append(types, et)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return types, nil
}

func (s *PostgresStorage) GetLocations(ctx context.Context) ([]model.Location, error) {
	const op = "infra.storage.postgres.GetLocations"

	rows, err := s.db.QueryContext(ctx, "SELECT id, name, route, features FROM locations")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var locations []model.Location
	for rows.Next() {
		var loc model.Location
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Route, &loc.Features); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return locations, nil
}

func (s *PostgresStorage) GetLocation(ctx context.Context, id int) (model.Location, error) {
	const op = "infra.storage.postgres.GetLocation"

	query := `
		SELECT id, name, route, features
		FROM locations
		WHERE id = $1;
	`

	var loc model.Location
	err := s.db.QueryRowContext(ctx, query, id).Scan(&loc.ID, &loc.Name, &loc.Route, &loc.Features)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Location{}, fmt.Errorf("%s: %w", op, storage.ErrLocationNotFound)
		}
		return model.Location{}, fmt.Errorf("%s: %w", op, err)
	}

	return loc, nil
}

func (s *PostgresStorage) GetEventType(ctx context.Context, id int) (model.EventType, error) {
	const op = "infra.storage.postgres.GetEventType"

	query := `
		SELECT id, name, description
		FROM event_types
		WHERE id = $1;
	`

	var et model.EventType
	err := s.db.QueryRowContext(ctx, query, id).Scan(&et.ID, &et.Name, &et.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.EventType{}, fmt.Errorf("%s: %w", op, storage.ErrEventTypeNotFound)
		}
		return model.EventType{}, fmt.Errorf("%s: %w", op, err)
	}

	return et, nil
}

func (s *PostgresStorage) GetOrchestraInfo(ctx context.Context, key string) (model.OrchestraInfo, error) {
	const op = "infra.storage.postgres.GetOrchestraInfo"
	stmt, err := s.db.PrepareContext(ctx, "SELECT key, value FROM orchestra_info WHERE key = $1;")
	if err != nil {
		return model.OrchestraInfo{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var info model.OrchestraInfo
	err = stmt.QueryRowContext(ctx, key).Scan(
		&info.Key,
		&info.Value,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OrchestraInfo{}, storage.ErrInfoNotFound
		}
		return model.OrchestraInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return info, nil
}

func (s *PostgresStorage) AddEventType(ctx context.Context, name, description string) (id int, err error) {
	const op = "infra.storage.postgres.AddEventType"

	query := `
		INSERT INTO event_types (name, description)
		VALUES ($1, $2)
		RETURNING id;
	`

	err = s.db.QueryRowContext(ctx, query, name, description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *PostgresStorage) AddLocation(ctx context.Context, name, route, features string) (id int, err error) {
	const op = "infra.storage.postgres.AddLocation"

	query := `
		INSERT INTO locations (name, route, features)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err = s.db.QueryRowContext(ctx, query, name, route, features).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *PostgresStorage) AddOrchestraInfo(ctx context.Context, key, value string) error {
	const op = "infra.storage.postgres.AddOrchestraInfo"

	query := `
		INSERT INTO orchestra_info (info_key, info_value)
		VALUES ($1, $2)
		ON CONFLICT (info_key)
		DO UPDATE SET info_value = EXCLUDED.info_value;
	`

	_, err := s.db.ExecContext(ctx, query, key, value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
