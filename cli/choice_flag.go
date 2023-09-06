package cli

import (
	"flag"
	"fmt"

	"github.com/urfave/cli"
)

type ChoiceFlag struct {
	cli.StringFlag
	ValidValues []string
}

func NewChoiceFlag(sf cli.StringFlag, validValues []string) *ChoiceFlag {
	return &ChoiceFlag{sf, validValues}
}

type ChoiceFlagValue struct {
	Flag *ChoiceFlag
}

func (value *ChoiceFlagValue) Apply(set *flag.FlagSet, val string) error {
	for _, validValue := range value.Flag.ValidValues {
		if val == validValue {
			return nil
		}
	}
	return fmt.Errorf("invalid value for %s. Valid values are: %v", value.Flag.Name, value.Flag.ValidValues)
}
