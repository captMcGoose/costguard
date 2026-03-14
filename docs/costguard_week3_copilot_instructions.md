
# CostGuard Week 3 – GitHub Pull Request Comment Integration

## Objective

Make CostGuard usable in CI by allowing it to **post cost estimation results directly as a comment on a GitHub Pull Request**.

At the end of Week 3, CostGuard should be able to:

1. Run in a CI environment
2. Detect that it is running inside a GitHub Pull Request
3. Generate a **Markdown report**
4. Post that report as a **comment on the PR** using `GITHUB_TOKEN`

This is the first step that makes CostGuard useful during real infrastructure reviews.

---

# Week 3 Success Criteria

Running CostGuard in a PR environment should create a PR comment like:

```markdown
## CostGuard Estimate

### Infrastructure Changes

| Action | Resource |
|------|------|
| CREATE | aws_db_instance.prod_db |
| CREATE | aws_nat_gateway.main |
| UPDATE | aws_ebs_volume.data |
| REPLACE | aws_instance.worker |

### Estimated Monthly Cost Impact

**+$228/month**

### Top Cost Drivers

| Resource | Monthly Cost |
|---------|-------------|
| aws_db_instance.prod_db | $180 |
| aws_nat_gateway.main | $32 |
| aws_ebs_volume.data | $8 |
| aws_instance.worker | $8 |

---

⚠️ Estimate based on supported AWS resources only.
```

---

# Scope for Week 3

## In Scope

- Generate a Markdown report from CostGuard results
- Detect GitHub Pull Request environment
- Post comment using GitHub REST API
- Use `GITHUB_TOKEN` for authentication
- Keep CLI usable locally

## Out of Scope

Do NOT implement yet:

- SaaS backend
- Slack alerts
- Policy enforcement
- Merge blocking
- Azure/GCP pricing
- GitHub App authentication

---

# Repository Structure Changes

Add a GitHub integration module:

```
internal/github/
    comment.go
    client.go
```

---

# Step 1 – Markdown Report Generator

Create:

```
internal/report/report.go
```

Implement a function:

```
func GenerateMarkdownReport(changes []terraform.ResourceChange, summary pricing.CostSummary) string
```

Responsibilities:

- Generate Markdown tables
- Include:
  - infrastructure changes
  - cost summary
  - cost drivers

Keep formatting simple and readable.

---

# Step 2 – GitHub Client

Create:

```
internal/github/client.go
```

Implement a minimal GitHub API client.

Required endpoint:

```
POST /repos/{owner}/{repo}/issues/{issue_number}/comments
```

Authentication:

```
Authorization: Bearer <GITHUB_TOKEN>
```

Function signature:

```
func PostPRComment(owner, repo string, prNumber string, body string) error
```

---

# Step 3 – Detect GitHub PR Environment

In CI environments (GitHub Actions), these environment variables are available:

```
GITHUB_REPOSITORY
GITHUB_TOKEN
GITHUB_EVENT_PATH
```

Parse `GITHUB_EVENT_PATH` to extract the pull request number.

Add helper:

```
func DetectPullRequest() (owner string, repo string, prNumber string, err error)
```

If not running inside a PR, skip commenting.

---

# Step 4 – Integrate into CLI

Modify:

```
cmd/costguard/main.go
```

After cost estimation:

1. Generate Markdown report
2. Detect PR environment
3. If PR detected:
   - post GitHub comment
4. Otherwise:
   - print report to stdout

Pseudo logic:

```
report := report.GenerateMarkdownReport(changes, summary)

if github.IsPullRequestEnvironment():
    github.PostPRComment(...)
else:
    fmt.Println(report)
```

---

# Step 5 – Add Tests

Add tests for:

```
internal/report/report_test.go
```

Test:

- Markdown formatting
- Table generation
- Cost driver output

GitHub API client can be tested using a mock HTTP server.

---

# Step 6 – Update README

Add section:

```
## GitHub Pull Request Integration

CostGuard can comment on Terraform pull requests with estimated cost impact.

It automatically detects GitHub Actions environments and posts a PR comment.
```

Include example screenshot or Markdown output.

---

# Suggested Implementation Order

1. Markdown report generator
2. GitHub API client
3. PR environment detection
4. CLI integration
5. Unit tests
6. README update

---

# Copilot Prompt

Paste this into Copilot Chat:

You are extending the CostGuard Go CLI project.

Current state:

- Terraform plan parser works
- AWS cost estimation works
- CLI prints resource changes and cost summary

Goal:

Add **GitHub Pull Request comment support**.

Requirements:

- Generate Markdown report
- Post comment using GitHub REST API
- Use `GITHUB_TOKEN`
- Detect PR context from GitHub environment variables
- Keep CLI working locally

Tasks:

1. Implement Markdown report generator
2. Implement GitHub comment client
3. Detect pull request environment
4. Integrate into CLI
5. Add tests for report generation

Do not add SaaS features yet.

Keep implementation simple and focused on MVP.

---

# Definition of Done

Week 3 is complete when:

- CostGuard generates Markdown report
- CLI prints Markdown locally
- GitHub PR comment is successfully posted in CI
- Tests pass
- No changes required to pricing or parser modules
