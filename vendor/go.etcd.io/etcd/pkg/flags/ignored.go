package flags

// IgnoredFlag encapsulates a flag that may have been previously valid but is
// now ignored. If an IgnoredFlag is set, a warning is printed and
// operation continues.
type IgnoredFlag struct {
	Name string
}

// IsBoolFlag is defined to allow the flag to be defined without an argument
func (f *IgnoredFlag) IsBoolFlag() bool {
	return true
}

func (f *IgnoredFlag) Set(s string) error {
	plog.Warningf(`flag "-%s" is no longer supported - ignoring.`, f.Name)
	return nil
}

func (f *IgnoredFlag) String() string {
	return ""
}
