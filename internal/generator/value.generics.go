package internal

import "github.com/cheekybits/genny/generic"

//go:generate genny -pkg=v8go -in=value.generics.go -out=../../value.gen.go gen "_t_=string,int,int32,int64,uint,uint32,uint64,*big.Int,bool,float64"

type _t_ generic.Type

// New_t_Value creates new Value of _t_
func New_t_Value(iso *Isolate, val _t_) (*Value, error) {
	return NewValue(iso, val)
}

// New_t_Valuer creates new Valuer of _t_. Same as New_t_Value
// except it is casted to interface.
func New_t_Valuer(iso *Isolate, val _t_) (Valuer, error) {
	return NewValue(iso, val)
}

// New_t_Values creates new list of _t_.
func New_t_Values(iso *Isolate, vals ..._t_) (Values, error) {
	vv := make(Values, len(vals))
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return vv, nil
}

// New_t_Valuers creates new list of _t_ Valuer.
func New_t_Valuers(iso *Isolate, vals ..._t_) ([]Valuer, error) {
	vv, err := New_t_Values(iso, vals...)
	return vv.Valuers(), err
}
