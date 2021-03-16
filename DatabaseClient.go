package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var DB *sql.DB

func OpenDatabaseConnection() {
	var err error

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	DB, err = sql.Open(os.Getenv("DB_DIALECT"), connectionString)
	if err != nil {
		panic(err)

	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}

func GetAllFaces() interface{} {
	OpenDatabaseConnection()

	sqlQuery := fmt.Sprintf("SELECT f.id, f.file_path FROM Face f WHERE f.status = 1")

	rows, err := DB.Query(sqlQuery)
	var faceList []Face

	if err == nil {
		var profileId int
		var filePath string

		for rows.Next() {
			face := Face{}
			_ = rows.Scan(&profileId, &filePath)
			face.ProfileId = profileId
			face.FilePath = filePath
			faceList = append(faceList, face)
		}
	}
	return faceList
}


func SaveAttendance(faceId int, date time.Time) error {
	sqlQuery := fmt.Sprintf("SELECT a.id FROM Attendance a WHERE a.face_id = %d AND DATE(a.date_time) = '%s' AND a.status = 1", faceId, date.Format("2006-01-02"))

	rows, err := DB.Query(sqlQuery)

	if err != nil {
		return err
	}

	if rows != nil {
		if !rows.Next() {
			sqlQuery = fmt.Sprintf("INSERT INTO Attendance (face_id, status, date_time) VALUES (%d, 1, '%s')", faceId, date.Format("2006-01-02 15:04:05"))
			_, err = DB.Exec(sqlQuery)
		}
	}
	err = rows.Close()
	return err
}