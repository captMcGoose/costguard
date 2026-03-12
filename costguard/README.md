# CostGuard

CostGuard is a DevOps tool that analyzes Terraform plans.

## Current status

Week 1 MVP:
- Parse Terraform plan JSON
- Detect changed resources
- Print resource actions

## Run

To analyze a Terraform plan JSON file, use the following command:

```bash
go run ./cmd/costguard analyze <path-to-plan-json>
```

### Example

To run the example provided in the project:

```bash
go run ./cmd/costguard analyze examples/plan.json
```

### Example Output

```text
Detected resource changes:

CREATE  aws_db_instance.prod_db
CREATE  aws_nat_gateway.main
UPDATE  aws_ebs_volume.data
REPLACE aws_instance.worker
```