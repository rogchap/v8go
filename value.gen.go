// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package v8go

import "math/big"

// NewStringValue creates new Value of string
func NewStringValue(iso *Isolate, val string) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewStringValueTemplate creates new ValueTemplate of string
func NewStringValueTemplate(iso *Isolate, val ...string) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewStringValuer creates new Valuer of string. Same as NewStringValue
// except it is casted to interface.
func NewStringValuer(iso *Isolate, val string) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewStringValuers creates new list of string Valuer.
func NewStringValuers(iso *Isolate, vals ...string) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}

// NewIntValue creates new Value of int
func NewIntValue(iso *Isolate, val int) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewIntValueTemplate creates new ValueTemplate of int
func NewIntValueTemplate(iso *Isolate, val ...int) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewIntValuer creates new Valuer of int. Same as NewIntValue
// except it is casted to interface.
func NewIntValuer(iso *Isolate, val int) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewIntValuers creates new list of int Valuer.
func NewIntValuers(iso *Isolate, vals ...int) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}

// NewInt32Value creates new Value of int32
func NewInt32Value(iso *Isolate, val int32) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewInt32ValueTemplate creates new ValueTemplate of int32
func NewInt32ValueTemplate(iso *Isolate, val ...int32) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewInt32Valuer creates new Valuer of int32. Same as NewInt32Value
// except it is casted to interface.
func NewInt32Valuer(iso *Isolate, val int32) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewInt32Valuers creates new list of int32 Valuer.
func NewInt32Valuers(iso *Isolate, vals ...int32) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}

// NewInt64Value creates new Value of int64
func NewInt64Value(iso *Isolate, val int64) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewInt64ValueTemplate creates new ValueTemplate of int64
func NewInt64ValueTemplate(iso *Isolate, val ...int64) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewInt64Valuer creates new Valuer of int64. Same as NewInt64Value
// except it is casted to interface.
func NewInt64Valuer(iso *Isolate, val int64) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewInt64Valuers creates new list of int64 Valuer.
func NewInt64Valuers(iso *Isolate, vals ...int64) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}

// NewUintValue creates new Value of uint
func NewUintValue(iso *Isolate, val uint) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewUintValueTemplate creates new ValueTemplate of uint
func NewUintValueTemplate(iso *Isolate, val ...uint) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewUintValuer creates new Valuer of uint. Same as NewUintValue
// except it is casted to interface.
func NewUintValuer(iso *Isolate, val uint) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewUintValuers creates new list of uint Valuer.
func NewUintValuers(iso *Isolate, vals ...uint) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}

// NewUint32Value creates new Value of uint32
func NewUint32Value(iso *Isolate, val uint32) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewUint32ValueTemplate creates new ValueTemplate of uint32
func NewUint32ValueTemplate(iso *Isolate, val ...uint32) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewUint32Valuer creates new Valuer of uint32. Same as NewUint32Value
// except it is casted to interface.
func NewUint32Valuer(iso *Isolate, val uint32) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewUint32Valuers creates new list of uint32 Valuer.
func NewUint32Valuers(iso *Isolate, vals ...uint32) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}

// NewUint64Value creates new Value of uint64
func NewUint64Value(iso *Isolate, val uint64) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewUint64ValueTemplate creates new ValueTemplate of uint64
func NewUint64ValueTemplate(iso *Isolate, val ...uint64) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewUint64Valuer creates new Valuer of uint64. Same as NewUint64Value
// except it is casted to interface.
func NewUint64Valuer(iso *Isolate, val uint64) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewUint64Valuers creates new list of uint64 Valuer.
func NewUint64Valuers(iso *Isolate, vals ...uint64) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}

// NewBigIntValue creates new Value of *big.Int
func NewBigIntValue(iso *Isolate, val *big.Int) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewBigIntValueTemplate creates new ValueTemplate of *big.Int
func NewBigIntValueTemplate(iso *Isolate, val ...*big.Int) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewBigIntValuer creates new Valuer of bigInt. Same as NewBigIntValue
// except it is casted to interface.
func NewBigIntValuer(iso *Isolate, val *big.Int) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewBigIntValuers creates new list of *big.Int Valuer.
func NewBigIntValuers(iso *Isolate, vals ...*big.Int) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}

// NewBoolValue creates new Value of bool
func NewBoolValue(iso *Isolate, val bool) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewBoolValueTemplate creates new ValueTemplate of bool
func NewBoolValueTemplate(iso *Isolate, val ...bool) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewBoolValuer creates new Valuer of bool. Same as NewBoolValue
// except it is casted to interface.
func NewBoolValuer(iso *Isolate, val bool) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewBoolValuers creates new list of bool Valuer.
func NewBoolValuers(iso *Isolate, vals ...bool) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}

// NewFloat64Value creates new Value of float64
func NewFloat64Value(iso *Isolate, val float64) (v *Value, err error) {
	return NewValue(iso, val)
}

// NewFloat64ValueTemplate creates new ValueTemplate of float64
func NewFloat64ValueTemplate(iso *Isolate, val ...float64) (vv *ValueTemplate, err error) {
	return NewValueTemplate(iso, val)
}

// NewFloat64Valuer creates new Valuer of float64. Same as NewFloat64Value
// except it is casted to interface.
func NewFloat64Valuer(iso *Isolate, val float64) (v Valuer, err error) {
	return NewValue(iso, val)
}

// NewFloat64Valuers creates new list of float64 Valuer.
func NewFloat64Valuers(iso *Isolate, vals ...float64) (vv []Valuer, err error) {
	for _, v := range vals {
		vx, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv = append(vv, vx)
	}
	return
}
