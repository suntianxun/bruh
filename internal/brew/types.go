// internal/brew/types.go
package brew

type PackageInfo struct {
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Version      string `json:"version"`
	IsCask       bool
	Dependencies []string `json:"dependencies,omitempty"`
}

type BrewJSONV2 struct {
	Formulae []PackageInfo `json:"formulae"`
	Casks    []PackageInfo `json:"casks"`
}