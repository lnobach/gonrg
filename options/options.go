package options

type Options []string

func (o Options) HasOption(name string) bool {

	for _, oe := range o {
		if oe == name {
			return true
		}
	}
	return false
}
