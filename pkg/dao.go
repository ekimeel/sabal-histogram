package histogram

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	DB           *sql.DB
	singletonDao *dao
	onceDao      sync.Once
)

type dao struct {
	insertStmt          *sql.Stmt
	updateStmt          *sql.Stmt
	selectByPointIdStmt *sql.Stmt
}

func getDao() *dao {
	onceDao.Do(func() {
		singletonDao = &dao{}
		var err error

		singletonDao.createTableIfNotExists()

		singletonDao.insertStmt, err = DB.Prepare(sqlInsert)
		if err != nil {
			panic(fmt.Sprintf("failed to prepare statement: %v", err))
		}

		singletonDao.updateStmt, err = DB.Prepare(sqlUpdate)
		if err != nil {
			panic(fmt.Sprintf("failed to prepare statement: %v", err))
		}

		singletonDao.selectByPointIdStmt, err = DB.Prepare(sqlSelectByPointId)
		if err != nil {
			panic(fmt.Sprintf("failed to prepare statement: %v", err))
		}

	})
	return singletonDao
}
func (dao *dao) createTableIfNotExists() {
	_, err := DB.Exec(sqlCreateTable)
	if err != nil {
		panic(fmt.Sprintf("failed to create table: %v", err))
	}
}

func (dao *dao) insert(hist *Histogram) (int64, error) {

	data, err := json.Marshal(hist.Histogram)
	if err != nil {
		log.WithField("plugin", PluginName).Errorf("failed to marshal histogram: %s", err)
	}

	result, err := dao.insertStmt.Exec(
		hist.PointId,
		time.Now(),
		hist.KeyCount,
		hist.ValueCount,
		data,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *dao) update(hist *Histogram) (int64, error) {

	data, err := json.Marshal(hist.Histogram)
	if err != nil {
		log.WithField("plugin", PluginName).Errorf("failed to marshal histogram: %s", err)
	}

	result, err := dao.updateStmt.Exec(
		time.Now(),
		hist.KeyCount,
		hist.ValueCount,
		data,
		hist.PointId,
	)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (dao *dao) selectByPointId(pointID uint32) (*Histogram, error) {

	row := dao.selectByPointIdStmt.QueryRow(pointID)
	var hist Histogram
	var data string
	err := row.Scan(
		&hist.PointId,
		&hist.LastUpdated,
		&hist.KeyCount,
		&hist.ValueCount,
		&data,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	err = json.Unmarshal([]byte(data), &hist.Histogram)

	if err != nil {
		log.WithField("plugin", PluginName).Errorf("failed to unmarshal histogram: %s", err)
	}
	return &hist, nil
}
