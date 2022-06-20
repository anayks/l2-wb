package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	for {
		s, err := ReadCommand()
		if err != nil {
			fmt.Println(fmt.Sprintf("Error on read command: %v", err))
		}

		cmd, err := ParsePipeCommands(s)

		if err != nil {
			fmt.Println(err)
			continue
		}

		output, err := cmd.execute("")

		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(output) > 0 {
			fmt.Printf("Output: %v", output)
		}
	}
}

func ParseCommands(s string) (Command, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("command is empty")
	}

	cmd := strings.Split(s, " ")
	for id, v := range cmd {
		cmd[id] = strings.TrimLeft(v, " ")
		cmd[id] = strings.TrimRight(v, " ")
	}

	result, err := NewCommand(cmd[0], cmd[1:])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ParsePipeCommands(s string) (Command, error) {
	var result CommandPipe
	pipe := strings.Split(s, " | ")
	if len(pipe) == 1 {
		cmd, err := ParseCommands(pipe[0])
		if err != nil {
			return nil, err
		}
		result = append(result, cmd)
		return &result, nil
	}

	for id, v := range pipe {
		cmd, err := ParseCommands(v)
		if err != nil {
			return nil, fmt.Errorf("Error at #%v command: %v", id, err)
		}
		result = append(result, cmd)
	}

	return &result, nil
}

func ReadCommand() (string, error) {
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	bresult, _, err := reader.ReadLine()

	if err != nil {
		if err == io.EOF {
			return "", nil
		}
		return "", err
	}

	result := string(bresult)
	return result, nil
}
