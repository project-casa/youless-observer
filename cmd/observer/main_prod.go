// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !dev

package main

import (
	"github.com/go-pogo/env"
	"github.com/go-pogo/env/envfile"
	"github.com/go-pogo/errors"
	"github.com/roeldev/youless-observer/app/observer-app"
)

func init() {
	unmarshalEnv = func(conf *observerapp.Config) (err error) {
		src := []env.Lookupper{env.System()}
		if f, fileErr := envfile.Open("/.env"); fileErr == nil {
			defer errors.AppendFunc(&err, f.Close)
			src = append(src, f)
		}

		err = env.NewDecoder(src...).Decode(conf)
		return
	}
}
