package main

import "fmt"

func main() {
	testBasicPanic()                       // recovers
	testGoRoutinePanic()                   // recovers
	testDeepGoRoutinePanic()               // panics
	testAsyncGoRoutinePanic()              // panics
	testAsyncGoRoutinePanicWithErrorChan() // recovers
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

func testAsyncGoRoutinePanicWithErrorChan() {
	defer handlePanic()

	done, errStream := func() (<-chan struct{}, <-chan error) {
		doneSignal := make(chan struct{})
		errStream := make(chan error)

		go func() {
			defer close(doneSignal)
			defer func() {
				if r := recover(); r != nil {
					errStream <- fmt.Errorf("%v", r)
				}
			}()
			panic("testAsyncGoRoutinePanicWithErrorChan panic")
		}()

		return doneSignal, errStream
	}()

	select {
	case <-done:
	case err := <-errStream:
		panic(err)
	}
}
