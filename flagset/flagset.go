package flagset

import (
	"github.com/urfave/cli/v2"
)

// FlagSet provides additional functionality on top of a collection of flags.
type FlagSet []cli.Flag

// Filter returns a new FlagSet that contains flags allowed by the provided Filter.
func (flags FlagSet) Filter(allow Filter) FlagSet {
	next := FlagSet{}

	for _, flag := range flags {
		if allow(flag) {
			next = append(next, flag)
		}
	}

	return next
}

// Filter allows the user to inspect the flag to determine if it should be in the resulting FlagSet.
type Filter func(flag cli.Flag) bool
