package metrics

import "sync"

func FanIn(chs ...chan []Metric) chan []Metric {
	var wg sync.WaitGroup
	outCh := make(chan []Metric)

	output := func(c chan []Metric) {
		for m := range c {
			outCh <- m
		}
		wg.Done()
	}

	wg.Add(len(chs))
	for _, c := range chs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}
