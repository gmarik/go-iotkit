package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/gmarik/go-iotkit/iotkit"

	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

var (
	logger = log.New(ioutil.Discard, "", log.LstdFlags)

	verbose = flag.Bool("verbose", false, "Log info to stdout")
	delay   = flag.Duration("delay", time.Duration(1*time.Second), "Delay between sensor reads")
	kind    = flag.String("kind", "", "sensor kind. Allowed: temp,pseudo, integer by default")
	pin     = flag.String("pin", "", "Analog pin-ID to value read from. Required.")
	maxSize = flag.Int("flush-size", 1000, "Observation batch size to submit")

	accountID   = flag.String("account-uuid", "", "Account ID component belongs to")
	componentID = flag.String("component-uuid", "", "Component Id being observed")
	deviceID    = flag.String("device-uuid", "", "Device uuid")
	deviceToken = flag.String("device-token", "", "Device Token")

	client *iotkit.Client = iotkit.NewClient(nil)
)

type AnalogReader interface {
	AnalogRead(string) (interface{}, error)
}

type AnalogReaderFunc func(string) (interface{}, error)

func (fn AnalogReaderFunc) AnalogRead(pin string) (interface{}, error) {
	return fn(pin)
}

func main() {
	flag.Parse()

	if *verbose {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	var (
		errc = make(chan error)
		obs  = make(chan iotkit.Observation)
		done = make(chan struct{})
	)

	go func() { errc <- SensorReader(done, obs, reader()) }()
	go func() { errc <- Collector(done, obs) }()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	// Block until a signal is received.

	select {
	case err := <-errc:
		if err != nil {
			log.Println("An error occured:", err)
		}

	case s := <-sig:
		logger.Println("Got signal:", s)
	}

	close(done)
}

func SensorReader(done <-chan struct{}, observations chan<- iotkit.Observation, rd AnalogReader) error {
	for {
		v, err := rd.AnalogRead(*pin)
		if err != nil {
			return err
		}
		logger.Println("Read:", v, err)

		o := iotkit.Observation{
			ComponentId: *componentID,
			On:          iotkit.Time(time.Now()),
			Value:       fmt.Sprintf("%v", v),
		}

		select {
		case <-done:
			return nil
		case observations <- o:
		}
	}

	return nil
}

func Collector(done <-chan struct{}, observations <-chan iotkit.Observation) error {
	batch := make([]iotkit.Observation, 0, 2*(*maxSize))

	for {
		if len(batch) >= *maxSize {
			logger.Println("Submitting batch; size: ", len(batch))
			if err := submit(batch); err != nil {
				return err
			}
			// emptying batch
			batch = batch[:0]
		}

		select {
		case <-done:
			return nil
		case o := <-observations:
			logger.Printf("Observation: %#v\n", o)
			batch = append(batch, o)
		}
	}
}

func reader() AnalogReader {
	var reader AnalogReader

	logger.Println("Reader for kind:", *kind)

	switch *kind {
	case "random":
		reader = AnalogReaderFunc(func(pin string) (interface{}, error) {
			return rand.Intn(1024), nil
		})
	case "temp":
		reader = AnalogReaderFunc(func(pin string) (interface{}, error) {
			device := edison.NewEdisonAdaptor("edison")
			v, err := device.AnalogRead(pin)
			if err != nil {
				return 0, err
			}
			return TempConvert(v), nil
		})
	default:
		device := edison.NewEdisonAdaptor("edison")
		reader = AnalogReaderFunc(func(pin string) (interface{}, error) {
			v, err := device.AnalogRead(pin)
			if err != nil {
				return 0, err
			}
			return v, nil
		})
	}

	// rate limiter
	return AnalogReaderFunc(func(pin string) (interface{}, error) {
		<-time.After(*delay)
		return reader.AnalogRead(pin)
	})
}

func submit(data []iotkit.Observation) error {
	batch := iotkit.ObservationBatch{
		AccountId: *accountID,
		On:        iotkit.Time(time.Now()),
		Data:      data,
	}

	_, err := client.ObservationAPI.Create(
		batch,
		iotkit.Device{
			*deviceID,
			*deviceToken,
		},
	)

	return err
}

func TempConvert(val int) float64 {
	const B = 3975
	v := float32(val)
	var (
		resistance = (1023.0 - v) * 10000.0 / v
		div        = math.Log(float64(resistance)/10000.0)/B + 1/298.15
	)

	return 1/div - 273.15
}
