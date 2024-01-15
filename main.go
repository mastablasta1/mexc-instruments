package main

/*
	Парсер API биржи bitget.com

	https://www.bitget.com/api-doc/common/release-note

	1. Бесконечный цикл
	2. Раз в пять минут сверяем актуальные валютные пары
	3. Сверяем

*/

import (
	sta "mexcstakan/stakan"
	sym "mexcstakan/symbols"
	tic "mexcstakan/tickers"
	"time"
)

func main() {

	//	Текущие данные по таймстампам
	var timeFive int64 = 0

	//	ОБновим данные по парам
	timeNow := time.Now().Unix()

	//	Класс обновления данных по парам
	sym.GetSymbolsUpdate(timeNow)

	for {

		if timeNow > timeFive {

			timeFive = timeNow + 30000 //	Задержка на 500 минут

			//	Обновление объема торгов за 24 часа
			tic.GetTickersUpdate(timeNow)

		}

		//	Обновим данные по данным в парах
		sta.StartParser()

		//	Тормознем на секунду
		//time.Sleep(time.Second)
		time.Sleep(60 * time.Second)

	}

}
