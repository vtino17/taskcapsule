package config

type Config struct {
	Version  int            `json:"version"`
	Defaults DefaultsConfig `json:"defaults"`
	Setup    []SetupCommand `json:"setup"`
	Services ServiceMap     `json:"services"`
	Checks   CheckMap       `json:"checks"`
}

type DefaultsConfig struct {
	BaseBranch              string `json:"baseBranch"`
	BranchPrefix            string `json:"branchPrefix"`
	GracefulShutdownSeconds int    `json:"gracefulShutdownSeconds"`
	HealthTimeoutSeconds    int    `json:"healthTimeoutSeconds"`
}

type SetupCommand struct {
	Command []string `json:"command"`
}

type ServiceConfig struct {
	Command            []string          `json:"command"`
	WorkingDirectory   string            `json:"workingDirectory"`
	Environment        map[string]string `json:"environment"`
	InheritEnvironment []string          `json:"inheritEnvironment"`
	Health             *HealthConfig     `json:"health"`
}

type HealthConfig struct {
	Type           string `json:"type"`
	URL            string `json:"url"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	ExpectedStatus int    `json:"expectedStatus"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
}

type CheckConfig struct {
	Command []string `json:"command"`
}

type ServiceMap map[string]ServiceConfig

type CheckMap map[string]CheckConfig
