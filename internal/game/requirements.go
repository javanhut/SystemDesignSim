package game

import (
	"fmt"
	"time"
)

type TaskStatus string

const (
	TaskStatusNotStarted TaskStatus = "not_started"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusSkipped    TaskStatus = "skipped"
)

type TaskProgress struct {
	Task        *Task
	Status      TaskStatus
	StartedAt   *time.Time
	CompletedAt *time.Time
	Notes       string
}

type RequirementTracker struct {
	Scenario         *Scenario
	TaskProgress     map[int]*TaskProgress
	BonusesCompleted []string
	StartTime        time.Time
	EndTime          *time.Time
}

func NewRequirementTracker(scenario *Scenario) *RequirementTracker {
	tracker := &RequirementTracker{
		Scenario:         scenario,
		TaskProgress:     make(map[int]*TaskProgress),
		BonusesCompleted: []string{},
		StartTime:        time.Now(),
	}

	for i := range scenario.Tasks {
		task := &scenario.Tasks[i]
		tracker.TaskProgress[task.Step] = &TaskProgress{
			Task:   task,
			Status: TaskStatusNotStarted,
		}
	}

	return tracker
}

func (rt *RequirementTracker) StartTask(step int) error {
	progress, exists := rt.TaskProgress[step]
	if !exists {
		return fmt.Errorf("task %d does not exist", step)
	}

	if progress.Status == TaskStatusCompleted {
		return fmt.Errorf("task %d is already completed", step)
	}

	now := time.Now()
	progress.Status = TaskStatusInProgress
	progress.StartedAt = &now

	return nil
}

func (rt *RequirementTracker) CompleteTask(step int, notes string) error {
	progress, exists := rt.TaskProgress[step]
	if !exists {
		return fmt.Errorf("task %d does not exist", step)
	}

	if progress.Status == TaskStatusCompleted {
		return fmt.Errorf("task %d is already completed", step)
	}

	now := time.Now()
	progress.Status = TaskStatusCompleted
	progress.CompletedAt = &now
	progress.Notes = notes

	return nil
}

func (rt *RequirementTracker) SkipTask(step int, reason string) error {
	progress, exists := rt.TaskProgress[step]
	if !exists {
		return fmt.Errorf("task %d does not exist", step)
	}

	if progress.Task.Mandatory {
		return fmt.Errorf("task %d is mandatory and cannot be skipped", step)
	}

	progress.Status = TaskStatusSkipped
	progress.Notes = reason

	return nil
}

func (rt *RequirementTracker) MarkBonusCompleted(bonusObjective string) {
	rt.BonusesCompleted = append(rt.BonusesCompleted, bonusObjective)
}

func (rt *RequirementTracker) GetTaskStatus(step int) TaskStatus {
	if progress, exists := rt.TaskProgress[step]; exists {
		return progress.Status
	}
	return TaskStatusNotStarted
}

func (rt *RequirementTracker) GetCompletedTasks() []*TaskProgress {
	completed := []*TaskProgress{}
	for _, progress := range rt.TaskProgress {
		if progress.Status == TaskStatusCompleted {
			completed = append(completed, progress)
		}
	}
	return completed
}

func (rt *RequirementTracker) GetPendingTasks() []*TaskProgress {
	pending := []*TaskProgress{}
	for _, progress := range rt.TaskProgress {
		if progress.Status == TaskStatusNotStarted || progress.Status == TaskStatusInProgress {
			pending = append(pending, progress)
		}
	}
	return pending
}

func (rt *RequirementTracker) GetNextTask() *TaskProgress {
	for i := 1; i <= len(rt.Scenario.Tasks); i++ {
		if progress, exists := rt.TaskProgress[i]; exists {
			if progress.Status == TaskStatusNotStarted {
				return progress
			}
		}
	}
	return nil
}

func (rt *RequirementTracker) GetProgress() (completed, total int) {
	total = 0
	completed = 0

	for _, progress := range rt.TaskProgress {
		if progress.Task.Mandatory {
			total++
			if progress.Status == TaskStatusCompleted {
				completed++
			}
		}
	}

	return completed, total
}

func (rt *RequirementTracker) AreAllMandatoryTasksCompleted() bool {
	for _, progress := range rt.TaskProgress {
		if progress.Task.Mandatory && progress.Status != TaskStatusCompleted {
			return false
		}
	}
	return true
}

func (rt *RequirementTracker) GetTimeSpent() time.Duration {
	if rt.EndTime != nil {
		return rt.EndTime.Sub(rt.StartTime)
	}
	return time.Since(rt.StartTime)
}

func (rt *RequirementTracker) Finish() {
	now := time.Now()
	rt.EndTime = &now
}

func (rt *RequirementTracker) GenerateReport() string {
	report := fmt.Sprintf("Requirement Tracking Report\n")
	report += fmt.Sprintf("==========================\n\n")
	report += fmt.Sprintf("Scenario: %s\n", rt.Scenario.CustomerName)
	report += fmt.Sprintf("Started: %s\n", rt.StartTime.Format(time.RFC3339))

	if rt.EndTime != nil {
		report += fmt.Sprintf("Completed: %s\n", rt.EndTime.Format(time.RFC3339))
		report += fmt.Sprintf("Duration: %s\n\n", rt.GetTimeSpent())
	} else {
		report += fmt.Sprintf("Duration: %s (ongoing)\n\n", rt.GetTimeSpent())
	}

	completed, total := rt.GetProgress()
	report += fmt.Sprintf("Mandatory Tasks: %d/%d completed (%.1f%%)\n\n",
		completed, total, float64(completed)/float64(total)*100)

	report += "Task Details:\n"
	for i := 1; i <= len(rt.Scenario.Tasks); i++ {
		if progress, exists := rt.TaskProgress[i]; exists {
			mandatoryStr := ""
			if progress.Task.Mandatory {
				mandatoryStr = " [MANDATORY]"
			}

			statusSymbol := ""
			switch progress.Status {
			case TaskStatusCompleted:
				statusSymbol = "[✓]"
			case TaskStatusInProgress:
				statusSymbol = "[▶]"
			case TaskStatusSkipped:
				statusSymbol = "[⊘]"
			case TaskStatusNotStarted:
				statusSymbol = "[ ]"
			}

			report += fmt.Sprintf("%s Step %d: %s%s\n",
				statusSymbol, progress.Task.Step, progress.Task.Title, mandatoryStr)

			if progress.Status == TaskStatusCompleted && progress.CompletedAt != nil {
				report += fmt.Sprintf("    Completed at: %s\n",
					progress.CompletedAt.Format(time.RFC3339))
			}

			if progress.Notes != "" {
				report += fmt.Sprintf("    Notes: %s\n", progress.Notes)
			}
		}
	}

	if len(rt.BonusesCompleted) > 0 {
		report += "\nBonus Objectives Completed:\n"
		for _, bonus := range rt.BonusesCompleted {
			report += fmt.Sprintf("  ✓ %s\n", bonus)
		}
	}

	return report
}

type RequirementValidator struct {
	tracker *RequirementTracker
}

func NewRequirementValidator(tracker *RequirementTracker) *RequirementValidator {
	return &RequirementValidator{
		tracker: tracker,
	}
}

func (rv *RequirementValidator) ValidateCompletion() (bool, []string) {
	errors := []string{}

	if !rv.tracker.AreAllMandatoryTasksCompleted() {
		errors = append(errors, "Not all mandatory tasks are completed")

		for _, progress := range rv.tracker.GetPendingTasks() {
			if progress.Task.Mandatory {
				errors = append(errors,
					fmt.Sprintf("  - Task %d: %s is not completed",
						progress.Task.Step, progress.Task.Title))
			}
		}
	}

	return len(errors) == 0, errors
}

func (rv *RequirementValidator) ValidateTask(step int, architecture map[string]interface{}) (bool, []string) {
	progress, exists := rv.tracker.TaskProgress[step]
	if !exists {
		return false, []string{"Task does not exist"}
	}

	warnings := []string{}

	switch progress.Task.Type {
	case TaskTypeInfrastructure:
		warnings = rv.validateInfrastructure(progress.Task, architecture)
	case TaskTypeNetworking:
		warnings = rv.validateNetworking(progress.Task, architecture)
	case TaskTypeDeployment:
		warnings = rv.validateDeployment(progress.Task, architecture)
	case TaskTypeConfiguration:
		warnings = rv.validateConfiguration(progress.Task, architecture)
	case TaskTypeOptimization:
		warnings = rv.validateOptimization(progress.Task, architecture)
	}

	return len(warnings) == 0, warnings
}

func (rv *RequirementValidator) validateInfrastructure(task *Task, architecture map[string]interface{}) []string {
	warnings := []string{}

	switch task.Title {
	case "Deploy Application Server", "Deploy Multiple API Servers":
		if components, ok := architecture["api_servers"].([]interface{}); !ok || len(components) == 0 {
			warnings = append(warnings, "No API servers deployed")
		}

	case "Configure Database":
		if components, ok := architecture["databases"].([]interface{}); !ok || len(components) == 0 {
			warnings = append(warnings, "No database configured")
		}

	case "Add Load Balancer":
		if components, ok := architecture["load_balancers"].([]interface{}); !ok || len(components) == 0 {
			warnings = append(warnings, "No load balancer deployed")
		}

	case "Add Caching Layer":
		if components, ok := architecture["caches"].([]interface{}); !ok || len(components) == 0 {
			warnings = append(warnings, "No cache layer deployed")
		}

	case "Implement CDN", "Deploy Global CDN":
		if components, ok := architecture["cdns"].([]interface{}); !ok || len(components) == 0 {
			warnings = append(warnings, "No CDN deployed")
		}
	}

	return warnings
}

func (rv *RequirementValidator) validateNetworking(task *Task, architecture map[string]interface{}) []string {
	warnings := []string{}

	return warnings
}

func (rv *RequirementValidator) validateDeployment(task *Task, architecture map[string]interface{}) []string {
	warnings := []string{}

	return warnings
}

func (rv *RequirementValidator) validateConfiguration(task *Task, architecture map[string]interface{}) []string {
	warnings := []string{}

	return warnings
}

func (rv *RequirementValidator) validateOptimization(task *Task, architecture map[string]interface{}) []string {
	warnings := []string{}

	if budget, ok := architecture["total_cost"].(float64); ok {
		if maxBudget, ok := architecture["max_budget"].(float64); ok {
			if budget > maxBudget {
				warnings = append(warnings,
					fmt.Sprintf("Cost $%.2f exceeds budget $%.2f", budget, maxBudget))
			}
		}
	}

	return warnings
}
