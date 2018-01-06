// Copyright 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/logic/cmdr"
)

type cmdHello struct {
	fs       *flag.FlagSet
	Greeting string
}

func (cmd *cmdHello) FlagSet() *flag.FlagSet {
	return cmd.fs
}

func (cmd *cmdHello) Help() string {
	return "Is it subcommands you're looking for?"
}

func (cmd *cmdHello) PositionalArguments() []cmdr.Argument {
	return []cmdr.Argument{
		cmdr.Argument{
			Name:        "target",
			Description: "who do we say it to",
			DefValue:    "world",
			Optional:    true,
		},
	}
}

func (cmd *cmdHello) Run(args []string) error {
	target := "world"
	if len(args) > 0 {
		target = strings.Join(args, " ")
	}
	fmt.Printf("%s, %s\n", cmd.Greeting, target)
	return nil
}

func init() {
	cmd := &cmdHello{
		fs: flag.NewFlagSet("hello", flag.ExitOnError),
	}
	cmd.fs.StringVar(&cmd.Greeting, "greeting", "hello",
		"Greeting to use")
	cmdr.Commands["hello"] = cmd
}
