package cmd

import (
	"fmt"
	"testing"
)

func TestAddCommandWithEmptyArguments(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand("new task")

	if len(cmd.commands) != 1 {
		t.Fatalf("Expected: %d, Received: %d", 1, len(cmd.commands))
	}
}

func TestAddCommandWithArguments(t *testing.T) {
	cmd := CreateCmd()

	err := cmd.AddCommand("new task",
		CommandArgumentSyntax{
			Name:     "Task name",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_INT,
		},
		CommandArgumentSyntax{
			Name:     "Task state",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_STR,
		},
		CommandArgumentSyntax{
			Name:     "Task owner",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_BOOL,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(cmd.commands) != 1 {
		t.Fatalf("Expected: %d, Received: %d", 1, len(cmd.commands))
	}

	if len(cmd.commands[0].Arguments) != 3 {
		t.Fatalf("Expected: %d, Received: %d", 3, len(cmd.commands[0].Arguments))
	}
}

func TestAddCommandWithInvalidArgumentTypes(t *testing.T) {
	cmd := CreateCmd()

	err := cmd.AddCommand("new task",
		CommandArgumentSyntax{
			Name:     "Task name",
			Required: false,
			Type:     COMMAND_ARGUMENT_TYPE_INT,
		},
		CommandArgumentSyntax{
			Name:     "Task state",
			Required: true,
			Type:     4,
		},
	)

	expectedErrorMessage := "Argument \"Task state\" have a invalid type \"4\""

	if err == nil {
		t.Fatal("Error expected but received nil")
	} else if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected: \"%v\", Received: \"%v\"", expectedErrorMessage, err.Error())
	}

	if len(cmd.commands) != 0 {
		t.Fatalf("Expected: %d, Received: %d", 0, len(cmd.commands))
	}
}

func TestAddCommandWithInvalidOptionalArguments(t *testing.T) {
	cmd := CreateCmd()

	err := cmd.AddCommand("new task",
		CommandArgumentSyntax{
			Name:     "Task name",
			Required: false,
			Type:     COMMAND_ARGUMENT_TYPE_INT,
		},
		CommandArgumentSyntax{
			Name:     "Task state",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_INT,
		},
	)

	expectedErrorMessage := "You cannot have required arguments after optionals"

	if err == nil {
		t.Fatal("Error expected but received nil")
	} else if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected: \"%v\", Received: \"%v\"", expectedErrorMessage, err.Error())
	}

	if len(cmd.commands) != 0 {
		t.Fatalf("Expected: %d, Received: %d", 0, len(cmd.commands))
	}
}

func TestAddCommandAfterCallingFinish(t *testing.T) {
	cmd := CreateCmd()

	err := cmd.AddCommand("new task",
		CommandArgumentSyntax{
			Name:     "Task name",
			Required: false,
			Type:     COMMAND_ARGUMENT_TYPE_INT,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	err = cmd.Finish()

	if err != nil {
		t.Fatal(err)
	}

	err = cmd.AddCommand("new task",
		CommandArgumentSyntax{
			Name:     "Task name",
			Required: false,
			Type:     COMMAND_ARGUMENT_TYPE_INT,
		},
	)

	expectedErrorMessage := "You cannot add commands after calling Finish method"

	if err == nil {
		t.Fatalf("Expected: \"%v\", Received: nil", expectedErrorMessage)
	} else if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected: \"%v\", Received: \"%v\"", expectedErrorMessage, err.Error())
	}
}

func TestCallingFinishTwice(t *testing.T) {
	cmd := CreateCmd()

	err := cmd.Finish()

	if err != nil {
		t.Fatalf("Expected:  nil, Received: \"%v\"", err.Error())
	}

	err = cmd.Finish()

	if err == nil {
		t.Fatalf("Expected: \"%v\", Received: nil", "You cannot call Finish more than once")
	}
}

func TestParseFromStringBeforeFinish(t *testing.T) {
	cmd := CreateCmd()

	_, err := cmd.ParseFromString("anything")

	if err == nil {
		t.Fatalf("Expected: \"%v\", Received: nil", "You cannot parse a command before finish the Cmd")
	} else if err.Error() != "You cannot parse a command before finish the Cmd" {
		t.Fatalf("Expected: \"%v\", Received: \"%v\"", "You cannot parse a command before finish the Cmd", err.Error())
	}
}

func TestParseFromStringWhenHavingInvalidCommand(t *testing.T) {
	cmd := CreateCmd()

	cmd.Finish()

	_, err := cmd.ParseFromString("anything")

	expectedErrorMessage := "\"anything\" is not a valid command"

	if err == nil {
		t.Fatalf("Expected: \"%v\", Received: nil", expectedErrorMessage)
	} else if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected: \"%v\", Received: \"%v\"", expectedErrorMessage, err.Error())
	}
}

func TestParseFromStringWhenHavingValidCommandWithNoArguments(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand("new task")

	cmd.Finish()

	command, err := cmd.ParseFromString("   new task   ")

	if err != nil {
		t.Fatalf("Expected: nil, Received: \"%v\"", err.Error())
	}

	if command.Name != "new task" {
		t.Fatalf("Invalid command name \"%v\"", command.Name)
	}

	if len(command.Arguments) != 0 {
		t.Fatalf("Invalid arguments length \"%d\"", len(command.Arguments))
	}
}

func TestParseFromStringWhenHasRequiredArgumentsNonProvided(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand(
		"delete task",
		CommandArgumentSyntax{
			Name:     "task id (int)",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_INT,
		},
	)

	cmd.Finish()

	_, err := cmd.ParseFromString("delete task")

	expectedErrorMessage := `"task id (int)" is required`

	if err == nil {
		t.Fatalf(`Expected: "%v", Received: nil`, expectedErrorMessage)
	} else if err.Error() != expectedErrorMessage {
		t.Fatalf(`Expected: "%v", Received: "%v"`, expectedErrorMessage, err.Error())
	}
}

func TestParseFromStringWhenHavingValidCommandWithBadIntArgument(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand(
		"delete task",
		CommandArgumentSyntax{
			Name:     "task id (int)",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_INT,
		},
	)

	cmd.Finish()

	_, err := cmd.ParseFromString("delete task two")

	expectedErrorMessage := `Invalid "task id (int)" argument. Cannot parse it as int`

	if err == nil {
		t.Fatalf(`Expected: "%v", Received: nil`, expectedErrorMessage)
	} else if err.Error() != expectedErrorMessage {
		t.Fatalf(`Expected: "%v", Received: "%v"`, expectedErrorMessage, err.Error())
	}
}

func TestParseFromStringWhenHavingValidCommandWithValidIntArgument(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand(
		"delete task",
		CommandArgumentSyntax{
			Name:     "task id (int)",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_INT,
		},
	)

	cmd.Finish()

	command, err := cmd.ParseFromString(`delete task 123`)

	if err != nil {
		t.Fatalf(`Expected: nil, Received: "%v"`, err.Error())
	}

	if command.Arguments[0].Value != 123 {
		t.Fatalf(`Expected: 123, Received: "%d"`, command.Arguments[0].Value)
	}
}

func TestParseFromStringWhenHavingValidCommandWithInvalidStringArgument(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand(
		"new task",
		CommandArgumentSyntax{
			Name:     "task name",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_STR,
		},
	)

	cmd.Finish()

	_, err := cmd.ParseFromString(`new task "this is the name of my task don't care about`)

	expectedErrorMessage := `Invalid argument "task name" format`

	if err == nil {
		t.Fatalf(`Expected: "%v", Received: nil`, expectedErrorMessage)
	} else if err.Error() != expectedErrorMessage {
		t.Fatalf(`Expected: "%v", Received: "%v"`, expectedErrorMessage, err.Error())
	}
}

func TestParseFromStringWhenHavingValidCommandWithValidStringWithoutQuotes(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand(
		"new task",
		CommandArgumentSyntax{
			Name:     "task name",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_STR,
		},
	)

	cmd.Finish()

	command, err := cmd.ParseFromString(`new task sdlkjsdlfkj2309dsfisdjfl`)

	if err != nil {
		t.Fatalf(`Expected: "%v", Received: nil`, err.Error())
	}

	expectedValue := "sdlkjsdlfkj2309dsfisdjfl"

	if command.Arguments[0].Value != expectedValue {
		t.Fatalf(`Expected: "%v", Received: "%v"`, expectedValue, command.Arguments[0].Value)
	}
}

func TestParseFromStringWhenHavingValidCommandWithValidTwoStringWithoutQuotes(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand(
		"new task",
		CommandArgumentSyntax{
			Name:     "task name",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_STR,
		},
		CommandArgumentSyntax{
			Name:     "task description",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_STR,
		},
	)

	cmd.Finish()

	command, err := cmd.ParseFromString(`new task sdlkjsdlfkj2309dsfisdjfl asdflksdfjsdflsdfhsdflk`)

	if err != nil {
		t.Fatalf(`Expected: "%v", Received: nil`, err.Error())
	}

	expectedValue := "sdlkjsdlfkj2309dsfisdjfl"

	if command.Arguments[0].Value != expectedValue {
		t.Fatalf(`Expected: "%v", Received: "%v"`, expectedValue, command.Arguments[0].Value)
	}

	expectedValue = "asdflksdfjsdflsdfhsdflk"

	if command.Arguments[1].Value != expectedValue {
		t.Fatalf(`Expected: "%v", Received: "%v"`, expectedValue, command.Arguments[1].Value)
	}
}

func TestParseFromStringWhenHavingValidCommandWithValidStringWithQuotes(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand(
		"new task",
		CommandArgumentSyntax{
			Name:     "task name",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_STR,
		},
	)

	cmd.Finish()

	command, err := cmd.ParseFromString(`new task "This is the name of my task"`)

	if err != nil {
		t.Fatalf(`Expected: "%v", Received: nil`, err.Error())
	}

	expectedValue := "This is the name of my task"

	if command.Arguments[0].Value != expectedValue {
		t.Fatalf(`Expected: "%v", Received: "%v"`, expectedValue, command.Arguments[0].Value)
	}
}

func TestParseFromStringWhenHavingValidCommandWithInvalidBooleanType(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand(
		"set completed",
		CommandArgumentSyntax{
			Name:     "state (boolean)",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_BOOL,
		},
	)

	cmd.Finish()

	_, err := cmd.ParseFromString(`set completed nono`)

	expectedErrorMessage := `Invalid argument "state (boolean)" type`

	if err == nil {
		t.Fatalf(`Expected: "%v", Received: nil`, expectedErrorMessage)
	} else if err.Error() != expectedErrorMessage {
		t.Fatalf(`Expected: "%v", Received: "%v"`, expectedErrorMessage, err.Error())
	}
}

func TestParseFromStringWhenHavingValidCommandWithValidBooleanType(t *testing.T) {
	cmd := CreateCmd()

	cmd.AddCommand(
		"set completed",
		CommandArgumentSyntax{
			Name:     "state (boolean)",
			Required: true,
			Type:     COMMAND_ARGUMENT_TYPE_BOOL,
		},
	)

	cmd.Finish()

	validTrue := []string{"true", "1", "t", "yes", "yeah", "y"}
	validFalse := []string{"false", "0", "f", "no", "not", "n"}

	for _, vt := range validTrue {
		command, err := cmd.ParseFromString(fmt.Sprintf("set completed %v", vt))

		if err != nil {
			t.Fatalf(`For: %v, Expected: nil, Received: "%v"`, vt, err.Error())
		}

		if command.Arguments[0].Value != true {
			t.Fatalf(`For: %v, Expected: true, Received: "%v"`, vt, command.Arguments[0].Value)
		}
	}

	for _, vf := range validFalse {
		command, err := cmd.ParseFromString(fmt.Sprintf("set completed %v", vf))

		if err != nil {
			t.Fatalf(`For: %v, Expected: nil, Received: "%v"`, vf, err.Error())
		}

		if command.Arguments[0].Value != false {
			t.Fatalf(`For: %v, Expected: false, Received: "%v"`, vf, command.Arguments[0].Value)
		}
	}
}
