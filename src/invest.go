package main

import (
	"log"
	//"math/rand"
)

func main() {
	var rate float64
	rate = 0.12
	times := 15
	var amountPerCycle float64
	amountPerCycle = 2000
	value := calCompound(rate, times, amountPerCycle)
	investAmount := float64(times) * amountPerCycle
	log.Println("invest:", investAmount, "get:", value)
	newValue := calRate(times, rate)
	log.Println("rate:", newValue)
	
}

func calRate(times int, rate float64) float64 {
	value := 1.0
	for i := 0; i < times; i++ {
		value = value * (1 + rate)
	}
	return value
}

func calCompound(rate float64, times int, amountPerCycle float64) float64 {
	value := 0.001
	//r := rand.New(rand.NewSource(99))
	for i := 0; i < times; i++ {
		//rateR := r.Float64() - 0.5 + rate
		//log.Println("rand:", rateR)
		value = (value + amountPerCycle) * (1 + rate)
	}
	return value
}