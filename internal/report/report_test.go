package report

import (
    "strings"
    "testing"

    "github.com/captMcGoose/costguard/internal/pricing"
    "github.com/captMcGoose/costguard/internal/terraform"
)

func TestGenerateMarkdownReport(t *testing.T) {
    changes := []terraform.ResourceChange{
        {Action: "create", Address: "aws_db_instance.prod_db"},
        {Action: "create", Address: "aws_nat_gateway.main"},
    }
    summary := pricing.CostSummary{
        TotalMonthly: 212,
        Drivers: []pricing.CostDriver{
            {Address: "aws_db_instance.prod_db", MonthlyCost: 180},
            {Address: "aws_nat_gateway.main", MonthlyCost: 32},
        },
    }

    out := GenerateMarkdownReport(changes, summary)
    if !strings.Contains(out, "## CostGuard Estimate") {
        t.Fatal("report missing header")
    }
    if !strings.Contains(out, "aws_db_instance.prod_db") || !strings.Contains(out, "aws_nat_gateway.main") {
        t.Fatal("report missing resources")
    }
    if !strings.Contains(out, "**+$212/month**") {
        t.Fatal("report missing cost impact")
    }
}
