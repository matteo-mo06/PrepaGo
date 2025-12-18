package db

import (
	"PrepaGo/models"
	"fmt"
)

func AddEntityAToEntityB(entityBId int, entityAId int) error {
	_, err := Conn.Exec(
		"INSERT INTO entity_b_entity_a (entity_b_id, entity_a_id) VALUES (?, ?)",
		entityBId, entityAId)
	if err != nil {
		return fmt.Errorf("addEntityAToEntityB : %v", err.Error())
	}
	return nil
}

func RemoveEntityAFromEntityB(entityBId int, entityAId int) error {
	result, err := Conn.Exec(
		"DELETE FROM entity_b_entity_a WHERE entity_b_id = ? AND entity_a_id = ?",
		entityBId, entityAId)
	if err != nil {
		return fmt.Errorf("removeEntityAFromEntityB : %v", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("removeEntityAFromEntityB rowsAffected : %v", err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("removeEntityAFromEntityB : relation not found")
	}

	return nil
}

func GetAllEntityAForEntityB(entityBId int) ([]models.EntityA, error) {
	var entities []models.EntityA

	rows, err := Conn.Query(`
		SELECT a.id, a.name, a.decimal_ex 
		FROM entity_a a
		INNER JOIN entity_b_entity_a ba ON a.id = ba.entity_a_id
		WHERE ba.entity_b_id = ?`,
		entityBId)
	if err != nil {
		return nil, fmt.Errorf("getAllEntityAForEntityB : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var entity models.EntityA
		err := rows.Scan(&entity.Id, &entity.Name, &entity.DecimalEx)
		if err != nil {
			return nil, fmt.Errorf("getAllEntityAForEntityB scan : %v", err.Error())
		}
		entities = append(entities, entity)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("getAllEntityAForEntityB rows : %v", err.Error())
	}

	return entities, nil
}

func CheckRelationExists(entityBId int, entityAId int) (bool, error) {
	var count int
	err := Conn.QueryRow(
		"SELECT COUNT(*) FROM entity_b_entity_a WHERE entity_b_id = ? AND entity_a_id = ?",
		entityBId, entityAId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checkRelationExists : %v", err.Error())
	}
	return count > 0, nil
}

func CountEntityAForEntityB(entityBId int) (int, error) {
	var count int
	err := Conn.QueryRow(
		"SELECT COUNT(*) FROM entity_b_entity_a WHERE entity_b_id = ?",
		entityBId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("countEntityAForEntityB : %v", err.Error())
	}
	return count, nil
}