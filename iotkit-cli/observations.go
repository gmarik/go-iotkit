package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gmarik/go-iotkit/iotkit"

	"github.com/google/subcommands"

	"golang.org/x/net/context"
)

type observationCreate struct {
	client *iotkit.Client

	iotkit.ObservationBatch
	iotkit.Observation
	iotkit.Device
}

func (*observationCreate) Name() string     { return "observation:submit" }
func (*observationCreate) Synopsis() string { return "Submit single observation" }
func (*observationCreate) Usage() string {
	return `observation:submit
	submits single observation
`
}

func (p *observationCreate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ObservationBatch.AccountId, "account-uuid", "", "account uuid identifier")

	f.StringVar(&p.Observation.ComponentId, "component-uuid", "", "component uuid identifier")
	f.StringVar(&p.Observation.Value, "value", "", "observation value")

	f.StringVar(&p.Device.ID, "device-uuid", "", "device uuid identifier")
	f.StringVar(&p.Device.Token, "device-token", "", "device auth JWT Token")
}

func (p *observationCreate) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	var now = iotkit.Time(time.Now())

	p.ObservationBatch.On = now
	p.Observation.On = now

	p.ObservationBatch.Data = []iotkit.Observation{p.Observation}

	_, err := p.client.Create(p.ObservationBatch, p.Device)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
