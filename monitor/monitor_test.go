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

package monitor

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/sync/errgroup"
)

// TestMuxIO tests muxIO
func TestMuxIO(t *testing.T) {
	tests := []struct {
		name       string
		inputs     []instruction
		initIdx    int
		outputsNum int
		wants      []string

		// Everytime string is written to the mux stdin, the output end
		// that received the string write backs to the string that is masked with
		// its index number. This is useful to check if writeback is written from the
		// expected output destination.
		wantsMaskedOutput string
	}{
		{
			name: "single output",
			inputs: []instruction{
				input("foo\nbar\n"),
				toggle(),
				input("1234"),
				toggle(),
				input("456"),
			},
			initIdx:           0,
			outputsNum:        1,
			wants:             []string{"foo\nbar\n1234456"},
			wantsMaskedOutput: `^0+$`,
		},
		{
			name: "multi output",
			inputs: []instruction{
				input("foo\nbar\n"),
				toggle(),
				input("12" + string([]rune{rune(1)}) + "34abc"),
				toggle(),
				input("456"),
			},
			initIdx:           0,
			outputsNum:        3,
			wants:             []string{"foo\nbar\n", "1234abc", "456"},
			wantsMaskedOutput: `^0+1+2+$`,
		},
		{
			name: "multi output with nonzero index",
			inputs: []instruction{
				input("foo\nbar\n"),
				toggle(),
				input("1234"),
				toggle(),
				input("456"),
			},
			initIdx:           1,
			outputsNum:        3,
			wants:             []string{"456", "foo\nbar\n", "1234"},
			wantsMaskedOutput: `^1+2+0+$`,
		},
		{
			name: "multi output many toggles",
			inputs: []instruction{
				input("foo\nbar\n"),
				toggle(),
				input("1234"),
				toggle(),
				toggle(),
				input("456"),
				toggle(),
				input("%%%%"),
				toggle(),
				toggle(),
				toggle(),
				input("aaaa"),
			},
			initIdx:           0,
			outputsNum:        3,
			wants:             []string{"foo\nbar\n456", "1234%%%%aaaa", ""},
			wantsMaskedOutput: `^0+1+0+1+$`,
		},
		{
			name: "enable disable",
			inputs: []instruction{
				input("foo\nbar\n"),
				toggle(),
				input("1234"),
				toggle(),
				input("456"),
				disable(2),
				input("%%%%"),
				enable(2),
				toggle(),
				toggle(),
				input("aaa"),
				disable(2),
				disable(1),
				input("1111"),
				toggle(),
				input("2222"),
				toggle(),
				input("3333"),
			},
			initIdx:           0,
			outputsNum:        3,
			wants:             []string{"foo\nbar\n%%%%111122223333", "1234", "456aaa"},
			wantsMaskedOutput: `^0+1+2+0+2+0+$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inBuf, end, in := newTestIn(t)
			var outBufs []*outBuf
			var outs []ioSetOutContext
			if tt.outputsNum != len(tt.wants) {
				t.Fatalf("wants != outputsNum")
			}
			for i := 0; i < tt.outputsNum; i++ {
				outBuf, out := newTestOut(t, i)
				outBufs = append(outBufs, outBuf)
				outs = append(outs, ioSetOutContext{out, nil, nil})
			}
			mio := newMuxIO(in, outs, tt.initIdx, func(prev int, res int) string { return "" })
			for _, i := range tt.inputs {
				// Add input to muxIO
				istr, writeback := i(mio)
				if _, err := end.stdin.Write([]byte(istr)); err != nil {
					t.Fatalf("failed to write data to stdin: %v", err)
				}

				// Wait for writeback of this input
				var eg errgroup.Group
				eg.Go(func() error {
					outbuf := make([]byte, len(writeback))
					if _, err := io.ReadAtLeast(end.stdout, outbuf, len(outbuf)); err != nil {
						return err
					}
					return nil
				})
				eg.Go(func() error {
					errbuf := make([]byte, len(writeback))
					if _, err := io.ReadAtLeast(end.stderr, errbuf, len(errbuf)); err != nil {
						return err
					}
					return nil
				})
				if err := eg.Wait(); err != nil {
					t.Fatalf("failed to wait for output: %v", err)
				}
			}

			// Close stdin on this muxIO
			end.stdin.Close()

			// Wait for all output ends reach EOF
			mio.waitClosed()

			// Close stdout/stderr as well
			in.Close()

			// Check if each output end received expected string
			<-inBuf.doneCh
			for i, o := range outBufs {
				<-o.doneCh
				if o.stdin != tt.wants[i] {
					t.Fatalf("output[%d]: got %q; wanted %q", i, o.stdin, tt.wants[i])
				}
			}

			// Check if expected string is returned from expected outputs
			if !regexp.MustCompile(tt.wantsMaskedOutput).MatchString(inBuf.stdout) {
				t.Fatalf("stdout: got %q; wanted %q", inBuf.stdout, tt.wantsMaskedOutput)
			}
			if !regexp.MustCompile(tt.wantsMaskedOutput).MatchString(inBuf.stderr) {
				t.Fatalf("stderr: got %q; wanted %q", inBuf.stderr, tt.wantsMaskedOutput)
			}
		})
	}
}

type instruction func(m *muxIO) (intput string, writeBackView string)

func input(s string) instruction {
	return func(m *muxIO) (string, string) {
		return s, strings.ReplaceAll(s, string([]rune{rune(1)}), "")
	}
}

func toggle() instruction {
	return func(m *muxIO) (string, string) {
		return string([]rune{rune(1)}) + "c", ""
	}
}

func enable(i int) instruction {
	return func(m *muxIO) (string, string) {
		m.enable(i)
		return "", ""
	}
}

func disable(i int) instruction {
	return func(m *muxIO) (string, string) {
		m.disable(i)
		return "", ""
	}
}

type inBuf struct {
	stdout string
	stderr string
	doneCh chan struct{}
}

func newTestIn(t *testing.T) (*inBuf, ioSetOut, ioSetIn) {
	ti := &inBuf{
		doneCh: make(chan struct{}),
	}
	gotOutR, gotOutW := io.Pipe()
	gotErrR, gotErrW := io.Pipe()
	outR, outW := io.Pipe()
	var eg errgroup.Group
	eg.Go(func() error {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(io.MultiWriter(gotOutW, buf), outR); err != nil {
			return err
		}
		ti.stdout = buf.String()
		return nil
	})
	errR, errW := io.Pipe()
	eg.Go(func() error {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(io.MultiWriter(gotErrW, buf), errR); err != nil {
			return err
		}
		ti.stderr = buf.String()
		return nil
	})
	go func() {
		eg.Wait()
		close(ti.doneCh)
	}()
	inR, inW := io.Pipe()
	return ti, ioSetOut{inW, gotOutR, gotErrR}, ioSetIn{inR, outW, errW}
}

type outBuf struct {
	idx    int
	stdin  string
	doneCh chan struct{}
}

func newTestOut(t *testing.T, idx int) (*outBuf, ioSetOut) {
	to := &outBuf{
		idx:    idx,
		doneCh: make(chan struct{}),
	}
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()
	errR, errW := io.Pipe()
	go func() {
		defer inR.Close()
		defer outW.Close()
		defer errW.Close()
		buf := new(bytes.Buffer)
		mw := io.MultiWriter(buf,
			writeMasked(outW, fmt.Sprintf("%d", to.idx)),
			writeMasked(errW, fmt.Sprintf("%d", to.idx)),
		)
		if _, err := io.Copy(mw, inR); err != nil {
			inR.CloseWithError(err)
			outW.CloseWithError(err)
			errW.CloseWithError(err)
			return
		}
		to.stdin = string(buf.Bytes())
		outW.Close()
		errW.Close()
		close(to.doneCh)
	}()
	return to, ioSetOut{inW, outR, errR}
}

func writeMasked(w io.Writer, s string) io.Writer {
	buf := make([]byte, 4096)
	pr, pw := io.Pipe()
	go func() {
		for {
			n, readErr := pr.Read(buf)
			if readErr != nil && readErr != io.EOF {
				pr.CloseWithError(readErr)
				return
			}
			var masked string
			for i := 0; i < n; i++ {
				masked += s
			}
			if _, err := w.Write([]byte(masked)); err != nil {
				pr.CloseWithError(err)
				return
			}
			if readErr == io.EOF {
				pr.Close()
				return
			}
		}
	}()
	return pw
}
