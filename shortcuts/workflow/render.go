package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

func renderTriageReport(w io.Writer, report TriageReport, format string) error {
	switch normalizeFormat(format) {
	case "json":
		return writeJSON(w, report)
	case "markdown":
		return writeTriageMarkdown(w, report)
	case "table":
		return writeTriageTable(w, report)
	default:
		return fmt.Errorf("unsupported workflow output format %q", format)
	}
}

func renderHealthResult(w io.Writer, result HealthResult, format string) error {
	switch normalizeFormat(format) {
	case "json":
		return writeJSON(w, result)
	case "markdown":
		return writeHealthMarkdown(w, result)
	case "table":
		return writeHealthTable(w, result)
	default:
		return fmt.Errorf("unsupported workflow output format %q", format)
	}
}

func normalizeFormat(format string) string {
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "" {
		return "json"
	}
	return format
}

func writeJSON(w io.Writer, data interface{}) error {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(encoded))
	return err
}

func writeTriageTable(w io.Writer, report TriageReport) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "NUMBER\tTYPE\tPRIORITY\tCONFIDENCE\tMISSING\tACTION"); err != nil {
		return err
	}
	for _, result := range report.Results {
		missing := "-"
		if len(result.MissingInformation) > 0 {
			missing = strings.Join(result.MissingInformation, ",")
		}
		if _, err := fmt.Fprintf(tw, "%d\t%s\t%s\t%d\t%s\t%s\n",
			result.Issue.Number,
			result.DetectedType,
			result.Priority,
			result.Confidence,
			missing,
			result.RecommendedAction,
		); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func writeHealthTable(w io.Writer, result HealthResult) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintf(tw, "REPOSITORY\tSCORE\tRISK\n%s\t%d\t%s\n\n", result.Repository, result.HealthScore, result.RiskLevel); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(tw, "METRIC\tSTATUS\tSCORE\tMAX\tREASON"); err != nil {
		return err
	}
	for _, metric := range result.Metrics {
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%d\t%d\t%s\n", metric.Name, metric.Status, metric.Score, metric.MaxScore, metric.Reason); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func writeTriageMarkdown(w io.Writer, report TriageReport) error {
	if _, err := fmt.Fprintf(w, "# Issue Triage Report\n\nRepository: `%s`\n\n", report.Repository); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| Issue | Type | Priority | Confidence | Action | Missing Information |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| --- | --- | --- | ---: | --- | --- |"); err != nil {
		return err
	}
	for _, result := range report.Results {
		missing := "-"
		if len(result.MissingInformation) > 0 {
			missing = strings.Join(result.MissingInformation, ", ")
		}
		if _, err := fmt.Fprintf(w, "| #%d | %s | %s | %d | %s | %s |\n",
			result.Issue.Number,
			result.DetectedType,
			result.Priority,
			result.Confidence,
			result.RecommendedAction,
			missing,
		); err != nil {
			return err
		}
	}
	return nil
}

func writeHealthMarkdown(w io.Writer, result HealthResult) error {
	if _, err := fmt.Fprintf(w, "# Repository Health Report\n\nRepository: `%s`\n\nHealth score: **%d**\n\nRisk level: **%s**\n\n", result.Repository, result.HealthScore, result.RiskLevel); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| Metric | Status | Score | Max | Reason |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| --- | --- | ---: | ---: | --- |"); err != nil {
		return err
	}
	for _, metric := range result.Metrics {
		if _, err := fmt.Fprintf(w, "| %s | %s | %d | %d | %s |\n", metric.Name, metric.Status, metric.Score, metric.MaxScore, metric.Reason); err != nil {
			return err
		}
	}
	if len(result.Recommendations) == 0 {
		return nil
	}
	if _, err := fmt.Fprintln(w, "\n## Recommendations"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	for _, recommendation := range result.Recommendations {
		if _, err := fmt.Fprintf(w, "- %s\n", recommendation); err != nil {
			return err
		}
	}
	return nil
}
