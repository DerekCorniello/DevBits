package database

import (
	"database/sql"
	"fmt"

	"backend/api/internal/types"
)

func QueryProject(id int) (*types.Project, error) {
	query := `SELECT id, name, description, status, likes, links, tags, owner, creation_date FROM Projects WHERE id = ?;`
	row := DB.QueryRow(query, id)
	var project types.Project
	var linksJSON, tagsJSON string

	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&project.Likes,
		&linksJSON,
		&tagsJSON,
		&project.Owner,
		&project.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := UnmarshalFromJSON(linksJSON, &project.Links); err != nil {
		return nil, err
	}
	if err := UnmarshalFromJSON(tagsJSON, &project.Tags); err != nil {
		return nil, err
	}

	return &project, nil
}

func QueryCreateProject(proj *types.Project) (int64, error) {
	linksJSON, err := MarshalToJSON(proj.Links)
	if err != nil {
		return -1, err
	}

	tagsJSON, err := MarshalToJSON(proj.Tags)
	if err != nil {
		return -1, err
	}

	query := `INSERT INTO Projects (name, description, status, links, tags, owner)
              VALUES (?, ?, ?, ?, ?, ?);`

	res, err := DB.Exec(query, proj.Name, proj.Description, proj.Status, string(linksJSON), string(tagsJSON), proj.Owner)
	if err != nil {
		return -1, fmt.Errorf("Failed to create project `%v`: %v", proj.Name, err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to ensure proj was created: %v", err)
	}

	return lastId, nil
}

func QueryDeleteProject(id int) (int16, error) {
	query := `DELETE from Projects WHERE id=?;`
	res, err := DB.Exec(query, id)
	if err != nil {
		return 400, fmt.Errorf("Failed to delete project `%v`: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return 404, fmt.Errorf("Deletion did not affect any records")
	} else if err != nil {
		return 500, fmt.Errorf("Failed to fetch affected rows: %v", err)
	}

	return 200, nil
}

func QueryUpdateProject(id int, updatedData map[string]interface{}) error {
	query := `UPDATE Projects SET `
	var args []interface{}

	queryParams, args, err := BuildUpdateQuery(updatedData)
	if err != nil {
		return fmt.Errorf("Error building query: %v", err)
	}

	query += queryParams + " WHERE id = ?"
	args = append(args, id)

	rowsAffected, err := ExecUpdate(query, args...)
	if err != nil {
		return fmt.Errorf("Error executing update query: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No project found with id `%d` to update", id)
	}

	return nil
}
