package main

import (
	"fmt"
	"lecture01_homework/fizzbuzz"
)

func main() {
	for num := 0; num <= 100; num++ {
		fmt.Println(fizzbuzz.FizzBuzz(num))
	}
}
