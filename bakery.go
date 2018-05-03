package main

import "fmt"
import "os"
import "math/rand"
import "time"
import "strconv"

type Order struct {
	num int
	sender chan int
}

func main() {

	startTime := time.Now()

	args := os.Args[1:] // get command arguments

	numServers, err1 := strconv.Atoi(args[0])
	numCustomers, err2 := strconv.Atoi(args[1])

	if err1 != nil {
		fmt.Println(err1)
		os.Exit(2)
	}
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(2)
	}

	serverLine := make(chan Order, numCustomers) // create channel for orders from customers
	trackingChannel := make(chan bool, numCustomers) // create channel to track finished customers

	// create servers
	for i := 0; i < numServers; i++ {
		go server(serverLine)
		fmt.Println("Server setup", i)
	}

	// create customers
	for i := 0; i < numCustomers; i++ {
		go customer(serverLine, trackingChannel)
		fmt.Println("Customer setup", i)
	}

	// wait for all customers to be served.
	for i := 0; i < numCustomers; i++ {
		<- trackingChannel
		fmt.Println("Customers served:", i)
		fmt.Println("--------------------------------")
	}

	timeElapsed := time.Since(startTime)

	fmt.Println("Took ", timeElapsed)
}


func fib(n int)(int) {
	switch n {
	case 0, 1:
		return n
	default:
		return fib(n-1) + fib(n-2)
	}
}


func server(c chan Order) {
	for {
		o := <-c // server takes an order
		o.sender <- fib(o.num) // server returns the processed order to customer
	}
}

func customer(c chan Order, t chan bool) {

	sleepTime := time.Duration(rand.Intn(10000)) * time.Millisecond

	time.Sleep(sleepTime) // sleep for a random duration

	resultChan := make(chan int)

	// make orders asking for fib of a random number
	myOrder := Order{sender: resultChan, num: rand.Intn(50)} 

	fmt.Println("Waited for", sleepTime)
	fmt.Println("Ordered", myOrder.num)
	fmt.Println("--------------------------------")

	c <- myOrder // put order in channel for server

	result := <- resultChan // get result from server

	fmt.Println("Fib:", result)

	t <- true // put confirmation in tracking channel

}
