package cmd

import (
	"bytes"
	"io"
	"sync"
)

type bufferedWriter struct {
	mutex         sync.Mutex
	bufferedInput bytes.Buffer
	writeCloser   io.WriteCloser
}

func newBufferedWriter(writeCloser io.WriteCloser) *bufferedWriter {
	return &bufferedWriter{
		writeCloser: writeCloser,
	}
}

func (w *bufferedWriter) Write(p []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.bufferedInput.Write(p)
}

func (w *bufferedWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	_, writeErr := w.bufferedInput.WriteTo(w.writeCloser)
	closeErr := w.writeCloser.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
