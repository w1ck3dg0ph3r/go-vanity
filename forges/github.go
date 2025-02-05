package forges

import "fmt"

type Github struct {
	host       string
	module     string
	home       string
	sourceDir  string
	sourceFile string
}

func NewGithub(fs []string) (Vanity, error) {
	if len(fs) != 3 {
		return Github{}, fmt.Errorf(`github vanity format should be "{module} {path} {ref}"`)
	}
	module, path, ref := fs[0], fs[1], fs[2]
	host := "https://github.com/"
	vcsPath := host + path
	return Github{
		host:       host,
		module:     module,
		home:       vcsPath,
		sourceDir:  vcsPath + "/tree/" + ref + "{/dir}",
		sourceFile: vcsPath + "/blob/" + ref + "{/dir}/{file}#L{line}",
	}, nil
}

func (v Github) Module() string        { return v.module }
func (v Github) VCSHome() string       { return v.home }
func (v Github) VCSSourceDir() string  { return v.sourceDir }
func (v Github) VCSSourceFile() string { return v.sourceFile }

func (v Github) String() string {
	return v.home + " -> " + v.module
}
