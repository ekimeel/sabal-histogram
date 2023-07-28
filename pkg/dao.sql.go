package tq

const (
	sqlCreateTable = `
		CREATE TABLE IF NOT EXISTS plugin_histogram (
			point_id INTEGER PRIMARY KEY,
			last_updated TIMESTAMP,
			key_count INT,
			value_count INT,
			histogram TEXT,
			FOREIGN KEY (point_id) REFERENCES point(id),
			UNIQUE(point_id)
		);
	`
	sqlInsert = `INSERT INTO plugin_histogram (
        point_id,
        last_updated,
        key_count,
    	value_count,
    	histogram
    )
    VALUES (?, ?, ?, ?, ?)`
	sqlUpdate = `UPDATE plugin_histogram
        SET last_updated = ?,
            key_count = ?,
            value_count = ?,
            histogram = ?
        WHERE point_id = ?`
	sqlSelectByPointId = `SELECT point_id, last_updated, key_count, value_count, histogram FROM plugin_histogram WHERE point_id = ?`
)
