package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRunOptionsBindFlagsDefaults(t *testing.T) {
	opts := newRunOptions()
	cmd := &cobra.Command{Use: "run"}
	opts.bindFlags(cmd)

	value, err := cmd.Flags().GetString("config")
	if err != nil {
		t.Fatalf("get config flag: %v", err)
	}
	if value != "./sithub.toml" {
		t.Fatalf("unexpected config default: %s", value)
	}
}

func TestRunOptionsOverrides(t *testing.T) {
	opts := newRunOptions()
	cmd := &cobra.Command{Use: "run"}
	opts.bindFlags(cmd)

	if err := cmd.Flags().Set("listen", "0.0.0.0"); err != nil {
		t.Fatalf("set listen: %v", err)
	}
	if err := cmd.Flags().Set("port", "1234"); err != nil {
		t.Fatalf("set port: %v", err)
	}
	if err := cmd.Flags().Set("entraid-client-id", "client-1"); err != nil {
		t.Fatalf("set entraid-client-id: %v", err)
	}
	if err := cmd.Flags().Set("spaces-config-file", "./spaces.yaml"); err != nil {
		t.Fatalf("set spaces-config-file: %v", err)
	}
	if err := cmd.Flags().Set("test-auth-enabled", "true"); err != nil {
		t.Fatalf("set test-auth-enabled: %v", err)
	}

	overrides := opts.overrides(cmd)
	if overrides["main.listen"] != "0.0.0.0" {
		t.Fatalf("listen override missing: %#v", overrides)
	}
	if overrides["main.port"] != 1234 {
		t.Fatalf("port override missing: %#v", overrides)
	}
	if overrides["entraid.client_id"] != "client-1" {
		t.Fatalf("entraid client override missing: %#v", overrides)
	}
	if overrides["spaces.config_file"] != "./spaces.yaml" {
		t.Fatalf("spaces config override missing: %#v", overrides)
	}
	if overrides["test_auth.enabled"] != true {
		t.Fatalf("test auth override missing: %#v", overrides)
	}
	if _, ok := overrides["log.level"]; ok {
		t.Fatalf("unexpected log level override: %#v", overrides)
	}
}

func TestRunOptionsOverridesEmpty(t *testing.T) {
	opts := newRunOptions()
	cmd := &cobra.Command{Use: "run"}
	opts.bindFlags(cmd)

	overrides := opts.overrides(cmd)
	if len(overrides) != 0 {
		t.Fatalf("expected no overrides, got %#v", overrides)
	}
}
