package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"

	"golang.org/x/net/context"
)

type uuidGen struct{}

func (*uuidGen) Name() string     { return "uuidgen" }
func (*uuidGen) Synopsis() string { return "generates a UUID" }
func (*uuidGen) Usage() string {
	return `uuidgen
	generates a UUID
`
}

func (p *uuidGen) SetFlags(f *flag.FlagSet) {}

func (p *uuidGen) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	_uuid, err := uuid()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return subcommands.ExitFailure
	}

	fmt.Fprintln(stdout, _uuid)

	return subcommands.ExitSuccess
}

// from http://stackoverflow.com/a/15130965
func uuid() (string, error) {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		return "", err
	}

	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid, nil
}
