package main

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os/exec"
)

type executionerCmd struct {
	program       string
	args          []string
	timeout       timeout
	preconditions preconditions
	streams       streamCollection
}

type timeout struct {
	milliseconds float64
	response     responseType
}

type responseType string

type preconditions struct {
	program string
	args    []string
	output  string
}

type streamCollection struct {
	stdin   stream
	stdout  stream
	stderr  stream
	logging stream
}

type stream struct {
	std    standardStreamType
	inline string
	file   string
}

type standardStreamType string

func main() {
	executioner, _ := ParseExecutioner("sample_command.yaml")
	cmd, _ := ExecutionerToCmd(executioner)

	fmt.Printf("executioner: %+v\n", executioner)
	fmt.Printf("cmd: %+v\n", cmd)
}

func RunScript(cmd exec.Cmd) (string, error) {
	return "", nil
}

func ExecutionerToCmd(executioner executionerCmd) (exec.Cmd, error) {
	return exec.Cmd{}, nil
}

func ParseExecutioner(file string) (executionerCmd, error) {
	viper.SetConfigFile(file)
	viper.ReadInConfig()

	viper.Get("timeout")
	timeout, _ := parseTimeout(viper.GetStringMap("timeout"))
	preconds, _ := parsePreconditions(viper.GetStringMap("preconditions"))
	streams, _ := parseStreamCollection(viper.GetStringMap("streams"))

	ecmd := executionerCmd{
		program:       viper.GetString("prog"),
		args:          viper.GetStringSlice("args"),
		timeout:       timeout,
		preconditions: preconds,
		streams:       streams,
	}
	return ecmd, nil
}

func parseTimeout(m map[string]interface{}) (timeout, error) {
	return timeout{
		milliseconds: cast.ToFloat64(m["milliseconds"]),
		response:     responseType(cast.ToString(m["response"])),
	}, nil
}

func parsePreconditions(m map[string]interface{}) (preconditions, error) {
	return preconditions{
		program: cast.ToString(m["prog"]),
		args:    cast.ToStringSlice(m["args"]),
		output:  cast.ToString(m["output"]),
	}, nil
}

func parseStreamCollection(m map[string]interface{}) (streamCollection, error) {
	stdin, _ := parseStream(cast.ToStringMapString(m["stdin"]))
	stdout, _ := parseStream(cast.ToStringMapString(m["stdout"]))
	stderr, _ := parseStream(cast.ToStringMapString(m["stderr"]))
	logging, _ := parseStream(cast.ToStringMapString(m["logging"]))

	return streamCollection{
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		logging: logging,
	}, nil
}

func parseStream(m map[string]string) (stream, error) {
	return stream{
		inline: m["inline"],
		std:    standardStreamType(m["standard"]),
		file:   m["file"],
	}, nil
}
