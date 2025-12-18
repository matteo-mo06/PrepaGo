package db

import (
	"PrepaGo/models"
	"database/sql"
	"fmt"
)

func GetAllEntityA() ([]models.EntityA, error) {
	var entities []models.EntityA

	rows, err := Conn.Query("SELECT id, name, decimal_ex FROM entity_a")
	if err != nil {
		return nil, fmt.Errorf("getAllEntityA : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var entity models.EntityA
		err := rows.Scan(&entity.Id, &entity.Name, &entity.DecimalEx)
		if err != nil {
			return nil, fmt.Errorf("getAllEntityA scan : %v", err.Error())
		}
		entities = append(entities, entity)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("getAllEntityA rows : %v", err.Error())
	}

	return entities, nil
}

func GetEntityAById(id int) (*models.EntityA, error) {
	var entity models.EntityA

	row := Conn.QueryRow("SELECT id, name, decimal_ex FROM entity_a WHERE id = ?", id)
	err := row.Scan(&entity.Id, &entity.Name, &entity.DecimalEx)
	if err != nil {
		return nil, fmt.Errorf("getEntityAById : %v", err.Error())
	}

	return &entity, nil
}

func CreateEntityA(entity models.EntityA) (int64, error) {
	var result sql.Result
	var err error

	if entity.Id == 0 {
		result, err = Conn.Exec(
			"INSERT INTO entity_a (name, decimal_ex) VALUES (?, ?)",
			entity.Name, entity.DecimalEx)
	} else {
		result, err = Conn.Exec(
			"INSERT INTO entity_a (id, name, decimal_ex) VALUES (?, ?, ?)",
			entity.Id, entity.Name, entity.DecimalEx)
	}

	if err != nil {
		return 0, fmt.Errorf("createEntityA : %v", err.Error())
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("createEntityA lastInsertId : %v", err.Error())
	}

	return lastId, nil
}

func UpdateEntityA(id int, entity models.EntityA) error {
	result, err := Conn.Exec(
		"UPDATE entity_a SET name = ?, decimal_ex = ? WHERE id = ?",
		entity.Name, entity.DecimalEx, id)
	if err != nil {
		return fmt.Errorf("updateEntityA : %v", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updateEntityA rowsAffected : %v", err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("updateEntityA : entity not found")
	}

	return nil
}

func DeleteEntityA(id int) error {
	result, err := Conn.Exec("DELETE FROM entity_a WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleteEntityA : %v", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteEntityA rowsAffected : %v", err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("deleteEntityA : entity not found")
	}

	return nil
}

func CheckEntityAExists(id int) (bool, error) {
	var count int
	err := Conn.QueryRow("SELECT COUNT(*) FROM entity_a WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checkEntityAExists : %v", err.Error())
	}
	return count > 0, nil
}