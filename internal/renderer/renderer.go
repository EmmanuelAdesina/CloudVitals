package renderer

import (
	"fmt"
	"strings"

	"github.com/EmmanuelAdesina/CloudVitals/internal/core"
)

const (
	red     = "\033[91m"
	green   = "\033[92m"
	yellow  = "\033[93m"
	blue    = "\033[94m"
	magenta = "\033[95m"
	cyan    = "\033[96m"
	white   = "\033[97m"
	bold    = "\033[1m"
	reset   = "\033[0m"
)

func color(sev core.Severity) string {
	switch sev {
	case core.Critical:
		return red
	case core.High:
		return magenta
	case core.Medium:
		return yellow
	case core.Low:
		return blue
	default:
		return white
	}
}

func icon(status string) string {
	switch status {
	case "pass":
		return green + "✓" + reset
	case "fail":
		return red + "✗" + reset
	case "error":
		return yellow + "!" + reset
	default:
		return white + "?" + reset
	}
}

func badge(sev core.Severity) string {
	c := color(sev)
	label := strings.ToUpper(string(sev))
	padding := 10 - len(label)
	return fmt.Sprintf("%s%s%s%s%s", c, bold, label, strings.Repeat(" ", padding), reset)
}

func boxScore(score int, status string) string {
	var statusColor string
	var emoji string
	switch {
	case score >= 90:
		statusColor = green
		emoji = "🛡️"
	case score >= 70:
		statusColor = yellow
		emoji = "⚠️"
	default:
		statusColor = red
		emoji = "🔥"
	}

	barWidth := 50
	filled := score * barWidth / 100
	empty := barWidth - filled

	bar := green + strings.Repeat("█", filled) + reset
	if score < 70 {
		bar = yellow + strings.Repeat("█", filled) + reset
	}
	if score < 40 {
		bar = red + strings.Repeat("█", filled) + reset
	}
	bar += strings.Repeat("░", empty)

	return fmt.Sprintf(`
%s┌────────────────────────────────────────────────────────────┐%s
%s│%s  %s CloudVitals Security Score %s                              %s│%s
%s│%s                                                            %s│%s
%s│%s  %s%d/100%s  %s[%s%s%s]%s                    %s%s%s  %s│%s
%s│%s                                                            %s│%s
%s│%s  Status: %s%s%s%s%s                                    %s│%s
%s└────────────────────────────────────────────────────────────┘%s`,
		cyan, reset,
		cyan, reset, bold, reset, cyan, reset,
		cyan, reset, cyan, reset,
		cyan, reset, bold, score, reset, bold, bar, reset, bold, reset, emoji, bold, reset, cyan, reset,
		cyan, reset, cyan, reset,
		cyan, reset, statusColor, bold, status, reset, strings.Repeat(" ", 42-len(status)), cyan, reset,
		cyan, reset)
}

func Render(results []core.CheckResult, score int) {
	var status string
	switch {
	case score >= 90:
		status = "HEALTHY"
	case score >= 70:
		status = "AT RISK"
	default:
		status = "CRITICAL"
	}

	fmt.Println(boxScore(score, status))

	fmt.Printf("\n%sCHECK RESULTS%s\n", bold+white, reset)
	fmt.Printf("%s────────────────────────────────────────────────────────────%s\n", white, reset)

	for _, r := range results {
		ic := icon(r.Status)
		bd := badge(r.Severity)
		name := r.CheckName
		if len(name) > 32 {
			name = name[:29] + "..."
		}
		count := ""
		if r.Status == "fail" && len(r.Findings) > 0 {
			count = fmt.Sprintf("  %s(%d)%s", red, len(r.Findings), reset)
		}

		fmt.Printf("  %s  %s  %s%s%s\n", ic, bd, white, name, count)
	}

	fmt.Printf("%s────────────────────────────────────────────────────────────%s\n\n", white, reset)

	failCount := 0
	for _, r := range results {
		if r.Status != "fail" || len(r.Findings) == 0 {
			continue
		}
		failCount++

		fmt.Printf("%s┌─ %sFINDING %d%s / %s%s%s ─────────────────────────────────────┐%s\n",
			red, bold, failCount, reset, color(r.Severity), strings.ToUpper(string(r.Severity)), reset, red, reset)

		for _, f := range r.Findings {
			fmt.Printf("  %sResource:%s  %s%s%s\n", cyan, reset, white, f.Resource, reset)
			fmt.Printf("  %sDetail:%s    %s%s%s\n", cyan, reset, white, f.Detail, reset)
			fmt.Printf("  %sFix:%s       %s%s%s\n", cyan, reset, yellow, f.Remediation, reset)
		}

		fmt.Printf("%s└────────────────────────────────────────────────────────────┘%s\n\n", red, reset)
	}

	if failCount == 0 {
		fmt.Printf("  %s%sAll checks passed. No action required.%s\n\n", green, bold, reset)
	} else {
		fmt.Printf("  %s%d finding(s) detected.%s Run the fix commands above.\n\n", red+bold, failCount, reset)
	}
}
