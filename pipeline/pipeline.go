package main

import "fmt"

// This is an example of syncronous channels forming a pipeline
// The channels are syncronous because they are not buffered
// Flow will ping-pong between the sender and receiver goroutines
// Each step of the pipeline creates an output channel that will be used by a spun off goroutine
// and then spins off the goroutine to process the input to the output channels

// NOTE: For a synchronous channel you need one goroutine pushing data and one pulling
// Otherwise you will have a deadlock.
// That is why it can't all be done in the main routine.

func sliceToChannel(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func double(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * 2
		}
		close(out)
	}()
	return out
}

func main() {
	nums := []int{1, 3, 5, 7}

	input := sliceToChannel(nums)
	stage1 := double(input)
	stage2 := square(stage1)

	for n := range stage2 {
		fmt.Println(n)
	}

	fmt.Println("Done")
}
