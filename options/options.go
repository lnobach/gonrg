package options

import "strings"

type Options []string

func (o Options) HasOption(name string) bool {

	for _, oe := range o {
		if oe == name {
			return true
		}
	}
	return false
}

func (o Options) String() string {
	return strings.Join(o, ",")
}
