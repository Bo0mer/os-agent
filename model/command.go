package model

type CommandRequest struct {
	Name           string            `json:"name"`
	Args           []string          `json:"args"`
	Env            map[string]string `json:"env"`
	UseIsolatedEnv bool              `json:"use_isolated_env"`
	WorkingDir     string            `json:"working_dir"`
	Input          string            `json:"input"`
}

type CommandResponse struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}
