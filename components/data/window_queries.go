package data

const (
	insertWindow = `
		INSERT INTO windows (name, working_directory) 
		VALUES (?, ?)
	`
	updateWindow = `
		UPDATE windows SET name = ?, working_directory = ? 
		WHERE id = ?
	`
	selectAllWindows = `
		SELECT id, name, working_directory 
		FROM windows
	`
	selectWindowByID = `
		SELECT id, name, working_directory 
		FROM windows 
		WHERE id = ?
	`
	deleteWindowByID = `
		DELETE FROM windows 
		WHERE id = ?
	`
	selectWindowPanes = `
		SELECT p.id, p.name, p.working_directory, p.split_direction, p.split_ratio 
		FROM windows_panes wp 
		JOIN windows w ON wp.window_id = w.id 
		JOIN panes p ON wp.pane_id = p.id
		WHERE w.id = ?
	`
	selectFirstWindow = `
		SELECT id, name, working_directory
		FROM windows
		LIMIT 1
	`
)
