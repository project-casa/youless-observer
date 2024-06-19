// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen-dotenv.go

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-pogo/buildinfo"
	"github.com/go-pogo/easytls"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/healthcheck"
	"github.com/go-pogo/healthcheck/healthclient"
	"github.com/go-pogo/serv"
	"github.com/roeldev/youless-logger/common/logging"
	youlessobserver "github.com/roeldev/youless-observer"
	"github.com/roeldev/youless-observer/app/observer-app"
)

var unmarshalEnv func(conf *observerapp.Config) error

func main() {
	var conf observerapp.Config
	errors.FatalOnErr(unmarshalEnv(&conf))

	log := logging.New(conf.Level, conf.WithTimestamp)

	if len(os.Args) >= 2 && os.Args[1] == "healthcheck" {
		runHealthCheck(log, conf.Server.Port, conf.Server.TLS)
		return
	}

	// collecting metrics is always enabled
	conf.Telemetry.Meter.Enabled = true
	if err := conf.Validate(); err != nil {
		log.Fatal().Err(err).Msg("invalid configuration")
	}

	runApp(log, conf)
}

func runApp(log *logging.Logger, conf observerapp.Config) {
	bld, err := buildinfo.New(youlessobserver.Version)
	if err != nil {
		log.Warn().Msg("unable to retrieve build info")
	} else {
		log.LogBuildInfo(bld,
			"github.com/roeldev/youless-client",
			"github.com/roeldev/youless-logger",
		)
	}

	app, err := observerapp.New(conf, log, bld)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create observer app")
	}

	runCtx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer stopFn()
		if err := app.Run(runCtx); err != nil {
			log.Err(err).Msg("error during run")
		}
	}()
	<-runCtx.Done()

	closeCtx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer stopFn()
		if err := app.Shutdown(closeCtx); err != nil {
			log.Err(err).Msg("error during shutdown")
		}
	}()
	<-closeCtx.Done()
}

func runHealthCheck(log *logging.Logger, port serv.Port, tls easytls.Config) {
	healthy, err := healthclient.New(
		healthclient.Config{
			TargetPort: uint16(port),
			TargetPath: healthcheck.PathPattern,
		},
		healthclient.WithTLSConfig(easytls.DefaultTLSConfig(), tls),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create healthcheck client")
	}

	stat, err := healthy.Request(context.Background())
	if err != nil {
		log.Err(err).Stringer("status", stat).Msg("healthcheck failed")
	} else {
		log.Info().Stringer("status", stat).Msg("health checked")
	}

	os.Exit(stat.ExitCode())
}
