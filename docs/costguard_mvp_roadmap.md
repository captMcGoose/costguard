# CostGuard MVP Roadmap (4 Week Plan)

## Goal

Build and launch a minimal viable version of **CostGuard**, a DevOps
tool that analyzes Terraform plans and comments on pull requests with
estimated cloud cost changes.

The objective of the MVP is to validate that:

> DevOps teams want cost visibility during Terraform pull request
> reviews.

The MVP should be intentionally small and focused.

**Core MVP capability:**

-   Parse a Terraform plan
-   Estimate AWS monthly cost changes
-   Comment the result on a pull request

No dashboard, SaaS platform, billing, or advanced features are required
at this stage.

------------------------------------------------------------------------

# Week 1 --- Terraform Plan Parser

## Objective

Build the core engine that can read and interpret Terraform plan JSON
output.

Terraform provides structured output via:

    terraform show -json plan.out

Your parser must extract relevant resource changes from this JSON.

------------------------------------------------------------------------

## Tasks

### 1. Generate example Terraform plans

Create a few Terraform test configurations that generate realistic
plans:

Examples:

-   EC2 instance
-   RDS instance
-   NAT Gateway
-   EBS volume

Generate plans:

    terraform plan -out plan.out
    terraform show -json plan.out > plan.json

Store these files in:

    examples/plan.json

------------------------------------------------------------------------

### 2. Parse Terraform plan JSON

Your parser should extract:

-   Resource type
-   Resource name
-   Action (`create`, `update`, `delete`)
-   Resource configuration values

Example structure to extract:

    resource_change {
      type: aws_db_instance
      name: prod_db
      change.actions: ["create"]
    }

------------------------------------------------------------------------

### 3. Create internal data structures

Define simple structs/models representing resources:

Example:

    ResourceChange {
      Type
      Name
      Action
      Attributes
    }

------------------------------------------------------------------------

### 4. Build CLI prototype

Create a command line tool:

    costguard analyze plan.json

Output example:

    Resources detected:

    + aws_db_instance.prod_db
    + aws_nat_gateway.main

------------------------------------------------------------------------

## Deliverables

At the end of Week 1:

-   Terraform plan parser working
-   CLI tool reading plan.json
-   List of resource changes printed
-   Repository structure created

------------------------------------------------------------------------

# Week 2 --- AWS Cost Estimation

## Objective

Add the ability to estimate monthly costs for detected AWS resources.

Start small and support only a few resources initially.

------------------------------------------------------------------------

## Tasks

### 1. Define pricing model

Start with static price mappings.

Example:

    db.m6i.large → $138/month
    nat_gateway → $32/month
    t3.medium → $30/month

Store pricing in a simple lookup table.

------------------------------------------------------------------------

### 2. Extract pricing attributes

Your parser must extract attributes required for pricing.

Examples:

EC2:

    instance_type

RDS:

    instance_class

EBS:

    size
    volume_type

------------------------------------------------------------------------

### 3. Implement cost calculator

Create a pricing module:

    pricing/
      aws.go

Example function:

    EstimateCost(resource ResourceChange) -> MonthlyCost

------------------------------------------------------------------------

### 4. Compute cost diff

Aggregate costs across all resources.

Example output:

    Estimated monthly cost change: +$240

    Top cost drivers:
    + aws_db_instance.prod_db → $180/month
    + aws_nat_gateway.main → $70/month

------------------------------------------------------------------------

## Deliverables

At the end of Week 2:

-   Basic AWS cost estimation
-   Cost aggregation
-   CLI output showing monthly cost change
-   Top cost-driving resources listed

------------------------------------------------------------------------

# Week 3 --- GitHub Pull Request Comment Bot

## Objective

Integrate CostGuard with GitHub so it can comment on pull requests
automatically.

------------------------------------------------------------------------

## Tasks

### 1. Learn GitHub API basics

Use GitHub REST API to create PR comments.

Endpoint:

    POST /repos/{owner}/{repo}/issues/{issue_number}/comments

------------------------------------------------------------------------

### 2. Format CostGuard comment

Example PR comment:

    ⚠️ CostGuard Estimate

    Monthly cost change: +$240

    Top cost drivers:
    + aws_db_instance.prod_db → $180/month
    + aws_nat_gateway → $70/month

Keep formatting simple and readable.

------------------------------------------------------------------------

### 3. Implement PR comment integration

Allow the CLI to post comments using:

    GITHUB_TOKEN

CLI usage example:

    costguard analyze plan.json --comment

------------------------------------------------------------------------

### 4. Improve comment formatting

Optional enhancements:

-   emojis
-   markdown tables
-   cost summary section

------------------------------------------------------------------------

## Deliverables

At the end of Week 3:

-   CostGuard can post PR comments
-   Output formatted clearly
-   CLI works in CI environments

------------------------------------------------------------------------

# Week 4 --- GitHub Action + Public Launch

## Objective

Package CostGuard into a GitHub Action so any repository can install it
easily.

------------------------------------------------------------------------

## Tasks

### 1. Create GitHub Action

Provide an action that runs CostGuard automatically.

Example usage:

    name: CostGuard

    on: pull_request

    jobs:
      costguard:
        runs-on: ubuntu-latest

        steps:
          - uses: actions/checkout@v3

          - uses: hashicorp/setup-terraform@v2

          - run: terraform init

          - run: terraform plan -out plan.out

          - run: terraform show -json plan.out > plan.json

          - run: costguard analyze plan.json

------------------------------------------------------------------------

### 2. Write a strong README

README should include:

-   project description
-   problem statement
-   example PR comment
-   installation instructions
-   roadmap

------------------------------------------------------------------------

### 3. Add GitHub topics

Recommended topics:

    terraform
    finops
    devops
    cloud-cost
    iac
    github-actions

------------------------------------------------------------------------

### 4. Launch publicly

Share CostGuard in:

-   r/devops
-   r/terraform
-   Hacker News
-   DevOps Discord communities
-   LinkedIn
-   Twitter/X

Example post:

> I built CostGuard --- a tool that comments on Terraform PRs with the
> estimated cloud cost change before merging.

------------------------------------------------------------------------

## Deliverables

At the end of Week 4:

-   Working GitHub Action
-   Public GitHub repo
-   Documentation complete
-   Initial public launch

------------------------------------------------------------------------

# What the MVP Does NOT Include

The following features are intentionally excluded from the MVP:

-   SaaS dashboard
-   billing or subscriptions
-   Slack integration
-   policy enforcement
-   Azure/GCP support
-   advanced pricing engine

These will come **after validation**.

------------------------------------------------------------------------

# Success Criteria

The MVP is successful if:

-   DevOps engineers install it
-   People star the GitHub repo
-   Teams use it in CI pipelines
-   Developers request new features

Metrics to track:

-   GitHub stars
-   GitHub installs
-   issues and feedback
-   community interest

------------------------------------------------------------------------

# Next Phase (Post-MVP)

If adoption is validated, the next stage includes:

-   FinOps guardrails
-   Merge blocking policies
-   Slack alerts
-   SaaS dashboard
-   multi-cloud support
-   team management
-   billing

These features will form the **CostGuard SaaS platform**.
