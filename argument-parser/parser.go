package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CommandArgumentType int

const (
	COMMAND_ARGUMENT_TYPE_INT  CommandArgumentType = iota
	COMMAND_ARGUMENT_TYPE_STR  CommandArgumentType = iota
	COMMAND_ARGUMENT_TYPE_BOOL CommandArgumentType = iota
)

type CommandArgumentSyntax struct {
	Name     string // just to better identify when error
	Required bool
	Type     CommandArgumentType
}

type CommandSyntax struct {
	Name      string                  // for example: new task
	Arguments []CommandArgumentSyntax // for example: "this is my task name"
}

type CommandArgument struct {
	Name  string
	Value any
	Type  CommandArgumentType
}

type Command struct {
	Name      string
	Arguments []CommandArgument
}

type ArgumentParser struct {
	commands    []CommandSyntax
	closed      bool
	commandsMap map[string]*CommandSyntax
}

func CreateArgumentParser() *ArgumentParser {
	return &ArgumentParser{
		commands: []CommandSyntax{},
	}
}

func (cmd *ArgumentParser) AddCommand(name string, arguments ...CommandArgumentSyntax) error {
	if cmd.closed {
		return errors.New("You cannot add commands after calling Finish method")
	}

	hasOptional := false

	for _, argument := range arguments {
		switch argument.Type {
		case COMMAND_ARGUMENT_TYPE_INT, COMMAND_ARGUMENT_TYPE_STR, COMMAND_ARGUMENT_TYPE_BOOL:
			break
		default:
			return errors.New(fmt.Sprintf("Argument \"%v\" have a invalid type \"%d\"", argument.Name, argument.Type))
		}

		if hasOptional && argument.Required {
			return errors.New("You cannot have required arguments after optionals")
		}

		hasOptional = hasOptional || !argument.Required
	}

	cmd.commands = append(cmd.commands, CommandSyntax{
		Name:      name,
		Arguments: arguments,
	})

	return nil
}

func (cmd *ArgumentParser) Finish() error {
	if cmd.closed {
		return errors.New("You cannot call Finish more than once")
	}

	cmd.closed = true
	cmd.commandsMap = make(map[string]*CommandSyntax, len(cmd.commands))

	for _, command := range cmd.commands {
		cmd.commandsMap[command.Name] = &command
	}

	return nil
}

func nextArgument(text string, argumentName string) (string, int, error) {
	pad := len(text)
	text = strings.TrimSpace(text)
	pad -= len(text)

	if len(text) == 0 {
		return "", 0, nil
	}

	var quotes []rune
	hasQuotes := false

	for index, ch := range text {
		if ch == '"' {
			hasQuotes = true

			if len(quotes) == 0 || quotes[len(quotes)-1] != ch {
				quotes = append(quotes, ch)
			} else {
				quotes = quotes[:len(quotes)-1]
			}
		}

		if ch == ' ' && len(quotes) == 0 {
			if hasQuotes {
				return text[1 : index-1], pad + index + 1, nil
			}
			return text[:index], pad + index + 1, nil
		}
	}

	if len(quotes) == 0 {
		if hasQuotes {
			return text[1 : len(text)-1], pad + len(text), nil
		}

		return text, pad + len(text), nil
	}

	return "", 0, errors.New(fmt.Sprintf(`Invalid argument "%v" format`, argumentName))
}

func parseArguments(text string, args []CommandArgumentSyntax) ([]CommandArgument, error) {
	text = strings.TrimSpace(text)

	var arguments []CommandArgument

	for _, arg := range args {
		argument, size, err := nextArgument(text, arg.Name)

		if err != nil {
			return nil, err
		}

		if len(argument) == 0 && arg.Required {
			return nil, errors.New(fmt.Sprintf(`"%v" is required`, arg.Name))
		}

		commandArgument := CommandArgument{
			Name: arg.Name,
			Type: arg.Type,
		}

		switch arg.Type {
		case COMMAND_ARGUMENT_TYPE_STR:
			commandArgument.Value = argument
			break
		case COMMAND_ARGUMENT_TYPE_INT:
			v, err := strconv.Atoi(argument)

			if err != nil {
				return nil, errors.New(fmt.Sprintf(`Invalid "%v" argument. Cannot parse it as int`, arg.Name))
			}

			commandArgument.Value = v
			break
		case COMMAND_ARGUMENT_TYPE_BOOL:
			switch argument {
			case "true", "1", "t", "yes", "yeah", "y":
				commandArgument.Value = true
				break
			case "false", "0", "f", "no", "not", "n":
				commandArgument.Value = false
				break
			default:
				return nil, errors.New(fmt.Sprintf(`Invalid argument "%v" type`, arg.Name))
			}
		default:
			break
		}

		arguments = append(arguments, commandArgument)
		text = text[size-1:]
	}

	return arguments, nil
}

func (cmd *ArgumentParser) ParseFromString(text string) (*Command, error) {
	if !cmd.closed {
		return nil, errors.New("You cannot parse a command before finish the ArgumentParser")
	}

	text = strings.TrimSpace(text)

	var command []rune

	for index, b := range text {
		command = append(command, b)

		if cmd, ok := cmd.commandsMap[string(command)]; ok {

			arguments, err := parseArguments(text[index+1:], cmd.Arguments)

			if err != nil {
				return nil, err
			}

			command := &Command{
				Name:      cmd.Name,
				Arguments: arguments,
			}

			return command, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("\"%v\" is not a valid command", text))
}
