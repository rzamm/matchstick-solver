package run

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestWorkers(t *testing.T) {
	found := make(chan *taskReturn)

	task := func(p *taskParams) {
		random := rand.Intn(5000) - 2500
		time.Sleep(time.Duration(3000+random) * time.Millisecond)
		fmt.Println(p.removeCombIndex)

		if p.removeCombIndex == 200 {
			found <- &taskReturn{}
		}
	}

	workers := Workers(task, found)

	go func() {
		for i := 0; i < 2000; i++ {
			workers <- &taskParams{
				removeCombIndex: i,
			}
		}
		close(workers)
	}()

	for range found {
		fmt.Println("Found!")
		break
	}
}
