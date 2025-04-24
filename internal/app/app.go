package app

import (
	"database/sql"
	"encoding/json"
	"github.com/Ilya-Repin/orchestra_api/internal/handler"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/metrics"
	"github.com/Ilya-Repin/orchestra_api/internal/infra/storage/postgres"
	"github.com/Ilya-Repin/orchestra_api/internal/openapi"
	"github.com/Ilya-Repin/orchestra_api/internal/service/auxiliary"
	"github.com/Ilya-Repin/orchestra_api/internal/service/events"
	"github.com/Ilya-Repin/orchestra_api/internal/service/members"
	"github.com/Ilya-Repin/orchestra_api/internal/service/registrations"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
)

type App struct {
	log                 *slog.Logger
	memberService       *members.Service
	eventService        *events.Service
	registrationService *registrations.Service
	auxService          *auxiliary.Service
	metrics             *metrics.Metrics // Храним метрики
}

func NewApp(log *slog.Logger, db *sql.DB, appMetrics *metrics.Metrics) *App {
	storage := postgres.New(db)

	return &App{
		log:                 log.With("component", "app"),
		memberService:       members.New(log, storage),
		eventService:        events.New(log, storage, storage),
		registrationService: registrations.New(log, storage, storage),
		auxService:          auxiliary.New(log, storage),
		metrics:             appMetrics,
	}
}

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/members", a.membersRoutes())
		r.Mount("/events", a.eventsRoutes())
		r.Mount("/locations", a.locRoutes())
		r.Mount("/types", a.eventTypeRoutes())
		r.Mount("/info", a.infoRoutes())
	})

	return r
}

func (a *App) membersRoutes() http.Handler {
	r := chi.NewRouter()

	membersHandler := handler.NewMembersHandler(a.log, a.memberService, a.metrics)

	r.Get("/", membersHandler.HandleGetMembers)
	r.Post("/", membersHandler.HandleCreateMember)
	r.Route("/{memberId}", func(r chi.Router) {
		r.Get("/", membersHandler.HandleGetMember)
		r.Put("/", membersHandler.HandleUpdateMemberProfile)
		r.Patch("/", membersHandler.HandleUpdateMemberStatus)
		r.Delete("/", membersHandler.HandleDeleteMember)
	})

	return r
}

func (a *App) locRoutes() http.Handler {
	r := chi.NewRouter()

	auxHandler := handler.NewAuxHandler(a.log, a.auxService, a.metrics)

	r.Get("/", auxHandler.HandleGetLocations)
	r.Post("/", auxHandler.HandleCreateLocation)

	return r
}

func (a *App) eventTypeRoutes() http.Handler {
	r := chi.NewRouter()

	auxHandler := handler.NewAuxHandler(a.log, a.auxService, a.metrics)

	r.Get("/", auxHandler.HandleGetEventTypes)
	r.Post("/", auxHandler.HandleCreateEventType)

	return r
}

func (a *App) infoRoutes() http.Handler {
	r := chi.NewRouter()

	auxHandler := handler.NewAuxHandler(a.log, a.auxService, a.metrics)

	r.Get("/", auxHandler.HandleGetOrchestraInfo)

	return r
}

func (a *App) eventsRoutes() http.Handler {
	r := chi.NewRouter()

	eventsHandler := handler.NewEventsHandler(a.log, a.eventService, a.metrics)
	registrationHandler := handler.NewRegistrationsHandler(a.log, a.registrationService, a.metrics)

	r.Get("/", eventsHandler.HandleGetEvents)
	r.Post("/", eventsHandler.HandleCreateEvent)
	r.Get("/upcoming", eventsHandler.HandleGetUpcomingEvents)
	r.Get("/available", eventsHandler.HandleGetAvailableEvents)
	r.Get("/registered", eventsHandler.HandleGetRegisteredEvents)
	//r.Post("/", a.handleCreateMember)
	r.Route("/{eventId}", func(r chi.Router) {
		r.Get("/", eventsHandler.HandleGetEvent)
		r.Put("/", eventsHandler.HandleUpdateEvent)
		r.Delete("/", eventsHandler.HandleDeleteEvent)
		r.Route("/registration", func(r chi.Router) {
			r.Get("/", registrationHandler.HandleCheckRegistration)
			r.Post("/", registrationHandler.HandleRegister)
			r.Delete("/", registrationHandler.HandleCancel)

		})
	})

	return r
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, openapi.ErrorResponse{Message: &message})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
