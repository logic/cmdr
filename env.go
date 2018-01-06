package cmdr

import "os"

// EnvVar defines a standard interface for triggering actions based on a
// defined environment variable, as well as associating a help text with
// it.
type EnvVar interface {
	// Trigger is called with the value of the environment variable,
	// and can return an error to be displayed to the user if something
	// goes wrong.
	Trigger(string) error

	// Help returns a one-line description of what this environment
	// variable does.
	Help() string
}

// EnvVars is a registry of our environment variable triggers.
var EnvVars = map[string]EnvVar{}

// ParseEnvironment walks the list of registered environment variables and
// calls their trigger functions.
func ParseEnvironment() {
	for name, action := range EnvVars {
		if value, ok := os.LookupEnv(name); ok {
			action.Trigger(value)
		}
	}
}
