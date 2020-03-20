package send

import "encoding/json"

func NewNotify(data []byte) (*Notify, error) {
	var alarm Notify
	if err := json.Unmarshal(data, &alarm); err != nil {
		return nil, err
	}
	return &alarm, nil
}

type Notify struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	*Content
}

func (n Notify) String() string {
	d, _ := json.Marshal(n)
	return string(d)
}
