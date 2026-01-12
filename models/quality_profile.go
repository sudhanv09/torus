package models

import (
	"encoding/json"
	"sudhanv09/torus/db"
	"sudhanv09/torus/engine"
)

func AddQualityProfile(p *engine.QualityProfile) error {
	res, _ := json.Marshal(p.PreferredResolutions)
	src, _ := json.Marshal(p.PreferredSources)
	cod, _ := json.Marshal(p.PreferredCodecs)
	blk, _ := json.Marshal(p.BlockedWords)
	grp, _ := json.Marshal(p.PreferredGroups)

	query := `INSERT INTO quality_profiles (name, preferred_resolutions, min_resolution, preferred_sources, preferred_codecs, blocked_words, preferred_groups, min_seeders) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	return db.DB.QueryRow(query,
		p.Name,
		string(res),
		p.MinResolution,
		string(src),
		string(cod),
		string(blk),
		string(grp),
		p.MinSeeders,
	).Scan(&p.ID)
}

func GetQualityProfiles() ([]engine.QualityProfile, error) {
	rows, err := db.DB.Query("SELECT id, name, preferred_resolutions, min_resolution, preferred_sources, preferred_codecs, blocked_words, preferred_groups, min_seeders FROM quality_profiles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []engine.QualityProfile
	for rows.Next() {
		var p engine.QualityProfile
		var res, src, cod, blk, grp string
		if err := rows.Scan(&p.ID, &p.Name, &res, &p.MinResolution, &src, &cod, &blk, &grp, &p.MinSeeders); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(res), &p.PreferredResolutions)
		json.Unmarshal([]byte(src), &p.PreferredSources)
		json.Unmarshal([]byte(cod), &p.PreferredCodecs)
		json.Unmarshal([]byte(blk), &p.BlockedWords)
		json.Unmarshal([]byte(grp), &p.PreferredGroups)
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func UpdateQualityProfile(p *engine.QualityProfile) error {
	res, _ := json.Marshal(p.PreferredResolutions)
	src, _ := json.Marshal(p.PreferredSources)
	cod, _ := json.Marshal(p.PreferredCodecs)
	blk, _ := json.Marshal(p.BlockedWords)
	grp, _ := json.Marshal(p.PreferredGroups)

	query := `UPDATE quality_profiles SET name = ?, preferred_resolutions = ?, min_resolution = ?, preferred_sources = ?, preferred_codecs = ?, blocked_words = ?, preferred_groups = ?, min_seeders = ? WHERE id = ?`
	_, err := db.DB.Exec(query,
		p.Name,
		string(res),
		p.MinResolution,
		string(src),
		string(cod),
		string(blk),
		string(grp),
		p.MinSeeders,
		p.ID,
	)
	return err
}

func DeleteQualityProfile(id int64) error {
	_, err := db.DB.Exec("DELETE FROM quality_profiles WHERE id = ?", id)
	return err
}
