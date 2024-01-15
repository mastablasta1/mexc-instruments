package base

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Job Connect
var Login = "work_work"
var Password = "nb2ZS{exLK]n"
var Host = "38.180.7.209"
var Base = "work_work"

// Обновляем данные в БД по символам
func InsertDataInMysql(query string) {

	// Подключение к базе данных MySQL
	db, err := sql.Open("mysql", Login+":"+Password+"@tcp("+Host+":3306)/"+Base)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Выполняем обновление записи
	result, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

	// Получаем количество обновленных строк
	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

}

// Делаем выборку из БД по символам
func GetDataFromMysql(limit int) ([]Mexcs, int) {

	// Подключение к базе данных MySQL
	db, err := sql.Open("mysql", Login+":"+Password+"@tcp("+Host+":3306)/"+Base)
	if err != nil {
		log.Fatal(err)

	}
	defer db.Close()

	query := "SELECT * FROM `mexcs` ORDER BY `updated_at` ASC LIMIT " + strconv.Itoa(limit) + ""
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)

	}
	defer rows.Close()

	var data []Mexcs

	for rows.Next() {
		var d Mexcs
		err := rows.Scan(
			&d.Symbol,
			&d.Volume,
			&d.QuoteVolume,
			&d.UpdatedAt,
			&d.AskOne,
			&d.AskDuo,
			&d.BidOne,
			&d.BidDuo,
			&d.Raznitca,
			&d.TradesCountOld,
			&d.TradesCountNew,
			&d.Active,
		)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, d)
	}

	return data, len(data)

}

// Обновляем данные в БД по данным
func InsertDataInMysqlDatas(query string) {

	// Подключение к базе данных MySQL
	db, err := sql.Open("mysql", Login+":"+Password+"@tcp("+Host+":3306)/"+Base)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Выполняем обновление записи
	result, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

	// Получаем количество обновленных строк
	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

}
