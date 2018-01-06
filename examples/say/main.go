// Copyright 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/logic/cmdr"
)

func main() {
	var timing bool
	cmdr.Global.BoolVar(&timing, "timing", false, "Display timing data")
	cmd := cmdr.Parse(os.Args)

	before := time.Now()
	if timing {
		fmt.Printf("Starting at %s\n", before)
	}

	err := cmd.Run()

	if timing {
		after := time.Now()
		fmt.Printf("Finished at %s\n", after)
		fmt.Printf("Took %s\n", after.Sub(before))
	}

	if err != nil {
		panic(err)
	}
}
