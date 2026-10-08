// Copyright 2026 Sonic Operations Ltd
// This file is part of the Sonic Client
//
// Sonic is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Sonic is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Sonic. If not, see <http://www.gnu.org/licenses/>.

package ethapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTraceTx_RestrictsTracersToWhitelist(t *testing.T) {
	const customJS = "{step:function(){},result:function(){return 0}}"

	for name, tracer := range map[string]string{
		"custom JS": customJS,
		"unknown":   "muxTracer",
	} {
		t.Run(name+" is rejected by default", func(t *testing.T) {
			api := NewPublicDebugAPI(nil, false)
			_, err := api.traceTx(context.Background(), nil, nil, nil, nil, &TraceConfig{Tracer: &tracer})
			require.ErrorContains(t, err, "not permitted")
		})
	}

	t.Run("custom JS is not rejected with allowJSTracers", func(t *testing.T) {
		tracer := customJS
		api := NewPublicDebugAPI(nil, true)
		_, err := api.traceTx(context.Background(), nil, nil, nil, nil, &TraceConfig{Tracer: &tracer})
		if err != nil {
			require.NotContains(t, err.Error(), "not permitted")
		}
	})
}
