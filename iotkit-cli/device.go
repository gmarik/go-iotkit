package main

import (
	"flag"
	"fmt"

	"github.com/gmarik/go-iotkit/iotkit"

	"github.com/google/subcommands"

	"golang.org/x/net/context"
)

type deviceActivate struct {
	client *iotkit.Client

	iotkit.Device
	iotkit.Account

	activationCode string
}

func (*deviceActivate) Name() string     { return "device:activate" }
func (*deviceActivate) Synopsis() string { return "activate device and return device token" }
func (*deviceActivate) Usage() string {
	return `
device:activate -account-uuid <aid> -device-uuid <did> -device-token <tok> -activation-code <code>
activate device and return device token
`
}

func (p *deviceActivate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.Account.ID, "account-uuid", "", "account uuid identifier")

	f.StringVar(&p.Device.ID, "device-uuid", "", "device uuid identifier")
	f.StringVar(&p.Device.Token, "device-token", "", "device token")
	f.StringVar(&p.activationCode, "activation-code", "", "device activation code")
}

func (p *deviceActivate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// TODO: generate Component UUID
	deviceToken, _, err := p.client.ActivateDevice(p.activationCode, p.Account, p.Device)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return subcommands.ExitFailure
	}

	fmt.Fprintf(stdout, "\nDevice Token: ", deviceToken)

	return subcommands.ExitSuccess
}
