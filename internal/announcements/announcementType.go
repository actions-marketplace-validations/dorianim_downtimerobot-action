package announcements

import (
	"bytes"
	"encoding/json"
)

// TaskState represents the state of task, moving through Created, Running then Finished or Errorred
type AnnouncementType int

const (
	Information = iota
	Warning
	Alert
)

func (s AnnouncementType) String() string {
	return toString[s]
}

var toString = map[AnnouncementType]string{
	Information: "Information",
	Warning:     "Warning",
	Alert:       "Alert",
}

var toID = map[string]AnnouncementType{
	"Information": Information,
	"Warning":     Warning,
	"Alert":       Alert,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AnnouncementType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AnnouncementType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toID[j]
	return nil
}
