// #############################################################################
// This file is part of the "clap" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

package clap

import (
	"fmt"
	"testing"
)

var intArg int
var int64Arg int64
var uintArg uint
var uint64Arg uint64
var boolArg bool
var float64Arg float64
var stringArg string

var defIntArg int

var intSubArg int
var int64SubArg int64

func createTestArgSet() *ArgSet {
	argSet := NewArgSet("command", "A test command.")

	argSet.AddIntArg("int", "i", &intArg, 0, true, "An int argument.")
	argSet.AddIntArg("dint", "d", &defIntArg, 54321, false, "A default int argument.")
	argSet.AddInt64Arg("int64", "l", &int64Arg, 0, true, "An int64 argument.")
	argSet.AddUIntArg("uint", "u", &uintArg, 0, true, "A uint argument.")
	argSet.AddUInt64Arg("uint64", "x", &uint64Arg, 0, true, "A uint64 argument.")
	argSet.AddBoolArg("bool", "b", &boolArg, false, true, "A bool argument.")
	argSet.AddFloat64Arg("float64", "f", &float64Arg, 0, true, "A float64 argument.")
	argSet.AddStringArg("string", "s", &stringArg, "empty", true, "A string argument.")

	return argSet
}

func addSubCmd(argSet *ArgSet) error {
	subCmd := NewArgSet("subcmd", "A test sub-command.")
	subCmd.AddIntArg("int", "i", &intSubArg, 0, true, "An int argument.")
	subCmd.AddInt64Arg("int64", "l", &int64SubArg, 0, true, "An int64 argument.")

	err := argSet.AddSubCommand(subCmd)
	if err !=  nil {
		return fmt.Errorf("Unable to add sub command.\n%s", err.Error())
	}

	return nil
}

func TestArgs(t *testing.T) {
	argSet := createTestArgSet()
	cmdLine := []string{
		"-int", "10", "-int64", "20", "-uint", "30", "-uint64", "40", "-bool",
		"-float64", "1.23", "-string", "hello"}
	cmdList, err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing:\n%s", err.Error())
		return
	}

	if cmdList[0] != "command" {
		t.Errorf("Bad command name. Expecting '%s'; found '%s'", "command", cmdList[0]);
	}

	if intArg != 10 {
		t.Errorf("Argument 'int' has value '%d'; expecting '%d'.", intArg, 10)
	}
	if defIntArg != 54321 {
		t.Errorf("Argument 'dint' has value '%d'; expecting '%d'.", defIntArg, 54321)
	}
	if int64Arg != 20 {
		t.Errorf("Argument 'int64' has value '%d'; expecting '%d'.", int64Arg, 20)
	}
	if uintArg != 30 {
		t.Errorf("Argument 'uint' has value '%d'; expecting '%d'.", uintArg, 30)
	}
	if uint64Arg != 40 {
		t.Errorf("Argument 'uint64' has value '%d'; expecting '%d'.", uint64Arg, 40)
	}
	if float64Arg != 1.23 {
		t.Errorf("Argument 'float64' has value '%f'; expecting '%f'.", float64Arg, 1.23)
	}
	if boolArg != true {
		t.Errorf("Argument 'bool' has value '%t'; expecting '%t'.", boolArg, true)
	}
	if stringArg != "hello" {
		t.Errorf("Argument 'string' has value '%s'; expecting '%s'.", stringArg, "hello")
	}
}

func TestArgsWithEqual(t *testing.T) {
	argSet := createTestArgSet()
	cmdLine := []string{
		"-int=10", "-int64=20", "-uint=30", "-uint64=40", "-bool=true",
		"-float64=1.23", "-string=hello"}
	cmdList, err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing:\n%s", err.Error())
		return
	}

	if cmdList[0] != "command" {
		t.Errorf("Bad command name");
	}

	if intArg != 10 {
		t.Errorf("Argument 'int' has value '%d'; expecting '%d'.", intArg, 10)
	}
	if int64Arg != 20 {
		t.Errorf("Argument 'int64' has value '%d'; expecting '%d'.", int64Arg, 20)
	}
	if uintArg != 30 {
		t.Errorf("Argument 'uint' has value '%d'; expecting '%d'.", uintArg, 30)
	}
	if uint64Arg != 40 {
		t.Errorf("Argument 'uint64' has value '%d'; expecting '%d'.", uint64Arg, 40)
	}
	if float64Arg != 1.23 {
		t.Errorf("Argument 'float64' has value '%f'; expecting '%f'.", float64Arg, 1.23)
	}
	if boolArg != true {
		t.Errorf("Argument 'bool' has value '%t'; expecting '%t'.", boolArg, true)
	}
	if stringArg != "hello" {
		t.Errorf("Argument 'string' has value '%s'; expecting '%s'.", stringArg, "hello")
	}
}

func TestShortArgs(t *testing.T) {
	argSet := createTestArgSet()
	cmdLine := []string{
		"-i", "10", "-l", "20", "-u", "30", "-x", "40", "-b", "-f", "1.23", "-s", "hello"}
	cmdList, err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing:\n%s", err.Error())
		return
	}

	if cmdList[0] != "command" {
		t.Errorf("Bad command name");
	}

	if intArg != 10 {
		t.Errorf("Argument 'int' has value '%d'; expecting '%d'.", intArg, 10)
	}
	if int64Arg != 20 {
		t.Errorf("Argument 'int64' has value '%d'; expecting '%d'.", int64Arg, 20)
	}
	if uintArg != 30 {
		t.Errorf("Argument 'uint' has value '%d'; expecting '%d'.", uintArg, 30)
	}
	if uint64Arg != 40 {
		t.Errorf("Argument 'uint64' has value '%d'; expecting '%d'.", uint64Arg, 40)
	}
	if float64Arg != 1.23 {
		t.Errorf("Argument 'float64' has value '%f'; expecting '%f'.", float64Arg, 1.23)
	}
	if boolArg != true {
		t.Errorf("Argument 'bool' has value '%t'; expecting '%t'.", boolArg, true)
	}
	if stringArg != "hello" {
		t.Errorf("Argument 'string' has value '%s'; expecting '%s'.", stringArg, "hello")
	}
}

func TestShortArgsWithEqual(t *testing.T) {
	argSet := createTestArgSet()
	cmdLine := []string{
		"-i=10", "-l=20", "-u=30", "-x=40", "-b=true", "-f=1.23", "-s=hello"}
	cmdList, err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing:\n%s", err.Error())
		return
	}

	if cmdList[0] != "command" {
		t.Errorf("Bad command name");
	}

	if intArg != 10 {
		t.Errorf("Argument 'int' has value '%d'; expecting '%d'.", intArg, 10)
	}
	if int64Arg != 20 {
		t.Errorf("Argument 'int64' has value '%d'; expecting '%d'.", int64Arg, 20)
	}
	if uintArg != 30 {
		t.Errorf("Argument 'uint' has value '%d'; expecting '%d'.", uintArg, 30)
	}
	if uint64Arg != 40 {
		t.Errorf("Argument 'uint64' has value '%d'; expecting '%d'.", uint64Arg, 40)
	}
	if float64Arg != 1.23 {
		t.Errorf("Argument 'float64' has value '%f'; expecting '%f'.", float64Arg, 1.23)
	}
	if boolArg != true {
		t.Errorf("Argument 'bool' has value '%t'; expecting '%t'.", boolArg, true)
	}
	if stringArg != "hello" {
		t.Errorf("Argument 'string' has value '%s'; expecting '%s'.", stringArg, "hello")
	}
}

func TestSubCommand(t *testing.T) {
	argSet := createTestArgSet()
	err := addSubCmd(argSet)
	if err != nil {
		t.Errorf(err.Error());
	}

	cmdLine := []string{"subcmd", "-i=10", "-l=20"}
	cmdList, err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf(err.Error())
	}

	if cmdList[0] != "subcmd" || cmdList[1] != "command" {
		t.Errorf("Expecting command list [subcmd, command]. Found '%s'", cmdList)
	}

	if intSubArg != 10 {
		t.Errorf("Argument 'int' to subcommand has value '%d'; Expecting 10", intSubArg)
	}

	if int64SubArg != 20 {
		t.Errorf("Argument 'int64' to subcommand has value '%d'; Expecting 20", int64SubArg)
	}
}

func TestCommandClearing(t *testing.T) {
	argSet := createTestArgSet()
	cmdLine := []string{
		"-i=10", "-d=12345", "-l=20", "-u=30", "-x=40", "-b=true", "-f=1.23", "-s=hello"}
	cmdList, err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing first time:\n%s", err.Error())
		return
	}

	err = argSet.Clear()
	if err != nil {
		t.Errorf("Error clearing arg set.\n%s", err.Error())
		return
	}

	if defIntArg != 54321 {
		t.Errorf(
			"Default arg 'dint' not reset after clearing. Got '%d'; expecting '%d'.",
			defIntArg, 54321)
	}

	cmdList, err = argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing second time:\n%s", err.Error())
		return
	}

	if cmdList[0] != "command" {
		t.Errorf("Bad command name");
	}

	if intArg != 10 {
		t.Errorf("Argument 'int' has value '%d'; expecting '%d'.", intArg, 10)
	}
	if defIntArg != 12345 {
		t.Errorf("Argument 'int' has value '%d'; expecting '%d'.", defIntArg, 12345)
	}
	if int64Arg != 20 {
		t.Errorf("Argument 'int64' has value '%d'; expecting '%d'.", int64Arg, 20)
	}
	if uintArg != 30 {
		t.Errorf("Argument 'uint' has value '%d'; expecting '%d'.", uintArg, 30)
	}
	if uint64Arg != 40 {
		t.Errorf("Argument 'uint64' has value '%d'; expecting '%d'.", uint64Arg, 40)
	}
	if float64Arg != 1.23 {
		t.Errorf("Argument 'float64' has value '%f'; expecting '%f'.", float64Arg, 1.23)
	}
	if boolArg != true {
		t.Errorf("Argument 'bool' has value '%t'; expecting '%t'.", boolArg, true)
	}
	if stringArg != "hello" {
		t.Errorf("Argument 'string' has value '%s'; expecting '%s'.", stringArg, "hello")
	}
}
