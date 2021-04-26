package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"
)

func main() {
	//hello := "Hello"
	//world := "world!"
	//num := strconv.Itoa(rand.Int())
	//words := []string{hello, num, world}
	//print_words(words)
	//loop(rand.Int())
	iterate_map()
	b := make([]byte, 10)
	_, _ = rand.Read(b)
	fmt.Println(b)
	_ = sha256.New()

	var Ball int
	table := make(chan int)
	go player(table)
	go player(table)
	table <- Ball
	time.Sleep(1 * time.Second)
	<-table
}

func player(table chan int) {
	for {
		ball := <-table
		ball++
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

//func print_words(words []string) {
//word_str := strings.Join(words, ",")
//fmt.Println(word_str)
//}

//func loop(num int) {
//iter_num := num % 1000
//sum := 0
//for i := 0; i < iter_num; i++ {
//sum += i
//}
//i := 0
//sum = 0
//for i < iter_num {
//sum += i
//i += 1
//}
//}

func iterate_map() {
	var m map[int]string = map[int]string{1: "One", 2: "Two", 3: "Three"}
	for k, v := range m {
		fmt.Printf("%d --> %s\n", k, v)
	}
}
