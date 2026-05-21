//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

var (
	client    *ink.Client
	workspace string
	project   string
	runID     string
)

func TestMain(m *testing.M) {
	loadDotEnv()

	apiKey := os.Getenv("INK_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "INK_API_KEY is required")
		os.Exit(1)
	}

	cfg := ink.Config{APIKey: apiKey}
	if u := os.Getenv("INK_API_URL"); u != "" {
		cfg.BaseURL = u
	}
	client = ink.NewClient(cfg)

	workspace = envOr("INK_WORKSPACE", "e2e")
	project = envOr("INK_PROJECT", "default")
	runID = envOr("GITHUB_RUN_ID", strconv.FormatInt(time.Now().Unix(), 10))

	os.Exit(m.Run())
}
