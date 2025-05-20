package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogWriter struct {
	*os.File
	date string
	name string
	path string
}

func New(path string, name string, flag int) {
	log.SetFlags(flag)
	logf := &LogWriter{path: path, name: name}
	logf.tofile()
	w := io.MultiWriter(logf, os.Stdout)
	log.SetOutput(w)
}

func (w *LogWriter) tofile() error {
	date := time.Now().Format("2006-01-02")
	if w.date != date {
		w.date = date

		fp := filepath.Join(w.path, fmt.Sprintf("%s-%s.log", w.name, w.date))
		f, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}

		w.File.Close()
		w.File = f
	}
	return nil
}

func (w *LogWriter) Writer(b []byte) {
	w.tofile()
	w.Write(b)
}
