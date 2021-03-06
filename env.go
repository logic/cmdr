// Copyright 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package cmdr

import "os"

// Variable defines a standard interface for triggering actions based on a
// defined environment variable, as well as associating a help text with
// it.
type Variable interface {
	// Trigger is called with the value of the environment variable,
	// and can return an error to be displayed to the user if something
	// goes wrong.
	Trigger(string) error

	// Help returns a one-line description of what this environment
	// variable does.
	Help() string
}

// Variables is a registry of our environment variable triggers.
var Variables = map[string]Variable{}

// ParseEnvironment walks the list of registered environment variables and
// calls their trigger functions.
func ParseEnvironment() {
	for name, action := range Variables {
		if value, ok := os.LookupEnv(name); ok {
			action.Trigger(value)
		}
	}
}
