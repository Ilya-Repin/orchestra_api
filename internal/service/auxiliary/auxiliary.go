package auxiliary

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/storage"
	"github.com/Ilya-Repin/orchestra_api/internal/model"
	"github.com/Ilya-Repin/orchestra_api/internal/service"
	"log/slog"
)

type Service struct {
	log        *slog.Logger
	auxStorage AuxStorage
}

type AuxStorage interface {
	GetEventTypes(ctx context.Context) ([]model.EventType, error)
	GetLocations(ctx context.Context) ([]model.Location, error)
	GetLocation(ctx context.Context, id int) (model.Location, error)
	GetEventType(ctx context.Context, id int) (model.EventType, error)
	GetOrchestraInfo(ctx context.Context, key string) (model.OrchestraInfo, error)
	AddEventType(ctx context.Context, name, description string) (int, error)
	AddLocation(ctx context.Context, name, route, features string) (int, error)
	AddOrchestraInfo(ctx context.Context, key, value string) error
}

func New(log *slog.Logger, storage AuxStorage) *Service {
	return &Service{log: log.With("component", "service"), auxStorage: storage}
}

func (s *Service) GetEventTypes(ctx context.Context) ([]model.EventType, error) {
	const op = "auxiliary.Service.GetEventTypes"
	log := s.log.With(slog.String("op", op))

	types, err := s.auxStorage.GetEventTypes(ctx)
	if err != nil {
		log.Error("failed to get event types", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return types, nil
}

func (s *Service) GetLocations(ctx context.Context) ([]model.Location, error) {
	const op = "auxiliary.Service.GetLocations"
	log := s.log.With(slog.String("op", op))

	locs, err := s.auxStorage.GetLocations(ctx)
	if err != nil {
		log.Error("failed to get locations", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return locs, nil
}

func (s *Service) GetLocation(ctx context.Context, id int) (model.Location, error) {
	const op = "auxiliary.Service.GetLocation"
	log := s.log.With(slog.String("op", op), slog.Int("id", id))

	loc, err := s.auxStorage.GetLocation(ctx, id)
	if err != nil {
		log.Error("failed to get location", "error", err)
		if errors.Is(err, storage.ErrLocationNotFound) {
			return model.Location{}, fmt.Errorf("%s: %w", op, service.ErrMetaNotFound)
		}
		return model.Location{}, fmt.Errorf("%s: %w", op, err)
	}

	return loc, nil
}

func (s *Service) GetEventType(ctx context.Context, id int) (model.EventType, error) {
	const op = "auxiliary.Service.GetEventType"
	log := s.log.With(slog.String("op", op), slog.Int("id", id))

	typeData, err := s.auxStorage.GetEventType(ctx, id)
	if err != nil {
		log.Error("failed to get event type", "error", err)
		if errors.Is(err, storage.ErrEventTypeNotFound) {
			return model.EventType{}, fmt.Errorf("%s: %w", op, service.ErrMetaNotFound)
		}
		return model.EventType{}, fmt.Errorf("%s: %w", op, err)
	}

	return typeData, nil
}

func (s *Service) GetOrchestraInfo(ctx context.Context, key string) (model.OrchestraInfo, error) {
	const op = "auxiliary.Service.GetOrchestraInfo"
	log := s.log.With(slog.String("op", op), slog.String("key", key))

	info, err := s.auxStorage.GetOrchestraInfo(ctx, key)
	if err != nil {
		if errors.Is(err, storage.ErrInfoNotFound) {
			log.Error("info not found", "key", key, "error", err)
			return model.OrchestraInfo{}, fmt.Errorf("%s: %w", op, service.ErrInfoNotFound)
		}
		log.Error("failed to get orchestra info", "error", err)
		return model.OrchestraInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return info, nil
}

func (s *Service) AddEventType(ctx context.Context, name, description string) (int, error) {
	const op = "auxiliary.Service.AddEventType"
	log := s.log.With(slog.String("op", op))

	id, err := s.auxStorage.AddEventType(ctx, name, description)
	if err != nil {
		log.Error("failed to add event type", "error", err)
		return 0, fmt.Errorf("%s: %w", op, service.ErrFailedToSaveMeta)
	}

	return id, nil
}

func (s *Service) AddLocation(ctx context.Context, name, route, features string) (int, error) {
	const op = "auxiliary.Service.AddLocation"
	log := s.log.With(slog.String("op", op))

	id, err := s.auxStorage.AddLocation(ctx, name, route, features)
	if err != nil {
		log.Error("failed to add location", "error", err)
		return 0, fmt.Errorf("%s: %w", op, service.ErrFailedToSaveMeta)
	}

	return id, nil
}

func (s *Service) AddOrchestraInfo(ctx context.Context, key, value string) error {
	const op = "auxiliary.Service.AddOrchestraInfo"
	log := s.log.With(slog.String("op", op), slog.String("key", key))

	err := s.auxStorage.AddOrchestraInfo(ctx, key, value)
	if err != nil {
		log.Error("failed to add orchestra info", "error", err)
		return fmt.Errorf("%s: %w", op, service.ErrFailedToSaveMeta)
	}

	return nil
}
