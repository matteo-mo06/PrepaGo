package db

import (
	"PrepaGo/models"
	"database/sql"
	"fmt"
)

func GetAllLibraries() ([]models.Library, error) {
	var libraries []models.Library

	rows, err := Conn.Query("SELECT id, owner_name, owner_password, is_premium, creation_year FROM libraries")
	if err != nil {
		return nil, fmt.Errorf("getAllLibraries : %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var library models.Library
		err := rows.Scan(&library.Id, &library.OwnerName, &library.OwnerPassword, &library.IsPremium, &library.CreationYear)
		if err != nil {
			return nil, fmt.Errorf("getAllLibraries scan : %v", err.Error())
		}
		libraries = append(libraries, library)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("getAllLibraries rows : %v", err.Error())
	}

	return libraries, nil
}

func GetLibraryById(id int) (*models.Library, error) {
	var library models.Library

	row := Conn.QueryRow(
		"SELECT id, owner_name, owner_password, is_premium, creation_year FROM libraries WHERE id = ?",
		id)
	err := row.Scan(&library.Id, &library.OwnerName, &library.OwnerPassword, &library.IsPremium, &library.CreationYear)
	if err != nil {
		return nil, fmt.Errorf("getLibraryById : %v", err.Error())
	}

	return &library, nil
}

func GetLibraryByOwnerName(ownerName string) (*models.Library, error) {
	var library models.Library

	row := Conn.QueryRow(
		"SELECT id, owner_name, owner_password, is_premium, creation_year FROM libraries WHERE owner_name = ?",
		ownerName)
	err := row.Scan(&library.Id, &library.OwnerName, &library.OwnerPassword, &library.IsPremium, &library.CreationYear)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getLibraryByOwnerName : %v", err.Error())
	}

	return &library, nil
}

func CreateLibrary(library models.Library) (int64, error) {
	var result sql.Result
	var err error

	if library.Id == 0 {
		result, err = Conn.Exec(
			"INSERT INTO libraries (owner_name, owner_password, is_premium, creation_year) VALUES (?, ?, ?, ?)",
			library.OwnerName, library.OwnerPassword, library.IsPremium, library.CreationYear)
	} else {
		result, err = Conn.Exec(
			"INSERT INTO libraries (id, owner_name, owner_password, is_premium, creation_year) VALUES (?, ?, ?, ?, ?)",
			library.Id, library.OwnerName, library.OwnerPassword, library.IsPremium, library.CreationYear)
	}

	if err != nil {
		return 0, fmt.Errorf("createLibrary : %v", err.Error())
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("createLibrary lastInsertId : %v", err.Error())
	}

	return lastId, nil
}

func UpdateLibraryPremium(id int, isPremium bool) error {
	result, err := Conn.Exec(
		"UPDATE libraries SET is_premium = ? WHERE id = ?",
		isPremium, id)
	if err != nil {
		return fmt.Errorf("updateLibraryPremium : %v", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updateLibraryPremium rowsAffected : %v", err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("updateLibraryPremium : library not found")
	}

	return nil
}

func DeleteLibrary(id int) error {
	result, err := Conn.Exec("DELETE FROM libraries WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleteLibrary : %v", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteLibrary rowsAffected : %v", err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("deleteLibrary : library not found")
	}

	return nil
}

func CheckLibraryExists(id int) (bool, error) {
	var count int
	err := Conn.QueryRow("SELECT COUNT(*) FROM libraries WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checkLibraryExists : %v", err.Error())
	}
	return count > 0, nil
}