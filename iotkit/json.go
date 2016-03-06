package iotkit

import (
	"fmt"
	"time"
)

// Time is a stdlib's Time wraper that implements custom json serialisation
// api expects it to be represented in milliseconds
type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	ts := fmt.Sprintf(`%d`, time.Time(t).Unix()*1000)
	return []byte(ts), nil
}
