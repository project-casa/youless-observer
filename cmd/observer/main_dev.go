// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build dev

package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-pogo/easytls"
	"github.com/go-pogo/env"
	"github.com/go-pogo/env/dotenv"
	"github.com/roeldev/youless-observer/app/observer-app"
)

func init() {
	unmarshalEnv = func(conf *observerapp.Config) error {
		_, dir, _, _ := runtime.Caller(1)
		dir = filepath.Dir(dir)

		ae, _ := dotenv.GetActiveEnvironmentOr(os.Args[1:], dotenv.Development)
		environ, err := dotenv.Read(dir, ae).Environ()
		if err != nil {
			return err
		}

		if err = env.Load(environ); err != nil {
			return err
		}
		if err = env.NewDecoder(environ).Decode(conf); err != nil {
			return err
		}

		caCertFile := conf.Server.TLS.CACertFile.String()
		prefixDir(dir, &caCertFile)
		conf.Server.TLS.CACertFile = easytls.CertificateFile(caCertFile)

		prefixDir(dir, &conf.Server.TLS.CertFile)
		prefixDir(dir, &conf.Server.TLS.KeyFile)
		return nil
	}
}

func prefixDir(dir string, ptr *string) {
	if ptr == nil || *ptr == "" || filepath.IsAbs(*ptr) {
		return
	}
	*ptr = filepath.Join(dir, *ptr)
}
