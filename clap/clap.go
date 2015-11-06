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

type Cmd struct {
	// Command name
	name string

	// Sub-commands
	subCmds map[string]*Cmd

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

// NewCmd creates a new command with name |name|.
// The description of the command (which is printed when the command is
// executed with '-h' or '--help' options) should be specified in
// |description|.
func NewCmd(name string, description string) *Cmd {
	cmd := new(Cmd)
	cmd.name = name
	cmd.description = description
	cmd.shouldRenderHelp = false
	cmd.parsed = false
	cmd.namedArgMap = make(map[string]*NamedArg)
	cmd.subCmds = make(map[string]*Cmd)

	cmd.AddBoolArg(
		"help", "h", &cmd.shouldRenderHelp, cmd.shouldRenderHelp,
		false, fmt.Sprintf("Print '%s' usage information.", name))

	return cmd
}

func (cmd *Cmd) Name() string {
	return cmd.name
}

func (cmd *Cmd) Description() string {
	return cmd.description
}

func (cmd *Cmd) AddSubCmd(subCmd *Cmd) error {
	subCmdName := subCmd.Name()
	_, exists := cmd.subCmds[subCmdName]
	if exists {
		return fmt.Errorf(
			"Sub-command with name '%s' already registered with '%s'.",
			subCmdName, cmd.name)
	}

	cmd.subCmds[subCmdName] = subCmd
	return nil
}

func (cmd *Cmd) addNamedArg(
	name, short, help, defValStr string, dest interface{}, required bool) {
	arg := newNamedArg(name, short, help, defValStr, dest, required)
	cmd.namedArgList = append(cmd.namedArgList, arg)
	cmd.namedArgMap[name] = arg
	cmd.namedArgMap[short] = arg
}

func (cmd *Cmd) AddIntArg(
	name string, short string, dest *int, def int, required bool, help string) {
	cmd.addNamedArg(name, short, help, fmt.Sprintf("%d", def), dest, required)
	*dest = def
}

func (cmd *Cmd) AddInt64Arg(
	name string, short string, dest *int64, def int64, required bool, help string) {
	cmd.addNamedArg(name, short, help, fmt.Sprintf("%d", def), dest, required)
	*dest = def
}

func (cmd *Cmd) AddUIntArg(
	name string, short string, dest *uint, def uint, required bool, help string) {
	cmd.addNamedArg(name, short, help, fmt.Sprintf("%d", def), dest, required)
	*dest = def
}

func (cmd *Cmd) AddUInt64Arg(
	name string, short string, dest *uint64, def uint64, required bool, help string) {
	cmd.addNamedArg(name, short, help, fmt.Sprintf("%d", def), dest, required)
	*dest = def
}

func (cmd *Cmd) AddFloat64Arg(
	name string, short string, dest *float64, def float64, required bool, help string) {
	cmd.addNamedArg(name, short, help, fmt.Sprintf("%f", def), dest, required)
	*dest = def
}

func (cmd *Cmd) AddBoolArg(
	name string, short string, dest *bool, def bool, required bool, help string) {
	cmd.addNamedArg(name, short, help, fmt.Sprintf("%t", def), dest, required)
	*dest = def
}

func (cmd *Cmd) AddStringArg(
	name string, short string, dest *string, def string, required bool, help string) {
	cmd.addNamedArg(name, short, help, fmt.Sprintf("%s", def), dest, required)
	*dest = def
}

func (cmd *Cmd) Parse(arguments []string) ([]string, error) {
	processedCmds := []string{cmd.name}

	if len(arguments) > 0 {
		subCmd, exists := cmd.subCmds[arguments[0]]
		if exists {
			subCmdList, err := subCmd.Parse(arguments[1:])
			return append(processedCmds, subCmdList...), err
		}
	}

	argCount := len(arguments)
	for i := 0; i < argCount; i++ {
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
				arg, exists = cmd.namedArgMap[name]
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
					if i >= argCount {
						err := fmt.Errorf(
							"Missing value for argument '%s'.", name)
						return processedCmds, err
					}
					valStr = arguments[i]
				case *bool:
					if i >= argCount {
						i -= 1;
						valStr = "true"
					} else {
						nextArgStr := arguments[i]
						_, err := strconv.ParseBool(nextArgStr)
						if err == nil {
							valStr = nextArgStr
						} else {
							i -= 1
							valStr = "true"
						}
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
				arg, exists = cmd.namedArgMap[name]
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
			cmd.argList = append(cmd.argList, Arg(argument))
		}
	}

	if !cmd.shouldRenderHelp {
		for _, arg := range cmd.namedArgList {
			if arg.required && !arg.set {
				err := fmt.Errorf("Required argument '%s' not specified.", arg.name)
				return processedCmds, err
			}
		}
	}

	return processedCmds, nil
}

func (cmd *Cmd) Args() []Arg {
	return cmd.argList
}

func (cmd *Cmd) Clear() error {
	cmd.argList = nil

	for _, namedArg := range cmd.namedArgList {
		err := namedArg.Reset()
		if err != nil {
			err = fmt.Errorf(
				"Unable to clear command '%s'.\n%s'", cmd.name, err.Error())
			return err
		}
	}

	for _, subCmd := range cmd.subCmds {
		err := subCmd.Clear()
		if err != nil {
			err = fmt.Errorf(
				"Unable to clear sub-command '%s' of command '%s'.\n%s",
				cmd.name, subCmd.name, err.Error())
			return err
		}
	}

	return nil
}

func (cmd *Cmd) ShouldRenderHelp() bool {
	return cmd.shouldRenderHelp
}

func (cmd *Cmd) RenderHelp() {
	fmt.Printf("%s\n\n", cmd.description)

	if len(cmd.subCmds) > 0 {
		fmt.Printf("Sub-commands:\n")
		for _, subCmd := range cmd.subCmds {
			fmt.Printf("     %s\n", subCmd.name)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("Options:\n")
	for _, arg := range cmd.namedArgList {
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
