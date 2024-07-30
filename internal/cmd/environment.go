package cmd

import (
	"errors"
	"strings"

	"github.com/KineticCommerce/kinetic-cli/internal/kinetic"
)

type environmentType string

const (
	environmentTypeLocal  environmentType = "local"
	environmentTypeDit    environmentType = "dit"
	environmentTypeStage  environmentType = "stage"
	environmentTypeProd   environmentType = "prod"
	environmentTypeProdEu environmentType = "prod-eu"
	environmentTypeAll    environmentType = "all"
)

var environmentTypeCompletionFunc = kinetic.FlagCompletionFunc([]string{
	string(environmentTypeDit),
	string(environmentTypeStage),
	string(environmentTypeProd),
	string(environmentTypeProdEu),
})

func (e *environmentType) Set(s string) error {
	switch strings.ToLower(s) {
	case "local":
		*e = environmentTypeLocal
	case "dit":
		*e = environmentTypeDit
	case "stage":
		*e = environmentTypeStage
	case "prod":
		*e = environmentTypeProd
	case "prod-eu":
		*e = environmentTypeProdEu
	case "all":
		*e = environmentTypeAll
	default:
		return errors.New("invalid or unsupported environment")
	}

	return nil
}

func (e environmentType) String() string { return string(e) }

func (e environmentType) Type() string { return "dit|stage|prod|prod-eu" }
