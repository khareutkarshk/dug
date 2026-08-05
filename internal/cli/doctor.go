package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
)

type checkStatus string

const (
	statusPass checkStatus = "pass"
	statusWarn checkStatus = "warn"
	statusFail checkStatus = "fail"
)

// CheckResult is one diagnostic check performed by `dug doctor`.
type CheckResult struct {
	Name    string      `json:"name"`
	Status  checkStatus `json:"status"`
	Message string      `json:"message"`
}

// DoctorReport is the full diagnostics report.
type DoctorReport struct {
	Config  string        `json:"config"`
	Passed  int           `json:"passed"`
	Warned  int           `json:"warned"`
	Failed  int           `json:"failed"`
	Checks  []CheckResult `json:"checks"`
	Elapsed string        `json:"elapsed"`
}

func (r *DoctorReport) add(name string, status checkStatus, message string) {
	r.Checks = append(r.Checks, CheckResult{
		Name:    name,
		Status:  status,
		Message: message,
	})

	switch status {
	case statusPass:
		r.Passed++
	case statusWarn:
		r.Warned++
	case statusFail:
		r.Failed++
	}
}

// Doctor runs local diagnostics against a configuration file.
//
// Flags:
//
//	-config   path to YAML config (default configs/edge.yaml)
//	-timeout  per-upstream probe timeout (default 3s)
//	-json     print diagnostics as JSON
func Doctor(args []string) error {
	return doctor(os.Stdout, os.Stderr, args)
}

func doctor(stdout, stderr io.Writer, args []string) error {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(stderr)

	configPath := fs.String(
		"config",
		"configs/edge.yaml",
		"Path to configuration file",
	)
	timeout := fs.Duration(
		"timeout",
		3*time.Second,
		"Timeout for each upstream probe",
	)
	jsonOut := fs.Bool("json", false, "Print diagnostics as JSON")

	if err := fs.Parse(args); err != nil {
		return reported(err)
	}

	start := time.Now()

	absPath, err := filepath.Abs(*configPath)
	if err != nil {
		absPath = *configPath
	}

	report := &DoctorReport{Config: absPath}

	cfg, err := config.Load(*configPath)
	if err != nil {
		report.add("config.load", statusFail, err.Error())
		return finishDoctor(stdout, stderr, report, start, *jsonOut)
	}
	report.add("config.load", statusPass, "configuration loaded")

	result := checkConfig(cfg)
	if !result.Ok() {
		for _, issue := range result.Issues {
			report.add("config.validate", statusFail, issue.String())
		}
		return finishDoctor(stdout, stderr, report, start, *jsonOut)
	}
	report.add("config.validate", statusPass, "configuration is valid")

	checkGatewayPort(report, cfg.Server.Port)

	if cfg.Server.TLS.Enabled {
		checkTLSFiles(report, cfg.Server.TLS)
	}

	client := &http.Client{Timeout: *timeout}
	checkUpstreams(report, client, cfg)

	return finishDoctor(stdout, stderr, report, start, *jsonOut)
}

func finishDoctor(stdout, stderr io.Writer, report *DoctorReport, start time.Time, jsonOut bool) error {
	report.Elapsed = time.Since(start).Round(time.Millisecond).String()

	if jsonOut {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			return err
		}
	} else {
		printDoctor(stdout, report)
	}

	if report.Failed > 0 {
		return reported(errors.New("doctor found failing checks"))
	}
	return nil
}

func printDoctor(w io.Writer, report *DoctorReport) {
	fmt.Fprintln(w, "Running DUG diagnostics...")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Config: %s\n", report.Config)
	fmt.Fprintln(w)

	currentSection := ""
	for _, check := range report.Checks {
		section := strings.SplitN(check.Name, ".", 2)[0]
		if section != currentSection {
			if currentSection != "" {
				fmt.Fprintln(w)
			}
			fmt.Fprintln(w, titleCase(section))
			currentSection = section
		}

		mark := "✓"
		switch check.Status {
		case statusWarn:
			mark = "!"
		case statusFail:
			mark = "✗"
		}
		fmt.Fprintf(w, "  %s %s\n", mark, check.Message)
	}

	fmt.Fprintln(w)
	fmt.Fprintf(
		w,
		"Summary: %d passed, %d warned, %d failed (%s)\n",
		report.Passed,
		report.Warned,
		report.Failed,
		report.Elapsed,
	)
}

func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func checkGatewayPort(report *DoctorReport, port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		report.add(
			"gateway.port",
			statusWarn,
			fmt.Sprintf("port %d already in use (gateway may already be running)", port),
		)
		return
	}
	_ = ln.Close()
	report.add("gateway.port", statusPass, fmt.Sprintf("port %d is available", port))
}

func checkTLSFiles(report *DoctorReport, tls config.TLSConfig) {
	if _, err := os.Stat(tls.CertFile); err != nil {
		report.add("tls.cert_file", statusFail, fmt.Sprintf("cert file missing: %s", tls.CertFile))
	} else {
		report.add("tls.cert_file", statusPass, fmt.Sprintf("cert file found: %s", tls.CertFile))
	}

	if _, err := os.Stat(tls.KeyFile); err != nil {
		report.add("tls.key_file", statusFail, fmt.Sprintf("key file missing: %s", tls.KeyFile))
	} else {
		report.add("tls.key_file", statusPass, fmt.Sprintf("key file found: %s", tls.KeyFile))
	}
}

func checkUpstreams(report *DoctorReport, client *http.Client, cfg *config.Config) {
	seen := make(map[string]struct{})

	for _, route := range cfg.Routes {
		for _, upstream := range route.Upstreams {
			if _, ok := seen[upstream.URL]; ok {
				continue
			}
			seen[upstream.URL] = struct{}{}
			probeUpstream(report, client, upstream.URL)
		}
	}
}

func probeUpstream(report *DoctorReport, client *http.Client, rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		report.add("upstream", statusFail, fmt.Sprintf("%s: invalid URL", rawURL))
		return
	}

	resp, err := client.Get(rawURL)
	if err != nil {
		report.add("upstream.reachability", statusFail, fmt.Sprintf("%s: unreachable (%v)", rawURL, err))
		return
	}
	_ = resp.Body.Close()
	report.add(
		"upstream.reachability",
		statusPass,
		fmt.Sprintf("%s: reachable (%d)", rawURL, resp.StatusCode),
	)

	healthURL := *u
	healthURL.Path = "/health"
	healthURL.RawQuery = ""
	healthURL.Fragment = ""

	resp, err = client.Get(healthURL.String())
	if err != nil {
		report.add(
			"upstream.health",
			statusWarn,
			fmt.Sprintf("%s: /health unreachable (optional)", healthURL.String()),
		)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		report.add("upstream.health", statusPass, fmt.Sprintf("%s: healthy", healthURL.String()))
		return
	}

	report.add(
		"upstream.health",
		statusWarn,
		fmt.Sprintf("%s: unhealthy status %d", healthURL.String(), resp.StatusCode),
	)
}
