package config

func DefaultTemplate() string {
	return `{
  "version": 1,
  "defaults": {
    "baseBranch": "main",
    "branchPrefix": "task/",
    "gracefulShutdownSeconds": 5,
    "healthTimeoutSeconds": 30
  },
  "setup": [],
  "services": {},
  "checks": {}
}
`
}
