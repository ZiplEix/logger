package logger

import "time"

type LogEntry struct {
	ID        int64     `db:"id" json:"id"`
	IPAddress string    `db:"ip_address" json:"ip_address"`
	Url       string    `db:"url" json:"url"`
	Action    string    `db:"action" json:"action"`
	Details   *string   `db:"details" json:"details,omitempty"` // nullable, omis si nil
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	UserAgent *string   `db:"user_agent" json:"user_agent,omitempty"` // nullable, omis si nil
	Status    *string   `db:"status" json:"status,omitempty"`         // nullable, omis si nil
	Latency   int64     `db:"latency" json:"latency"`
}
