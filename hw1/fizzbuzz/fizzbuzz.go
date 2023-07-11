package fizzbuzz

import "strconv"

func FizzBuzz(num int) string {
	var S string
	if num%3 == 0 {
		S += "Fizz"
	}
	if num%5 == 0 {
		S += "Buzz"
	}
	if S == "" {
		S += strconv.Itoa(num)
	}
	return S
}
