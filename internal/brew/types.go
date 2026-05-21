// internal/brew/types.go
package brew

import (
	"encoding/json"
)

type PackageInfo struct {
	Name           string   `json:"-"`
	Desc           string   `json:"desc"`
	CurrentVersion string   `json:"-"`
	LatestVersion  string   `json:"-"`
	IsOutdated     bool     `json:"-"`
	IsCask         bool     `json:"-"`
	Dependencies   []string `json:"dependencies,omitempty"`
}

type rawCask struct {
	Token     string   `json:"token"`
	Name      []string `json:"name"`
	Desc      string   `json:"desc"`
	Version   string   `json:"version"`
	Installed string   `json:"installed"`
	Outdated  bool     `json:"outdated"`
}

type rawFormula struct {
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Dependencies []string `json:"dependencies,omitempty"`
	Versions     struct {
		Stable string `json:"stable"`
	} `json:"versions"`
	Installed []struct {
		Version string `json:"version"`
	} `json:"installed"`
	Outdated bool `json:"outdated"`
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
		current := ""
		if len(f.Installed) > 0 {
			current = f.Installed[0].Version
		}
		b.Formulae = append(b.Formulae, PackageInfo{
			Name:           f.Name,
			Desc:           f.Desc,
			CurrentVersion: current,
			LatestVersion:  f.Versions.Stable,
			IsOutdated:     f.Outdated,
			IsCask:         false,
			Dependencies:   f.Dependencies,
		})
	}

	for _, c := range raw.Casks {
		name := c.Token
		b.Casks = append(b.Casks, PackageInfo{
			Name:           name,
			Desc:           c.Desc,
			CurrentVersion: c.Installed,
			LatestVersion:  c.Version,
			IsOutdated:     c.Outdated,
			IsCask:         true,
		})
	}

	return nil
}

type BrewJSONV2 struct {
	Formulae []PackageInfo
	Casks    []PackageInfo
}