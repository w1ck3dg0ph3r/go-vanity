package forges

import "fmt"

type Gitlab struct {
	host       string
	module     string
	home       string
	sourceDir  string
	sourceFile string
}

func NewGitlab(fs []string) (Vanity, error) {
	if len(fs) != 3 {
		return Gitlab{}, fmt.Errorf(`gitlab vanity format should be "{module} {path} {ref}"`)
	}
	module, path, ref := fs[0], fs[1], fs[2]
	host := "https://gitlab.com/"
	vcsPath := host + path
	return Gitlab{
		host:       host,
		module:     module,
		home:       vcsPath,
		sourceDir:  vcsPath + "/-/tree/" + ref + "{/dir}?ref_type=heads",
		sourceFile: vcsPath + "/-/blob/" + ref + "{/dir}/{file}?ref_type=heads#L{line}",
	}, nil
}

func (v Gitlab) Module() string        { return v.module }
func (v Gitlab) VCSHome() string       { return v.home }
func (v Gitlab) VCSSourceDir() string  { return v.sourceDir }
func (v Gitlab) VCSSourceFile() string { return v.sourceFile }

func (v Gitlab) String() string {
	return v.home + " -> " + v.module
}
