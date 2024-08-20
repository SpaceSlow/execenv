package metrics

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_compress(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "medium text",
			data: []byte(`Dolores doloribus est ut qui cumque.
Ratione repudiandae placeat est voluptas sequi aliquam illum.
Voluptates a veniam quidem quidem explicabo doloribus. 
Explicabo perferendis voluptas harum eveniet et dolorem. 
Error tempore perspiciatis sit perferendis. 
At voluptatem soluta quod esse. Ipsa ducimus dolores et quam. 
Voluptates voluptas ea blanditiis placeat dolorem. 
Exercitationem animi deserunt repellat assumenda eveniet quam reprehenderit.
Consequatur fugiat et sequi.`),
		},
		{
			name: "large text",
			data: []byte(`Consequatur ut velit officia repudiandae quas. 
Quam voluptatem voluptatibus nesciunt et ut nobis. 
Molestiae placeat et non atque error omnis.
Accusamus neque quasi consequatur necessitatibus nihil iure.
Exercitationem minima amet tempore ratione aperiam aut in.
Nulla dolor omnis molestiae ex optio.

Maiores dolores placeat sunt odit quidem. 
Distinctio nesciunt rerum porro et. 
Quis quibusdam aut qui. Qui culpa earum sit.
Consequatur consequuntur autem tempore. Autem ex nesciunt minima officia quisquam natus. 
Et odio pariatur et. Dolorem consequatur voluptas nihil necessitatibus. Ut ipsam aut maiores adipisci similique omnis.

Aperiam dolore est molestiae non qui vel. Delectus nostrum quaerat dolor fugiat. 
Unde qui omnis reprehenderit beatae esse. Aliquid dolores tempora cum.
Repellat ducimus quam ut quasi occaecati. Id sint voluptatum non libero ipsum doloremque. 
Nulla rerum at consequatur deleniti officia ut voluptate enim. Aliquid dolorum rerum qui. 
Vel ratione voluptatem commodi voluptatem illo. Exercitationem sunt omnis voluptates consequatur.

Facilis reiciendis magnam explicabo quod repellendus et fuga cumque. Repudiandae sequi ut eos aliquid nemo.
Eum id aliquid delectus ipsum magni qui. Tempore non ab excepturi ut. Minima id est incidunt quaerat qui ratione. 
Qui omnis commodi blanditiis molestiae ut. Non ex ut nesciunt. Et repellat perferendis eos maiores ratione.
Eius vel aut possimus ipsa omnis. Consequatur quae dicta magnam incidunt consequuntur doloremque dolor. 
Omnis cupiditate tempore sit corporis nam et tempore a.
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compress(tt.data)
			require.NoError(t, err)
			assert.Less(t, len(got), len(tt.data))
		})
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func generateRandomMetric() Metric {
	mType := MetricType(rand.Intn(2) + 1)

	metric := Metric{
		Type: mType,
		Name: randStringBytes(rand.Intn(10) + 10),
	}

	switch mType {
	case Counter:
		metric.Value = rand.Int63n(1000)
	case Gauge:
		metric.Value = rand.Float64()
	}
	return metric
}

func Test_fanIn(t *testing.T) {
	metricsChs := make([]chan []Metric, 10)
	for i := range metricsChs {
		metricsChs[i] = make(chan []Metric, 1)
	}

	expected := make([]Metric, 100)
	for i := range expected {
		expected[i] = generateRandomMetric()
	}

	for i := 0; i < len(metricsChs); i++ {
		metricsChs[i] <- expected[i*10 : i*10+10]
		close(metricsChs[i])
	}

	got := make([]Metric, 0)
	for slice := range fanIn(metricsChs...) {
		got = append(got, slice...)
	}

	assert.ElementsMatch(t, expected, got)
}
