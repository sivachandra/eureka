// #############################################################################
// This file is part of the "clap" package of the "Eureka" project.
// It is distributed under the MIT License. Refer to the LICENSE file for more
// information.
//
// Website: http://www.github.com/sivachandra/eureka
// #############################################################################

// The clap package provides Command Line Argument Parsing facilities.
// It is a built as a facade over the standard library package 'flag'.
package clap

import (
	"testing"
)

var intArg int
var int64Arg int64
var uintArg uint
var uint64Arg uint64
var boolArg bool
var float64Arg float64
var stringArg string

func createTestArgSet() *ArgSet {
	argSet := NewArgSet("command", "A test command.")

	argSet.AddIntArg("int", "i", &intArg, 0, true, "An int argument.")
	argSet.AddInt64Arg("int64", "l", &int64Arg, 0, true, "An int64 argument.")
	argSet.AddUIntArg("uint", "u", &uintArg, 0, true, "A uint argument.")
	argSet.AddUInt64Arg("uint64", "x", &uint64Arg, 0, true, "A uint64 argument.")
	argSet.AddBoolArg("bool", "b", &boolArg, false, true, "A bool argument.")
	argSet.AddFloat64Arg("float64", "f", &float64Arg, 0, true, "A float64 argument.")
	argSet.AddStringArg("string", "s", &stringArg, "empty", true, "A string argument.")

	return argSet
}

func TestArgs(t *testing.T) {
	argSet := createTestArgSet()
	cmdLine := []string{
		"-int", "10", "-int64", "20", "-uint", "30", "-uint64", "40", "-bool",
		"-float64", "1.23", "-string", "hello"}
	err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing:\n%s", err.Error())
		return
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

func TestArgsWithEqual(t *testing.T) {
	argSet := createTestArgSet()
	cmdLine := []string{
		"-int=10", "-int64=20", "-uint=30", "-uint64=40", "-bool=true",
		"-float64=1.23", "-string=hello"}
	err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing:\n%s", err.Error())
		return
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
	err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing:\n%s", err.Error())
		return
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
	err := argSet.Parse(cmdLine)
	if err != nil {
		t.Errorf("Error while parsing:\n%s", err.Error())
		return
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
