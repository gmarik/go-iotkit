package main

import (
	"flag"
	"fmt"

	"github.com/gmarik/go-iotkit/iotkit"

	"github.com/google/subcommands"

	"golang.org/x/net/context"
)

type componentTypeCreate struct {
	client *iotkit.Client

	iotkit.Account
	iotkit.ComponentType
}

func (*componentTypeCreate) Name() string     { return "component-type:create" }
func (*componentTypeCreate) Synopsis() string { return "create component type" }
func (*componentTypeCreate) Usage() string {
	return `
device:component:create -account-uuid <aid> -device-uuid <did> -device-token <tok> -device-type <type> -device-name <name> [-component-uuid <uuid>]
	creates new component for specified account and device
	-component-uuid gets auto generated if not specified
`
}

func (p *componentTypeCreate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.Account.ID, "account-uuid", "", "account uuid identifier")
	f.StringVar(&p.Account.Token, "account-token", "", "account authorization JWT token")

	f.StringVar(&p.ComponentType.Dimension, "dimension", "", "")
	f.StringVar(&p.ComponentType.Version, "version", "", "")
	f.StringVar(&p.ComponentType.Type, "type", "", "")
	f.StringVar(&p.ComponentType.DataType, "data-type", "", "")
	f.StringVar(&p.ComponentType.Format, "format", "", "")

	f.Float64Var(&p.ComponentType.Min, "min", 0.0, "")
	f.Float64Var(&p.ComponentType.Max, "max", 0.0, "")

	f.StringVar(&p.ComponentType.MeasureUnit, "measureunit", "", "")
	f.StringVar(&p.ComponentType.Display, "display", "", "")
}

func (p *componentTypeCreate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	c, _, err := p.client.CreateComponentType(p.Account, p.ComponentType)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return subcommands.ExitFailure
	}

	fmt.Fprintf(stdout, "\nComponent created: %#v", c)

	return subcommands.ExitSuccess
}
