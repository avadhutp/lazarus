package main

import (
	"fmt"
	"github.com/avadhutp/lazarus/geddit"
)

func main() {
	lst := geddit.Get()

	fmt.Println(lst)
}
