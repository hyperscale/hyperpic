// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/hyperscale/hyperpic/cmd/hyperpic/app"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := app.Run(); err != nil {
		if errors.Cause(err) == context.Canceled {
			log.Debug().Err(err).Msg("ignore error since context is cancelled")
		} else {
			log.Fatal().Err(err).Msg("hyperlens run failed")
		}
	}
}
