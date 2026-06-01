// internal/brew/client.go
package brew

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
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

func RunStreaming(args []string, outChan chan<- string, askPass chan<- struct{}, answerPass <-chan string) error {
	cmd := exec.Command("brew", args...)
	cmd.Env = append(os.Environ(), "HOMEBREW_NO_AUTO_UPDATE=1", "HOMEBREW_NO_ENV_HINTS=1")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	defer ptmx.Close()

	go func() {
		buf := make([]byte, 1024)
		var lineBuilder strings.Builder
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				chunk := string(buf[:n])
				lineBuilder.WriteString(chunk)

				str := lineBuilder.String()
				if strings.Contains(str, "Password:") {
					askPass <- struct{}{}
					pass := <-answerPass
					ptmx.Write([]byte(pass + "\n"))
					lineBuilder.Reset() // Clear buffer
					continue
				}

				// Extract lines
				for {
					idx := strings.IndexAny(str, "\r\n")
					if idx == -1 {
						break
					}
					line := str[:idx]
					str = str[idx+1:]
					lineBuilder.Reset()
					lineBuilder.WriteString(str)
					
					trimmed := strings.TrimSpace(line)
					if trimmed != "" && !strings.Contains(trimmed, "Sorry, try again") {
						outChan <- trimmed
					}
				}
			}
			if err != nil {
				break
			}
		}
	}()

	err = cmd.Wait()
	return err
}

func Install(name string, isCask bool, outChan chan<- string, askPass chan<- struct{}, answerPass <-chan string) error {
	args := []string{"install"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	return RunStreaming(args, outChan, askPass, answerPass)
}

func Uninstall(name string, isCask bool, outChan chan<- string, askPass chan<- struct{}, answerPass <-chan string) error {
	args := []string{"uninstall"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	return RunStreaming(args, outChan, askPass, answerPass)
}

func Upgrade(name string, isCask bool, outChan chan<- string, askPass chan<- struct{}, answerPass <-chan string) error {
	args := []string{"upgrade"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	return RunStreaming(args, outChan, askPass, answerPass)
}

func UpgradeAll(outChan chan<- string, askPass chan<- struct{}, answerPass <-chan string) error {
	args := []string{"upgrade"}
	return RunStreaming(args, outChan, askPass, answerPass)
}

func Reinstall(name string, isCask bool, outChan chan<- string, askPass chan<- struct{}, answerPass <-chan string) error {
	args := []string{"reinstall"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	return RunStreaming(args, outChan, askPass, answerPass)
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