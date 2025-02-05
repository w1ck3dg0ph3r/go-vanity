package forges

import "fmt"

type Forgejo struct {
	host       string
	module     string
	home       string
	sourceDir  string
	sourceFile string
}

func NewForgejo(fs []string) (Vanity, error) {
	if len(fs) != 4 {
		return Forgejo{}, fmt.Errorf(`forgejo vanity format should be "{module} {host} {path} {ref}"`)
	}
	module, host, path, ref := fs[0], fs[1], fs[2], fs[3]
	vcsPath := "https://" + host + "/" + path
	vcsDir := vcsPath + "/src/" + ref
	return Forgejo{
		host:       host,
		module:     module,
		home:       vcsPath,
		sourceDir:  vcsDir + "{/dir}",
		sourceFile: vcsDir + "{/dir}/{file}#L{line}",
	}, nil
}

func (v Forgejo) Module() string        { return v.module }
func (v Forgejo) VCSHome() string       { return v.home }
func (v Forgejo) VCSSourceDir() string  { return v.sourceDir }
func (v Forgejo) VCSSourceFile() string { return v.sourceFile }

func (v Forgejo) String() string {
	return v.home + " -> " + v.module
}
