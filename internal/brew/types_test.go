package brew

import (
	"encoding/json"
	"testing"
)

func TestBrewJSONV2Unmarshal(t *testing.T) {
	jsonData := []byte(`{
		"formulae": [
			{
				"name": "wget",
				"desc": "Internet file retriever",
				"dependencies": ["libidn2", "openssl@3"],
				"versions": {
					"stable": "1.21.4"
				},
				"installed": [
					{
						"version": "1.21.3"
					}
				],
				"outdated": true
			}
		],
		"casks": [
			{
				"token": "google-chrome",
				"name": ["Google Chrome"],
				"desc": "Web browser",
				"version": "120.0.6099.216",
				"installed": "120.0.6099.216",
				"outdated": false
			}
		]
	}`)

	var parsed BrewJSONV2
	err := json.Unmarshal(jsonData, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(parsed.Formulae) != 1 {
		t.Errorf("Expected 1 formula, got %d", len(parsed.Formulae))
	} else {
		f := parsed.Formulae[0]
		if f.Name != "wget" {
			t.Errorf("Expected name 'wget', got '%s'", f.Name)
		}
		if f.IsCask {
			t.Errorf("Expected IsCask to be false")
		}
		if f.LatestVersion != "1.21.4" {
			t.Errorf("Expected LatestVersion '1.21.4', got '%s'", f.LatestVersion)
		}
		if f.CurrentVersion != "1.21.3" {
			t.Errorf("Expected CurrentVersion '1.21.3', got '%s'", f.CurrentVersion)
		}
		if !f.IsOutdated {
			t.Errorf("Expected IsOutdated to be true")
		}
	}

	if len(parsed.Casks) != 1 {
		t.Errorf("Expected 1 cask, got %d", len(parsed.Casks))
	} else {
		c := parsed.Casks[0]
		if c.Name != "google-chrome" {
			t.Errorf("Expected name 'google-chrome', got '%s'", c.Name)
		}
		if !c.IsCask {
			t.Errorf("Expected IsCask to be true")
		}
		if c.LatestVersion != "120.0.6099.216" {
			t.Errorf("Expected LatestVersion '120.0.6099.216', got '%s'", c.LatestVersion)
		}
		if c.CurrentVersion != "120.0.6099.216" {
			t.Errorf("Expected CurrentVersion '120.0.6099.216', got '%s'", c.CurrentVersion)
		}
		if c.IsOutdated {
			t.Errorf("Expected IsOutdated to be false")
		}
	}
}
