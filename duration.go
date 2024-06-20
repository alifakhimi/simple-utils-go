package simutils

import (
	"encoding/json"
	"time"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	if p, err := time.ParseDuration(str); err != nil {
		return err
	} else {
		*d = Duration{p}
	}

	return nil
}
