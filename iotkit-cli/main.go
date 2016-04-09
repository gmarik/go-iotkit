package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gmarik/go-iotkit/iotkit"

	"github.com/google/subcommands"
	"golang.org/x/net/context"
)

var (
	stdin  = os.Stdin
	stdout = os.Stdout
	stderr = os.Stderr
)

func main() {
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			log.Printf("\nRedirect to: %s %s", req.Method, req.URL.String())
			// fmt.Printf("\n%#v\n\n\n%#v\n", req, via)
			return nil
		},
	}
	// client := iotkit.NewClient(&iotkit.Dumper{httpClient})
	client := iotkit.NewClient(httpClient)

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&authLogin{client: client}, "")
	subcommands.Register(&deviceComponentCreate{client: client}, "")
	subcommands.Register(&deviceActivate{client: client}, "")
	subcommands.Register(&observationCreate{client: client}, "")
	subcommands.Register(&uuidGen{}, "")
	// subcommands.Register(&componentTypeCreate{client: client}, "")

	flag.Parse()

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
