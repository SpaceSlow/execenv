package execenv_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/SpaceSlow/execenv/cmd/routers"
	"github.com/SpaceSlow/execenv/cmd/storages"
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
	res, err := http.Post("http://localhost:8080/update/counter/foo/42", "text/plain", nil)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	res, err = http.Post("http://localhost:8080/update/gauge/bar/42.42", "text/plain", nil)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

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

	res, err = http.Post("http://localhost:8080/updates/", "application/json", jsonMetricsReader)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

	// Получим отправленные ранее метрики с сервера
	response, err := http.Get("http://localhost:8080/")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Выведем полученные метрики
	metricsResponse, _ := io.ReadAll(response.Body)

	metricSlice := sort.StringSlice(strings.Split(string(metricsResponse), "\n"))
	metricSlice.Sort()
	fmt.Println(strings.Join(metricSlice, "\n"))

	// Output:
	// GolangDayMonth = 10.11 (gauge)
	// GolangYear = 2009 (counter)
	// bar = 42.42 (gauge)
	// foo = 42 (counter)
}
