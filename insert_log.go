package logger

import "log"

var logQueue chan LogEntry

func startLoggerWorker(bufferSize int) {
	logQueue = make(chan LogEntry, bufferSize)

	go func() {
		for entry := range logQueue {
			if err := insertLog(entry); err != nil {
				log.Printf("failed to insert log entry %v", err)
			}
		}
	}()
}

func insertLog(entry LogEntry) error {
	_, err := db.Exec(`
		INSERT INTO logs (ip_address, url, action, details, timestamp, user_agent, status, latency)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, entry.IPAddress, entry.Url, entry.Action, entry.Details, entry.Timestamp, entry.UserAgent, entry.Status, entry.Latency)

	return err
}
