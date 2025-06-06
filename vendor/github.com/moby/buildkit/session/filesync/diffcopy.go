package filesync

import (
	"bufio"
	"context"
	io "io"
	"os"
	"time"

	"github.com/moby/buildkit/util/bklog"

	"github.com/pkg/errors"
	"github.com/tonistiigi/fsutil"
	fstypes "github.com/tonistiigi/fsutil/types"
	"google.golang.org/grpc"
)

type Stream interface {
	Context() context.Context
	SendMsg(m any) error
	RecvMsg(m any) error
}

func newStreamWriter(stream grpc.ClientStream) io.WriteCloser {
	wc := &streamWriterCloser{ClientStream: stream}
	return &bufferedWriteCloser{Writer: bufio.NewWriter(wc), Closer: wc}
}

type bufferedWriteCloser struct {
	*bufio.Writer
	io.Closer
}

func (bwc *bufferedWriteCloser) Close() error {
	if err := bwc.Flush(); err != nil {
		return errors.WithStack(err)
	}
	return bwc.Closer.Close()
}

type streamWriterCloser struct {
	grpc.ClientStream
}

func (wc *streamWriterCloser) Write(dt []byte) (int, error) {
	// grpc-go has a 4MB limit on messages by default. Split large messages
	// so we don't get close to that limit.
	const maxChunkSize = 3 * 1024 * 1024
	if len(dt) > maxChunkSize {
		n1, err := wc.Write(dt[:maxChunkSize])
		if err != nil {
			return n1, err
		}
		dt = dt[maxChunkSize:]
		var n2 int
		if n2, err = wc.Write(dt); err != nil {
			return n1 + n2, err
		}
		return n1 + n2, nil
	}

	if err := wc.SendMsg(&BytesMessage{Data: dt}); err != nil {
		// SendMsg return EOF on remote errors
		if errors.Is(err, io.EOF) {
			if err := errors.WithStack(wc.RecvMsg(struct{}{})); err != nil {
				return 0, err
			}
		}
		return 0, errors.WithStack(err)
	}
	return len(dt), nil
}

func (wc *streamWriterCloser) Close() error {
	if err := wc.CloseSend(); err != nil {
		return errors.WithStack(err)
	}
	// block until receiver is done
	var bm BytesMessage
	if err := wc.RecvMsg(&bm); !errors.Is(err, io.EOF) {
		return errors.WithStack(err)
	}
	return nil
}

func recvDiffCopy(ds grpc.ClientStream, dest string, cu CacheUpdater, progress progressCb, differ fsutil.DiffType, filter, metadataOnlyFilter func(string, *fstypes.Stat) bool) (err error) {
	st := time.Now()
	defer func() {
		bklog.G(ds.Context()).Debugf("diffcopy took: %v", time.Since(st))
	}()
	var cf fsutil.ChangeFunc
	var ch fsutil.ContentHasher
	if cu != nil {
		cu.MarkSupported(true)
		cf = cu.HandleChange
		ch = cu.ContentHasher()
	}
	defer func() {
		// tracing wrapper requires close trigger even on clean eof
		if err == nil {
			ds.CloseSend()
		}
	}()
	return errors.WithStack(fsutil.Receive(ds.Context(), ds, dest, fsutil.ReceiveOpt{
		NotifyHashed:  cf,
		ContentHasher: ch,
		ProgressCb:    progress,
		Filter:        fsutil.FilterFunc(filter),
		Differ:        differ,
		MetadataOnly:  metadataOnlyFilter,
	}))
}

func syncTargetDiffCopy(ds grpc.ServerStream, dest string) error {
	if err := os.MkdirAll(dest, 0700); err != nil {
		return errors.Wrapf(err, "failed to create synctarget dest dir %s", dest)
	}
	return errors.WithStack(fsutil.Receive(ds.Context(), ds, dest, fsutil.ReceiveOpt{
		Merge: true,
		Filter: func() func(string, *fstypes.Stat) bool {
			uid := os.Getuid()
			gid := os.Getgid()
			return func(p string, st *fstypes.Stat) bool {
				st.Uid = uint32(uid)
				st.Gid = uint32(gid)
				return true
			}
		}(),
	}))
}

func writeTargetFile(ds grpc.ServerStream, wc io.WriteCloser) error {
	var bm BytesMessage
	for {
		bm.Data = bm.Data[:0]
		if err := ds.RecvMsg(&bm); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return errors.WithStack(err)
		}
		if _, err := wc.Write(bm.Data); err != nil {
			return errors.WithStack(err)
		}
	}
}
