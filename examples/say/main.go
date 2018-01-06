package main

import (
	"fmt"
	"os"
	"time"

	"github.com/logic/cmdr"
)

func main() {
	before := time.Now()
	var timing bool
	cmdr.Global.BoolVar(&timing, "timing", false, "Display timing data")

	cmdr.Parse(os.Args)

	if timing {
		after := time.Now()
		fmt.Printf("Took %s\n", after.Sub(before))
	}
}
