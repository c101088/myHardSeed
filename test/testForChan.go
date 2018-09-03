package  main

import (
	"fmt"
)

func main(){
	chan1 := make(chan int)
	chan2 := make(chan int)
	go fun1(chan1)


	for i := 0;i <10;i++{
		createWorker(i,chan1,chan2)
	}
	for{
		<-chan2
	}
//	time.Sleep(1*time.Second)
}

func fun1(chan1 chan int){
	for i := 0;i<100;i++{
		chan1<- i
	}
}

func fun2(x int){
	fmt.Printf("I am func2 %v\n",x)
}

func createWorker(i int,chan1 chan int,chan2 chan int )  {
	go func() {
		for{
			x := <-chan1
			fmt.Printf("This is worker %v is working %v\n",i,x)
			chan2<-1
		}
	}()
}