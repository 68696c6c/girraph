package workflow

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"

	"github.com/68696c6c/girraph"
)

type Plan interface {
	copy() (Plan, error)
	handleInput(input Input)
	handleDecision(inputType ConditionType, decision girraph.Graph[Workflow]) error
	areNodesComplete(nodes []girraph.Graph[Workflow]) (bool, error)
	addTask(task *Task)
	getTaskState(task *Task) (State, bool)

	JSON() (string, error)
	GetTaskStateByType(taskType TaskType) State
	GetTasks() planTasks
	GetDecisions() planDecisions
	GetInputs() planInputs
}

type plan struct {
	Tasks     planTasks
	Decisions planDecisions
	Inputs    planInputs
}

// A request to update the plan state.
type Input struct {
	Type     ConditionType
	TaskType TaskType
	State    State
}

func GetInitialPlan(workflow girraph.Graph[Workflow]) Plan {
	initialTasks := make(planTasks)
	for _, node := range GetInitialTasks(workflow) {
		meta := node.GetMeta()
		if meta.GetType() == TaskNode {
			task := meta.GetTask()
			if task != nil {
				initialTasks[task.Type] = StateTodo
			}
		}
	}
	return &plan{
		Tasks:     initialTasks,
		Decisions: make(planDecisions),
		Inputs:    planInputs{},
	}
}

func GetUpdatedPlan(currentPlan Plan, workflow girraph.Graph[Workflow], input Input) (Plan, error) {
	newPlan, err := currentPlan.copy()
	if err != nil {
		return nil, err
	}

	// Record the input and update the state of the task in the plan, if it is present.
	newPlan.handleInput(input)

	// Evaluate workflow decisions that have conditions matching the input type.
	workflowDecisions := GetDecisionsByInputType(workflow, input.Type)
	for _, node := range workflowDecisions {
		err = newPlan.handleDecision(input.Type, node)
		if err != nil {
			return nil, err
		}
	}

	// If the input task was completed or cancelled, find the task parent in the workflow and check to see if it is unlocked now.
	// TODO: treating modules the same as tasks means an additional input is needed to mark the module as complete, essentially forcing an external review step.  Confirm desired behavior.
	if stateIsCompletedOrCancelled(input.State) {
		for _, inputTaskParent := range GetTaskParents(workflow, input.TaskType) {
			inputTaskMeta := inputTaskParent.GetMeta()
			parentComplete, err := newPlan.areNodesComplete(inputTaskParent.GetChildren())
			if err != nil {
				return nil, err
			}
			if parentComplete {
				newPlan.addTask(inputTaskMeta.GetTask())
			}
		}
	}

	return newPlan, nil
}

func PlanFromJSON(input []byte) (Plan, error) {
	result := &plan{}
	err := json.Unmarshal(input, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type planTasks map[TaskType]State

type planDecisions map[DecisionType]ConditionType

type planInputs []Input

func (p *plan) JSON() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (p *plan) GetTaskStateByType(taskType TaskType) State {
	return p.Tasks[taskType]
}

func (p *plan) GetTasks() planTasks {
	return p.Tasks
}

func (p *plan) GetDecisions() planDecisions {
	return p.Decisions
}

func (p *plan) GetInputs() planInputs {
	return p.Inputs
}

func (p *plan) copy() (Plan, error) {
	newPlan := &plan{}
	err := copier.Copy(newPlan, p)
	if err != nil {
		return nil, err
	}
	return newPlan, nil
}

func (p *plan) getTaskState(task *Task) (State, bool) {
	if task == nil {
		return StateNil, false
	}
	state, ok := p.Tasks[task.Type]
	if ok {
		return state, true
	}
	return StateNil, false
}

func (p *plan) getDecisionOutcome(decision *Decision) (ConditionType, bool) {
	if decision == nil {
		return "", false
	}
	outcome, ok := p.Decisions[decision.Type]
	if ok {
		return outcome, true
	}
	return "", false
}

func (p *plan) addTask(task *Task) {
	// TODO: validate that we aren't adding a duplicate task.
	if task != nil {
		p.Tasks[task.Type] = StateTodo
	}
}

func (p *plan) handleInput(input Input) {

	// Add the input to the plan input log.
	// TODO: probably need to add some kind of timestamp to the input log.
	p.Inputs = append(p.Inputs, input)

	// If the input task was in the plan, update the task state.
	_, ok := p.Tasks[input.TaskType]
	if ok {
		p.Tasks[input.TaskType] = input.State
	}
}

func (p *plan) handleDecision(inputType ConditionType, node girraph.Graph[Workflow]) error {
	if node == nil {
		return errors.New("nil decision provided")
	}
	conditionNodes, err := GetConditionalChildren(node, inputType)
	if err != nil {
		return err
	}
	if len(conditionNodes) > 0 {
		decision := node.GetMeta().GetDecision()

		// Make sure the decision hasn't already been made.
		_, exists := p.Decisions[decision.Type]
		if exists {
			msg := fmt.Sprintf("duplicate decision: %v", decision.Type)
			return errors.New(msg)
		}

		// Record the decision outcome in the plan.
		p.Decisions[decision.Type] = inputType

		// Add the outcome tasks to the plan.
		for _, conditionNode := range conditionNodes {
			p.addTask(conditionNode.GetMeta().GetTask())
		}
	}
	return nil
}

func (p *plan) areNodesComplete(nodes []girraph.Graph[Workflow]) (bool, error) {
	complete := true

	for _, child := range nodes {
		meta := child.GetMeta()
		switch meta.GetType() {

		// Check the task state in the plan.
		case TaskNode:
			planTaskState, _ := p.getTaskState(meta.GetTask())
			if !stateIsCompletedOrCancelled(planTaskState) {
				complete = false
			}
			break

		// Get the outcome of the plan decision and check their states.
		case DecisionNode:
			childDecision := meta.GetDecision()
			input, _ := p.getDecisionOutcome(childDecision)
			conditionNodes, err := GetConditionalChildren(child, input)
			if err != nil {
				return false, err
			}
			for _, conditionChild := range conditionNodes {
				conditionMeta := conditionChild.GetMeta()
				switch conditionMeta.GetType() {

				// Check the task state in the plan.
				case TaskNode:
					planTaskState, _ := p.getTaskState(conditionMeta.GetTask())
					if !stateIsCompletedOrCancelled(planTaskState) {
						complete = false
					}
					break

				// If the plan has recorded the decision, it's complete.
				case DecisionNode:
					_, exists := p.getDecisionOutcome(conditionMeta.GetDecision())
					if !exists {
						complete = false
					}
					break

				default:
					break
				}
			}
			break

		default:
			break
		}
	}

	return complete, nil
}
