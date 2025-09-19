package main

import (
	"fmt"
	"math"
)

// Fungsi untuk menghitung faktorial
func factorial(n int) uint64 {
	if n == 0 {
		return 1
	}
	result := uint64(1)
	for i := 2; i <= n; i++ {
		result *= uint64(i)
	}
	return result
}

// Fungsi untuk membulatkan ke atas
func calcValue(n int) uint64 {
	num := float64(factorial(n))
	den := math.Pow(2, float64(n))
	result := math.Ceil(num / den)
	return uint64(result)
}

func main() {
	for i := 0; i <= 10; i++ {
		fmt.Println("f(", i, ") =", calcValue(i))
	}
}
