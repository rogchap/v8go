package internal

type Isolate struct{}
type Value struct{}
type ValueTemplate struct{}
type Valuer interface{}

type Values []*Value

func (vv Values) Valuers() []Valuer { return nil }

func NewValue(Isolate, val interface{}) (*Value, error)                 { return nil, nil }
func NewValueTemplate(Isolate, val interface{}) (*ValueTemplate, error) { return nil, nil }
