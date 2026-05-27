package renderer

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/reliastra/cloudvitals/internal/core"
)

func Render(results []core.CheckResult, score int) {
	fmt.Println()
	fmt.Println("===================================================")
	fmt.Printf("  CloudVitals Security Score: %d/100\n", score)

	switch {
	case score >= 90:
		fmt.Println("  Status: HEALTHY")
	case score >= 70:
		fmt.Println("  Status: AT RISK")
	default:
		fmt.Println("  Status: CRITICAL")
	}
	fmt.Println("===================================================")
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "SEVERITY\tCHECK\t\tSTATUS\tFINDINGS\t")
	for _, r := range results {
		icon := "PASS"
		if r.Status == "fail" {
			icon = "FAIL"
		} else if r.Status == "error" {
			icon = "ERR "
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t\n", r.Severity, r.CheckName, icon, len(r.Findings))
	}
	w.Flush()

	for _, r := range results {
		if r.Status != "fail" || len(r.Findings) == 0 {
			continue
		}
		fmt.Println()
		fmt.Printf("--- %s (%s) ---\n", r.CheckName, r.Severity)
		for _, f := range r.Findings {
			fmt.Printf("  Resource: %s\n", f.Resource)
			fmt.Printf("  Detail:   %s\n", f.Detail)
			fmt.Printf("  Fix:      %s\n", f.Remediation)
		}
	}
	fmt.Println()
}
