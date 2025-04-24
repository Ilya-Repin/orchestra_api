package handler

import (
	"encoding/json"
	"errors"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/metrics"
	"github.com/Ilya-Repin/orchestra_api/internal/openapi"
	"github.com/Ilya-Repin/orchestra_api/internal/service"
	"github.com/Ilya-Repin/orchestra_api/internal/service/auxiliary"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type AuxHandler struct {
	log        *slog.Logger
	auxService *auxiliary.Service
	metrics    *metrics.Metrics
}

func NewAuxHandler(log *slog.Logger, as *auxiliary.Service, metrics *metrics.Metrics) *AuxHandler {
	return &AuxHandler{log: log, auxService: as, metrics: metrics}
}

func (ah *AuxHandler) HandleGetEventTypes(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auxiliary.HandleGetEventTypes"

	log := ah.log.With(slog.String("op", op))
	ctx := r.Context()

	eventTypes, err := ah.auxService.GetEventTypes(ctx)
	if err != nil {
		log.Error("failed to get event types", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to get event types")
		ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()

		return
	}

	var typeResponses []openapi.EventTypeResponse
	for _, e := range eventTypes {
		id := int32(e.ID)

		typeResponses = append(typeResponses, openapi.EventTypeResponse{
			Id:          &id,
			Name:        &e.Name,
			Description: &e.Description,
		})
	}

	writeJSON(w, http.StatusOK, typeResponses)
	ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (ah *AuxHandler) HandleGetLocations(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auxiliary.HandleGetLocations"

	log := ah.log.With(slog.String("op", op))
	ctx := r.Context()

	readLocations, err := ah.auxService.GetLocations(ctx)
	if err != nil {
		log.Error("failed to get locations", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to get locations")
		ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()

		return
	}

	var locResponses []openapi.LocationResponse
	for _, m := range readLocations {
		id := int32(m.ID)

		locResponses = append(locResponses, openapi.LocationResponse{
			Id:       &id,
			Name:     &m.Name,
			Route:    &m.Route,
			Features: &m.Features,
		})
	}

	writeJSON(w, http.StatusOK, locResponses)
	ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (ah *AuxHandler) HandleCreateEventType(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auxiliary.HandleCreateEventType"
	ctx := r.Context()

	var req openapi.NewEventTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.log.Warn("failed to decode request", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid request body")
		ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	if req.GetName() == "" || req.GetDescription() == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	id, err := ah.auxService.AddEventType(ctx, req.GetName(), req.GetDescription())
	if err != nil {
		ah.log.Error("failed to add event type", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to add event type")
		ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	writeJSON(w, http.StatusCreated, id)
	ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "201").Inc()
}

func (ah *AuxHandler) HandleCreateLocation(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auxiliary.HandleCreateLocation"
	ctx := r.Context()

	var req openapi.NewLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.log.Warn("failed to decode request", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid request body")
		ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	if req.GetName() == "" || req.GetRoute() == "" || req.GetFeatures() == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	id, err := ah.auxService.AddLocation(ctx, req.GetName(), req.GetRoute(), req.GetFeatures())
	if err != nil {
		ah.log.Error("failed to add location", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to add location")
		ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	writeJSON(w, http.StatusCreated, id)
	ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "201").Inc()
}

func (ah *AuxHandler) HandleGetOrchestraInfo(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.events.HandleGetEvent"

	log := ah.log.With(slog.String("op", op))
	ctx := r.Context()

	key := chi.URLParam(r, "key")

	info, err := ah.auxService.GetOrchestraInfo(ctx, key)
	if err != nil {
		if errors.Is(err, service.ErrInfoNotFound) {
			log.Error("info not found", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusNotFound, "info not found")
			ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}
		log.Error("failed to get event", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to get info")
		ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	infoResponse := openapi.OrchestraInfoResponse{
		Key:   &info.Key,
		Value: &info.Value,
	}

	writeJSON(w, http.StatusOK, infoResponse)
	ah.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}
