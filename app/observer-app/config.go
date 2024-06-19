// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package observerapp

import (
	"github.com/go-pogo/telemetry"
	youlessclient "github.com/roeldev/youless-client"
	"github.com/roeldev/youless-logger/common/server"
	youlessobserver "github.com/roeldev/youless-observer"
	"github.com/rs/zerolog"
)

//goland:noinspection GoUnusedConst
const ErrInvalidConfig = youlessclient.ErrInvalidConfig

// Config is the configuration for App.
type Config struct {
	Level         zerolog.Level `env:"LOG_LEVEL" default:"debug"`
	WithTimestamp bool          `env:"LOG_TIMESTAMP" default:"true"`
	Server        server.Config `env:",include"`
	Telemetry     telemetry.Config
	YouLess       youlessclient.Config `env:"YOULESS,include"`
	//Mqtt          mqtt.Config             `env:",include"`

	Observer struct {
		youlessobserver.MeterReadingRegisterer
		youlessobserver.PhaseReadingRegisterer
	} `env:",include"`
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return err
	}
	if err := c.YouLess.Validate(); err != nil {
		return err
	}
	return nil
}
