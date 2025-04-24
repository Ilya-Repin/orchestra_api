package handler

import (
	"encoding/json"
	"errors"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/metrics"
	"github.com/Ilya-Repin/orchestra_api/internal/openapi"
	"github.com/Ilya-Repin/orchestra_api/internal/service"
	"github.com/Ilya-Repin/orchestra_api/internal/service/events"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type EventsHandler struct {
	log          *slog.Logger
	eventService *events.Service
	metrics      *metrics.Metrics
}

func NewEventsHandler(log *slog.Logger, es *events.Service, metrics *metrics.Metrics) *EventsHandler {
	return &EventsHandler{log: log, eventService: es, metrics: metrics}
}

func (eh *EventsHandler) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.events.HandleCreateEvent"
	log := eh.log.With(slog.String("op", op))
	ctx := r.Context()

	var req openapi.NewEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("invalid request body", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid request body")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	eventID, err := eh.eventService.AddEvent(ctx, req.GetTitle(), req.GetDescription(), int(req.GetEventType()), req.GetEventDate(), int(req.GetLocation()), int(req.GetCapacity()))
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			log.Error("event not found", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusNotFound, "event not found")
			eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}
		log.Error("failed to add event", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to add event")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	writeJSON(w, http.StatusCreated, eventID)
}

func (eh *EventsHandler) HandleGetEvents(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.events.HandleGetEvents"

	log := eh.log.With(slog.String("op", op))
	ctx := r.Context()

	eventTypeStr := r.URL.Query().Get("type")
	dateFromStr := r.URL.Query().Get("date_from")
	dateToStr := r.URL.Query().Get("date_to")

	var (
		eventType *int
		begin     *time.Time
		end       *time.Time
	)

	if eventTypeStr != "" {
		et, err := strconv.Atoi(eventTypeStr)
		if err != nil {
			log.Error("invalid event type", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusBadRequest, "invalid event type")
			eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
			return
		}
		eventType = &et
	}

	if dateFromStr != "" {
		t, err := time.Parse(time.RFC3339, dateFromStr)
		if err != nil {
			log.Error("invalid date_from", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusBadRequest, "invalid date_from")
			eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
			return
		}
		begin = &t
	}

	if dateToStr != "" {
		t, err := time.Parse(time.RFC3339, dateToStr)
		if err != nil {
			log.Error("invalid date_to", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusBadRequest, "invalid date_to")
			eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
			return
		}
		end = &t
	}

	readEvents, err := eh.eventService.GetEvents(ctx, eventType, begin, end)

	if err != nil {
		log.Error("failed to get events", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to get events")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()

		return
	}

	var eventResponses []openapi.EventResponse
	for _, e := range readEvents {
		id := int32(e.ID)
		eventTypeId := int32(e.EventType.ID)
		locId := int32(e.Location.ID)
		capacity := int32(e.Capacity)
		eventResponses = append(eventResponses, openapi.EventResponse{
			Id:          &id,
			Title:       &e.Title,
			Description: &e.Description,
			EventType:   &eventTypeId,
			EventDate:   &e.EventDate,
			Location:    &locId,
			Capacity:    &capacity,
			CreatedAt:   &e.CreatedAt,
			UpdatedAt:   &e.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, eventResponses)
	eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (eh *EventsHandler) HandleGetUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.events.HandleGetUpcomingEvents"

	log := eh.log.With(slog.String("op", op))
	ctx := r.Context()

	readEvents, err := eh.eventService.GetUpcomingEvents(ctx)

	if err != nil {
		log.Error("failed to get upcoming events", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to get upcoming events")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()

		return
	}

	var eventResponses []openapi.EventResponse
	for _, e := range readEvents {
		id := int32(e.ID)
		eventTypeId := int32(e.EventType.ID)
		locId := int32(e.Location.ID)
		capacity := int32(e.Capacity)
		eventResponses = append(eventResponses, openapi.EventResponse{
			Id:          &id,
			Title:       &e.Title,
			Description: &e.Description,
			EventType:   &eventTypeId,
			EventDate:   &e.EventDate,
			Location:    &locId,
			Capacity:    &capacity,
			CreatedAt:   &e.CreatedAt,
			UpdatedAt:   &e.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, eventResponses)
	eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (eh *EventsHandler) HandleGetAvailableEvents(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.events.HandleGetUpcomingEvents"

	log := eh.log.With(slog.String("op", op))
	ctx := r.Context()

	memberID, err := uuid.Parse(r.URL.Query().Get("memberId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "wrong format memberId")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	readEvents, err := eh.eventService.GetAvailableEvents(ctx, memberID)

	if err != nil {
		if errors.Is(err, service.ErrMemberNotFound) {
			writeError(w, http.StatusNotFound, "member not found")
			eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}

		log.Error("failed to get available events", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to get available events")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	var eventResponses []openapi.EventResponse
	for _, e := range readEvents {
		id := int32(e.ID)
		eventTypeId := int32(e.EventType.ID)
		locId := int32(e.Location.ID)
		capacity := int32(e.Capacity)
		eventResponses = append(eventResponses, openapi.EventResponse{
			Id:          &id,
			Title:       &e.Title,
			Description: &e.Description,
			EventType:   &eventTypeId,
			EventDate:   &e.EventDate,
			Location:    &locId,
			Capacity:    &capacity,
			CreatedAt:   &e.CreatedAt,
			UpdatedAt:   &e.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, eventResponses)
	eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (eh *EventsHandler) HandleGetRegisteredEvents(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.events.HandleGetRegisteredEvents"

	log := eh.log.With(slog.String("op", op))
	ctx := r.Context()

	memberID, err := uuid.Parse(r.URL.Query().Get("memberId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "wrong format memberId")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	readEvents, err := eh.eventService.GetRegisteredEvents(ctx, memberID)

	if err != nil {
		if errors.Is(err, service.ErrMemberNotFound) {
			writeError(w, http.StatusNotFound, "member not found")
			eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}

		log.Error("failed to get registered events", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusInternalServerError, "failed to get registered events")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	var eventResponses []openapi.EventResponse
	for _, e := range readEvents {
		id := int32(e.ID)
		eventTypeId := int32(e.EventType.ID)
		locId := int32(e.Location.ID)
		capacity := int32(e.Capacity)
		eventResponses = append(eventResponses, openapi.EventResponse{
			Id:          &id,
			Title:       &e.Title,
			Description: &e.Description,
			EventType:   &eventTypeId,
			EventDate:   &e.EventDate,
			Location:    &locId,
			Capacity:    &capacity,
			CreatedAt:   &e.CreatedAt,
			UpdatedAt:   &e.UpdatedAt,
		})
	}

	writeJSON(w, http.StatusOK, eventResponses)
	eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (eh *EventsHandler) HandleGetEvent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.events.HandleGetEvent"

	log := eh.log.With(slog.String("op", op))
	ctx := r.Context()

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Error("invalid event id", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid event id")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	e, err := eh.eventService.GetEvent(ctx, eventID)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			log.Error("event not found", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusNotFound, "event not found")
			eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}
		log.Error("failed to get event", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to get event")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}
	id := int32(e.ID)
	eventTypeId := int32(e.EventType.ID)
	locId := int32(e.Location.ID)
	capacity := int32(e.Capacity)
	eventResponse := openapi.EventResponse{
		Id:          &id,
		Title:       &e.Title,
		Description: &e.Description,
		EventType:   &eventTypeId,
		EventDate:   &e.EventDate,
		Location:    &locId,
		Capacity:    &capacity,
		CreatedAt:   &e.CreatedAt,
		UpdatedAt:   &e.UpdatedAt,
	}

	writeJSON(w, http.StatusOK, eventResponse)
	eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (eh *EventsHandler) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.events.HandleUpdateEvent"
	log := eh.log.With(slog.String("op", op))
	ctx := r.Context()

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Error("invalid event id", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid event id")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	var req openapi.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("invalid request body", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid request body")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	err = eh.eventService.UpdateEvent(ctx, eventID, req.GetTitle(), req.GetDescription(), int(req.GetEventType()), req.GetEventDate(), int(req.GetLocation()), int(req.GetCapacity()))
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			log.Error("event not found", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusNotFound, "event not found")
			eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}
		log.Error("failed to update event", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to update event")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	writeJSON(w, http.StatusOK, eventID)
	eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}

func (eh *EventsHandler) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.events.HandleDeleteEvent"
	log := eh.log.With(slog.String("op", op))
	ctx := r.Context()

	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Error("invalid event id", slog.String("op", op), slog.Any("err", err))
		writeError(w, http.StatusBadRequest, "invalid event id")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "400").Inc()
		return
	}

	err = eh.eventService.DeleteEvent(ctx, eventID)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			log.Error("event not found", slog.String("op", op), slog.Any("err", err))
			writeError(w, http.StatusNotFound, "event not found")
			eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "404").Inc()
			return
		}
		log.Error("failed to delete event", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "failed to delete event")
		eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "500").Inc()
		return
	}

	w.WriteHeader(http.StatusNoContent)
	eh.metrics.ApiRequestsTotal.WithLabelValues(r.Method, "200").Inc()
}
