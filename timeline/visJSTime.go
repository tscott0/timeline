package timeline

import "time"

type visJSTime struct {
	time.Time
}

func (t visJSTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.Format("2006-01-02 15:04:05") + `"`), nil
}
