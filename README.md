# Spinet

Spinet is a configurable task trigger like ifttt

## Build
``` bash
make
```
## Run task
### Standalone Mode
``` bash
# Simple task
bin/spinet task -f examples/simple.yaml

# Start example api server
go run examples/api/example.go

# Start a simple task which use example api server
bin/spinet task -f examples/api-simple.yaml

# Task with custom app
bin/spinet -a examples/apps/fetch-and-submit.yaml task -f examples/custom-app.yaml
```

### Cluster Mode
```bash
# Start cluster
bin/spinet

# Create namespace
bin/spictl create namespace example

# Apply task/app to cluster
bin/spictl apply -f example/simple.yaml
bin/spictl apply -f example/apps/fetch-and-submit.yaml

# Get resource from cluster
bin/spictl get apps
bin/spictl get namespaces
bin/spictl -n default get tasks
```

### Logging Level
```bash
bin/spinet -v N # 0 < N < 10
```

# Create a task

## Flow

### Task

Triggered by triggers → Init dictionary → Process Inputs → Check Conditions → Process Outputs

### Input

Process Dependencies →  Execute App → Process Mapper → Check Conditions

### Output

Process Dependencies → Execute App

## Definition

```yaml
kind: task
name: name
namespace: namespace
dictionary:
  key: value
triggers:
  - type: timer
    options:
      period: 60
  - type: hook
    options:
      name: example
      mapper:
        data: "$.__app__" # set hole app data to dictionary.data
inputs:
  - app: simple
    mapper:
      data: "$.content" # set app.content to dictionary.data
    conditions:
      - operator: eq # compare app.content to apple
        values:
          - "$.content"
          - "apple"
conditions:
  - operator: and # nested conditions
    conditions:
      - operator: eq # compare dictionary.data to apple
        values:
          - "$.data"
          - "apple"
      - operator: eq # compare dictionary.key to value
        values:
          - "$.key"
          - "value"
outputs:
  - app: simple
    options:
      content: "$.data" # print dictionary.data
```

### Meta

Name & Namespace

### Dictionary

Global data dictionary for task, every app can read & write this dictionary.

### Trigger

Define how a task be trigged, spinet has two types of trigger:

Timer: Run a task periodically

Hook: Run a task when receive an HTTP request

### Input & Output

Define the process of a task. Each input or output should specify an app, app will run with its options, and return its data. Optionally, you can config mapper & conditions for each input.

Mapper: The data can merge into task's dictionary by mapper.

Conditions: If conditions return false, task will skip outputs.

Both input and output could have a list of dependencies, and it will be executed before input & output execution. Each dependency is also an input or output.

### Conditions

Conditions can appear in inputs and top-level of tasks.

When in inputs, it will be executed after the app's execution.

When in top-level, it will be executed after all inputs completed.

Each input's condition or top-level condition return false, task will skip all outputs.

## Values

### Types

- Constant
- Variable
- Template
- Map
- Build-in Variable (Not implemented)

```yaml
# Constant
- 10
- "text"
- type: constant
  value: "text"

# Variable
- "$.data"
- "$.__dict__.data"
- "$.__app__.data"
- "$.__magic__.__index__"
- "$.data.0.key"
- type: variable
  value: "data"
- type: variable
  value: ["data", 0, "key"]

# Template
- "${Data is: {{ .data }}}"
- type: template
  value: "Data is: {{ .data }}"

# Map
- data:
    key: "$.data"
- type: map
  value:
    data:
      key: "$.data"

# Build-in
- "$DATE"
- "$TIME"
```

# Custom App
Coming soon