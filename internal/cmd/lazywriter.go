package cmd

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// A lazyWriter that buffers its input and writes everything at once on close.
type lazyWriter struct {
	mutex         sync.Mutex
	bufferedInput bytes.Buffer
	openFunc      func() (io.WriteCloser, error)
}

func newLazyWriter(openFunc func() (io.WriteCloser, error)) *lazyWriter {
	return &lazyWriter{
		openFunc: openFunc,
	}
}

func (w *lazyWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.openFunc == nil {
		return fmt.Errorf("lazyWriter already closed")
	}
	openFunc := w.openFunc
	w.openFunc = nil
	if w.bufferedInput.Len() == 0 {
		return nil
	}
	writeCloser, openErr := openFunc()
	if openErr != nil {
		return openErr
	}
	_, writeErr := w.bufferedInput.WriteTo(writeCloser)
	closeErr := writeCloser.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

func (w *lazyWriter) Write(p []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.openFunc == nil {
		return 0, fmt.Errorf("lazyWriter already closed")
	}
	return w.bufferedInput.Write(p)
}
