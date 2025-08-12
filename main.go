package main

import "fmt"

func main() {
	testBasicPanic()
	testGoRoutinePanic()
	testDeepGoRoutinePanic()
	testAsyncGoRoutinePanic()
}

func handlePanic() {
	if r := recover(); r != nil {
		err := fmt.Errorf("recovered from panic: %v", r)
		println(err.Error())
	}
}

func testBasicPanic() {
	defer handlePanic()
	panic("testBasicPanic panic")
}

func testGoRoutinePanic() {
	doneSignal := make(chan struct{})

	go func() {
		defer close(doneSignal)
		defer handlePanic()
		panic("testGoRoutinePanic panic")
	}()

	<-doneSignal
}

func testDeepGoRoutinePanic() {
	doneSignal := make(chan struct{})

	defer handlePanic()
	go func() {
		defer close(doneSignal)
		panic("testDeepGoRoutinePanic panic")
	}()

	<-doneSignal
}

func testAsyncGoRoutinePanic() {
	defer handlePanic()

	done := func() <-chan struct{} {
		doneSignal := make(chan struct{})

		go func() {
			defer close(doneSignal)
			panic("testAsyncGoRoutinePanic panic")
		}()

		return doneSignal
	}()

	<-done
}
