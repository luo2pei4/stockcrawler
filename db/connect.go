package db

import (
	"database/sql"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Connection 数据库
var db *sql.DB
var tx *sql.Tx

// NewConnection 新建连接
func NewConnection(name, driver, dsn string) error {

	var err error
	db, err = sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	fmt.Printf("Create %s connection successful.\n", name)

	return nil
}

// GetDBNowTime 获取数据库系统时间
func GetDBNowTime() string {

	rows, err := db.Query("select current_time() from dual")

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	nowTime := ""

	rows.Next()
	rows.Scan(&nowTime)

	return nowTime
}

// Select 根据传入的SQL文, 返回rows对象指针
func Select(sql string) (result *sql.Rows, err error) {

	result, err = db.Query(sql)

	if err != nil {
		return nil, err
	}

	return
}

// Execute 执行传入的SQL, 向数据库写入数据.
func Execute(sql string) (lastInsertID, rowsAffected int64, err error) {

	stmt, err := db.Prepare(sql)
	result, err := stmt.Exec()
	defer stmt.Close()

	if err != nil {
		return 0, 0, err
	}

	lastInsertID, _ = result.LastInsertId()
	rowsAffected, _ = result.RowsAffected()

	return
}

func TxBegin() error {
	var err error
	tx, err = db.Begin()
	if err != nil {
		return err
	}
	return nil
}

func TxCommit() error {
	err := tx.Commit()
	return err
}

func Rollback() error {
	err := tx.Rollback()
	return err
}
