package wsrouter

type Event struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
}

func (e *Event) GetID() string {
	return e.Code
}

func (e *Event) GetData() interface{} {
	return e.Data
}

func CreateEvent(code string, data interface{}) *Event {
	return &Event{
		Code: code,
		Data: data,
	}
}

func CreateErrorEvent(code string, err string) *Event {
	return &Event{
		Code: code,
		Data: err,
	}
}
