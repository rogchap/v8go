package v8go

// Values as a list of v8go Values.
type Values []*Value

func NewValues(iso *Isolate, vals ...interface{}) (Values, error) {
	vv := make(Values, len(vals))
	for i, v := range vals {
		v, err := NewValue(iso, v)
		if err != nil {
			return nil, err
		}
		vv[i] = v
	}
	return vv, nil
}

// ValuesToValuer converts actual value to Valuer interface.
func (vv Values) Valuers() []Valuer {
	xx := make([]Valuer, len(vv))
	for _, v := range vv {
		xx = append(xx, v)
	}
	return xx
}

func (vv Values) Get(i int) (*Value, bool) {
	if i < len(vv) {
		return vv[i], true
	}
	return nil, false
}

// ArgsValidate will validate arguments based on given conditions.
// Error returned can be either *ArgCondErr or any other.
// In case of ArgCondErr Error() will report given message.
func (vv Values) Validate(conds ...ValueCond) error {
	for _, c := range conds {
		if err := c(vv); err != nil {
			return err
		}
	}
	return nil
}

// ValueCondErr is an error indicating validation of arguments
// failed
type ValueCondErr struct{ message string }

func (err *ValueCondErr) Error() string { return err.message }

// NewValueCondErr creates new argument error.
func NewValueCondErr(message string) *ValueCondErr {
	return &ValueCondErr{message}
}

// ValueCond is a condition returning argument error or any other.
// Useful for validation of arguments.
type ValueCond func(Values) error

// ValueCondLen asserts excact number of arguments.
func ValueCondLen(message string, l int) ValueCond {
	return func(vv Values) error {
		if len(vv) != l {
			return NewValueCondErr(message)
		}
		return nil
	}
}

// ValueCondMinLen asserts minimum of arguments.
func ValueCondMinLen(message string, l int) ValueCond {
	return func(vv Values) error {
		if len(vv) < l {
			return NewValueCondErr(message)
		}
		return nil
	}
}

// ValueCondTypeOf asserts any given with logical OR.
func ValueCondTypeOf(message string, i int, xx ...func(*Value) bool) ValueCond {
	return func(vv Values) error {
		arg, ok := vv.Get(i)
		if !ok {
			return NewValueCondErr(message)
		}
		for _, x := range xx {
			ok := x(arg)
			if ok {
				return nil
			}
		}
		return NewValueCondErr(message)
	}
}

// ValueCondOptionalType asserts given type at given index. If argument is null or undefined
// it will not be validated.
func ValueCondOptionalType(message string, i int, x func(*Value) bool) ValueCond {
	return func(vv Values) error {
		arg, ok := vv.Get(i)
		if !ok {
			return nil
		}
		if arg.IsNullOrUndefined() {
			return nil
		}
		ok = x(arg)
		if !ok {
			return NewValueCondErr(message)
		}
		return nil
	}
}
