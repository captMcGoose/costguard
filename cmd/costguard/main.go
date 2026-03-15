package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"

    "github.com/captMcGoose/costguard/internal/github"
    "github.com/captMcGoose/costguard/internal/pricing"
    "github.com/captMcGoose/costguard/internal/report"
    "github.com/captMcGoose/costguard/internal/terraform"
)

const version = "v1.0.0"

func main() {
    if len(os.Args) < 2 {
        usage()
        os.Exit(2)
    }

    cmd := os.Args[1]
    switch cmd {
    case "analyze":
        fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
        debug := fs.Bool("debug", false, "enable debug logs")
        if err := fs.Parse(os.Args[2:]); err != nil {
            fmt.Fprintln(os.Stderr, "usage: costguard analyze [--debug] <path-to-plan.json>")
            os.Exit(2)
        }

        args := fs.Args()
        if len(args) != 1 {
            fmt.Fprintln(os.Stderr, "usage: costguard analyze [--debug] <path-to-plan.json>")
            os.Exit(2)
        }

        path := args[0]
        abs, err := filepath.Abs(path)
        if err == nil {
            path = abs
        }

        rcs, err := terraform.ParsePlanFile(path)
        if err != nil {
            if os.IsNotExist(err) {
                fmt.Fprintf(os.Stderr, "Error: missing plan file %q\n", path)
            } else {
                fmt.Fprintf(os.Stderr, "Error: unable to parse Terraform plan file: %v\n", err)
            }
            os.Exit(1)
        }

        summary := pricing.CalculateCostDiff(rcs, *debug)

        sign := "+"
        total := summary.TotalMonthly
        if total < 0 {
            sign = "-"
            total = -total
        }
        fmt.Println("💰 CostGuard Estimate")
        fmt.Println("")
        fmt.Printf("Estimated Monthly Cost Impact: %s$%.0f/month\n\n", sign, total)

        fmt.Println("Top Cost Drivers")
        fmt.Println("----------------")
        if len(summary.Drivers) == 0 {
            fmt.Println("No supported resources found.")
        } else {
            for _, d := range summary.Drivers {
                fmt.Printf("%-25s $%.0f\n", d.Address, d.MonthlyCost)
            }
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

    case "version":
        fmt.Printf("CostGuard %s\n", version)
    default:
        usage()
        os.Exit(2)
    }
}

func usage() {
    fmt.Fprintln(os.Stderr, "Usage:")
    fmt.Fprintln(os.Stderr, "  costguard analyze [--debug] <path-to-plan.json>")
    fmt.Fprintln(os.Stderr, "  costguard version")
}