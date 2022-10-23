package main

import (
	"fmt"
	"github.com/ykallan/memorydb"
)

func main() {
	set := memorydb.NewSet()
	set.Add(1)
	set.Add(1)
	set.Add(1)
	set.Add(1)
	set.Add(1)
	//list := set.ToList()
	fmt.Println(set.ToString())
	set.Print()
}
