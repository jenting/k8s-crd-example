package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/jenting/k8s-crd-example/pkg/apis/health/v1"
)

func Test_testDefaultMethod(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(healthHandlerFunc()))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/health")
	assert.Nil(err)
	res.Body.Close()
	assert.Equal(res.StatusCode, http.StatusOK)

	res, err = http.Post(ts.URL+"/health", "", nil)
	assert.Nil(err)
	res.Body.Close()
	assert.Equal(res.StatusCode, http.StatusMethodNotAllowed)
}

func Test_enablePostMethod(t *testing.T) {
	assert := assert.New(t)

	logger := zerolog.New(ioutil.Discard)
	ts := httptest.NewServer(http.HandlerFunc(healthHandlerFunc()))
	defer ts.Close()

	res, err := http.Post(ts.URL+"/health", "", nil)
	assert.Nil(err)
	res.Body.Close()
	assert.Equal(res.StatusCode, http.StatusMethodNotAllowed)

	h := &HealthHandler{logger: &logger}
	obj := &v1.Health{
		Spec: v1.HealthSpec{
			Action: http.MethodPost,
			Switch: true,
		},
	}
	h.ObjectUpdated(obj)

	res, err = http.Post(ts.URL+"/health", "", nil)
	assert.Nil(err)
	res.Body.Close()
	assert.Equal(res.StatusCode, http.StatusOK)
}

func Test_deleteObject(t *testing.T) {
	assert := assert.New(t)

	logger := zerolog.New(ioutil.Discard)
	ts := httptest.NewServer(http.HandlerFunc(healthHandlerFunc()))
	defer ts.Close()

	h := &HealthHandler{logger: &logger}
	obj := &v1.Health{
		Spec: v1.HealthSpec{
			Action: http.MethodPost,
			Switch: true,
		},
	}
	h.ObjectUpdated(obj)

	res, err := http.Post(ts.URL+"/health", "", nil)
	assert.Nil(err)
	res.Body.Close()
	assert.Equal(res.StatusCode, http.StatusOK)

	h.ObjectDeleted(&v1.Health{})

	res, err = http.Post(ts.URL+"/health", "", nil)
	assert.Nil(err)
	res.Body.Close()
	assert.Equal(res.StatusCode, http.StatusMethodNotAllowed)
}

func Test_serverShutdown(t *testing.T) {
	stopCh := make(chan struct{})
	defer close(stopCh)

	logger := zerolog.New(ioutil.Discard)
	server := &http.Server{}
	h := &HealthHandler{logger: &logger, server: server}

	go h.Run(stopCh)

	timeout := time.NewTimer(time.Second * 3)
	defer timeout.Stop()

	<-timeout.C

	// stop server
	stopCh <- struct{}{}
}
