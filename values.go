package v8go

// Values as a list of v8go Values.
type Values []*Value

// ValuesToValuer converts actual value to Valuer interface.
func (vv Values) Valuers() []Valuer {
	xx := make([]Valuer, len(vv))
	for _, v := range vv {
		xx = append(xx, v)
	}
	return xx
}
