// Copyright 2022 Docker Buildx authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracing

import (
	"context"
	"os"
	"strings"

	"github.com/moby/buildkit/util/tracing/detect"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TraceCurrentCommand(ctx context.Context, name string) (context.Context, func(error), error) {
	tp, err := detect.TracerProvider()
	if err != nil {
		return context.Background(), nil, err
	}
	ctx, span := tp.Tracer("").Start(ctx, name, trace.WithAttributes(
		attribute.String("command", strings.Join(os.Args, " ")),
	))

	return ctx, func(err error) {
		if err != nil {
			span.RecordError(err)
		}
		span.End()

		detect.Shutdown(context.TODO())
	}, nil
}
