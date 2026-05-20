// internal/brew/client.go
package brew

import (
	"encoding/json"
	"os/exec"
)

func GetInstalled() ([]PackageInfo, error) {
	cmd := exec.Command("brew", "info", "--json=v2", "--installed")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var parsed BrewJSONV2
	if err := json.Unmarshal(out, &parsed); err != nil {
		return nil, err
	}

	var all []PackageInfo
	for _, f := range parsed.Formulae {
		f.IsCask = false
		all = append(all, f)
	}
	for _, c := range parsed.Casks {
		c.IsCask = true
		all = append(all, c)
	}
	return all, nil
}