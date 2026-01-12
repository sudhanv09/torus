package models

import (
	"database/sql"
	"sudhanv09/torus/db"
)

func AddTrackedItem(item *TrackedItem) error {
	query := `INSERT INTO tracked_items (type, external_id, title, year, poster_url, backdrop_path, overview, genres, path, status) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
			  RETURNING id, added_at`
	err := db.DB.QueryRow(query, item.Type, item.ExternalID, item.Title, item.Year, item.PosterURL, item.BackdropPath, item.Overview, item.Genres, item.Path, item.Status).
		Scan(&item.ID, &item.AddedAt)
	return err
}

func GetTrackedItems() ([]TrackedItem, error) {
	rows, err := db.DB.Query("SELECT id, type, external_id, title, year, poster_url, backdrop_path, overview, genres, path, status, added_at FROM tracked_items ORDER BY added_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TrackedItem
	for rows.Next() {
		var item TrackedItem
		if err := rows.Scan(&item.ID, &item.Type, &item.ExternalID, &item.Title, &item.Year, &item.PosterURL, &item.BackdropPath, &item.Overview, &item.Genres, &item.Path, &item.Status, &item.AddedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func GetTrackedItemByID(id int64) (*TrackedItem, error) {
	var item TrackedItem
	query := "SELECT id, type, external_id, title, year, poster_url, backdrop_path, overview, genres, path, status, added_at FROM tracked_items WHERE id = ?"
	err := db.DB.QueryRow(query, id).Scan(&item.ID, &item.Type, &item.ExternalID, &item.Title, &item.Year, &item.PosterURL, &item.BackdropPath, &item.Overview, &item.Genres, &item.Path, &item.Status, &item.AddedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &item, err
}

func GetTrackedItemsByType(item_type string) ([]TrackedItem, error) {
	rows, err := db.DB.Query("SELECT id, type, external_id, title, year, poster_url, backdrop_path, overview, genres, path, status, added_at FROM tracked_items WHERE type = ? ORDER BY added_at DESC", item_type)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TrackedItem
	for rows.Next() {
		var item TrackedItem
		if err := rows.Scan(&item.ID, &item.Type, &item.ExternalID, &item.Title, &item.Year, &item.PosterURL, &item.BackdropPath, &item.Overview, &item.Genres, &item.Path, &item.Status, &item.AddedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func UpdateTrackedItemStatus(id int64, status string) error {
	_, err := db.DB.Exec("UPDATE tracked_items SET status = ? WHERE id = ?", status, id)
	return err
}

func DeleteTrackedItem(id int64) error {
	_, err := db.DB.Exec("DELETE FROM tracked_items WHERE id = ?", id)
	return err
}

