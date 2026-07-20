package spec

import "time"

const ManifestVersion = 2

type ManagedFile struct {
	Path       string `json:"path"`
	SHA256     string `json:"sha256"`
	Customized bool   `json:"customized,omitempty"`
}

type Manifest struct {
	Version     int           `json:"version"`
	Generator   string        `json:"generator"`
	GeneratedAt time.Time     `json:"generatedAt"`
	Config      Config        `json:"config"`
	Files       []ManagedFile `json:"files"`
}
