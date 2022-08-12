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

package logutil

import (
	"bytes"
	"io"
	"sync"

	"github.com/sirupsen/logrus"
)

func Pause(l *logrus.Logger) func() {
	// initialize formatter with original terminal settings
	l.Formatter.Format(logrus.NewEntry(l))

	bw := newBufferedWriter(l.Out)
	l.SetOutput(bw)
	return func() {
		bw.resume()
	}
}

type bufferedWriter struct {
	mu  sync.Mutex
	buf *bytes.Buffer
	w   io.Writer
}

func newBufferedWriter(w io.Writer) *bufferedWriter {
	return &bufferedWriter{
		buf: bytes.NewBuffer(nil),
		w:   w,
	}
}

func (bw *bufferedWriter) Write(p []byte) (int, error) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	if bw.buf == nil {
		return bw.w.Write(p)
	}
	return bw.buf.Write(p)
}

func (bw *bufferedWriter) resume() {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	if bw.buf == nil {
		return
	}
	io.Copy(bw.w, bw.buf)
	bw.buf = nil
}
