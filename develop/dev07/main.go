package main

import (
	"fmt"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	joinChannels(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v", time.Since(start))

}

func joinChannels(ch ...<-chan interface{}) (r chan interface{}) {
	ready := make(chan struct{})

	for _, v := range ch {
		go func(cin <-chan interface{}) {
			v, ok := <-cin
			if !ok {
				ready <- struct{}{}
			}
			r <- v
		}(v)
	}

	<-ready

	return
}
