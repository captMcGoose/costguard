package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

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
    default:
        usage()
        os.Exit(2)
    }
}

func usage() {
    fmt.Fprintln(os.Stderr, "Usage:")
    fmt.Fprintln(os.Stderr, "  costguard analyze <path-to-plan.json>")
}