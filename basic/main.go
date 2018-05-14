package main

import (
	"strconv"
	"time"
	"fmt"
)

func sample(ch chan string) {
	for i := 0; i < 19; i++ {
		ch <- "I'm sample1 num:" + strconv.Itoa(i)
		time.Sleep(1* time.Second)
	}
}

func sample2(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
		time.Sleep(2 * time.Second)
	}
}

func main() {
	ch1 := make(chan string, 3)
	ch2 := make(chan int, 5)
	for i := 0 ; i < 10; i++ {
		go sample(ch1)
		go sample2(ch2)
	}

	for {
		select {
		case str, ch1Check := <-ch1:
			if !ch1Check {
				fmt.Println("ch1Check false")
			} else {
				fmt.Println(str)
			}
		case p, ch2Check := <-ch2:
			if !ch2Check {
				fmt.Println("ch2Check false")
			} else {
				fmt.Println(p)
			}


		}
		fmt.Printf("ch1 has %d items", len(ch1))
		fmt.Printf("ch2 has %d items", len(ch2))
		time.Sleep(1 * time.Second)
	}


}

