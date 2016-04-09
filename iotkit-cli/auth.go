package main

import (
	"flag"
	"fmt"

	"github.com/gmarik/go-iotkit/iotkit"

	"github.com/google/subcommands"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
)

type authLogin struct {
	client *iotkit.Client

	username string
}

func (*authLogin) Name() string     { return "auth:login" }
func (*authLogin) Synopsis() string { return "get the authenticate token" }
func (*authLogin) Usage() string {
	return `auth:login -user <email>
	asks for a password, requests auth token and prints it
`
}

func (p *authLogin) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.username, "user", "", "user account's email address")
}

func (p *authLogin) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.username == "" {
		f.Usage()
		return subcommands.ExitUsageError
	}

	fmt.Fprint(stdout, "Password:")
	buf, err := terminal.ReadPassword(int(stdin.Fd()))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return subcommands.ExitFailure
	}

	tok, _, err := p.client.CreateToken(p.username, string(buf))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return subcommands.ExitFailure
	}

	fmt.Fprintln(stdout, "\nUser token:", tok)

	return subcommands.ExitSuccess
}
