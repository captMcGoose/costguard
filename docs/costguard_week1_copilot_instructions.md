# CostGuard Week 1 Build Instructions for GitHub Copilot

## Objective

Build the **Week 1 MVP target** for CostGuard:

> A CLI that reads a Terraform plan JSON file and prints the infrastructure resources that were created, updated, or deleted.

This is **not** the full product yet.  
Week 1 is only about parsing Terraform plan JSON and extracting relevant resource changes.

---

## What success looks like

By the end of Week 1, the project should be able to do this:

```bash
costguard analyze examples/plan.json
```

And print output like:

```text
Detected resource changes:

CREATE  aws_db_instance.prod_db
CREATE  aws_nat_gateway.main
UPDATE  aws_ebs_volume.data
DELETE  aws_instance.old_worker
```

---

## Scope for Week 1

### In scope

- Create a Go CLI named `costguard`
- Read a Terraform plan JSON file from disk
- Parse Terraform `resource_changes`
- Extract:
  - resource type
  - resource name
  - resource address
  - action (`create`, `update`, `delete`)
- Print a readable summary to stdout
- Add at least one example Terraform plan JSON file
- Add unit tests for parsing logic

### Out of scope

Do **not** build any of this yet:

- cost estimation
- AWS pricing integration
- GitHub PR comments
- GitHub Action packaging
- dashboards
- SaaS backend
- billing
- multi-cloud support

---

## Tech choices

Use the following stack:

- **Language:** Go
- **Project type:** CLI
- **Input format:** Terraform plan JSON generated from:
  - `terraform plan -out plan.out`
  - `terraform show -json plan.out > plan.json`

Keep dependencies minimal. Prefer the Go standard library unless a small dependency is clearly justified.

---

## Repository structure to create

Create this structure:

```text
costguard/
├── cmd/
│   └── costguard/
│       └── main.go
├── internal/
│   └── terraform/
│       ├── parser.go
│       ├── models.go
│       └── parser_test.go
├── examples/
│   └── plan.json
├── README.md
└── go.mod
```

---

## Step 1: initialize the Go project

Create a Go module for the repo.

Suggested module name:

```bash
go mod init github.com/<your-github-username>/costguard
```

Then make sure the CLI entrypoint is:

```text
cmd/costguard/main.go
```

---

## Step 2: define the Terraform models

Create Go structs that map only the Terraform JSON fields needed for Week 1.

You do **not** need to model the full Terraform JSON schema.

### Minimum JSON fields to support

Top-level:

- `resource_changes`

For each resource change:

- `address`
- `mode`
- `type`
- `name`
- `change.actions`

### Suggested internal structs

Use something close to this shape:

```go
type Plan struct {
    ResourceChanges []TerraformResourceChange `json:"resource_changes"`
}

type TerraformResourceChange struct {
    Address string         `json:"address"`
    Mode    string         `json:"mode"`
    Type    string         `json:"type"`
    Name    string         `json:"name"`
    Change  TerraformChange `json:"change"`
}

type TerraformChange struct {
    Actions []string `json:"actions"`
}
```

Also create a simplified internal output model:

```go
type ResourceChange struct {
    Address string
    Type    string
    Name    string
    Action  string
}
```

---

## Step 3: implement the parser

Create a parser in:

```text
internal/terraform/parser.go
```

### Responsibilities

The parser should:

1. Open and read a JSON file from disk
2. Unmarshal the Terraform plan JSON
3. Iterate over `resource_changes`
4. Convert Terraform actions into a simplified action string
5. Return a list of `ResourceChange`

### Function signature suggestion

```go
func ParsePlanFile(path string) ([]ResourceChange, error)
```

---

## Step 4: normalize Terraform actions

Terraform actions can sometimes appear as arrays like:

- `["create"]`
- `["update"]`
- `["delete"]`
- `["delete", "create"]`

For Week 1, normalize actions into these labels:

- `create`
- `update`
- `delete`
- `replace`

### Suggested rules

- `["create"]` -> `create`
- `["update"]` -> `update`
- `["delete"]` -> `delete`
- `["delete", "create"]` -> `replace`
- anything unknown -> `unknown`

Create a helper function for this.

Suggested signature:

```go
func normalizeAction(actions []string) string
```

---

## Step 5: build the CLI command

Implement a basic CLI in:

```text
cmd/costguard/main.go
```

### Required command behavior

Support this command:

```bash
costguard analyze <path-to-plan-json>
```

Example:

```bash
costguard analyze examples/plan.json
```

### CLI requirements

- Validate arguments
- Show a helpful usage message if arguments are missing
- Call `ParsePlanFile`
- Print resource changes in a readable format
- Exit non-zero on error

### Suggested output format

```text
Detected resource changes:

CREATE  aws_db_instance.prod_db
CREATE  aws_nat_gateway.main
UPDATE  aws_ebs_volume.data
REPLACE aws_instance.worker
```

Use uppercase action labels in the CLI output, but keep lowercase in the internal model.

---

## Step 6: add an example plan file

Create:

```text
examples/plan.json
```

This should be a small but realistic Terraform plan JSON fixture.

It only needs enough structure to support the parser.

Include at least these cases:

- one `create`
- one `update`
- one `delete`
- one `replace`

### Minimal sample shape

The JSON can look like this:

```json
{
  "resource_changes": [
    {
      "address": "aws_db_instance.prod_db",
      "mode": "managed",
      "type": "aws_db_instance",
      "name": "prod_db",
      "change": {
        "actions": ["create"]
      }
    },
    {
      "address": "aws_nat_gateway.main",
      "mode": "managed",
      "type": "aws_nat_gateway",
      "name": "main",
      "change": {
        "actions": ["create"]
      }
    },
    {
      "address": "aws_ebs_volume.data",
      "mode": "managed",
      "type": "aws_ebs_volume",
      "name": "data",
      "change": {
        "actions": ["update"]
      }
    },
    {
      "address": "aws_instance.worker",
      "mode": "managed",
      "type": "aws_instance",
      "name": "worker",
      "change": {
        "actions": ["delete", "create"]
      }
    }
  ]
}
```

---

## Step 7: add tests

Create unit tests in:

```text
internal/terraform/parser_test.go
```

### Test coverage required

Add tests for:

1. parsing a valid plan file
2. handling missing file path
3. handling invalid JSON
4. normalizing actions:
   - create
   - update
   - delete
   - replace
   - unknown

### Suggested test names

- `TestParsePlanFile`
- `TestParsePlanFile_FileNotFound`
- `TestParsePlanFile_InvalidJSON`
- `TestNormalizeAction`

---

## Step 8: update the README

Write a simple README with:

- what CostGuard is
- what Week 1 does
- how to run the CLI
- example output

### Suggested README structure

```md
# CostGuard

CostGuard is a DevOps tool that analyzes Terraform plans.

## Current status

Week 1 MVP:
- parse Terraform plan JSON
- detect changed resources
- print resource actions

## Run

go run ./cmd/costguard analyze examples/plan.json
```

---

## Coding guidelines for Copilot

Follow these implementation rules:

- Keep the code small and readable
- Prefer explicit structs over generic maps where possible
- Do not over-engineer abstractions
- Do not add interfaces unless clearly necessary
- Use clear error messages
- Keep functions focused and testable
- Prefer deterministic tests with local fixtures
- Do not build future Week 2 features yet

---

## Suggested implementation order

Build in this order:

1. create repo structure
2. initialize Go module
3. add Terraform JSON models
4. implement `normalizeAction`
5. implement `ParsePlanFile`
6. create `examples/plan.json`
7. implement CLI in `main.go`
8. add parser tests
9. update README

---

## Exact prompt to give GitHub Copilot Chat

Use the prompt below in Copilot Chat.

---

### Copilot prompt

You are helping me build Week 1 of a Go CLI project called CostGuard.

Goal:
Build a CLI command `costguard analyze <plan.json>` that reads a Terraform plan JSON file and prints detected resource changes.

Constraints:
- Use Go
- Keep dependencies minimal
- Do not implement pricing, GitHub comments, SaaS features, or dashboards
- Focus only on parsing Terraform `resource_changes`

Project structure to create:

```text
costguard/
├── cmd/
│   └── costguard/
│       └── main.go
├── internal/
│   └── terraform/
│       ├── parser.go
│       ├── models.go
│       └── parser_test.go
├── examples/
│   └── plan.json
├── README.md
└── go.mod
```

Requirements:
1. Create Go structs for the Terraform JSON fields:
   - resource_changes
   - address
   - mode
   - type
   - name
   - change.actions
2. Create a simplified internal model `ResourceChange`
3. Implement:
   - `ParsePlanFile(path string) ([]ResourceChange, error)`
   - `normalizeAction(actions []string) string`
4. Normalize actions:
   - create
   - update
   - delete
   - replace for `["delete", "create"]`
   - unknown otherwise
5. Implement CLI:
   - `costguard analyze <path>`
   - print readable output
   - exit non-zero on failure
6. Create `examples/plan.json` with create, update, delete, and replace cases
7. Add unit tests for parser and action normalization
8. Keep the code simple and production-quality

After generating code, also explain how to run the CLI and tests locally.

---

## Definition of done

Week 1 is complete when all of the following are true:

- `go run ./cmd/costguard analyze examples/plan.json` works
- the CLI prints detected resource changes
- parsing logic is covered by tests
- invalid input returns clear errors
- repo structure is clean and easy to extend in Week 2

---

## Notes for Week 2

Do not implement Week 2 yet, but keep the design ready for it.

Week 2 will add:

- AWS pricing lookup
- monthly cost estimation
- cost diff summary

That means the output of Week 1 should already produce clean `ResourceChange` objects that can later feed a pricing module.
