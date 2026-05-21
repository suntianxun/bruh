// internal/brew/types.go
package brew

import (
	"encoding/json"
)

type PackageInfo struct {
	Name         string `json:"-"` // We'll handle this specially
	Desc         string `json:"desc"`
	Version      string `json:"-"` // We'll handle this specially
	IsCask       bool
	Dependencies []string `json:"dependencies,omitempty"`
}

// Casks have name as an array of strings
type rawCask struct {
	Token   string   `json:"token"`
	Name    []string `json:"name"`
	Desc    string   `json:"desc"`
	Version string   `json:"version"`
}

// Formulae have name as a string, and version is nested
type rawFormula struct {
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Dependencies []string `json:"dependencies,omitempty"`
	Versions     struct {
		Stable string `json:"stable"`
	} `json:"versions"`
}

type rawBrewJSONV2 struct {
	Formulae []rawFormula `json:"formulae"`
	Casks    []rawCask    `json:"casks"`
}

func (b *BrewJSONV2) UnmarshalJSON(data []byte) error {
	var raw rawBrewJSONV2
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, f := range raw.Formulae {
		b.Formulae = append(b.Formulae, PackageInfo{
			Name:         f.Name,
			Desc:         f.Desc,
			Version:      f.Versions.Stable,
			IsCask:       false,
			Dependencies: f.Dependencies,
		})
	}

	for _, c := range raw.Casks {
		name := c.Token
		if len(c.Name) > 0 {
			name = c.Name[0]
		}
		b.Casks = append(b.Casks, PackageInfo{
			Name:    name,
			Desc:    c.Desc,
			Version: c.Version,
			IsCask:  true,
		})
	}

	return nil
}

type BrewJSONV2 struct {
	Formulae []PackageInfo
	Casks    []PackageInfo
}