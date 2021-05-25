package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"
)

func main() {
	iterate_map()
	b := make([]byte, 10)
	_, _ = rand.Read(b)
	fmt.Println(b)
	_ = sha256.New()


	var Score float64
	var Score2 float32
	

	var Ball int
	table := make(chan int)
	go player(table)
	go player(table)
	table <- Ball
	Score = Score + 1.05
	Score2 = Score2 + 1.0000000000072
	time.Sleep(1 * time.Second)
	<-table

	mm := make(map[int]bool)
	mm[0] = false
	mm[1] = true
	for k, v := range mm {
		fmt.Println(k, "-->", v)
	}
}

func player(table chan int) {
	for {
		ball := <-table
		ball++
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

func iterate_map() {
	var m map[int]string = map[int]string{1: "One", 2: "Two", 3: "Three"}
	for k, v := range m {
		fmt.Printf("%d --> %s\n", k, v)
	}
}
