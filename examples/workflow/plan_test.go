package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/68696c6c/girraph"
)

const (
	initialPlanSnapshot      = `{"Tasks":{"Task G":"todo","Task N":"todo"},"Decisions":{},"Inputs":[]}`
	taskCompleteSnapshot     = `{"Tasks":{"Task G":"todo","Task N":"complete"},"Decisions":{},"Inputs":[{"Type":"Task N Resolved","TaskType":"Task N","State":"complete"}]}`
	decisionCompleteSnapshot = `{"Tasks": {"Task G": "complete","Task H": "todo","Task N": "todo"},"Decisions": {"Is Condition A Met?": "Condition A Is Met"},"Inputs": [{"Type": "Condition A Is Met","TaskType": "Task G","State": "complete"}]}`
)

func TestGetInitialPlan_Snapshot(t *testing.T) {
	result := GetInitialPlan(getPlanFixture())

	expected, err := PlanFromJSON([]byte(initialPlanSnapshot))
	require.Nil(t, err)

	assert.Equal(t, expected, result)
}

func TestGetUpdatedPlan_Tasks_Snapshot(t *testing.T) {
	workflow := getPlanFixture()
	initialPlan := GetInitialPlan(workflow)

	result, err := GetUpdatedPlan(initialPlan, workflow, Input{
		Type:     TaskNResolved,
		TaskType: TaskN,
		State:    StateComplete,
	})
	require.Nil(t, err)

	expected, err := PlanFromJSON([]byte(taskCompleteSnapshot))
	require.Nil(t, err)

	assert.Equal(t, expected, result)
}

func TestGetUpdatedPlan_Tasks_StatusChanged(t *testing.T) {
	workflow := getPlanFixture()
	initialPlan := GetInitialPlan(workflow)

	updatedPlan, err := GetUpdatedPlan(initialPlan, workflow, Input{
		Type:     TaskNResolved,
		TaskType: TaskN,
		State:    StateComplete,
	})
	require.Nil(t, err)

	result := updatedPlan.GetTaskStateByType(TaskN)
	assert.Equal(t, StateComplete, result)
}

func TestGetUpdatedPlan_Tasks_UnlockParent(t *testing.T) {
	workflow := getPlanFixture()
	initialPlan := GetInitialPlan(workflow)

	updatedPlan, err := GetUpdatedPlan(initialPlan, workflow, Input{
		Type:     TaskNResolved,
		TaskType: TaskN,
		State:    StateComplete,
	})
	require.Nil(t, err)

	result := updatedPlan.GetTaskStateByType(TaskG)
	assert.Equal(t, StateTodo, result)
}

func TestGetUpdatedPlan_Inputs_Record(t *testing.T) {
	workflow := getPlanFixture()
	initialPlan := GetInitialPlan(workflow)

	updatedPlan, err := GetUpdatedPlan(initialPlan, workflow, Input{
		Type:     TaskNResolved,
		TaskType: TaskN,
		State:    StateComplete,
	})
	require.Nil(t, err)

	assert.Len(t, updatedPlan.GetInputs(), 1)
}

func TestGetUpdatedPlan_Decisions_Snapshot(t *testing.T) {
	workflow := getPlanFixture()
	initialPlan := GetInitialPlan(workflow)

	result, err := GetUpdatedPlan(initialPlan, workflow, Input{
		Type:     ConditionAIsMet,
		TaskType: TaskG,
		State:    StateComplete,
	})
	require.Nil(t, err)

	expected, err := PlanFromJSON([]byte(decisionCompleteSnapshot))
	require.Nil(t, err)

	assert.Equal(t, expected, result)
}

func TestGetUpdatedPlan_Decisions_Record(t *testing.T) {
	workflow := getPlanFixture()
	initialPlan := GetInitialPlan(workflow)

	updatedPlan, err := GetUpdatedPlan(initialPlan, workflow, Input{
		Type:     ConditionAIsMet,
		TaskType: TaskG,
		State:    StateComplete,
	})
	require.Nil(t, err)

	assert.Len(t, updatedPlan.GetDecisions(), 1)
}

func TestGetUpdatedPlan_Decisions_AddTasks(t *testing.T) {
	workflow := getPlanFixture()
	initialPlan := GetInitialPlan(workflow)

	updatedPlan, err := GetUpdatedPlan(initialPlan, workflow, Input{
		Type:     ConditionAIsMet,
		TaskType: TaskG,
		State:    StateComplete,
	})
	require.Nil(t, err)

	result := updatedPlan.GetTaskStateByType(TaskH)
	assert.Equal(t, StateTodo, result)
}

func TestWorkflow_JSON(t *testing.T) {
	graph := getPlanFixture()

	// Convert it to JSON.
	expected, err := graph.JSON()
	require.Nil(t, err)

	// Convert back to a graph.
	graphFromJSON, err := girraph.GraphFromJSON[*workflow](expected)
	assert.Equal(t, string(TaskA), graphFromJSON.GetMeta().GetName())
	require.Nil(t, err)

	// Convert back to JSON.
	result, err := graphFromJSON.JSON()
	require.Nil(t, err)

	// Should match the original JSON.
	assert.Equal(t, expected, result)
}

func getPlanFixture() girraph.Graph[Workflow] {
	multiParentTask := MakeTask(TaskN)
	return MakeTask(TaskA).SetChildren([]girraph.Graph[Workflow]{
		MakeTask(TaskB).SetChildren([]girraph.Graph[Workflow]{
			MakeTask(TaskE).SetChildren([]girraph.Graph[Workflow]{
				MakeTask(TaskF).SetChildren([]girraph.Graph[Workflow]{
					multiParentTask,
				}),
			}),
		}),
		MakeTask(TaskC).SetChildren([]girraph.Graph[Workflow]{
			multiParentTask,
			MakeTask(TaskG),
		}),
		MakeTask(TaskD).SetChildren([]girraph.Graph[Workflow]{
			MakeDecision(IsConditionAMet).SetChildren([]girraph.Graph[Workflow]{
				MakeCondition(ConditionAIsMet).SetChildren([]girraph.Graph[Workflow]{
					MakeTask(TaskH).SetChildren([]girraph.Graph[Workflow]{
						MakeTask(TaskI),
					}),
					MakeDecision(IsConditionBMet).SetChildren([]girraph.Graph[Workflow]{
						MakeCondition(ConditionBIsMet).SetChildren([]girraph.Graph[Workflow]{
							MakeTask(TaskJ),
						}),
					}),
				}),
				MakeCondition(ConditionAIsNotMet).SetChildren([]girraph.Graph[Workflow]{
					MakeDecision(IsConditionCMet).SetChildren([]girraph.Graph[Workflow]{
						MakeCondition(ConditionCIsMet).SetChildren([]girraph.Graph[Workflow]{
							MakeTask(TaskK),
							MakeTask(TaskL).SetChildren([]girraph.Graph[Workflow]{
								MakeTask(TaskM),
							}),
						}),
					}),
				}),
			}),
		}),
	})
}

const (
	IsConditionAMet DecisionType = "Is Condition A Met?"
	IsConditionBMet DecisionType = "Is Condition B Met?"
	IsConditionCMet DecisionType = "Is Condition B Met?"
)

const (
	ConditionAIsMet    ConditionType = "Condition A Is Met"
	ConditionAIsNotMet ConditionType = "Condition A Is Not Met"
	ConditionBIsMet    ConditionType = "Condition B Is Met"
	ConditionCIsMet    ConditionType = "Condition C Is Met"
	TaskNResolved      ConditionType = "Task N Resolved"
)

const (
	TaskA TaskType = "Task A"
	TaskB TaskType = "Task B"
	TaskC TaskType = "Task C"
	TaskD TaskType = "Task D"
	TaskE TaskType = "Task E"
	TaskF TaskType = "Task F"
	TaskG TaskType = "Task G"
	TaskH TaskType = "Task H"
	TaskI TaskType = "Task I"
	TaskJ TaskType = "Task J"
	TaskK TaskType = "Task K"
	TaskL TaskType = "Task L"
	TaskM TaskType = "Task M"
	TaskN TaskType = "Task N"
)
