// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package observerapp

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/zerologr"
	"github.com/go-pogo/buildinfo"
	"github.com/go-pogo/healthcheck"
	"github.com/go-pogo/telemetry"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

func setupTelemetry(conf telemetry.Config, log *zerolog.Logger, bld *buildinfo.BuildInfo, reg healthcheck.Registerer) (*telemetry.Telemetry, error) {
	telem := telemetry.NewBuilder(conf).Global().WithDefaultExporter()
	if bld != nil {
		telem.TracerProvider.WithBuildInfo(bld.Internal())
	}

	if reg == nil {
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			log.Err(err).Msg("otel error")
		}))
	} else {
		health := &otelHealthCheck{timeout: conf.ExporterOTLP.TimeoutDuration()}
		reg.Register("otel.otlp_exporter", health)

		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			log.Err(err).Msg("otel error")
			go health.markError()
		}))
	}

	zl := log.Level(zerolog.DebugLevel)
	otel.SetLogger(zerologr.New(&zl))

	return telem.Build()
}

var _ healthcheck.HealthChecker = (*otelHealthCheck)(nil)

type otelHealthCheck struct {
	mut     sync.RWMutex
	lastErr time.Time
	timeout time.Duration
}

func (o *otelHealthCheck) CheckHealth(_ context.Context) healthcheck.Status {
	o.mut.RLock()
	defer o.mut.RUnlock()

	if time.Since(o.lastErr) <= o.timeout {
		return healthcheck.StatusUnhealthy
	}
	return healthcheck.StatusHealthy
}

func (o *otelHealthCheck) markError() {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.lastErr = time.Now()
}
