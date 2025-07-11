package logger

import (
	"io"
	"log"
	"os"
)

type MLogger struct {
	logFile *os.File
}

func NewLogger(filePath string) *MLogger {
	var newLog MLogger
	if filePath != "" {
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalf("Không thể mở file log: %v", err)
		}
		newLog.logFile = logFile

		// Ghi log ra cả terminal (stdout) và file
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
	}

	return &newLog
}

func (l *MLogger) Log(format string, v ...any) {
	log.Printf(format, v...)
}

func (l *MLogger) Close() {
	l.logFile.Close()
	log.SetOutput(os.Stderr)
}
