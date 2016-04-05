package main

import (
	"flag"
	"fmt"

	"github.com/gmarik/go-iotkit/iotkit"

	"github.com/google/subcommands"

	"golang.org/x/net/context"
)

type deviceComponentCreate struct {
	client *iotkit.Client

	iotkit.Device
	iotkit.Account
	iotkit.Component
}

func (*deviceComponentCreate) Name() string     { return "device:component:create" }
func (*deviceComponentCreate) Synopsis() string { return "create deveice component" }
func (*deviceComponentCreate) Usage() string {
	return `
	device:component:create -account-uuid <aid> -device-uuid <did> -device-token <tok> -type <type> -name <name> -uuid <uuid>
	creates new component for specified account and device
`
}

func (p *deviceComponentCreate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.Account.ID, "account-uuid", "", "account uuid identifier")

	f.StringVar(&p.Device.ID, "device-uuid", "", "device uuid identifier")
	f.StringVar(&p.Device.Token, "device-token", "", "device JWT token")

	f.StringVar(&p.Component.ID, "uuid", "", "component uuid identifier")
	f.StringVar(&p.Component.Type, "type", "custom.v1.0", "component type:  [custom.v1.0|temperature.v1.0|humidity.v1.0|powerswitch.v1.0]")
	f.StringVar(&p.Component.Name, "name", "", "component name")
}

func (p *deviceComponentCreate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// TODO: generate Component UUID
	_, err := p.client.CreateComponent(p.Component, p.Device, p.Account)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return subcommands.ExitFailure
	}

	fmt.Fprintf(stdout, "\nComponent[%s] created.", p.Component.ID)

	return subcommands.ExitSuccess
}
