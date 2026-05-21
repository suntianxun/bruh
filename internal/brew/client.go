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

func Search(query string) ([]PackageInfo, error) {
	if query == "" {
		return nil, nil
	}
	// Brew search returns simple text lines, not json
	cmd := exec.Command("brew", "search", query)
	out, err := cmd.Output()
	if err != nil {
		// brew search returns non-zero if no results found
		return nil, nil
	}

	var results []PackageInfo
	lines := strings.Split(string(out), "\n")
	
	isCaskSection := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "==> Formulae" {
			isCaskSection = false
			continue
		}
		if line == "==> Casks" {
			isCaskSection = true
			continue
		}
		
		// It's a package name
		results = append(results, PackageInfo{
			Name: line,
			Desc: "Press enter to view details", // We don't get descriptions from basic search
			IsCask: isCaskSection,
		})
	}
	return results, nil
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