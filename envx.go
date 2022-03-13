package envx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// A EnvSet represents a set of defined flags. The zero value of a EnvSet
// has no name and has ContinueOnError error handling.
//
// Flag names must be unique within a EnvSet. An attempt to define a flag whose
// name is already in use will cause a panic.
type EnvSet struct {
	prefix string
	parsed bool
	envs   map[string]Value
}

// NewEnvSet returns a new, empty flag set with the specified name and
// error handling property. If the name is not empty, it will be printed
// in the default usage message and in error messages.
func NewEnvSet(prefix string) *EnvSet {
	return &EnvSet{
		envs:   make(map[string]Value),
		prefix: strings.ToUpper(prefix),
	}
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the EnvSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help or -h were set but not defined.
func (f *EnvSet) Parse(envs []string) error {
	f.parsed = true
	for _, env := range envs {
		idx := strings.Index(env, "=")
		if idx == -1 {
			continue
		}

		name, value := env[:idx], env[idx+1:]
		env, ok := f.envs[name]
		if ok {
			if err := env.Set(value); err != nil {
				return fmt.Errorf("cannot set value for %s: %w", name, err)
			}
		}
	}
	return nil
}

// IsParsed reports whether f.Parse has been called.
func (f *EnvSet) IsParsed() bool {
	return f.parsed
}

func (f *EnvSet) getName(s string) string {
	if f.prefix == "" {
		return strings.ToUpper(s)
	}
	return f.prefix + "_" + strings.ToUpper(s)
}

// // VisitAll visits the flags in lexicographical order, calling fn for each.
// // It visits all flags, even those not set.
// func (f *EnvSet) VisitAll(fn func(*Flag)) {
// 	for _, flag := range sortFlags(f.formal) {
// 		fn(flag)
// 	}
// }

// // Visit visits the flags in lexicographical order, calling fn for each.
// // It visits only those flags that have been set.
// func (f *EnvSet) Visit(fn func(*Flag)) {
// 	for _, flag := range sortFlags(f.actual) {
// 		fn(flag)
// 	}
// }

// // Lookup returns the Flag structure of the named flag, returning nil if none exists.
// func (f *EnvSet) Lookup(name string) *Flag {
// 	return f.formal[name]
// }

// // Set the value of the env.
// func (f *EnvSet) Set(name, value string) error {
// 	flag, ok := f.formal[name]
// 	if !ok {
// 		return fmt.Errorf("no such flag -%v", name)
// 	}
// 	err := flag.Value.Set(value)
// 	if err != nil {
// 		return err
// 	}
// 	if f.actual == nil {
// 		f.actual = make(map[string]*Flag)
// 	}
// 	f.actual[name] = flag
// 	return nil
// }

// // PrintDefaults prints, to standard error unless configured otherwise, the
// // default values of all defined command-line flags in the set. See the
// // documentation for the global function PrintDefaults for more information.
// func (f *EnvSet) PrintDefaults() {
// 	f.VisitAll(func(flag *Flag) {
// 		var b strings.Builder
// 		fmt.Fprintf(&b, "  -%s", flag.Name) // Two spaces before -; see next two comments.
// 		name, usage := UnquoteUsage(flag)
// 		if len(name) > 0 {
// 			b.WriteString(" ")
// 			b.WriteString(name)
// 		}
// 		// Boolean flags of one ASCII letter are so common we
// 		// treat them specially, putting their usage on the same line.
// 		if b.Len() <= 4 { // space, space, '-', 'x'.
// 			b.WriteString("\t")
// 		} else {
// 			// Four spaces before the tab triggers good alignment
// 			// for both 4- and 8-space tab stops.
// 			b.WriteString("\n    \t")
// 		}
// 		b.WriteString(strings.ReplaceAll(usage, "\n", "\n    \t"))

// 		if !isZeroValue(flag, flag.DefValue) {
// 			if _, ok := flag.Value.(*stringValue); ok {
// 				// put quotes on the value
// 				fmt.Fprintf(&b, " (default %q)", flag.DefValue)
// 			} else {
// 				fmt.Fprintf(&b, " (default %v)", flag.DefValue)
// 			}
// 		}
// 		fmt.Fprint(f.Output(), b.String(), "\n")
// 	})
// }

// // defaultUsage is the default function to print a usage message.
// func (f *EnvSet) defaultUsage() {
// 	if f.name == "" {
// 		fmt.Fprintf(f.Output(), "Usage:\n")
// 	} else {
// 		fmt.Fprintf(f.Output(), "Usage of %s:\n", f.name)
// 	}
// 	f.PrintDefaults()
// }

// NOTE: Usage is not just defaultUsage(CommandLine)
// because it serves (via godoc flag Usage) as the example
// for how to write your own usage function.

// Usage prints a usage message documenting all defined command-line flags
// to CommandLine's output, which by default is os.Stderr.
// It is called when an error occurs while parsing flags.
// The function is a variable that may be changed to point to a custom function.
// By default it prints a simple header and calls PrintDefaults; for details about the
// format of the output and how to control it, see the documentation for PrintDefaults.
// Custom usage functions may choose to exit the program; by default exiting
// happens anyway as the command line's error handling strategy is set to
// ExitOnError.
// var Usage = func() {
// 	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
// 	PrintDefaults()
// }

// // NFlag returns the number of flags that have been set.
// func (f *EnvSet) NFlag() int { return len(f.actual) }

// // Arg returns the i'th argument. Arg(0) is the first remaining argument
// // after flags have been processed. Arg returns an empty string if the
// // requested element does not exist.
// func (f *EnvSet) Arg(i int) string {
// 	if i < 0 || i >= len(f.args) {
// 		return ""
// 	}
// 	return f.args[i]
// }

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func (f *EnvSet) Bool(p *bool, name string, value bool, usage string) {
	f.Var(newBoolValue(value, p), name, usage)
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func (f *EnvSet) Int(p *int, name string, value int, usage string) {
	f.Var(newIntValue(value, p), name, usage)
}

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func (f *EnvSet) Int64(p *int64, name string, value int64, usage string) {
	f.Var(newInt64Value(value, p), name, usage)
}

// UintVar defines a uint flag with specified name, default value, and usage string.
// The argument p points to a uint variable in which to store the value of the flag.
func (f *EnvSet) Uint(p *uint, name string, value uint, usage string) {
	f.Var(newUintValue(value, p), name, usage)
}

// Uint64Var defines a uint64 flag with specified name, default value, and usage string.
// The argument p points to a uint64 variable in which to store the value of the flag.
func (f *EnvSet) Uint64(p *uint64, name string, value uint64, usage string) {
	f.Var(newUint64Value(value, p), name, usage)
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (f *EnvSet) String(p *string, name string, value string, usage string) {
	f.Var(newStringValue(value, p), name, usage)
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func (f *EnvSet) Float64(p *float64, name string, value float64, usage string) {
	f.Var(newFloat64Value(value, p), name, usage)
}

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
// The flag accepts a value acceptable to time.ParseDuration.
func (f *EnvSet) Duration(p *time.Duration, name string, value time.Duration, usage string) {
	f.Var(newDurationValue(value, p), name, usage)
}

// Func defines a flag with the specified name and usage string.
// Each time the flag is seen, fn is called with the value of the flag.
// If fn returns a non-nil error, it will be treated as a flag value parsing error.
func (f *EnvSet) Func(name, usage string, fn func(string) error) {
	f.Var(funcValue(fn), name, usage)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (f *EnvSet) Var(value Value, name string, usage string) {
	// Flag must not begin "-" or contain "=".
	// if strings.HasPrefix(name, "-") {
	// 	panic(f.sprintf("flag %q begins with -", name))
	// } else if strings.Contains(name, "=") {
	// 	panic(f.sprintf("flag %q contains =", name))
	// }

	flag := &Flag{
		Name:     f.getName(name),
		Usage:    usage,
		Value:    value,
		DefValue: value.String(),
	}
	// _, alreadythere := f.formal[name]
	// if alreadythere {
	// 	var msg string
	// 	// if f.name == "" {
	// 	msg = f.sprintf("flag redefined: %s", name)
	// 	// } else {
	// 	// 	msg = f.sprintf("%s flag redefined: %s", f.name, name)
	// 	// }
	// 	panic(msg) // Happens only if flags are declared with identical names
	// }
	// if f.formal == nil {
	// 	f.formal = make(map[string]*Flag)
	// }
	// f.formal[name] = flag
	f.envs[flag.Name] = flag.Value
}

// sprintf formats the message, prints it to output, and returns it.
func (f *EnvSet) sprintf(format string, a ...interface{}) string {
	msg := fmt.Sprintf(format, a...)
	// fmt.Fprintln(f.Output(), msg)
	return msg
}

// failf prints to standard error a formatted error and usage message and
// returns the error.
func (f *EnvSet) failf(format string, a ...interface{}) error {
	msg := f.sprintf(format, a...)
	// f.usage()
	return errors.New(msg)
}

// usage calls the Usage method for the flag set if one is specified,
// or the appropriate default usage function otherwise.
// func (f *EnvSet) usage() {
// 	if f.Usage == nil {
// 		// f.defaultUsage()
// 	} else {
// 		f.Usage()
// 	}
// }

// // parseOne parses one flag. It reports whether a flag was seen.
// func (f *EnvSet) parseOne() (bool, error) {
// 	if len(f.args) == 0 {
// 		return false, nil
// 	}
// 	s := f.args[0]
// 	if len(s) < 2 || s[0] != '-' {
// 		return false, nil
// 	}
// 	numMinuses := 1
// 	if s[1] == '-' {
// 		numMinuses++
// 		if len(s) == 2 { // "--" terminates the flags
// 			f.args = f.args[1:]
// 			return false, nil
// 		}
// 	}
// 	name := s[numMinuses:]
// 	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
// 		return false, f.failf("bad flag syntax: %s", s)
// 	}

// 	// it's a flag. does it have an argument?
// 	f.args = f.args[1:]
// 	hasValue := false
// 	value := ""
// 	for i := 1; i < len(name); i++ { // equals cannot be first
// 		if name[i] == '=' {
// 			value = name[i+1:]
// 			hasValue = true
// 			name = name[0:i]
// 			break
// 		}
// 	}
// 	m := f.formal
// 	flag, alreadythere := m[name] // BUG
// 	if !alreadythere {
// 		if name == "help" || name == "h" { // special case for nice help message.
// 			// f.usage()
// 			return false, errors.New("ErrHelp")
// 		}
// 		return false, f.failf("flag provided but not defined: -%s", name)
// 	}

// 	if fv, ok := flag.Value.(*boolValue); ok { //&& fv.IsBoolFlag() { // special case: doesn't need an arg
// 		if hasValue {
// 			if err := fv.Set(value); err != nil {
// 				return false, f.failf("invalid boolean value %q for -%s: %v", value, name, err)
// 			}
// 		} else {
// 			if err := fv.Set("true"); err != nil {
// 				return false, f.failf("invalid boolean flag %s: %v", name, err)
// 			}
// 		}
// 	} else {
// 		// It must have a value, which might be the next argument.
// 		if !hasValue && len(f.args) > 0 {
// 			// value is the next arg
// 			hasValue = true
// 			value, f.args = f.args[0], f.args[1:]
// 		}
// 		if !hasValue {
// 			return false, f.failf("flag needs an argument: -%s", name)
// 		}
// 		if err := flag.Value.Set(value); err != nil {
// 			return false, f.failf("invalid value %q for flag -%s: %v", value, name, err)
// 		}
// 	}
// 	if f.actual == nil {
// 		f.actual = make(map[string]*Flag)
// 	}
// 	f.actual[name] = flag
// 	return true, nil
// }

// func commandLineUsage() {
// 	Usage()
// }

// // Init sets the name and error handling property for a flag set.
// // By default, the zero EnvSet uses an empty name and the
// // ContinueOnError error handling policy.
// func (f *EnvSet) Init(name string, errorHandling ErrorHandling) {
// 	f.name = name
// 	f.errorHandling = errorHandling
// }

// isZeroValue determines whether the string represents the zero
// value for a flag.
func isZeroValue(flag *Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(flag.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return value == z.Interface().(Value).String()
}
