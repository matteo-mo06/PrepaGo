package db

import (
	"PrepaGo/models"
	"database/sql"
	"fmt"
)

func GetAllEntityB() ([]models.EntityB, error) {
	var entities []models.EntityB

	rows, err := Conn.Query("SELECT id, name, max_varchar_ex, bool_ex, int_ex FROM entity_b")
	if err != nil {
		return nil, fmt.Errorf("getAllEntityB : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var entity models.EntityB
		err := rows.Scan(&entity.Id, &entity.Name, &entity.MaxVarcharEx, &entity.BoolEx, &entity.IntEx)
		if err != nil {
			return nil, fmt.Errorf("getAllEntityB scan : %v", err.Error())
		}
		entities = append(entities, entity)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("getAllEntityB rows : %v", err.Error())
	}

	return entities, nil
}

func GetEntityBById(id int) (*models.EntityB, error) {
	var entity models.EntityB

	row := Conn.QueryRow(
		"SELECT id, name, max_varchar_ex, bool_ex, int_ex FROM entity_b WHERE id = ?",
		id)
	err := row.Scan(&entity.Id, &entity.Name, &entity.MaxVarcharEx, &entity.BoolEx, &entity.IntEx)
	if err != nil {
		return nil, fmt.Errorf("getEntityBById : %v", err.Error())
	}

	return &entity, nil
}

func CreateEntityB(entity models.EntityB) (int64, error) {
	var result sql.Result
	var err error

	if entity.Id == 0 {
		result, err = Conn.Exec(
			"INSERT INTO entity_b (name, max_varchar_ex, bool_ex, int_ex) VALUES (?, ?, ?, ?)",
			entity.Name, entity.MaxVarcharEx, entity.BoolEx, entity.IntEx)
	} else {
		result, err = Conn.Exec(
			"INSERT INTO entity_b (id, name, max_varchar_ex, bool_ex, int_ex) VALUES (?, ?, ?, ?, ?)",
			entity.Id, entity.Name, entity.MaxVarcharEx, entity.BoolEx, entity.IntEx)
	}

	if err != nil {
		return 0, fmt.Errorf("createEntityB : %v", err.Error())
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("createEntityB lastInsertId : %v", err.Error())
	}

	return lastId, nil
}

func UpdateEntityBBoolEx(id int, boolEx bool) error {
	result, err := Conn.Exec(
		"UPDATE entity_b SET bool_ex = ? WHERE id = ?",
		boolEx, id)
	if err != nil {
		return fmt.Errorf("updateEntityBBoolEx : %v", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updateEntityBBoolEx rowsAffected : %v", err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("updateEntityBBoolEx : entity not found")
	}

	return nil
}

func DeleteEntityB(id int) error {
	result, err := Conn.Exec("DELETE FROM entity_b WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleteEntityB : %v", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteEntityB rowsAffected : %v", err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("deleteEntityB : entity not found")
	}

	return nil
}

func CheckEntityBExists(id int) (bool, error) {
	var count int
	err := Conn.QueryRow("SELECT COUNT(*) FROM entity_b WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checkEntityBExists : %v", err.Error())
	}
	return count > 0, nil
}