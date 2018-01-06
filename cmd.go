package cmdr

import (
	"flag"
	"fmt"
	"os"
	"sort"
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

	// PositionalArguments returns an ordered array of any positional
	// arguments that the command requires.
	PositionalArguments() []Argument
}

// Argument defines a single positional argument
type Argument struct {
	Name        string
	Description string
	DefValue    string
	Optional    bool
}

func (arg *Argument) String() string {
	if arg.Optional {
		return fmt.Sprintf("[%s]", arg.Name)
	}
	return fmt.Sprintf("%s", arg.Name)
}

// Global is our global flagset.
var Global = flag.NewFlagSet("_global", flag.ExitOnError)

// Commands are all defined subcommands and their flagsets.
var Commands = map[string]Command{}

// Help displays either a partial or full help text for our command
// and all subcommands.
func Help(full bool) {
	flag.Usage()

	names := make([]string, len(Commands))
	i := 0
	for k := range Commands {
		names[i] = k
		i++
	}
	sort.Strings(names)

	Global.PrintDefaults()

	fmt.Println("\nSubcommands:")
	for _, name := range names {
		pArgs := Commands[name].PositionalArguments()
		if full {
			fmt.Printf("\n%s - %s\n", name, Commands[name].Help())
			Commands[name].FlagSet().PrintDefaults()
			if pArgs != nil {
				for _, arg := range pArgs {
					out := "  " + arg.String()
					if len(out) < 4 {
						out += "\t"
					} else {
						out += "\n    \t"
					}
					out += arg.Description
					if !arg.Optional || arg.DefValue != "" {
						out += " ("
						if arg.DefValue != "" {
							out += "default \"" + arg.DefValue + "\""
						}
						if !arg.Optional {
							if arg.DefValue != "" {
								out += ", "
							}
							out += "required"
						}
						out += ")"
					}
					fmt.Println(out)
				}
			}
		} else {
			out := "  " + name
			Commands[name].FlagSet().VisitAll(func(f *flag.Flag) {
				out += " [-" + f.Name
				name, _ := flag.UnquoteUsage(f)
				if name != "" {
					out += " " + name
				}
				out += "]"
			})
			if pArgs != nil {
				for _, arg := range pArgs {
					out += " " + arg.String()
				}
			}
			out += "\n    \t" + Commands[name].Help()
			fmt.Println(out)
		}
	}

	if len(Variables) > 0 {
		out := "\nEnvironment variables:"
		if full {
			out += "\n"
		}
		for name, action := range Variables {
			out += "\n  " + name + "\n    \t" + action.Help()
		}
		fmt.Println(out)
	}
	os.Exit(1)
}

// Parse takes a list of command-line arguments (typically os.Args), parses the
// global arguments, then checks to see if there is a subcommand to execute.
func Parse(args []string) {
	ParseEnvironment()

	var shortHelp bool
	Global.BoolVar(&shortHelp, "help", false,
		"Print all subcommands")

	var longHelp bool
	Global.BoolVar(&longHelp, "long-help", false,
		"Print full help for all subcommands")

	if err := Global.Parse(args[1:]); err != nil {
		panic(err)
	}
	args = Global.Args()

	if longHelp {
		Help(true)
	}
	if shortHelp || len(args) < 1 {
		Help(false)
	}

	if fs, ok := Commands[args[0]]; ok {
		fs.FlagSet().Parse(args[1:])
		if err := fs.Run(fs.FlagSet().Args()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("No such subcommand '%s'.\n", args[0])
		Help(false)
	}
}
