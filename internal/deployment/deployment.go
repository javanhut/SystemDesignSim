package deployment

import (
	"fmt"
	"time"

	"github.com/javanhut/systemdesignsim/internal/components/config"
)

type Deployment struct {
	ID                string
	Name              string
	ApplicationName   string
	Version           string
	Runtime           *config.Runtime
	Region            string
	AvailabilityZones []string
	Strategy          DeploymentStrategy
	Status            DeploymentStatus
	Instances         []*Instance
	Config            *DeploymentConfig
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeployedAt        *time.Time
}

type DeploymentStrategy string

const (
	DeploymentStrategyRolling   DeploymentStrategy = "rolling"
	DeploymentStrategyBlueGreen DeploymentStrategy = "blue-green"
	DeploymentStrategyCanary    DeploymentStrategy = "canary"
	DeploymentStrategyRecreate  DeploymentStrategy = "recreate"
	DeploymentStrategyAllAtOnce DeploymentStrategy = "all-at-once"
)

type DeploymentStatus string

const (
	DeploymentStatusPending     DeploymentStatus = "pending"
	DeploymentStatusInProgress  DeploymentStatus = "in-progress"
	DeploymentStatusSucceeded   DeploymentStatus = "succeeded"
	DeploymentStatusFailed      DeploymentStatus = "failed"
	DeploymentStatusRollingBack DeploymentStatus = "rolling-back"
	DeploymentStatusRolledBack  DeploymentStatus = "rolled-back"
)

type Instance struct {
	ID               string
	DeploymentID     string
	InstanceType     string
	AvailabilityZone string
	PrivateIP        string
	PublicIP         string
	Status           InstanceStatus
	HealthStatus     HealthStatus
	LaunchedAt       time.Time
	Version          string
}

type InstanceStatus string

const (
	InstanceStatusPending     InstanceStatus = "pending"
	InstanceStatusRunning     InstanceStatus = "running"
	InstanceStatusStopping    InstanceStatus = "stopping"
	InstanceStatusStopped     InstanceStatus = "stopped"
	InstanceStatusTerminating InstanceStatus = "terminating"
	InstanceStatusTerminated  InstanceStatus = "terminated"
)

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

type DeploymentConfig struct {
	MinHealthyInstances    int
	MaxBatchSize           int
	HealthCheckGracePeriod time.Duration
	RollbackOnFailure      bool
	AutoScaling            *AutoScalingConfig
	EnvironmentVariables   map[string]string
	CodeSource             *CodeSource
	HealthCheck            *HealthCheckConfig
}

type AutoScalingConfig struct {
	Enabled            bool
	MinInstances       int
	MaxInstances       int
	DesiredInstances   int
	ScaleUpThreshold   float64
	ScaleDownThreshold float64
	CooldownPeriod     time.Duration
}

type CodeSource struct {
	Type       CodeSourceType
	Repository string
	Branch     string
	Tag        string
	Path       string
}

type CodeSourceType string

const (
	CodeSourceTypeGit       CodeSourceType = "git"
	CodeSourceTypeS3        CodeSourceType = "s3"
	CodeSourceTypeContainer CodeSourceType = "container"
	CodeSourceTypeZip       CodeSourceType = "zip"
)

type HealthCheckConfig struct {
	Protocol           string
	Port               int
	Path               string
	IntervalSeconds    int
	TimeoutSeconds     int
	HealthyThreshold   int
	UnhealthyThreshold int
}

func NewDeployment(id, name, appName, version string, runtime *config.Runtime, region string) *Deployment {
	return &Deployment{
		ID:              id,
		Name:            name,
		ApplicationName: appName,
		Version:         version,
		Runtime:         runtime,
		Region:          region,
		Strategy:        DeploymentStrategyRolling,
		Status:          DeploymentStatusPending,
		Instances:       []*Instance{},
		Config:          NewDefaultDeploymentConfig(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func NewDefaultDeploymentConfig() *DeploymentConfig {
	return &DeploymentConfig{
		MinHealthyInstances:    1,
		MaxBatchSize:           1,
		HealthCheckGracePeriod: 60 * time.Second,
		RollbackOnFailure:      true,
		AutoScaling: &AutoScalingConfig{
			Enabled:            false,
			MinInstances:       1,
			MaxInstances:       10,
			DesiredInstances:   2,
			ScaleUpThreshold:   70.0,
			ScaleDownThreshold: 30.0,
			CooldownPeriod:     300 * time.Second,
		},
		EnvironmentVariables: make(map[string]string),
		CodeSource: &CodeSource{
			Type:       CodeSourceTypeGit,
			Repository: "",
			Branch:     "main",
		},
		HealthCheck: &HealthCheckConfig{
			Protocol:           "HTTP",
			Port:               8080,
			Path:               "/health",
			IntervalSeconds:    30,
			TimeoutSeconds:     5,
			HealthyThreshold:   2,
			UnhealthyThreshold: 3,
		},
	}
}

func (d *Deployment) AddInstance(instanceType, az string) *Instance {
	instance := &Instance{
		ID:               fmt.Sprintf("i-%s-%d", d.ID, len(d.Instances)),
		DeploymentID:     d.ID,
		InstanceType:     instanceType,
		AvailabilityZone: az,
		Status:           InstanceStatusPending,
		HealthStatus:     HealthStatusUnknown,
		LaunchedAt:       time.Now(),
		Version:          d.Version,
	}
	d.Instances = append(d.Instances, instance)
	return instance
}

func (d *Deployment) GetHealthyInstances() []*Instance {
	healthy := []*Instance{}
	for _, instance := range d.Instances {
		if instance.HealthStatus == HealthStatusHealthy {
			healthy = append(healthy, instance)
		}
	}
	return healthy
}

func (d *Deployment) GetRunningInstances() []*Instance {
	running := []*Instance{}
	for _, instance := range d.Instances {
		if instance.Status == InstanceStatusRunning {
			running = append(running, instance)
		}
	}
	return running
}

func (d *Deployment) Start() error {
	if d.Status != DeploymentStatusPending {
		return fmt.Errorf("deployment must be in pending state to start")
	}

	d.Status = DeploymentStatusInProgress
	d.UpdatedAt = time.Now()

	return nil
}

func (d *Deployment) Complete() {
	d.Status = DeploymentStatusSucceeded
	now := time.Now()
	d.DeployedAt = &now
	d.UpdatedAt = now
}

func (d *Deployment) Fail(reason string) {
	d.Status = DeploymentStatusFailed
	d.UpdatedAt = time.Now()
}

func (d *Deployment) Rollback() error {
	if d.Status != DeploymentStatusFailed {
		return fmt.Errorf("can only rollback failed deployments")
	}

	d.Status = DeploymentStatusRollingBack
	d.UpdatedAt = time.Now()

	return nil
}

func (instance *Instance) Start() {
	instance.Status = InstanceStatusRunning
	instance.PrivateIP = fmt.Sprintf("10.0.%d.%d",
		(len(instance.ID)%254)+1,
		(len(instance.DeploymentID)%254)+1)
}

func (instance *Instance) Stop() {
	instance.Status = InstanceStatusStopping
}

func (instance *Instance) Terminate() {
	instance.Status = InstanceStatusTerminated
}

func (instance *Instance) UpdateHealth(healthy bool) {
	if healthy {
		instance.HealthStatus = HealthStatusHealthy
	} else {
		instance.HealthStatus = HealthStatusUnhealthy
	}
}

type DeploymentPipeline struct {
	ID          string
	Name        string
	Stages      []*PipelineStage
	Deployments []*Deployment
	Status      PipelineStatus
}

type PipelineStage struct {
	Name          string
	Type          StageType
	Configuration map[string]string
	Status        StageStatus
}

type StageType string

const (
	StageTypeSource   StageType = "source"
	StageTypeBuild    StageType = "build"
	StageTypeTest     StageType = "test"
	StageTypeDeploy   StageType = "deploy"
	StageTypeApproval StageType = "approval"
)

type StageStatus string

const (
	StageStatusPending   StageStatus = "pending"
	StageStatusRunning   StageStatus = "running"
	StageStatusSucceeded StageStatus = "succeeded"
	StageStatusFailed    StageStatus = "failed"
	StageStatusSkipped   StageStatus = "skipped"
)

type PipelineStatus string

const (
	PipelineStatusIdle      PipelineStatus = "idle"
	PipelineStatusRunning   PipelineStatus = "running"
	PipelineStatusSucceeded PipelineStatus = "succeeded"
	PipelineStatusFailed    PipelineStatus = "failed"
)

func NewDeploymentPipeline(id, name string) *DeploymentPipeline {
	return &DeploymentPipeline{
		ID:          id,
		Name:        name,
		Stages:      []*PipelineStage{},
		Deployments: []*Deployment{},
		Status:      PipelineStatusIdle,
	}
}

func (p *DeploymentPipeline) AddStage(name string, stageType StageType) *PipelineStage {
	stage := &PipelineStage{
		Name:          name,
		Type:          stageType,
		Configuration: make(map[string]string),
		Status:        StageStatusPending,
	}
	p.Stages = append(p.Stages, stage)
	return stage
}

func (p *DeploymentPipeline) Execute() error {
	p.Status = PipelineStatusRunning

	for _, stage := range p.Stages {
		stage.Status = StageStatusRunning

		time.Sleep(100 * time.Millisecond)

		stage.Status = StageStatusSucceeded
	}

	p.Status = PipelineStatusSucceeded
	return nil
}

type DeploymentPreset struct {
	Name        string
	Description string
	Strategy    DeploymentStrategy
	Config      *DeploymentConfig
}

var DeploymentPresets = map[string]*DeploymentPreset{
	"simple": {
		Name:        "Simple Deployment",
		Description: "Single instance deployment for development",
		Strategy:    DeploymentStrategyAllAtOnce,
		Config: &DeploymentConfig{
			MinHealthyInstances:    0,
			MaxBatchSize:           1,
			HealthCheckGracePeriod: 30 * time.Second,
			RollbackOnFailure:      false,
			AutoScaling: &AutoScalingConfig{
				Enabled:          false,
				MinInstances:     1,
				MaxInstances:     1,
				DesiredInstances: 1,
			},
			EnvironmentVariables: make(map[string]string),
		},
	},
	"rolling": {
		Name:        "Rolling Deployment",
		Description: "Zero-downtime rolling deployment",
		Strategy:    DeploymentStrategyRolling,
		Config: &DeploymentConfig{
			MinHealthyInstances:    1,
			MaxBatchSize:           1,
			HealthCheckGracePeriod: 60 * time.Second,
			RollbackOnFailure:      true,
			AutoScaling: &AutoScalingConfig{
				Enabled:            true,
				MinInstances:       2,
				MaxInstances:       10,
				DesiredInstances:   2,
				ScaleUpThreshold:   70.0,
				ScaleDownThreshold: 30.0,
			},
			EnvironmentVariables: make(map[string]string),
		},
	},
	"blue-green": {
		Name:        "Blue/Green Deployment",
		Description: "Deploy to new environment, switch traffic instantly",
		Strategy:    DeploymentStrategyBlueGreen,
		Config: &DeploymentConfig{
			MinHealthyInstances:    2,
			MaxBatchSize:           100,
			HealthCheckGracePeriod: 120 * time.Second,
			RollbackOnFailure:      true,
			AutoScaling: &AutoScalingConfig{
				Enabled:            true,
				MinInstances:       2,
				MaxInstances:       20,
				DesiredInstances:   4,
				ScaleUpThreshold:   70.0,
				ScaleDownThreshold: 30.0,
			},
			EnvironmentVariables: make(map[string]string),
		},
	},
	"canary": {
		Name:        "Canary Deployment",
		Description: "Gradual rollout with traffic shifting",
		Strategy:    DeploymentStrategyCanary,
		Config: &DeploymentConfig{
			MinHealthyInstances:    3,
			MaxBatchSize:           1,
			HealthCheckGracePeriod: 180 * time.Second,
			RollbackOnFailure:      true,
			AutoScaling: &AutoScalingConfig{
				Enabled:            true,
				MinInstances:       3,
				MaxInstances:       20,
				DesiredInstances:   5,
				ScaleUpThreshold:   60.0,
				ScaleDownThreshold: 20.0,
			},
			EnvironmentVariables: make(map[string]string),
		},
	},
}

func GetDeploymentPreset(name string) *DeploymentPreset {
	if preset, exists := DeploymentPresets[name]; exists {
		return preset
	}
	return DeploymentPresets["simple"]
}

func GetDeploymentPresetNames() []string {
	return []string{"simple", "rolling", "blue-green", "canary"}
}

func CreateDeploymentFromPreset(id, name, appName, version string, runtime *config.Runtime, region, presetName string) *Deployment {
	preset := GetDeploymentPreset(presetName)

	deployment := NewDeployment(id, name, appName, version, runtime, region)
	deployment.Strategy = preset.Strategy
	deployment.Config = preset.Config

	return deployment
}

type DeploymentMetrics struct {
	DeploymentID       string
	TotalInstances     int
	HealthyInstances   int
	UnhealthyInstances int
	AverageCPU         float64
	AverageMemory      float64
	RequestsPerSecond  int
	ErrorRate          float64
	Timestamp          time.Time
}

func (d *Deployment) GetMetrics() *DeploymentMetrics {
	healthy := 0
	unhealthy := 0

	for _, instance := range d.Instances {
		if instance.HealthStatus == HealthStatusHealthy {
			healthy++
		} else if instance.HealthStatus == HealthStatusUnhealthy {
			unhealthy++
		}
	}

	return &DeploymentMetrics{
		DeploymentID:       d.ID,
		TotalInstances:     len(d.Instances),
		HealthyInstances:   healthy,
		UnhealthyInstances: unhealthy,
		Timestamp:          time.Now(),
	}
}
