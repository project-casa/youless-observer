// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package observerapp

import (
	"context"

	"github.com/go-pogo/buildinfo"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/errors/errgroup"
	"github.com/go-pogo/healthcheck"
	"github.com/go-pogo/telemetry"
	youlessclient "github.com/roeldev/youless-client"
	"github.com/roeldev/youless-logger/common/logging"
	"github.com/roeldev/youless-logger/common/server"
	"github.com/roeldev/youless-observer"
	"github.com/rs/zerolog"
)

const (
	ErrCreateHealthCheck errors.Msg = "failed to create health checker"
	ErrSetupTelemetry    errors.Msg = "failed to set up telemetry"
	ErrCreateServer      errors.Msg = "failed to create server"
	ErrCreateClient      errors.Msg = "failed to create client"
	ErrCreateObserver    errors.Msg = "failed to create observer"

	Name = "observer"
)

// App is the youless-observer application which observes the YouLess device
// using [youless.Client]. It also runs a server to expose health status and
// build info endpoints.
type App struct {
	telem    *telemetry.Telemetry
	health   *healthcheck.Checker
	server   *server.Server
	client   *youlessclient.Client
	observer *youlessobserver.Observer
}

func New(conf Config, log *logging.Logger, bld *buildinfo.BuildInfo) (*App, error) {
	var app App
	var err error

	app.health, err = healthcheck.New()
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateHealthCheck)
	}

	// set up telemetry
	if conf.Telemetry.ServiceName == "" {
		conf.Telemetry.ServiceName = "youless-" + Name
	}
	telem, err := setupTelemetry(conf.Telemetry, &log.Logger, bld, app.health)
	if err != nil {
		return nil, errors.Wrap(err, ErrSetupTelemetry)
	}

	// set up server
	app.server, err = server.New(Name, conf.Server, log, telem,
		server.WithBuildInfo(bld),
		server.WithHealthChecker(app.health),
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateServer)
	}

	// set up YouLess client
	app.client, err = youlessclient.NewClient(conf.YouLess,
		youlessclient.WithLogger(log),
		youlessclient.WithTracerProvider(telem.TracerProvider()),
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateClient)
	}

	// set up observer
	app.observer, err = youlessobserver.NewObserver(telem.MeterProvider(),
		youlessobserver.WithLogger(&observerLogger{log.Logger}),
		youlessobserver.WithMeterReading(conf.Observer.MeterReadingRegisterer, app.client),
		youlessobserver.WithPhaseReading(conf.Observer.PhaseReadingRegisterer, app.client),
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateObserver)
	}

	app.observer.RegisterHealthCheckers(app.health)
	return &app, nil
}

// Run the app by starting the internal observer and server.
func (app *App) Run(ctx context.Context) error {
	if err := app.observer.Start(); err != nil {
		return err
	}
	return app.server.Run(ctx)
}

// Shutdown the app by stopping the internal observer and server.
func (app *App) Shutdown(ctx context.Context) error {
	// shutdown server first before shutting down other services
	serverErr := app.server.Shutdown(ctx)

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return app.observer.Stop()
	})
	wg.Go(func() error {
		if err := app.telem.ForceFlush(ctx); err != nil {
			return err
		}

		return app.telem.Shutdown(ctx)
	})

	return errors.Append(serverErr, wg.Wait())
}

var _ youlessobserver.Logger = (*observerLogger)(nil)

type observerLogger struct{ zl zerolog.Logger }

func (l *observerLogger) LogRegister(name string) {
	l.zl.Debug().
		Str("name", name).
		Msg("observer register")
}

func (l *observerLogger) LogObserverStart() { l.zl.Info().Msg("observer starting") }

func (l *observerLogger) LogObserverStop() { l.zl.Info().Msg("observer stopped") }
