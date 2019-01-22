package main

import (
	"github.com/rs/zerolog"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(objOld, objNew interface{})
}

// HealthHandler is a sample implementation of Handler
type HealthHandler struct {
	logger *zerolog.Logger
}

// Init handles any handler initialization
func (t *HealthHandler) Init() error {
	t.logger.Info().Msg("HealthHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *HealthHandler) ObjectCreated(obj interface{}) {
	t.logger.Info().Msg("HealthHandler.ObjectCreated")
}

// ObjectDeleted is called when an object is deleted
func (t *HealthHandler) ObjectDeleted(obj interface{}) {
	t.logger.Info().Msg("HealthHandler.ObjectDeleted")
}

// ObjectUpdated is called when an object is updated
func (t *HealthHandler) ObjectUpdated(objOld, objNew interface{}) {
	t.logger.Info().Msg("HealthHandler.ObjectUpdated")
}
