package handler

import (
	"errors"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/metrics"
	"github.com/Ilya-Repin/orchestra_api/internal/service"
	"github.com/Ilya-Repin/orchestra_api/internal/service/registrations"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strconv"
)

type RegistrationsHandler struct {
	log        *slog.Logger
	regService *registrations.Service
	metrics    *metrics.Metrics
}

func NewRegistrationsHandler(log *slog.Logger, rs *registrations.Service, metrics *metrics.Metrics) *RegistrationsHandler {
	return &RegistrationsHandler{log: log, regService: rs, metrics: metrics}
}

func (rh *RegistrationsHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.registrations.HandleRegister"

	log := rh.log.With(slog.String("op", op))
	ctx := r.Context()

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Error("invalid event id", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid event id")
		rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	memberID, err := uuid.Parse(r.URL.Query().Get("memberId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "wrong format memberId")
		rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	status, err := rh.regService.RegisterForEvent(ctx, memberID, eventID)
	if err != nil {
		if errors.Is(err, service.ErrMemberNotFound) {
			log.Error("member not found", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusNotFound, "member not found")
			rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
		} else if errors.Is(err, service.ErrMemberNotApproved) {
			log.Error("member not approved", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusBadRequest, "member not approved")
			rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		} else if errors.Is(err, service.ErrRegAlreadyExists) {
			log.Error("registration already exists", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusBadRequest, "registration already exists")
			rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		} else {
			log.Error("failed to register", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusInternalServerError, "failed to register")
			rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		}
		return
	}

	writeJSON(w, http.StatusCreated, status)
	rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "201").Inc()
	rh.metrics.EventRegistrationsTotal.WithLabelValues("registered").Inc()
}

func (rh *RegistrationsHandler) HandleCancel(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.registrations.HandleCancel"

	log := rh.log.With(slog.String("op", op))
	ctx := r.Context()

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Error("invalid event id", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid event id")
		rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	memberID, err := uuid.Parse(r.URL.Query().Get("memberId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "wrong format memberId")
		rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	status, err := rh.regService.CancelRegistration(ctx, memberID, eventID)
	if err != nil {
		if errors.Is(err, service.ErrRegNotFound) {
			log.Error("registration not found", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusNotFound, "registration not found")
			rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
		} else {
			log.Error("failed to cancel", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusInternalServerError, "failed to cancel")
			rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		}
		return
	}

	writeJSON(w, http.StatusNoContent, status)
	rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "204").Inc()
	rh.metrics.EventRegistrationsTotal.WithLabelValues("cancelled").Inc()
}

func (rh *RegistrationsHandler) HandleCheckRegistration(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.registrations.HandleCheckRegistration"

	log := rh.log.With(slog.String("op", op))
	ctx := r.Context()

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Error("invalid event id", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid event id")
		rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	memberID, err := uuid.Parse(r.URL.Query().Get("memberId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "wrong format memberId")
		rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	status, err := rh.regService.GetRegistrationStatus(ctx, memberID, eventID)
	if err != nil {
		if errors.Is(err, service.ErrRegNotFound) {
			log.Error("registration not found", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusNotFound, "registration not found")
			rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
		} else {
			log.Error("failed to check registration", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusInternalServerError, "failed to check registration")
			rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		}
		return
	}

	writeJSON(w, http.StatusOK, status)
	rh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}
