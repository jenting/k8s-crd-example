package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	v1 "github.com/jenting/k8s-crd-example/pkg/apis/health/v1"
)

// Handler interface contains the methods that are required
type Handler interface {
	Run(stopCh <-chan struct{}) error
	ObjectDeleted(obj interface{})
	ObjectUpdated(obj interface{})
}

// HealthHandler is a sample implementation of Handler
type HealthHandler struct {
	logger *zerolog.Logger
	server *http.Server
	put    bool
}

var methodEnabled = map[string]bool{
	http.MethodHead:    false,
	http.MethodPost:    false,
	http.MethodPut:     false,
	http.MethodPatch:   false,
	http.MethodDelete:  false,
	http.MethodConnect: false,
	http.MethodOptions: false,
	http.MethodTrace:   false,
}

// Handle is a function that can be registered to a route to handle HTTP requests.
type Handle func(http.ResponseWriter, *http.Request)

func healthHandlerFunc() Handle {
	return func(w http.ResponseWriter, r *http.Request) {
		if http.MethodGet == r.Method {
			fmt.Fprintf(w, "Health method %s\n", r.Method)
		} else {
			if methodEnabled[r.Method] {
				fmt.Fprintf(w, "Health method %s\n", r.Method)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}
	}
}

// Run handles any handler initialization
func (h *HealthHandler) Run(stopCh <-chan struct{}) error {
	h.logger.Info().Msg("HealthHandler: run")

	http.HandleFunc("/health", healthHandlerFunc())

	go func() {
		<-stopCh

		// Gracefully shutdown server.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.server.Shutdown(ctx); err != nil {
			h.logger.Info().Msg("http server shutdown failed")
		}
		h.logger.Info().Msg("server down")
	}()

	if err := h.server.ListenAndServe(); err != nil {
		h.logger.Error().Err(err).Msgf("http server error")
		return err
	}

	return nil
}

// ObjectDeleted is called when an object is deleted
func (h *HealthHandler) ObjectDeleted(obj interface{}) {
	h.logger.Info().Msg("HealthHandler: object deleted")

	for k := range methodEnabled {
		methodEnabled[k] = false
	}
}

// ObjectUpdated is called when an object is updated
func (h *HealthHandler) ObjectUpdated(obj interface{}) {
	h.logger.Info().Msg("HealthHandler: object updated")

	methodEnabled[obj.(*v1.Health).Spec.Action] = obj.(*v1.Health).Spec.Switch

	h.logger.Info().Msgf("Method %s %t", obj.(*v1.Health).Spec.Action, obj.(*v1.Health).Spec.Switch)
}
