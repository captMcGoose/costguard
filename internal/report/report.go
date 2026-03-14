package report

import (
    "fmt"
    "strings"
    "github.com/captMcGoose/costguard/internal/pricing"
    "github.com/captMcGoose/costguard/internal/terraform"
)

func GenerateMarkdownReport(changes []terraform.ResourceChange, summary pricing.CostSummary) string {
    var b strings.Builder
    b.WriteString("## CostGuard Estimate\n\n")
    b.WriteString("### Infrastructure Changes\n\n")
    b.WriteString("| Action | Resource |\n")
    b.WriteString("|------|------|\n")
    for _, c := range changes {
        b.WriteString(fmt.Sprintf("| %s | %s |\n", strings.ToUpper(c.Action), c.Address))
    }
    b.WriteString("\n### Estimated Monthly Cost Impact\n\n")
    b.WriteString(fmt.Sprintf("**+$%.0f/month**\n\n", summary.TotalMonthly))
    b.WriteString("### Top Cost Drivers\n\n")
    if len(summary.Drivers) == 0 {
        b.WriteString("No supported resources priced.\n\n")
    } else {
        b.WriteString("| Resource | Monthly Cost |\n")
        b.WriteString("|---------|-------------|\n")
        for _, d := range summary.Drivers {
            b.WriteString(fmt.Sprintf("| %s | $%.0f |\n", d.Address, d.MonthlyCost))
        }
        b.WriteString("\n")
    }
    b.WriteString("---\n\n⚠️ Estimate based on supported AWS resources only.\n")
    return b.String()
}
