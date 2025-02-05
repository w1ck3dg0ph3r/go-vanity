package forges

import (
	"fmt"
	"strings"
)

type Vanity interface {
	Module() string
	VCSHome() string
	VCSSourceDir() string
	VCSSourceFile() string
	fmt.Stringer
}

type Forge = func([]string) (Vanity, error)

var forges = map[string]Forge{
	"forgejo": NewForgejo,
	"github":  NewGithub,
	"gitlab":  NewGitlab,
}

func Parse(s string) (Vanity, error) {
	fs := strings.Fields(s)
	if len(fs) == 0 {
		return nil, fmt.Errorf("invalid vanity: %q", s)
	}
	if parse, ok := forges[fs[0]]; ok {
		return parse(fs[1:])
	} else {
		return nil, fmt.Errorf("unknown forge: %q", fs[0])
	}
}
