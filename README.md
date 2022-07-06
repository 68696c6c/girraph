# Girraph
A proof-of-concept graph library using generics.  Two types of graphs are supported: "trees", where each node has only
one parent node and "graphs", for lack of a better word, where each node can have multiple parents.  In both types of
graph, nodes can have many child nodes.  Both types of graphs can be converted to and from JSON.


## Examples
The `examples/filesystem` package models a filesystem as a tree, with each node being a directory.

The `examples/workflow` package models a simple workflow engine.  A workflow is a set of tasks that must be completed in
a specific order.  Each node in the workflow can represent a task, decision, or condition.  A task represents some work
that needs to be done.  A decision represents a set of conditional tasks that are only unlocked when some condition is
met.  A workflow is stateless, describing only the work to be done.  The current state of a workflow is expressed as a 
"plan".

For example, a workflow might describe the process of fulfilling an order, with tasks like "collect payment info", 
"charge the customer's card", "create a shipping label", etc.  When a new order is made, the workflow would be used to 
generate a "plan" for how to fulfill that specific order.  The plan is a kind of state machine that describes the state
of an order, that is, which tasks are ready to be done and which tasks have been completed.  Inputs change the state of
the plan, completing tasks or fulfilling conditions to unlock more tasks until the workflow is complete.


## Unit tests
Each example includes unit tests covering the important functionality.  The tests can be ran using the Makefile:
```shell
make test
```
