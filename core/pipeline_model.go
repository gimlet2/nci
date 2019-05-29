package core

type Pipeline struct {
	Name    string   `yaml:"name"`
	Kind    string   `yaml:"kind"`
	Steps   []Step   `yaml:"steps"`
	Volumes []Volume `yaml:"volumes"`
}

type Step struct {
	Name         string            `yaml:"name"`
	Image        string            `yaml:"image"`
	Commands     []string          `yaml:"commands"`
	When         When              `yaml:"when"`
	Environment  map[string]string `yaml:"environment"`
	VolumeMounts []VolumeMount     `yaml:"volumes"`
}

type When struct {
	Event  []string `yaml:"event"`
	Branch []string `yaml:"branch"`
}

type Volume struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type VolumeMount struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}
