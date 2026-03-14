
# CostGuard Week 4 – GitHub Action Packaging & Public Launch

## Objective

Make CostGuard easy for other engineers to install and use by packaging it as a reusable GitHub Action and preparing the repository for a public launch.

By the end of Week 4:

- CostGuard can be installed in any repo with one GitHub Action step
- The README clearly explains installation
- The project is structured and documented for external users
- A minimal release version is ready

This step turns CostGuard from a prototype into a usable developer tool.

---

# Week 4 Success Criteria

A user should be able to install CostGuard in their repository using:

- uses: captMcGoose/costguard@v1

And the Action should:

1. Run CostGuard
2. Analyze the Terraform plan
3. Post a PR comment with cost impact

No manual cloning or building required by the user.

---

# Scope for Week 4

## In Scope

- Create a reusable GitHub Action
- Package CostGuard CLI for Action use
- Add action.yml
- Simplify installation workflow
- Improve README for public usage
- Create an example workflow snippet

## Out of Scope

Do NOT implement yet:

- SaaS backend
- dashboards
- Slack integration
- policy enforcement
- multi-cloud pricing
- advanced caching
- comment update logic

---

# Repository Changes

Add GitHub Action metadata:

action.yml

Optional:

.github/workflows/test-action.yml

---

# Step 1 – Create GitHub Action Metadata

Create file:

action.yml

Example structure:

name: "CostGuard"
description: "Analyze Terraform changes and comment estimated cloud cost impact on pull requests."

inputs:
  plan-file:
    description: "Path to Terraform plan JSON"
    required: true
    default: "plan.json"

runs:
  using: "composite"
  steps:
    - run: go build -o costguard ./cmd/costguard
      shell: bash

    - run: ./costguard analyze ${{ inputs.plan-file }}
      shell: bash

This allows users to call CostGuard directly from workflows.

---

# Step 2 – Add Example Workflow

Add documentation snippet in README.

Example usage:

name: Terraform CostGuard

on:
  pull_request:

jobs:
  costguard:
    runs-on: ubuntu-latest

    permissions:
      contents: read
      pull-requests: write

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Terraform Plan
        run: |
          terraform init
          terraform plan -out plan.out
          terraform show -json plan.out > plan.json

      - name: Run CostGuard
        uses: captMcGoose/costguard@v1
        with:
          plan-file: plan.json
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

---

# Step 3 – Improve README

Update README with sections:

Overview

Explain what CostGuard does:

CostGuard analyzes Terraform pull requests and comments the estimated cloud cost impact directly on the PR.

Installation

Provide the GitHub Action snippet.

Example Output

Include example PR comment.

Supported Resources

List current support:

- aws_instance
- aws_db_instance
- aws_nat_gateway
- aws_ebs_volume

Limitations

Explain:

- estimates only
- limited resource coverage
- static pricing

---

# Step 4 – Version Tag

Create a GitHub tag for the first usable version.

git tag v1
git push origin v1

GitHub Actions require version tags for stable references.

---

# Step 5 – Add Action Test Workflow

Create:

.github/workflows/test-action.yml

Purpose: ensure the Action builds correctly.

Example:

name: Test Action

on:
  push:

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - run: go build ./cmd/costguard

      - run: go test ./...

---

# Suggested Implementation Order

1. create action.yml
2. verify CLI runs through action locally in repo workflow
3. update README
4. add example workflow snippet
5. add test workflow
6. tag version v1

---

# Copilot Prompt

You are extending the CostGuard Go CLI project.

Current state:

- Terraform plan parser works
- AWS cost estimation works
- PR comment integration works

Goal:

Package CostGuard as a reusable GitHub Action so users can install it easily.

Tasks:

1. Create action.yml using a composite GitHub Action
2. Build the CostGuard CLI inside the action
3. Run costguard analyze using a provided plan file
4. Add README installation instructions
5. Add example workflow snippet
6. Add a basic test workflow for CI

Constraints:

- Do not modify parser or pricing logic
- Keep the Action simple
- Assume Go is available via setup-go
- Use GITHUB_TOKEN for PR comments

---

# Definition of Done

Week 4 is complete when:

- uses: captMcGoose/costguard@v1 works
- CostGuard runs in GitHub Actions
- PR comments are posted
- README contains installation instructions
- Tests pass
