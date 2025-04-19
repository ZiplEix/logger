package logger

func GetAllLogs() ([]LogEntry, error) {
	rows, err := db.Query(`
		SELECT id, ip_address, url, action, details, timestamp, user_agent, status, latency
		FROM logs
		ORDER BY timestamp DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var entry LogEntry
		err := rows.Scan(
			&entry.ID,
			&entry.IPAddress,
			&entry.Url,
			&entry.Action,
			&entry.Details,
			&entry.Timestamp,
			&entry.UserAgent,
			&entry.Status,
			&entry.Latency,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, entry)
	}

	return logs, rows.Err()
}

func GetLogById(id int64) (LogEntry, error) {
	row := db.QueryRow(`
		SELECT id, ip_address, url, action, details, timestamp, user_agent, status, latency
		FROM logs
		WHERE id = ?
	`, id)

	var entry LogEntry
	err := row.Scan(
		&entry.ID,
		&entry.IPAddress,
		&entry.Url,
		&entry.Action,
		&entry.Details,
		&entry.Timestamp,
		&entry.UserAgent,
		&entry.Status,
		&entry.Latency,
	)
	if err != nil {
		return LogEntry{}, err
	}

	return entry, nil
}
