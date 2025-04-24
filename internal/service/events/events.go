package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/storage"
	"github.com/Ilya-Repin/orchestra_api/internal/model"
	"github.com/Ilya-Repin/orchestra_api/internal/service"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Service struct {
	log           *slog.Logger
	eventStorage  EventStorage
	memberStorage MemberStorage
}

type MemberStorage interface {
	CheckIsApproved(ctx context.Context, id uuid.UUID) (bool, error)
}

type EventStorage interface {
	GetEvents(ctx context.Context, eventType *int, begin, end *time.Time) ([]model.Event, error)
	GetUpcomingEvents(ctx context.Context) ([]model.Event, error)
	GetAvailableEvents(ctx context.Context, memberID uuid.UUID) ([]model.Event, error)
	GetRegisteredEvents(ctx context.Context, memberID uuid.UUID) ([]model.Event, error)
	GetEvent(ctx context.Context, id int) (model.Event, error)
	AddEvent(ctx context.Context, title, description string, evType int, evDate time.Time, location int, capacity int) (int, error)
	DeleteEvent(ctx context.Context, id int) error
	UpdateEvent(ctx context.Context, id int, title, description string, evType int, evDate time.Time, location int, capacity int) error
}

func New(log *slog.Logger, eventStorage EventStorage, memberStorage MemberStorage) *Service {
	return &Service{log: log.With("component", "service"), eventStorage: eventStorage, memberStorage: memberStorage}
}

func (s *Service) AddEvent(ctx context.Context, title, description string, evType int, evDate time.Time, location int, capacity int) (int, error) {
	const op = "events.Service.AddEvent"

	log := s.log.With(slog.String("op", op))
	log.Info("adding new event")

	id, err := s.eventStorage.AddEvent(ctx, title, description, evType, evDate, location, capacity)
	if err != nil {
		log.Error("failed to add event", "error", err)
		return 0, fmt.Errorf("%s: %w", op, service.ErrFailedToAdd)
	}

	log.Info("event added", "id", id)
	return id, nil
}

func (s *Service) GetEvent(ctx context.Context, id int) (model.Event, error) {
	const op = "events.Service.GetEvent"

	log := s.log.With(slog.String("op", op))
	log.Info("getting event", "id", id)

	event, err := s.eventStorage.GetEvent(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			log.Error("event not found", "error", err)
			return model.Event{}, fmt.Errorf("%s: %w", op, service.ErrEventNotFound)
		}
		log.Error("failed to get event", "error", err)
		return model.Event{}, fmt.Errorf("%s: %w", op, service.ErrFailedToGet)
	}

	return event, nil
}

func (s *Service) GetEvents(ctx context.Context, eventType *int, begin, end *time.Time) ([]model.Event, error) {
	const op = "events.Service.GetEvents"

	log := s.log.With(slog.String("op", op))
	log.Info("getting events")

	events, err := s.eventStorage.GetEvents(ctx, eventType, begin, end)
	if err != nil {
		log.Error("failed to get events", "error", err)
		return nil, fmt.Errorf("%s: %w", op, service.ErrFailedToGetEvents)
	}

	return events, nil
}

func (s *Service) GetUpcomingEvents(ctx context.Context) ([]model.Event, error) {
	const op = "events.Service.GetUpcomingEvents"

	log := s.log.With(slog.String("op", op))
	log.Info("getting upcoming events")

	events, err := s.eventStorage.GetUpcomingEvents(ctx)
	if err != nil {
		log.Error("failed to get events", "error", err)
		return nil, fmt.Errorf("%s: %w", op, service.ErrFailedToGetUpcoming)
	}

	return events, nil
}

func (s *Service) GetAvailableEvents(ctx context.Context, memberID uuid.UUID) ([]model.Event, error) {
	const op = "events.Service.GetAvailableEvents"
	log := s.log.With(slog.String("op", op), slog.String("member_id", memberID.String()))

	approved, err := s.memberStorage.CheckIsApproved(ctx, memberID)
	if err != nil {
		if errors.Is(err, storage.ErrMemberNotFound) {
			log.Error("member not found", "error", err)
			return nil, fmt.Errorf("%s: %w", op, service.ErrMemberNotFound)
		}
		log.Error("failed to check approval", "error", err)
		return nil, fmt.Errorf("%s: %w", op, service.ErrFailedToGetAvailable)
	}

	if !approved {
		log.Warn("member is not approved")
		return nil, service.ErrMemberNotApproved
	}

	events, err := s.eventStorage.GetAvailableEvents(ctx, memberID)
	if err != nil {
		log.Error("failed to get available events", "error", err)
		return nil, fmt.Errorf("%s: %w", op, service.ErrFailedToGetAvailable)
	}

	log.Info("fetched available events", "count", len(events))
	return events, nil
}

func (s *Service) GetRegisteredEvents(ctx context.Context, memberID uuid.UUID) ([]model.Event, error) {
	const op = "events.Service.GetRegisteredEvents"
	log := s.log.With(slog.String("op", op), slog.String("member_id", memberID.String()))

	approved, err := s.memberStorage.CheckIsApproved(ctx, memberID)
	if err != nil {
		if errors.Is(err, storage.ErrMemberNotFound) {
			log.Error("member not found", "error", err)
			return nil, fmt.Errorf("%s: %w", op, service.ErrMemberNotFound)
		}
		log.Error("failed to check approval", "error", err)
		return nil, fmt.Errorf("%s: %w", op, service.ErrFailedToGetRegistered)
	}

	if !approved {
		log.Warn("member is not approved")
		return nil, service.ErrMemberNotApproved
	}

	events, err := s.eventStorage.GetRegisteredEvents(ctx, memberID)
	if err != nil {
		log.Error("failed to get registered events", "error", err)
		return nil, fmt.Errorf("%s: %w", op, service.ErrFailedToGetRegistered)
	}

	log.Info("fetched registered events", "count", len(events))
	return events, nil
}

func (s *Service) DeleteEvent(ctx context.Context, id int) error {
	const op = "events.Service.DeleteEvent"

	log := s.log.With(slog.String("op", op))
	log.Info("deleting event", "id", id)

	err := s.eventStorage.DeleteEvent(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			log.Error("event not found", "error", err)
			return fmt.Errorf("%s: %w", op, service.ErrEventNotFound)
		}

		log.Error("failed to delete event", "error", err)
		return fmt.Errorf("%s: %w", op, service.ErrFailedToDelete)
	}

	return nil
}

func (s *Service) UpdateEvent(ctx context.Context, id int, title, description string, evType int, evDate time.Time, location int, capacity int) error {
	const op = "events.Service.UpdateEvent"

	log := s.log.With(slog.String("op", op))
	log.Info("updating event", "id", id)

	err := s.eventStorage.UpdateEvent(ctx, id, title, description, evType, evDate, location, capacity)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			log.Error("event not found", "error", err)
			return fmt.Errorf("%s: %w", op, service.ErrEventNotFound)
		}

		log.Error("failed to update event", "error", err)
		return fmt.Errorf("%s: %w", op, service.ErrFailedToUpdate)
	}

	log.Info("event updated successfully", "id", id)
	return nil
}
