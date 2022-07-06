package workflow

import "github.com/68696c6c/girraph"

// Find all task children of the provided node that are not locked behind a decision and have no pre-requisites tasks.
func GetInitialTasks(workflow girraph.Graph[Workflow]) []girraph.Graph[Workflow] {
	var query girraph.QueryResult[girraph.Graph[Workflow]]
	findInitialTasks(&query, workflow)
	return query.GetNodes()
}

func GetDecisionsByInputType(workflow girraph.Graph[Workflow], inputType ConditionType) []girraph.Graph[Workflow] {
	var found girraph.QueryResult[girraph.Graph[Workflow]]
	queryDecisionsByInputType(&found, workflow, inputType)
	return found.GetNodes()
}

func GetTaskParents(workflow girraph.Graph[Workflow], taskType TaskType) []girraph.Graph[Workflow] {
	var found girraph.QueryResult[girraph.Graph[Workflow]]
	findTasksByType(&found, workflow, []TaskType{taskType})
	var result []girraph.Graph[Workflow]
	for _, node := range found.GetNodes() {
		result = append(result, node.GetParents()...)
	}
	return result
}

func findInitialTasks(query *girraph.QueryResult[girraph.Graph[Workflow]], node girraph.Graph[Workflow]) {
	children := node.GetChildren()
	if children != nil && len(children) > 0 {
		for _, child := range children {
			meta := child.GetMeta()
			if meta.GetType() == TaskNode {
				findInitialTasks(query, child)
			}
		}
	} else if node.GetMeta().GetType() == TaskNode {
		query.AddNode(node)
	}
}

func findTasksByType(query *girraph.QueryResult[girraph.Graph[Workflow]], node girraph.Graph[Workflow], types []TaskType) {
	task := node.GetMeta().GetTask()
	if task != nil {
		if taskTypeInTaskTypes(task.Type, types) {
			query.AddNode(node)
		}
	}
	for _, child := range node.GetChildren() {
		findTasksByType(query, child, types)
	}
}

func queryDecisionsByInputType(query *girraph.QueryResult[girraph.Graph[Workflow]], node girraph.Graph[Workflow], inputType ConditionType) {
	for _, child := range node.GetChildren() {
		meta := child.GetMeta()
		if meta.GetType() == ConditionNode && meta.GetCondition().Type == inputType {
			query.AddNode(node)
		}
		queryDecisionsByInputType(query, child, inputType)
	}
}

func taskTypeInTaskTypes(taskType TaskType, types []TaskType) bool {
	for _, t := range types {
		if t == taskType {
			return true
		}
	}
	return false
}
