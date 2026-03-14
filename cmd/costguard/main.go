package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/captMcGoose/costguard/internal/github"
    "github.com/captMcGoose/costguard/internal/pricing"
    "github.com/captMcGoose/costguard/internal/report"
    "github.com/captMcGoose/costguard/internal/terraform"
)

func main() {
    if len(os.Args) < 2 {
        usage()
        os.Exit(2)
    }

    cmd := os.Args[1]
    switch cmd {
    case "analyze":
        if len(os.Args) != 3 {
            fmt.Fprintln(os.Stderr, "usage: costguard analyze <path-to-plan.json>")
            os.Exit(2)
        }
        path := os.Args[2]
        abs, err := filepath.Abs(path)
        if err == nil {
            path = abs
        }
        rcs, err := terraform.ParsePlanFile(path)
        if err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }

        // print header (remove redundant '\n' to satisfy vet)
        fmt.Println("Detected resource changes:")
        for _, r := range rcs {
            fmt.Printf("%-7s %s\n", strings.ToUpper(r.Action), r.Address)
        }

        summary := pricing.CalculateCostDiff(rcs)
        sign := "+"
        if summary.TotalMonthly < 0 {
            sign = "-"
        }
        fmt.Printf("\nEstimated monthly cost change: %s$%.0f\n", sign, summary.TotalMonthly)

        if len(summary.Drivers) > 0 {
            fmt.Println("\nTop cost drivers:")
            for _, d := range summary.Drivers {
                fmt.Printf("+ %s → $%.0f/month\n", d.Address, d.MonthlyCost)
            }
        } else {
            fmt.Println("\nTop cost drivers: pricing unavailable")
        }

        reportBody := report.GenerateMarkdownReport(rcs, summary)
        owner, repo, prNumber, prErr := github.DetectPullRequest()
        if prErr == nil {
            if err := github.PostPRComment(owner, repo, prNumber, reportBody); err != nil {
                fmt.Fprintf(os.Stderr, "warning: failed to post PR comment: %v\n", err)
            } else {
                fmt.Println("Posted cost report to GitHub PR comment.")
            }
        } else {
            fmt.Println("PR report mode not available:", prErr)
            fmt.Println("Generated report:")
            fmt.Println(reportBody)
        }
    default:
        usage()
        os.Exit(2)
    }
}

func usage() {
    fmt.Fprintln(os.Stderr, "Usage:")
    fmt.Fprintln(os.Stderr, "  costguard analyze <path-to-plan.json>")
}