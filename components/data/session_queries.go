package data

const (
	insertSessions = `
		INSERT INTO sessions (name, working_directory) 
		VALUES (?, ?)
	`
	updateSession = `
		UPDATE sessions 
		SET name = ?, working_directory = ? 
		WHERE id = ?
	`
	selectAllSessions = `
		SELECT id, name 
		FROM sessions
		ORDER BY ID
	`
	selectSessionByID = `
		SELECT id, name, working_directory 
		FROM sessions 
		WHERE id = ?
	`
	deleteSessionByID = `
		DELETE FROM sessions 
		WHERE id = ?
	`
	selectAllSessionWindows = `
		SELECT w.id, w.name, w.working_directory
		FROM Session_Windows sw
		JOIN Windows w ON sw.window_id = w.id
		WHERE sw.session_id = ?
	`
	selectSessionWorkingDirectory = `
		SELECT working_directory 
		FROM Sessions 
		WHERE id = ?
	`
	selectFirstSession = `
		SELECT s.id, s.name, s.working_directory 
		FROM Sessions s 
		LIMIT 1
	`
)
