package data

const (
	insertPane = `
		INSERT INTO 
		Panes (name, working_directory, split_direction, split_ratio)
		VALUES (?, ?, ?, ?)

	`
	updatePane = `
		UPDATE Panes 
		SET name = ?, working_directory = ?, split_direction = ?, split_ratio = ? 
		WHERE id = ?
	`
	selectAllPanes = `
		SELECT id, name, working_directory, split_direction, split_ratio
		FROM Panes
	`
	selectPaneByID = `
		SELECT id, name, working_directory, split_direction, split_ratio
		FROM Panes
		WHERE id = ?
	`
	deletePaneByID = `
		DELETE FROM Panes 
		WHERE id = ?
	`
)
