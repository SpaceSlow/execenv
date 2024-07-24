package execenv_test

import (
	"bytes"
	"fmt"
	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
	"io"
	"log"
	"net/http"
)

func Example() {
	// Запускаем сервер для хранения метрик
	go func() {
		// Объявляем хранилище метрик в памяти
		storage := storages.NewMemStorage()
		// Объявляем шлюз запросов
		router := routers.MetricRouter(storage)
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatal(err)
		}
	}()

	// Отправляем на сервер метрики в формате "http://localhost:8080/update/{type}/{name}/{value}"
	http.Post("http://localhost:8080/update/counter/foo/42", "text/plain", nil)
	http.Post("http://localhost:8080/update/gauge/bar/42.42", "text/plain", nil)

	// Можно также отправить метрики на сервер batch-запросом в формате json
	jsonMetric := []byte(`[
		{
			"id": "GolangDayMonth",
			"type": "gauge",
			"value": 10.11
		},
		{
			"id": "GolangYear",
			"type": "counter",
			"delta": 2009
		}
	]`)
	jsonMetricsReader := bytes.NewReader(jsonMetric)

	_, err := http.Post("http://localhost:8080/updates/", "application/json", jsonMetricsReader)
	if err != nil {
		log.Println("json error", err)
	}

	// Получим отправленные ранее метрики с сервера
	response, err := http.Get("http://localhost:8080/")
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)

	// Выведем полученные метрики
	metrics, _ := io.ReadAll(response.Body)
	fmt.Println(string(metrics))

	// Output:
	// foo = 42 (counter)
	// GolangYear = 2009 (counter)
	// bar = 42.42 (gauge)
	// GolangDayMonth = 10.11 (gauge)
}
