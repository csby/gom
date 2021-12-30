package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var stdWriter io.Writer = os.Stdout

type LogWriter struct {
	sync.Mutex

	folder string
	year   int
	month  time.Month
	day    int
	file   *os.File
}

func (s *LogWriter) Write(p []byte) (n int, err error) {
	if len(p) < 1 {
		return 0, nil
	}

	writer := s.getLogger()
	if writer != nil {
		s.Lock()
		defer s.Unlock()
		return writer.Write(p)
	}

	return 0, nil
}

func (s *LogWriter) Close() error {
	s.Lock()
	defer s.Unlock()

	if s.file != nil {
		err := s.file.Close()
		s.file = nil
		return err
	}

	return nil
}

func (s *LogWriter) Info(v ...interface{}) string {
	return s.output("INFO", fmt.Sprint(v...))
}

func (s *LogWriter) Error(v ...interface{}) string {
	return s.output("ERROR", fmt.Sprint(v...))
}

func (s *LogWriter) output(l string, m string) string {
	now := time.Now()
	str := fmt.Sprintf("%s %d --- %s", l, os.Getpid(), m)
	msg := fmt.Sprintf("%s %s", now.Format("2006-01-02 15:04:05.000"), str)

	s.Write([]byte(fmt.Sprintln(msg)))

	return str
}

func (s *LogWriter) getLogger() io.Writer {
	if s.folder != "" {
		now := time.Now()
		if s.year != now.Year() || s.month != now.Month() || s.day != now.Day() || s.file == nil {
			s.Lock()
			defer s.Unlock()

			s.year = now.Year()
			s.month = now.Month()
			s.day = now.Day()
			if s.file != nil {
				s.file.Close()
				s.file = nil
			}

			os.MkdirAll(s.folder, 0777)
			fileName := fmt.Sprintf("%d-%d-%d.log", now.Year(), int(now.Month()), now.Day())
			filePath := filepath.Join(s.folder, fileName)
			file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return nil
			}
			s.file = file
		}
	}

	if s.file != nil {
		return s.file
	}

	return stdWriter
}
