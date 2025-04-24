package members

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/storage"
	"github.com/Ilya-Repin/orchestra_api/internal/model"
	"github.com/Ilya-Repin/orchestra_api/internal/service"
	"github.com/google/uuid"
	"log/slog"
)

type Service struct {
	log           *slog.Logger
	memberStorage MemberStorage
}

type MemberStorage interface {
	AddMember(ctx context.Context, fullName, email, phone string) (uuid.UUID, error)
	GetMember(ctx context.Context, id uuid.UUID) (model.Member, error)
	GetMembers(ctx context.Context) ([]model.Member, error)
	GetMembersWithStatus(ctx context.Context, status model.MemberStatus) ([]model.Member, error)
	DeleteMember(ctx context.Context, id uuid.UUID) error
	UpdateMember(ctx context.Context, id uuid.UUID, fullName, email, phone string) error
	UpdateMemberStatus(ctx context.Context, id uuid.UUID, status model.MemberStatus) error
}

func New(log *slog.Logger, storage MemberStorage) *Service {
	return &Service{log: log.With("component", "service"), memberStorage: storage}
}

func (s *Service) AddMember(ctx context.Context, fullName, email, phone string) (uuid.UUID, error) {
	const op = "members.Service.AddMember"

	log := s.log.With(slog.String("op", op))
	log.Info("adding new member")

	id, err := s.memberStorage.AddMember(ctx, fullName, email, phone)
	if err != nil {
		log.Error("failed to save member", "error", err)

		switch {
		case errors.Is(err, storage.ErrEmailDuplicate):
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, service.ErrEmailDuplicate)
		case errors.Is(err, storage.ErrPhoneDuplicate):
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, service.ErrPhoneDuplicate)
		case errors.Is(err, storage.ErrInvalidEmail):
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, service.ErrInvalidEmail)
		case errors.Is(err, storage.ErrInvalidPhone):
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, service.ErrInvalidPhone)
		default:
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, service.ErrFailedToAddMember)
		}
	}

	log.Info("member added", "id", id)

	return id, nil
}

func (s *Service) GetMember(ctx context.Context, id uuid.UUID) (model.Member, error) {
	const op = "members.Service.GetMember"

	log := s.log.With(slog.String("op", op), slog.String("id", id.String()))
	log.Info("getting member")

	member, err := s.memberStorage.GetMember(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrMemberNotFound) {
			log.Warn("member not found", "error", err)
			return model.Member{}, fmt.Errorf("%s: %w", op, service.ErrMemberNotFound)
		}

		log.Error("failed to get member", "error", err)
		return model.Member{}, fmt.Errorf("%s: %w", op, service.ErrFailedToGetMember)
	}

	log.Info("member retrieved", "member", member)
	return member, nil
}

func (s *Service) GetMembers(ctx context.Context, status model.MemberStatus) (members []model.Member, err error) {
	const op = "members.Service.GetAllMembers"

	log := s.log.With(slog.String("op", op))
	log.Info("getting all members")

	if len(status) == 0 {
		members, err = s.memberStorage.GetMembers(ctx)
		if err != nil {
			log.Error("failed to get members", "error", err)
			return nil, fmt.Errorf("%s: %w", op, service.ErrFailedToGetMembers)
		}
	} else {

		if status != model.StatusDeclined && status != model.StatusApproved && status != model.StatusPending {
			return nil, service.ErrUnknownStatus
		}

		members, err = s.memberStorage.GetMembersWithStatus(ctx, model.MemberStatus(status))
		if err != nil {
			log.Error("failed to get members", "error", err)
			return nil, fmt.Errorf("%s: %w", op, service.ErrFailedToGetMembers)
		}
	}

	log.Info("members retrieved", "len", len(members))

	return members, nil
}

func (s *Service) DeleteMember(ctx context.Context, id uuid.UUID) error {
	const op = "members.Service.DeleteMember"

	log := s.log.With(slog.String("op", op), slog.String("id", id.String()))
	log.Info("deleting member")

	err := s.memberStorage.DeleteMember(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrMemberNotFound) {
			log.Warn("member not found", "error", err)
			return fmt.Errorf("%s: %w", op, service.ErrMemberNotFound)
		}

		log.Error("failed to delete member", "error", err)
		return fmt.Errorf("%s: %w", op, service.ErrFailedToDeleteMember)
	}

	log.Info("member deleted")
	return nil
}

func (s *Service) UpdateMember(ctx context.Context, id uuid.UUID, fullName, email, phone string) error {
	const op = "members.Service.UpdateMember"

	log := s.log.With(slog.String("op", op), slog.String("id", id.String()))
	log.Info("updating member")

	err := s.memberStorage.UpdateMember(ctx, id, fullName, email, phone)
	if err != nil {
		if errors.Is(err, storage.ErrMemberNotFound) {
			log.Warn("member not found", "error", err)
			return fmt.Errorf("%s: %w", op, service.ErrMemberNotFound)
		}

		log.Error("failed to update member", "error", err)
		return fmt.Errorf("%s: %w", op, service.ErrFailedToUpdateMember)
	}

	log.Info("member updated")
	return nil
}

func (s *Service) UpdateMemberStatus(ctx context.Context, id uuid.UUID, status model.MemberStatus) error {
	const op = "members.Service.UpdateMemberStatus"

	log := s.log.With(slog.String("op", op), slog.String("id", id.String()), slog.String("status", string(status)))
	log.Info("updating member status")

	if status != model.StatusDeclined && status != model.StatusApproved && status != model.StatusPending {
		return service.ErrUnknownStatus
	}

	err := s.memberStorage.UpdateMemberStatus(ctx, id, status)
	if err != nil {
		if errors.Is(err, storage.ErrMemberNotFound) {
			log.Warn("member not found", "error", err)
			return fmt.Errorf("%s: %w", op, service.ErrMemberNotFound)
		}

		log.Error("failed to update member status", "error", err)
		return fmt.Errorf("%s: %w", op, service.ErrFailedToUpdateMemStatus)
	}

	log.Info("member status updated")
	return nil
}
