package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/hsiaoairplane/k8s-crd/pkg/apis/health/v1"
)

// Handler interface contains the methods that are required
type Handler interface {
	Run(stopCh <-chan struct{}) error
	// ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(obj interface{})
}

// HealthHandler is a sample implementation of Handler
type HealthHandler struct {
	logger *zerolog.Logger
	server *http.Server
	put    bool
}

// Run handles any handler initialization
func (h *HealthHandler) Run(stopCh <-chan struct{}) error {
	h.logger.Info().Msg("HealthHandler: run")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "Health method %s\n", r.Method)
		case http.MethodHead:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodPost:
			if h.put {
				fmt.Fprintf(w, "Health method %s\n", r.Method)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		case http.MethodPut:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodPatch:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodDelete:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodConnect:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodOptions:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodTrace:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	if err := h.server.ListenAndServe(); err != nil {
		h.logger.Error().Err(err).Msgf("http server error")
		return err
	}

	return nil
}

// ObjectDeleted is called when an object is deleted
func (h *HealthHandler) ObjectDeleted(obj interface{}) {
	h.logger.Info().Msg("HealthHandler: object deleted")

	// Gracefully shutdown server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.Info().Msg("http server shutdown failed")
	}

	h.logger.Info().Msg("server down")
}

// ObjectUpdated is called when an object is updated
func (h *HealthHandler) ObjectUpdated(obj interface{}) {
	h.logger.Info().Msg("HealthHandler: object updated")

	h.put = obj.(*v1.Health).Spec.Switch

	h.logger.Info().Msgf("put on/off: %v", h.put)
}
