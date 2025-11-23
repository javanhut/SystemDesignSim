package config

type Runtime struct {
	ID             string
	Name           string
	DisplayName    string
	Language       string
	Version        string
	MemoryOverhead int
	StartupTimeMs  int
	Description    string
	DefaultPort    int
	SupportedPorts []int
}

var Runtimes = map[string]*Runtime{
	"nodejs-18": {
		ID:             "nodejs-18",
		Name:           "nodejs",
		DisplayName:    "Node.js 18 LTS",
		Language:       "JavaScript",
		Version:        "18.x",
		MemoryOverhead: 128,
		StartupTimeMs:  1000,
		Description:    "JavaScript runtime - Fast, event-driven",
		DefaultPort:    3000,
		SupportedPorts: []int{3000, 8080, 8000},
	},
	"nodejs-20": {
		ID:             "nodejs-20",
		Name:           "nodejs",
		DisplayName:    "Node.js 20 LTS",
		Language:       "JavaScript",
		Version:        "20.x",
		MemoryOverhead: 128,
		StartupTimeMs:  900,
		Description:    "Latest Node.js - Better performance",
		DefaultPort:    3000,
		SupportedPorts: []int{3000, 8080, 8000},
	},
	"python-39": {
		ID:             "python-39",
		Name:           "python",
		DisplayName:    "Python 3.9",
		Language:       "Python",
		Version:        "3.9",
		MemoryOverhead: 256,
		StartupTimeMs:  1500,
		Description:    "Python runtime - Great for data processing",
		DefaultPort:    8000,
		SupportedPorts: []int{8000, 8080, 5000},
	},
	"python-311": {
		ID:             "python-311",
		Name:           "python",
		DisplayName:    "Python 3.11",
		Language:       "Python",
		Version:        "3.11",
		MemoryOverhead: 256,
		StartupTimeMs:  1300,
		Description:    "Latest Python - 25% faster than 3.9",
		DefaultPort:    8000,
		SupportedPorts: []int{8000, 8080, 5000},
	},
	"go-120": {
		ID:             "go-120",
		Name:           "go",
		DisplayName:    "Go 1.20",
		Language:       "Go",
		Version:        "1.20",
		MemoryOverhead: 64,
		StartupTimeMs:  500,
		Description:    "Go runtime - Fast, compiled, low memory",
		DefaultPort:    8080,
		SupportedPorts: []int{8080, 8000, 3000},
	},
	"go-121": {
		ID:             "go-121",
		Name:           "go",
		DisplayName:    "Go 1.21",
		Language:       "Go",
		Version:        "1.21",
		MemoryOverhead: 64,
		StartupTimeMs:  450,
		Description:    "Latest Go - Improved performance",
		DefaultPort:    8080,
		SupportedPorts: []int{8080, 8000, 3000},
	},
	"java-17": {
		ID:             "java-17",
		Name:           "java",
		DisplayName:    "Java 17 LTS",
		Language:       "Java",
		Version:        "17",
		MemoryOverhead: 512,
		StartupTimeMs:  3000,
		Description:    "Java runtime - Enterprise-grade",
		DefaultPort:    8080,
		SupportedPorts: []int{8080, 8000, 9000},
	},
	"java-21": {
		ID:             "java-21",
		Name:           "java",
		DisplayName:    "Java 21 LTS",
		Language:       "Java",
		Version:        "21",
		MemoryOverhead: 512,
		StartupTimeMs:  2500,
		Description:    "Latest Java LTS - Virtual threads",
		DefaultPort:    8080,
		SupportedPorts: []int{8080, 8000, 9000},
	},
	"dotnet-6": {
		ID:             "dotnet-6",
		Name:           "dotnet",
		DisplayName:    ".NET 6",
		Language:       "C#",
		Version:        "6.0",
		MemoryOverhead: 384,
		StartupTimeMs:  2000,
		Description:    ".NET runtime - Cross-platform",
		DefaultPort:    5000,
		SupportedPorts: []int{5000, 8080, 8000},
	},
	"dotnet-8": {
		ID:             "dotnet-8",
		Name:           "dotnet",
		DisplayName:    ".NET 8",
		Language:       "C#",
		Version:        "8.0",
		MemoryOverhead: 384,
		StartupTimeMs:  1800,
		Description:    "Latest .NET - Better performance",
		DefaultPort:    5000,
		SupportedPorts: []int{5000, 8080, 8000},
	},
	"ruby-32": {
		ID:             "ruby-32",
		Name:           "ruby",
		DisplayName:    "Ruby 3.2",
		Language:       "Ruby",
		Version:        "3.2",
		MemoryOverhead: 256,
		StartupTimeMs:  2000,
		Description:    "Ruby runtime - Developer-friendly",
		DefaultPort:    3000,
		SupportedPorts: []int{3000, 8080, 9292},
	},
	"php-82": {
		ID:             "php-82",
		Name:           "php",
		DisplayName:    "PHP 8.2",
		Language:       "PHP",
		Version:        "8.2",
		MemoryOverhead: 128,
		StartupTimeMs:  800,
		Description:    "PHP runtime - Web-focused",
		DefaultPort:    8000,
		SupportedPorts: []int{8000, 8080, 9000},
	},
	"rust-173": {
		ID:             "rust-173",
		Name:           "rust",
		DisplayName:    "Rust 1.73",
		Language:       "Rust",
		Version:        "1.73",
		MemoryOverhead: 32,
		StartupTimeMs:  300,
		Description:    "Rust runtime - Extremely fast, safe",
		DefaultPort:    8080,
		SupportedPorts: []int{8080, 8000, 3000},
	},
}

func GetRuntime(id string) *Runtime {
	if rt, exists := Runtimes[id]; exists {
		return rt
	}
	return Runtimes["nodejs-18"]
}

func GetRuntimesByLanguage(language string) []*Runtime {
	runtimes := []*Runtime{}
	for _, rt := range Runtimes {
		if rt.Language == language {
			runtimes = append(runtimes, rt)
		}
	}
	return runtimes
}

func GetAllRuntimes() []*Runtime {
	runtimes := []*Runtime{}
	for _, rt := range Runtimes {
		runtimes = append(runtimes, rt)
	}
	return runtimes
}

func GetRuntimeIDs() []string {
	ids := []string{
		"nodejs-18", "nodejs-20",
		"python-39", "python-311",
		"go-120", "go-121",
		"java-17", "java-21",
		"dotnet-6", "dotnet-8",
		"ruby-32", "php-82", "rust-173",
	}
	return ids
}

func GetRuntimeLanguages() []string {
	return []string{
		"JavaScript", "Python", "Go", "Java",
		"C#", "Ruby", "PHP", "Rust",
	}
}

type DeploymentStrategy string

const (
	DeploymentRolling   DeploymentStrategy = "rolling"
	DeploymentBlueGreen DeploymentStrategy = "blue-green"
	DeploymentCanary    DeploymentStrategy = "canary"
	DeploymentRecreate  DeploymentStrategy = "recreate"
)

type DeploymentConfig struct {
	Strategy          DeploymentStrategy
	HealthCheckPath   string
	HealthCheckPort   int
	MinHealthyPercent int
	MaxSurge          int
	RollbackOnFailure bool
}

var DefaultDeploymentConfig = &DeploymentConfig{
	Strategy:          DeploymentRolling,
	HealthCheckPath:   "/health",
	HealthCheckPort:   8080,
	MinHealthyPercent: 50,
	MaxSurge:          100,
	RollbackOnFailure: true,
}

func GetDeploymentStrategies() []DeploymentStrategy {
	return []DeploymentStrategy{
		DeploymentRolling,
		DeploymentBlueGreen,
		DeploymentCanary,
		DeploymentRecreate,
	}
}

type EnvironmentVariable struct {
	Key   string
	Value string
}

type ApplicationConfig struct {
	Runtime          *Runtime
	Port             int
	EnvironmentVars  []EnvironmentVariable
	DeploymentConfig *DeploymentConfig
	AutoScaling      *AutoScalingConfig
	HealthCheck      *HealthCheckConfig
}

type AutoScalingConfig struct {
	Enabled           bool
	MinInstances      int
	MaxInstances      int
	TargetCPU         int
	TargetMemory      int
	ScaleUpCooldown   int
	ScaleDownCooldown int
}

var DefaultAutoScalingConfig = &AutoScalingConfig{
	Enabled:           false,
	MinInstances:      1,
	MaxInstances:      10,
	TargetCPU:         70,
	TargetMemory:      80,
	ScaleUpCooldown:   60,
	ScaleDownCooldown: 300,
}

type HealthCheckConfig struct {
	Enabled            bool
	Path               string
	Port               int
	IntervalSeconds    int
	TimeoutSeconds     int
	HealthyThreshold   int
	UnhealthyThreshold int
}

var DefaultHealthCheckConfig = &HealthCheckConfig{
	Enabled:            true,
	Path:               "/health",
	Port:               8080,
	IntervalSeconds:    30,
	TimeoutSeconds:     5,
	HealthyThreshold:   2,
	UnhealthyThreshold: 3,
}

func NewApplicationConfig(runtimeID string) *ApplicationConfig {
	runtime := GetRuntime(runtimeID)
	return &ApplicationConfig{
		Runtime:          runtime,
		Port:             runtime.DefaultPort,
		EnvironmentVars:  []EnvironmentVariable{},
		DeploymentConfig: DefaultDeploymentConfig,
		AutoScaling:      DefaultAutoScalingConfig,
		HealthCheck:      DefaultHealthCheckConfig,
	}
}
