package workflow

import (
	"errors"

	"github.com/68696c6c/girraph"
)

type DecisionType string

type Decision struct {
	Type DecisionType
}

func MakeDecision(t DecisionType) girraph.Graph[Workflow] {
	result := MakeWorkflow()
	result.GetMeta().SetDecision(&Decision{
		Type: t,
	})
	return result
}

type ConditionType string

type Condition struct {
	Type ConditionType
}

func MakeCondition(t ConditionType) girraph.Graph[Workflow] {
	result := MakeWorkflow()
	result.GetMeta().SetCondition(&Condition{
		Type: t,
	})
	return result
}

type TaskType string

type Task struct {
	Type TaskType
}

func MakeTask(t TaskType) girraph.Graph[Workflow] {
	result := MakeWorkflow()
	result.GetMeta().SetTask(&Task{
		Type: t,
	})
	return result
}

type NodeType string

const (
	ConditionNode NodeType = "condition"
	DecisionNode  NodeType = "decision"
	TaskNode      NodeType = "task"
)

type Workflow interface {
	SetName(string) Workflow
	GetName() string
	SetCondition(*Condition) Workflow
	GetCondition() *Condition
	SetDecision(*Decision) Workflow
	GetDecision() *Decision
	SetTask(*Task) Workflow
	GetTask() *Task
	GetType() NodeType
}

type workflow struct {
	Name      string
	Condition *Condition
	Decision  *Decision
	Task      *Task
	Type      NodeType
}

func MakeWorkflow() girraph.Graph[Workflow] {
	return girraph.MakeGraph[Workflow]().SetMeta(&workflow{})
}

func (n *workflow) SetName(name string) Workflow {
	n.Name = name
	return n
}

func (n *workflow) GetName() string {
	return n.Name
}

func (n *workflow) SetCondition(condition *Condition) Workflow {
	n.Type = ConditionNode
	n.Condition = condition
	n.Decision = nil
	n.Task = nil
	n.SetName(string(condition.Type))
	return n
}

func (n *workflow) GetCondition() *Condition {
	return n.Condition
}

func (n *workflow) SetDecision(decision *Decision) Workflow {
	n.Type = DecisionNode
	n.Condition = nil
	n.Decision = decision
	n.Task = nil
	n.SetName(string(decision.Type))
	return n
}

func (n *workflow) GetDecision() *Decision {
	return n.Decision
}

func (n *workflow) SetTask(task *Task) Workflow {
	n.Type = TaskNode
	n.Condition = nil
	n.Decision = nil
	n.Task = task
	n.SetName(string(task.Type))
	return n
}

func (n *workflow) GetTask() *Task {
	return n.Task
}

func (n *workflow) GetType() NodeType {
	return n.Type
}

func GetConditionalChildren(node girraph.Graph[Workflow], input ConditionType) ([]girraph.Graph[Workflow], error) {
	if node.GetMeta().GetType() != DecisionNode {
		return nil, errors.New("node is not a decision, only decisions have conditional children")
	}
	for _, child := range node.GetChildren() {
		meta := child.GetMeta()
		if meta.GetType() == ConditionNode {
			condition := meta.GetCondition()
			if condition.Type == input {
				return child.GetChildren(), nil
			}
		}
	}
	return []girraph.Graph[Workflow]{}, nil
}
