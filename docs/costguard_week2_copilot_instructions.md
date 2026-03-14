
# CostGuard Week 2 Build Instructions for GitHub Copilot

## Objective

Extend the Week 1 MVP by adding **basic AWS monthly cost estimation** for detected Terraform resource changes.

Week 1 produced a CLI that parses Terraform `plan.json` and returns normalized `ResourceChange` objects.

Week 2 must add:

- A pricing module
- Basic monthly cost estimation
- A cost diff summary printed in the CLI

This is still **not the full SaaS product**.

---

# Week 2 Success Criteria

Running:

```bash
costguard analyze examples/plan.json
```

Should produce output like:

```text
Detected resource changes:

CREATE  aws_db_instance.prod_db
CREATE  aws_nat_gateway.main
UPDATE  aws_ebs_volume.data
REPLACE aws_instance.worker

Estimated monthly cost change: +$240

Top cost drivers:
+ aws_db_instance.prod_db → $180/month
+ aws_nat_gateway.main → $70/month
```

Important: numbers can be approximate for MVP.

---

# Scope for Week 2

## In Scope

- Implement a **pricing module**
- Support cost estimation for a **small set of AWS resources**
- Calculate **monthly cost estimates**
- Produce **cost delta summary**
- Display **top cost drivers**
- Integrate pricing with the CLI

## Out of Scope

Do **not** implement yet:

- GitHub PR comments
- GitHub Actions
- SaaS backend
- dashboards
- billing
- Slack integration
- Azure/GCP pricing
- AWS pricing APIs

For MVP we will use **static pricing values**.

---

# Supported AWS Resources (MVP)

Only implement pricing for these resources:

| Terraform Resource | Pricing Basis |
|--------------------|---------------|
| aws_instance | instance_type |
| aws_db_instance | instance_class |
| aws_nat_gateway | flat monthly |
| aws_ebs_volume | size * price per GB |

Any unsupported resource should return:

```text
pricing unavailable
```

and should not break the CLI.

---

# Repository Structure Changes

Add a pricing module:

```text
costguard/
├── internal/
│   ├── terraform/
│   │   parser.go
│   │   models.go
│   │   parser_test.go
│   │
│   └── pricing/
│       aws.go
│       catalog.go
│       calculator.go
│       calculator_test.go
```

---

# Step 1: Define Pricing Models

Create pricing models in:

```text
internal/pricing/catalog.go
```

Define a **static price catalog**.

Example values (approximate):

```go
var EC2Pricing = map[string]float64{
    "t3.micro": 8.0,
    "t3.small": 15.0,
    "t3.medium": 30.0,
}

var RDSPricing = map[string]float64{
    "db.t3.micro": 15.0,
    "db.t3.small": 30.0,
    "db.m6i.large": 180.0,
}

const NATGatewayMonthly = 32.0

const EBSPricePerGB = 0.08
```

These numbers are approximate for MVP only.

---

# Step 2: Pricing Input Extraction

Add logic that extracts the attributes required for pricing.

Examples:

### aws_instance

Need:

```
instance_type
```

### aws_db_instance

Need:

```
instance_class
```

### aws_ebs_volume

Need:

```
size
```

Terraform attributes are located inside:

```
change.after
```

Extend the Terraform models if necessary to expose these values.

---

# Step 3: Implement Cost Estimator

Create:

```text
internal/pricing/calculator.go
```

Implement:

```go
func EstimateMonthlyCost(change terraform.ResourceChange, attrs map[string]interface{}) (float64, error)
```

Responsibilities:

- Identify supported resource types
- Extract required attributes
- Look up pricing
- Return estimated monthly cost

Unsupported resources should return:

```
0, ErrUnsupportedResource
```

---

# Step 4: Aggregate Cost Changes

Add a function:

```go
func CalculateCostDiff(changes []terraform.ResourceChange) CostSummary
```

Define a model:

```go
type CostSummary struct {
    TotalMonthly float64
    Drivers []CostDriver
}

type CostDriver struct {
    Address string
    MonthlyCost float64
}
```

Only include drivers with non-zero cost.

Sort drivers descending by cost.

---

# Step 5: Integrate with CLI

Modify:

```text
cmd/costguard/main.go
```

After printing detected resources, add:

1. cost estimation
2. summary output

Example CLI output:

```text
Detected resource changes:

CREATE  aws_db_instance.prod_db
CREATE  aws_nat_gateway.main
UPDATE  aws_ebs_volume.data
REPLACE aws_instance.worker

Estimated monthly cost change: +$240

Top cost drivers:
+ aws_db_instance.prod_db → $180/month
+ aws_nat_gateway.main → $70/month
```

If no priced resources exist:

```text
No supported resources detected for cost estimation.
```

---

# Step 6: Add Tests

Create:

```text
internal/pricing/calculator_test.go
```

Tests required:

- EC2 instance pricing
- RDS instance pricing
- NAT gateway pricing
- EBS volume pricing
- unsupported resource handling
- cost aggregation and sorting

Use simple fixtures instead of requiring full Terraform JSON.

---

# Step 7: Update README

Add a **Week 2 section** describing:

- supported AWS resources
- example CLI output
- known limitations

Example:

```md
## Cost Estimation (Week 2)

CostGuard now estimates monthly cloud cost impact for common AWS resources.

Supported resources:
- EC2 instances
- RDS instances
- NAT Gateways
- EBS volumes
```

---

# Coding Guidelines

Follow these rules:

- Keep Week 1 parser untouched except for exposing needed attributes
- Keep pricing logic isolated in `internal/pricing`
- Avoid large dependencies
- Prefer deterministic unit tests
- Fail gracefully on unsupported resources
- Do not implement AWS API pricing yet

---

# Suggested Implementation Order

1. create pricing catalog
2. implement attribute extraction
3. implement `EstimateMonthlyCost`
4. implement cost aggregation
5. integrate CLI output
6. add unit tests
7. update README

---

# Copilot Prompt

Use this prompt in Copilot Chat.

---

You are helping extend the CostGuard Go CLI project.

Current state:
- Terraform parser works
- CLI prints resource changes

Goal:
Add **basic AWS monthly cost estimation**.

Constraints:

- Support only:
  - aws_instance
  - aws_db_instance
  - aws_nat_gateway
  - aws_ebs_volume
- Use static pricing values
- Do not add AWS API integrations
- Keep pricing logic inside `internal/pricing`
- CLI must print cost summary and top drivers

Tasks:

1. Create pricing catalog with static prices
2. Extract pricing attributes from Terraform resources
3. Implement `EstimateMonthlyCost`
4. Implement `CalculateCostDiff`
5. Integrate cost summary into CLI output
6. Add unit tests
7. Update README

Keep the code simple and focused on MVP requirements.

After generating code:

- explain how to run tests
- show example CLI output
- confirm which resources are supported

---

# Definition of Done

Week 2 is complete when:

- CLI prints cost estimation summary
- top cost drivers are shown
- supported AWS resources produce estimates
- unsupported resources fail gracefully
- tests pass
