package cmdr

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

// Command defines a standard interface for retrieving a flagset for a
// subcommand, for running the command (if it matches), and for retrieving
// help text.
type Command interface {
	// FlagSet returns a fully-populated flag set for the subcommand.
	FlagSet() *flag.FlagSet

	// Run is called with the remaining list of arguments after parsing
	// the flag set, and performs the action tied to the command. An
	// error can be returned if something goes wrong, which will be
	// presented to the user.
	Run([]string) error

	// Help returns a one-line description of what this command does.
	Help() string
}

// CommandGlobal is our global flagset.
var CommandGlobal = flag.NewFlagSet("_global", flag.ExitOnError)

// SubCommands are all defined subcommands and their flagsets.
var SubCommands = map[string]Command{}

// GlobalHelp displays either a partial or full help text for our command
// and subcommands.
func GlobalHelp(full bool) {
	flag.Usage()

	names := make([]string, len(SubCommands))
	i := 0
	for k := range SubCommands {
		names[i] = k
		i++
	}
	sort.Strings(names)

	CommandGlobal.PrintDefaults()

	fmt.Println("\nSubcommands:")
	for _, name := range names {
		if full {
			fmt.Printf("\n%s - %s\n", name, SubCommands[name].Help())
			SubCommands[name].FlagSet().PrintDefaults()
		} else {
			fmt.Printf("  %-10s %s\n", name, SubCommands[name].Help())
		}
	}

	if len(EnvVars) != 0 {
		fmt.Println("\nValid environment variables:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		for name, action := range EnvVars {
			fmt.Fprintf(w, "  %s\t%s\t\n", name, action.Help())
		}
		w.Flush()
	}
	os.Exit(1)
}

// ParseCommand takes a list of command-line arguments (typically os.Args),
// parses the global arguments, then checks to see if there is a subcommand
// to execute.
func ParseCommand(args []string) {
	var shortHelp bool
	CommandGlobal.BoolVar(&shortHelp, "help", false,
		"Print all subcommands")
	var longHelp bool
	CommandGlobal.BoolVar(&longHelp, "long-help", false,
		"Print full help for all subcommands")

	if err := CommandGlobal.Parse(args[1:]); err != nil {
		panic(err)
	}
	args = CommandGlobal.Args()

	if longHelp {
		GlobalHelp(true)
	}
	if shortHelp || len(args) < 1 {
		GlobalHelp(false)
	}

	if fs, ok := SubCommands[args[0]]; ok {
		fs.FlagSet().Parse(args[1:])
		if err := fs.Run(fs.FlagSet().Args()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("No such subcommand '%s'.\n", args[0])
		GlobalHelp(false)
	}
}
