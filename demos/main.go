package main

import (
	"fmt"
	"github.com/ykallan/memorydb"
)

func main() {
	ms := memorydb.NewWithLock()
	//2407669  lock
	//2118867  unlock
	for i := 0; i < 20; i++ {
		ms.Set(i, 100)
	}
	fmt.Println(ms.GetAll())

}
