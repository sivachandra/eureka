// #############################################################################
// This file is part of the "clap" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

// The clap package provides Command Line Argument Parsing facilities.
package clap

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Arg string

type NamedArg struct {
	name string
	short string
	help string
	defValStr string
	dest interface{}
	required bool
	set bool
}

func (namedArg *NamedArg) Reset() error {
	namedArg.set = false
	if !namedArg.required {
		var valid bool
		var err error
		switch namedArg.dest.(type) {
		case *int:
			var ptr *int
			ptr, valid = namedArg.dest.(*int)
			if valid {
				var int64Val int64
				int64Val, err = strconv.ParseInt(namedArg.defValStr, 0, 0)
				if err == nil {
					*ptr = int(int64Val)
				}
			}
		case *uint:
			var ptr *uint
			ptr, valid = namedArg.dest.(*uint)
			if valid {
				var uint64Val uint64
				uint64Val, err = strconv.ParseUint(namedArg.defValStr, 0, 0)
				if err == nil {
					*ptr = uint(uint64Val)
				}
			}
		case *int64:
			var ptr *int64
			ptr, valid = namedArg.dest.(*int64)
			if valid {
				*ptr, err = strconv.ParseInt(namedArg.defValStr, 0, 64)
			}
		case *uint64:
			var ptr *uint64
			ptr, valid = namedArg.dest.(*uint64)
			if valid {
				*ptr, err = strconv.ParseUint(namedArg.defValStr, 0, 64)
			}
		case *float64:
			var ptr *float64
			ptr, valid = namedArg.dest.(*float64)
			if valid {
				*ptr, err = strconv.ParseFloat(namedArg.defValStr, 64)
			}
		case *bool:
			var ptr *bool
			ptr, valid = namedArg.dest.(*bool)
			if valid {
				*ptr, err = strconv.ParseBool(namedArg.defValStr)
			}
		case *string:
			var ptr *string
			ptr, valid = namedArg.dest.(*string)
			if valid {
				*ptr = namedArg.defValStr
			}
		default:
			err := fmt.Errorf(
				"Unexpected argument type while resetting named arg '%s'.",
				namedArg.name)
			return err
		}

		if !valid {
			err = fmt.Errorf(
				"Unable to cast to argument type for arg '%s'.", namedArg.name)
			return err
		}
		if err != nil {
			err = fmt.Errorf(
				"Error while resetting named arg '%s' to default value.\n%s",
				namedArg.name,
				err.Error())
			return err
		}
	}

	return nil
}

func newNamedArg(name, short, help, defValStr string, dest interface{}, required bool) *NamedArg {
	arg := new(NamedArg)
	arg.name = name
	arg.short = short
	arg.help = help
	arg.defValStr = defValStr
	arg.dest = dest
	arg.required = required
	arg.set = false

	return arg
}

type ArgSet struct {
	// Command name
	name string

	// Sub-commands
	subCommands map[string]*ArgSet

	// Command description
	description string

	// Mapping from arg names to args.
	// This is used during parsing.
	namedArgMap map[string]*NamedArg

	// List of all named args.
	namedArgList []*NamedArg

	// List of unnamed arguments to the command.
	// This is populated while parsing.
	argList []Arg

	// Indicates whether -h or --help was specified during parsing.
	shouldRenderHelp bool

	// Indicates whether the Parse method was called and that it was
	// successfull.
	parsed bool
}

// NewArgSet creates a new argument set for a command given by |name|. The
// description of the command (which is printed along when the command is
// executed with '-h' or '--help' options) should be specified in
// |description|.
func NewArgSet(name string, description string) *ArgSet {
	argSet := new(ArgSet)
	argSet.name = name
	argSet.description = description
	argSet.shouldRenderHelp = false
	argSet.parsed = false
	argSet.namedArgMap = make(map[string]*NamedArg)
	argSet.subCommands = make(map[string]*ArgSet)

	argSet.AddBoolArg(
		"help", "h", &argSet.shouldRenderHelp, argSet.shouldRenderHelp,
		false, fmt.Sprintf("Print '%s' usage information.", name))

	return argSet
}

func (argSet *ArgSet) Name() string {
	return argSet.name
}

func (argSet *ArgSet) AddSubCommand(subArgSet *ArgSet) error {
	subCommandName := subArgSet.Name()
	_, exists := argSet.subCommands[subCommandName]
	if exists {
		return fmt.Errorf(
			"Sub-command with name '%s' already registered with '%s'.",
			subCommandName, argSet.name)
	}

	argSet.subCommands[subCommandName] = subArgSet
	return nil
}

func (argSet *ArgSet) addNamedArg(
	name, short, help, defValStr string, dest interface{}, required bool) {
	arg := newNamedArg(name, short, help, defValStr, dest, required)
	argSet.namedArgList = append(argSet.namedArgList, arg)
	argSet.namedArgMap[name] = arg
	argSet.namedArgMap[short] = arg
}

func (argSet *ArgSet) AddIntArg(
	name string, short string, dest *int, def int, required bool, help string) {
	argSet.addNamedArg(name, short, help, fmt.Sprintf("%d", def), dest, required)
	*dest = def
}

func (argSet *ArgSet) AddInt64Arg(
	name string, short string, dest *int64, def int64, required bool, help string) {
	argSet.addNamedArg(name, short, help, fmt.Sprintf("%d", def), dest, required)
	*dest = def
}

func (argSet *ArgSet) AddUIntArg(
	name string, short string, dest *uint, def uint, required bool, help string) {
	argSet.addNamedArg(name, short, help, fmt.Sprintf("%d", def), dest, required)
	*dest = def
}

func (argSet *ArgSet) AddUInt64Arg(
	name string, short string, dest *uint64, def uint64, required bool, help string) {
	argSet.addNamedArg(name, short, help, fmt.Sprintf("%d", def), dest, required)
	*dest = def
}

func (argSet *ArgSet) AddFloat64Arg(
	name string, short string, dest *float64, def float64, required bool, help string) {
	argSet.addNamedArg(name, short, help, fmt.Sprintf("%f", def), dest, required)
	*dest = def
}

func (argSet *ArgSet) AddBoolArg(
	name string, short string, dest *bool, def bool, required bool, help string) {
	argSet.addNamedArg(name, short, help, fmt.Sprintf("%t", def), dest, required)
	*dest = def
}

func (argSet *ArgSet) AddStringArg(
	name string, short string, dest *string, def string, required bool, help string) {
	argSet.addNamedArg(name, short, help, fmt.Sprintf("%s", def), dest, required)
	*dest = def
}

func (argSet *ArgSet) Parse(arguments []string) ([]string, error) {
	processedCmds := []string{argSet.name}

	if len(arguments) > 0 {
		subCommand, exists := argSet.subCommands[arguments[0]]
		if exists {
			subCommandList, err := subCommand.Parse(arguments[1:])
			return append(subCommandList, argSet.name), err
		}
	}

	for i := 0; i < len(arguments); i++ {
		argument := arguments[i]
		if strings.HasPrefix(argument, "-") {
			// A named argument can be specified in the following ways:
			//     -name value
			//     --name value
			//     -name=value
			//     --name=value
			// If it were a bool value argument, the value can be omitted to
			// imply a value of 'true':
			//     -name
			//     --name

			stripped := argument[1:]
			if strings.HasPrefix(stripped, "-") {
				stripped = stripped[1:]
			}

			var arg *NamedArg
			var valStr string

			indexOfEqual := strings.Index(stripped, "=")
			if indexOfEqual < 0 {
				// The stripped argument is the name if there is no "=".
				name := stripped
				var exists bool
				arg, exists = argSet.namedArgMap[name]
				if !exists {
					err := fmt.Errorf("Unknown argument '%s'.", name)
					return processedCmds, err
				}

				// If the argument is of bool type, then the next argument
				// can be a string which can be parsed error free by
				// strconv.ParseBool, or can be unspecified to mean 'true'.
				i += 1
				switch arg.dest.(type)  {
				default:
					valStr = arguments[i]
				case *bool:
					nextArgStr := arguments[i]
					_, err := strconv.ParseBool(nextArgStr)
					if err == nil {
						valStr = nextArgStr
					} else {
						i -= 1
						valStr = "true"
					}
				}
			} else if indexOfEqual == 0 {
				// This is an error
				err := fmt.Errorf(
					"Probably missing an argument name in '%s'.", argument)
				return processedCmds, err
			} else {
				name := stripped[0:indexOfEqual]
				valStr = stripped[indexOfEqual + 1:]
				var exists bool
				arg, exists = argSet.namedArgMap[name]
				if !exists {
					err := fmt.Errorf("Unknown argument '%s'.", name)
					return processedCmds, err
				}
			}

			var err error
			var valid bool
			switch arg.dest.(type) {
			case *int:
				var ptr *int
				ptr, valid = arg.dest.(*int)
				if valid {
					var int64Val int64
					int64Val, err = strconv.ParseInt(valStr, 0, 0)
					if err == nil {
						*ptr = int(int64Val)
					}
				}
			case *uint:
				var ptr *uint
				ptr, valid = arg.dest.(*uint)
				if valid {
					var uint64Val uint64
					uint64Val, err = strconv.ParseUint(valStr, 0, 0)
					if err == nil {
						*ptr = uint(uint64Val)
					}
				}
			case *int64:
				var ptr *int64
				ptr, valid = arg.dest.(*int64)
				if valid {
					*ptr, err = strconv.ParseInt(valStr, 0, 64)
				}
			case *uint64:
				var ptr *uint64
				ptr, valid = arg.dest.(*uint64)
				if valid {
					*ptr, err = strconv.ParseUint(valStr, 0, 64)
				}
			case *float64:
				var ptr *float64
				ptr, valid = arg.dest.(*float64)
				if valid {
					*ptr, err = strconv.ParseFloat(valStr, 64)
				}
			case *bool:
				var ptr *bool
				ptr, valid = arg.dest.(*bool)
				if valid {
					*ptr, err = strconv.ParseBool(valStr)
				}
			case *string:
				var ptr *string
				ptr, valid = arg.dest.(*string)
				if valid {
					*ptr = valStr
				}
			default:
				err := fmt.Errorf("Unexpected argument type while parsing.")
				return processedCmds, err
			}

			if !valid {
				err := fmt.Errorf("Unable to perform type assertion while parsing.")
				return processedCmds, err
			}
			if err != nil {
				err := fmt.Errorf(
					"Error parsing value of argument '%s'.\n%s", err.Error())
				return processedCmds, err
			}

			if arg.required {
				arg.set = true
			}
		} else {
			// This is not a named argument.
			argSet.argList = append(argSet.argList, Arg(argument))
		}
	}

	for _, arg := range argSet.namedArgList {
		if arg.required && !arg.set {
			err := fmt.Errorf("Required argument '%s' not specified.", arg.name)
			return processedCmds, err
		}
	}

	if argSet.shouldRenderHelp {
		argSet.RenderHelp()
		os.Exit(0)
	}

	return processedCmds, nil
}

func (argSet *ArgSet) Args() []Arg {
	return argSet.argList
}

func (argSet *ArgSet) Clear() error {
	argSet.argList = nil

	for _, namedArg := range argSet.namedArgList {
		err := namedArg.Reset()
		if err != nil {
			err = fmt.Errorf(
				"Unable to clear arg set '%s'.\n%s'", argSet.name, err.Error())
			return err
		}
	}

	for _, subCmd := range argSet.subCommands {
		err := subCmd.Clear()
		if err != nil {
			err = fmt.Errorf(
				"Unable to clear arg set '%s'. Error clearing sub cmd '%s'.\n%s",
				err.Error())
			return err
		}
	}

	return nil
}

func (argSet *ArgSet) ShouldRenderHelp() bool {
	return argSet.shouldRenderHelp
}

func (argSet *ArgSet) RenderHelp() {
	fmt.Printf("%s\n\n", argSet.description)
	fmt.Printf("Options:\n")
	for _, arg := range argSet.namedArgList {
		fmt.Printf("  -%s,  --%s\n", arg.short, arg.name)
		if arg.required {
			fmt.Printf("     Required argument.\n")
		} else {
			fmt.Printf("     Default value: %s\n", arg.defValStr)
		}
		usage := strings.Replace(arg.help, "\n", "\n     ", -1)
		fmt.Printf("     %s\n", usage)
	}
}
