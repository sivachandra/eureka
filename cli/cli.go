// #############################################################################
// This file is part of the "cli" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

// The cli package provides facilities for building cmd line interfaces.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

import (
	"eureka/clap"
)

type ResponseType int

const (
	RspEndSession = ResponseType(0)
	RspPrintMsg = ResponseType(1)
	RspGetUsrInput = ResponseType(2)
)

type CmdResponse struct {
	Type ResponseType
	Msg string
}

type CmdHandler interface {
	Run(input chan string, output chan *CmdResponse) error
}

type CLI struct {
	banner string
	prompt string
	cmds map[string]*clap.Cmd
	cmdHandlers map[string]CmdHandler
}

func NewCLI(banner string, prompt string) *CLI {
	cli := new(CLI)
	cli.banner = banner
	cli.prompt = prompt
	cli.cmds = make(map[string]*clap.Cmd)
	cli.cmdHandlers = make(map[string]CmdHandler)

	cli.addQuitCmd()
	cli.addHelpCmd()

	return cli
}

func (cli *CLI) AddCmd(cmd *clap.Cmd, handler CmdHandler) error {
	name := cmd.Name()
	_, exists := cli.cmds[name]
	if exists {
		return fmt.Errorf("Command with name '%s' already registered with CLI.", name)
	}

	cli.cmds[name] = cmd
	cli.cmdHandlers[name] = handler

	return nil
}

func (cli *CLI) MainLoop() {
	fmt.Print(cli.banner)
	fmt.Println("")

	var endSession bool = false
	for !endSession {
		fmt.Printf("%s ", cli.prompt)
		cmdName, args, err := cli.readCmd()
		if err != nil {
			fmt.Printf("%s", err.Error())
			continue
		}
		if len(cmdName) == 0 {
			continue
		}

		cmd, exists := cli.cmds[cmdName]
		if !exists {
			fmt.Printf("Unknown command '%s'.\n", cmdName)
			continue
		}

		if args != nil {
			_, err = cmd.Parse(args)
			if err != nil {
				fmt.Printf(
					"Error parsing arguments to command '%s'.\n%s",
					cmd.Name(), err.Error())
				continue
			}
		}
		if cmd.ShouldRenderHelp() {
			cmd.RenderHelp()
			cmd.Clear()
			continue
		}

		handler, exists := cli.cmdHandlers[cmdName]
		if !exists {
			fmt.Printf("Handler for command '%s' not found.\n", cmdName)
			continue
		}

		cmdInput := make(chan string)
		cmdOutput := make(chan *CmdResponse)
		go handler.Run(cmdInput, cmdOutput)

		for true {
			response, done := <-cmdOutput
			if response != nil {
				switch response.Type {
				case RspGetUsrInput:
					fmt.Printf("%s\n", response.Msg)
					input := cli.readInput()
					cmdInput <- input
				case RspPrintMsg:
					fmt.Printf("%s\n", response.Msg)
				case RspEndSession:
					endSession = true
				default:
					fmt.Printf("Unexpected response from command handler.\n")
				}
			}

			if done || endSession {
				break
			}
		}

		cmd.Clear()
	}
}

func (cli *CLI) readInput() string {
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	return stdin.Text()
}

func parseCmdStr(input string) ([]string, error) {
	cmdStr := strings.TrimSpace(input)
	var current []byte
	var args []string
	var parsingQuotedArg bool = false // Indicates if we in the midst of parsing a quoted arg.
	var doubleBackSlash bool = false // Indicates if a back slach was parsed immediately before.
	for i, char := range cmdStr {
		if parsingQuotedArg {
			if char == '"' {
				if cmdStr[i - 1] == '\\' && !doubleBackSlash {
					current[len(current) - 1] = byte(char)
				} else {
					args = append(args, string(current))
					current = nil
					parsingQuotedArg = false
					doubleBackSlash = false // It could arleady be false.
				}
			} else if char == '\\' {
				if cmdStr[i - 1] == '\\' && !doubleBackSlash {
					current[len(current) - 1] = byte(char)
					doubleBackSlash = true
				} else {
					current = append(current, byte(char))
					doubleBackSlash = false // It could already be false.
				}
			} else {
				current = append(current, byte(char))
				doubleBackSlash = false // It could already be false.
			}
		} else {
			if char == '"' {
				parsingQuotedArg = true
			} else if unicode.IsSpace(rune(char)) {
				if current != nil {
					args = append(args, string(current))
					current = nil
				}
			} else {
				current = append(current, byte(char))
			}
		}
	}

	if len(current) != 0 {
		args = append(args, string(current))
	}

	if parsingQuotedArg {
		return nil, fmt.Errorf("Invalid command syntax.")
	}

	return args, nil
}

func (cli *CLI) readCmd() (string, []string, error) {
	input := cli.readInput()
	if len(input) == 0 {
		return input, nil, nil
	}

        args, err := parseCmdStr(input)
        if err != nil {
		return "", nil, err
	}

	if len(args) == 0 {
		return "", nil, fmt.Errorf("Incorrectly parsed '%s'.", input)
	}

	if len(args) == 1 {
		return args[0], nil, nil
	} else {
		return args[0], args[1:], nil
	}
}

type quitCmdHandler struct {
}

func (handler *quitCmdHandler) Run(input chan string, output chan *CmdResponse) error {
	rsp := new(CmdResponse)
	rsp.Type = RspEndSession

	output <- rsp
	close(output)
	return nil
}

func (cli *CLI) addQuitCmd() {
	cmd := clap.NewCmd("quit", "End current session.")
	cli.AddCmd(cmd, new(quitCmdHandler))
}

type helpCmdHandler struct {
	cli *CLI
}

func (handler *helpCmdHandler) Run(input chan string, output chan *CmdResponse) error {
	rsp := new(CmdResponse)
	rsp.Msg = "List of available commands:\n\n"
	for name, cmd := range(handler.cli.cmds) {
		desc := cmd.Description()
		desc = strings.Replace(desc, "\n", "\n    ", -1)
		rsp.Msg += name + " -- " + desc + "\n"
	}

	rsp.Type = RspPrintMsg
	output <- rsp
	close(output)
	return nil
}

func (cli *CLI) addHelpCmd() {
	cmd := clap.NewCmd("help", "Show help message.")

	handler := new(helpCmdHandler)
	handler.cli = cli

	cli.AddCmd(cmd, handler)
}
