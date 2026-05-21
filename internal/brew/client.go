// internal/brew/client.go
package brew

import (
	"encoding/json"
	"os/exec"
	"strings"
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

func SearchRemote(query string) ([]PackageInfo, error) {
	if query == "" {
		return nil, nil
	}
	cmd := exec.Command("brew", "search", query)
	out, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	var names []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "==>") {
			continue
		}
		names = append(names, line)
	}

	if len(names) == 0 {
		return nil, nil
	}
	if len(names) > 20 {
		names = names[:20] // limit to avoid slow queries
	}

	args := []string{"info", "--json=v2"}
	args = append(args, names...)
	cmd = exec.Command("brew", args...)
	outInfo, err := cmd.Output()
	if err != nil {
		// some casks might fail if we don't specify --cask, but ignore errors and parse what we can
	}

	var parsed BrewJSONV2
	if err := json.Unmarshal(outInfo, &parsed); err != nil {
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

func Search(query string) ([]PackageInfo, error) {
	// Not used anymore, replaced by SearchRemote
	return SearchRemote(query)
}

func Install(name string, isCask bool) error {
	args := []string{"install"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	
	cmd := exec.Command("brew", args...)
	return cmd.Run()
}

func Uninstall(name string, isCask bool) error {
	args := []string{"uninstall"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	
	cmd := exec.Command("brew", args...)
	return cmd.Run()
}

func Upgrade(name string, isCask bool) error {
	args := []string{"upgrade"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	
	cmd := exec.Command("brew", args...)
	return cmd.Run()
}

func Reinstall(name string, isCask bool) error {
	args := []string{"reinstall"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	
	cmd := exec.Command("brew", args...)
	return cmd.Run()
}

func GetOutdated() ([]PackageInfo, error) {
	cmd := exec.Command("brew", "outdated", "--json=v2")
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