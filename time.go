package dns

import (
	"encoding/json"
	"time"
)

const weirdApitimeLayout = "2006-01-02 15:04:05.999 -0700 MST"

type HetznerTime time.Time

func (ht *HetznerTime) MarshalJSON() ([]byte, error) {
	asString := time.Time(*ht).Format(weirdApitimeLayout)
	return json.Marshal(asString)
}

func (ht *HetznerTime) UnmarshalJSON(b []byte) error {
	var asString string
	err := json.Unmarshal(b, &asString)
	if err != nil {
		return err
	}

	if asString == "" {
		return nil
	}

	value, err := time.Parse(weirdApitimeLayout, asString)
	if err != nil {
		return err
	}

	*ht = HetznerTime(value)
	return nil
}
