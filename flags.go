package nxboundary

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

var errWrongAllowedTags = errors.New("allowedTags flag must be of form tag:tag1,tag2")

var FlagAllowedTags = "allowedTags"

func flags(config *Config) flag.FlagSet {
	fs := flag.FlagSet{}
	fs.Var(stringMap(config.DepConstraints), FlagAllowedTags, "a tag that's allowed to import from the specified tags form tag|tag1,tag2")

	return fs
}

type stringMap map[string][]string

func (v stringMap) Set(val string) error {
	lastPipe := strings.LastIndex(val, "|")

	// Make sure there's at least one colon, and it's not the first or last character
	if lastPipe <= 1 || lastPipe == len(val)-1 {
		return errWrongAllowedTags
	}

	v[val[:lastPipe]] = strings.Split(val[lastPipe+1:], ",")

	return nil
}

func (v stringMap) String() string {
	return fmt.Sprintf("%v", (map[string][]string)(v))
}
