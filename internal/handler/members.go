package handler

import (
	"encoding/json"
	"errors"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/metrics"
	"github.com/Ilya-Repin/orchestra_api/internal/model"
	"github.com/Ilya-Repin/orchestra_api/internal/openapi"
	"github.com/Ilya-Repin/orchestra_api/internal/service"
	"github.com/Ilya-Repin/orchestra_api/internal/service/members"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type MembersHandler struct {
	log           *slog.Logger
	memberService *members.Service
	metrics       *metrics.Metrics
}

func NewMembersHandler(log *slog.Logger, memberService *members.Service, metrics *metrics.Metrics) *MembersHandler {
	return &MembersHandler{
		log:           log,
		memberService: memberService,
		metrics:       metrics,
	}
}

func (mh *MembersHandler) HandleGetMembers(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.members.HandleGetMembers"

	log := mh.log.With(slog.String("op", op))
	ctx := r.Context()

	status := r.URL.Query().Get("status")

	readMembers, err := mh.memberService.GetMembers(ctx, model.MemberStatus(status))
	if err != nil {
		if errors.Is(err, service.ErrUnknownStatus) {
			log.Error("unknown status", slog.String("op", op), slog.String("status", status), slog.Any("err", err))
			writeError(w, http.StatusBadRequest, "unknown status")
			mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		} else {
			log.Error("failed to get members", slog.String("op", op), slog.String("status", status), slog.Any("err", err))
			writeError(w, http.StatusInternalServerError, "failed to get readMembers")
			mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		}
		return
	}

	var memberResponses []openapi.MemberResponse
	for _, m := range readMembers {
		id := m.ID.String()
		statusStr := string(m.Status)

		memberResponses = append(memberResponses, openapi.MemberResponse{
			Id:        &id,
			FullName:  m.FullName,
			Email:     m.Email,
			Phone:     m.Phone,
			Status:    &statusStr,
			CreatedAt: &m.CreatedAt,
			UpdatedAt: &m.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, memberResponses)
	mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (mh *MembersHandler) HandleCreateMember(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.members.HandleCreateMember"
	ctx := r.Context()

	var req openapi.NewMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mh.log.Warn("failed to decode request", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid request body")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	if req.GetFullName() == "" || req.GetEmail() == "" || req.GetPhone() == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	id, err := mh.memberService.AddMember(ctx, req.GetFullName(), req.GetEmail(), req.GetPhone())
	if err != nil {
		mh.log.Error("failed to create member", slog.String("op", op), slog.Any("err", err))

		code := "400"

		switch {
		case errors.Is(err, service.ErrEmailDuplicate):
			writeError(w, http.StatusBadRequest, "email already exists")
		case errors.Is(err, service.ErrPhoneDuplicate):
			writeError(w, http.StatusBadRequest, "phone number already exists")
		case errors.Is(err, service.ErrInvalidEmail):
			writeError(w, http.StatusBadRequest, "invalid email format")
		case errors.Is(err, service.ErrInvalidPhone):
			writeError(w, http.StatusBadRequest, "invalid phone number format")
		default:
			code = "500"
			writeError(w, http.StatusInternalServerError, "failed to create member")
		}

		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, code).Inc()
		return
	}

	writeJSON(w, http.StatusCreated, id.String())
}

func (mh *MembersHandler) HandleGetMember(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.members.HandleGetMember"
	ctx := r.Context()

	memberID, err := uuid.Parse(chi.URLParam(r, "memberId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "wrong format memberId")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	member, err := mh.memberService.GetMember(ctx, memberID)
	if err != nil {
		if errors.Is(err, service.ErrMemberNotFound) {
			writeError(w, http.StatusNotFound, "member not found")
			mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}
		mh.log.Error("failed to get member", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to get member")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	writeJSON(w, http.StatusOK, member)
	mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (mh *MembersHandler) HandleUpdateMemberProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.members.HandleUpdateMemberProfile"

	memberID, err := uuid.Parse(chi.URLParam(r, "memberId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "wrong format memberId")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}
	var req openapi.UpdateMemberProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	if req.GetFullName() == "" || req.GetEmail() == "" || req.GetPhone() == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	err = mh.memberService.UpdateMember(r.Context(), memberID, req.GetFullName(), req.GetEmail(), req.GetPhone())
	if err != nil {
		if errors.Is(err, service.ErrMemberNotFound) {
			writeError(w, http.StatusNotFound, "member not found")
			mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}
		mh.log.Error("failed to update member", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to update member")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	writeJSON(w, http.StatusOK, memberID)
	mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (mh *MembersHandler) HandleUpdateMemberStatus(w http.ResponseWriter, r *http.Request) {
	const op = "members.member.HandleUpdateMemberStatus"

	memberID, err := uuid.Parse(chi.URLParam(r, "memberId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "wrong format memberId")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	var req openapi.UpdateMemberStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	err = mh.memberService.UpdateMemberStatus(r.Context(), memberID, model.MemberStatus(req.GetStatus()))
	if err != nil {
		if errors.Is(err, service.ErrMemberNotFound) {
			writeError(w, http.StatusNotFound, "member not found")
			mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}
		mh.log.Error("failed to update member status", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to update member status")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	writeJSON(w, http.StatusOK, memberID)
	mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
	mh.metrics.UserStatusDecisionsTotal.WithLabelValues(req.GetStatus()).Inc()
}

func (mh *MembersHandler) HandleDeleteMember(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.members.HandleDeleteMember"

	memberID, err := uuid.Parse(chi.URLParam(r, "memberId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "wrong format memberId")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	err = mh.memberService.DeleteMember(r.Context(), memberID)
	if err != nil {
		if errors.Is(err, service.ErrMemberNotFound) {
			writeError(w, http.StatusNotFound, "member not found")
			mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}
		mh.log.Error("failed to delete member", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to delete member")
		mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	w.WriteHeader(http.StatusNoContent)
	mh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "204").Inc()
}
