package registrations

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/storage"
	"github.com/Ilya-Repin/orchestra_api/internal/service"
	"github.com/google/uuid"
	"log/slog"
)

type Service struct {
	log           *slog.Logger
	regStorage    RegStorage
	memberStorage MemberStorage
}

type MemberStorage interface {
	CheckIsApproved(ctx context.Context, id uuid.UUID) (bool, error)
}

type RegStorage interface {
	RegisterForEvent(ctx context.Context, memberID uuid.UUID, eventID int) (string, error)
	CancelRegistration(ctx context.Context, memberID uuid.UUID, eventID int) (string, error)
	GetRegistrationStatus(ctx context.Context, memberID uuid.UUID, eventID int) (string, error)
}

func New(log *slog.Logger, regStorage RegStorage, memberStorage MemberStorage) *Service {
	return &Service{
		log:           log.With("component", "service"),
		regStorage:    regStorage,
		memberStorage: memberStorage,
	}
}

func (s *Service) RegisterForEvent(ctx context.Context, memberID uuid.UUID, eventID int) (string, error) {
	const op = "registrations.Service.RegisterForEvent"
	log := s.log.With(slog.String("op", op), slog.String("member_id", memberID.String()), slog.Int("event_id", eventID))

	log.Info("attempting registration")

	approved, err := s.memberStorage.CheckIsApproved(ctx, memberID)
	if err != nil {
		if errors.Is(err, storage.ErrMemberNotFound) {
			log.Error("member not found", "error", err)
			return "", fmt.Errorf("%s: %w", op, service.ErrMemberNotFound)
		}
		log.Error("failed to check member approval", "error", err)
		return "", fmt.Errorf("%s: %w", op, service.ErrRegistrationFailed)
	}

	if !approved {
		log.Error("registration denied: member not approved")
		return "", fmt.Errorf("%s: %w", op, service.ErrMemberNotApproved)
	}

	status, err := s.regStorage.RegisterForEvent(ctx, memberID, eventID)
	if err != nil {
		if errors.Is(err, storage.ErrEventFull) {
			log.Error("event full", "error", err)
			return "", fmt.Errorf("%s: %w", op, service.ErrEventFull)
		}
		if errors.Is(err, storage.ErrEventNotFound) {
			log.Error("event not found", "error", err)
			return "", fmt.Errorf("%s: %w", op, service.ErrEventNotFound)
		}
		if errors.Is(err, storage.ErrRegAlreadyExists) {
			log.Error("registration already exists", "error", err)
			return "", fmt.Errorf("%s: %w", op, service.ErrRegAlreadyExists)
		}

		log.Error("failed to register", "error", err)
		return "", fmt.Errorf("%s: %w", op, service.ErrRegistrationFailed)
	}

	log.Info("registration successful")
	return status, nil
}

func (s *Service) CancelRegistration(ctx context.Context, memberID uuid.UUID, eventID int) (string, error) {
	const op = "registrations.Service.CancelRegistration"
	log := s.log.With(slog.String("op", op), slog.String("member_id", memberID.String()), slog.Int("event_id", eventID))

	log.Info("cancelling registration")

	status, err := s.regStorage.CancelRegistration(ctx, memberID, eventID)
	if err != nil {
		if errors.Is(err, storage.ErrRegNotFound) {
			log.Error("registration not found", "error", err)
			return "", fmt.Errorf("%s: %w", op, service.ErrRegNotFound)
		}

		log.Error("failed to cancel", "error", err)
		return "", fmt.Errorf("%s: %w", op, service.ErrCancellationFailed)
	}

	log.Info("cancellation successful")
	return status, nil
}

func (s *Service) GetRegistrationStatus(ctx context.Context, memberID uuid.UUID, eventID int) (string, error) {
	const op = "registrations.Service.GetRegistrationStatus"
	log := s.log.With(slog.String("op", op), slog.String("member_id", memberID.String()), slog.Int("event_id", eventID))

	log.Info("fetching registration status")

	status, err := s.regStorage.GetRegistrationStatus(ctx, memberID, eventID)
	if err != nil {
		if errors.Is(err, storage.ErrRegNotFound) {
			log.Error("registration not found", "error", err)
			return "", fmt.Errorf("%s: %w", op, service.ErrRegNotFound)
		}

		log.Error("failed to get registration status", "error", err)
		return "", fmt.Errorf("%s: %w", op, service.ErrStatusCheckFailed)
	}

	log.Info("member is registered on event")
	return status, nil
}
