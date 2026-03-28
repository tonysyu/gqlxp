package gqlfmt

import "strings"

// IncludeOptions specifies which optional sections to include in output
type IncludeOptions struct {
	Usages     bool
	ReturnType bool
}

// ParseIncludeOptions parses a comma-separated string of include options
func ParseIncludeOptions(include string) IncludeOptions {
	var opts IncludeOptions
	if include == "" {
		return opts
	}

	for _, opt := range strings.Split(include, ",") {
		switch strings.TrimSpace(strings.ToLower(opt)) {
		case "usages":
			opts.Usages = true
		case "return-type":
			opts.ReturnType = true
		}
	}
	return opts
}
