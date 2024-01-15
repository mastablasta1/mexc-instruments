package stakan

/*
https://bybit-exchange.github.io/docs/v5/market/orderbook
аски и биды перепутаны в выдаче
*/
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	b "mexcstakan/base"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Получим обновления
/*
	принимает текущее время
*/

func StartParser() {

	//	Начальный замер времени
	//timeOne := time.Now().Unix()

	//	Сколько записей берем
	limit := 10

	//	Карта с символами для парсинга, и количество потоков
	symbols, count := b.GetDataFromMysql(limit)

	// Создаем канал для получения данных
	ch := make(chan b.Mexcs, limit)

	var wg sync.WaitGroup

	// Запуск горутин
	for i := 0; i < count; i++ {

		wg.Add(1)
		go getStakanData(i, ch, &wg, symbols[i])

	}

	// Создаем слайс для сбора данных

	var queryUpdArray []string

	// Запускаем горутину для чтения данных из канала и добавления их в слайс

	go func() {
		for value := range ch {

			symbol := value.Symbol
			volume := strconv.FormatFloat(value.Volume, 'f', -1, 64)
			quoteVolume := strconv.FormatFloat(value.QuoteVolume, 'f', -1, 64)
			updatedAt := strconv.FormatInt(value.UpdatedAt, 10)
			askOne := strconv.FormatFloat(value.AskOne, 'f', -1, 64)
			askDuo := strconv.FormatFloat(value.AskDuo, 'f', -1, 64)
			bidOne := strconv.FormatFloat(value.BidOne, 'f', -1, 64)
			bidDuo := strconv.FormatFloat(value.BidDuo, 'f', -1, 64)
			raznitca := strconv.FormatFloat(value.Raznitca, 'f', -1, 64)
			tradesCountOld := strconv.Itoa(value.TradesCountOld)
			tradesCountNew := strconv.Itoa(value.TradesCountNew)

			queryUpdArray = append(queryUpdArray, "('"+symbol+"','"+volume+"','"+quoteVolume+"','"+updatedAt+"','"+askOne+"','"+askDuo+"','"+bidOne+"','"+bidDuo+"','"+raznitca+"','"+tradesCountOld+"','"+tradesCountNew+"','1')")

		}
	}()

	// Ожидание завершения всех горутин
	wg.Wait()

	// Закрываем канал, чтобы избежать блокировки
	close(ch)

	/*
		//	Итоговый замер времени
		timeTwo := time.Now().Unix()

		//	Итоговый вывод данных
		fmt.Println("=================")
		fmt.Println("Количество символов: ", count)
		fmt.Println("Начало, конец, Разница: ", timeOne, timeTwo, timeTwo-timeOne)
	*/
	//	Отправим данные для занесения в БД
	megaUpdateString := strings.Join(queryUpdArray, ",")
	updateData(megaUpdateString)

	//	ЗАснем на секунду
	time.Sleep(1 * time.Second)

}

func getStakanData(id int, ch chan b.Mexcs, wg *sync.WaitGroup, mysqlDataStruct b.Mexcs) {

	defer wg.Done()

	//	Глубина стакана по каждому направлению
	//var depth int = 5

	//	Сформируем ссылку на парсинг
	url := "https://api.mexc.com/api/v3/depth?symbol=" + mysqlDataStruct.Symbol

	resp, err := http.Get(url)
	if err != nil {

		log.Fatalf("Ошибка при выполнении GET-запроса: %s", err)

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		log.Fatalf("Ошибка при чтении данных ответа: %s", err)

	}

	//	Разберем JSON в структуру ошибки
	var stakanDataStruct StakanDataStruct

	err = json.Unmarshal(body, &stakanDataStruct)
	if err != nil {

		//	Ошибка разбора
		time.Sleep(1 * time.Second)

	}

	/*
		Если в стакане есть данные по текушему времени, значит торги по этой паре идут
		Если в предложения в Аске или в Бидах отсутствуют то все окау
	*/

	if len(stakanDataStruct.Asks) > 0 {
		if len(stakanDataStruct.Bids) > 0 {

			/*
				ask — цена покупки
				Buy - покупка

				bid — цена продажи
				Sell - продажа
			*/

			//  Берем первые данные в стакане с обеих сторн
			mysqlDataStruct.BidOne, _ = strconv.ParseFloat(stakanDataStruct.Bids[0][0], 64)
			mysqlDataStruct.BidDuo, _ = strconv.ParseFloat(stakanDataStruct.Bids[0][1], 64)

			mysqlDataStruct.AskOne, _ = strconv.ParseFloat(stakanDataStruct.Asks[0][0], 64)
			mysqlDataStruct.AskDuo, _ = strconv.ParseFloat(stakanDataStruct.Asks[0][1], 64)

			fmt.Println("Разница: BidOne: ", mysqlDataStruct.BidOne)
			fmt.Println("Разница: AskOne: ", mysqlDataStruct.AskOne)

			mysqlDataStruct.Raznitca = (mysqlDataStruct.BidOne/mysqlDataStruct.AskOne - 1) * 100
			fmt.Println("Разница: ", mysqlDataStruct.Raznitca)
			fmt.Println("------------")

		} else {

			mysqlDataStruct.BidOne = 0
			mysqlDataStruct.BidDuo = 0
			mysqlDataStruct.Raznitca = 0

		}
	} else {

		mysqlDataStruct.AskOne = 0
		mysqlDataStruct.AskDuo = 0
		mysqlDataStruct.Raznitca = 0

	}

	//	Время обновления данных
	mysqlDataStruct.UpdatedAt = time.Now().Unix()
	mysqlDataStruct.Active = 1

	//	Вернем полученное в канал
	ch <- mysqlDataStruct

}

// Апдейтим данные
func updateData(megaUpdateString string) {

	//	Занесем свежие данные и включим символы, где есть данные
	query2 := "REPLACE INTO `mexcs`(`symbol`,`volume`,`quote_volume`,`updated_at`,`ask_one`,`ask_duo`,`bid_one`,`bid_duo`,`raznitca`,`trades_count_old`,`trades_count_new`,`active`) VALUES" + megaUpdateString //	Запрос создадим
	b.InsertDataInMysql(query2)

}
