// Package main provides the SitHub server CLI.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thorstenkramm/sithub/internal/config"
	"github.com/thorstenkramm/sithub/internal/startup"
)

func main() {
	opts := newRunOptions()

	rootCmd := &cobra.Command{
		Use:   "sithub",
		Short: "SitHub server",
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the SitHub server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.LoadWithOverrides(opts.configPath, opts.overrides(cmd))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			if err := startup.Run(cmd.Context(), cfg); err != nil {
				return fmt.Errorf("run server: %w", err)
			}
			return nil
		},
	}

	opts.bindFlags(runCmd)
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type runOptions struct {
	configPath           string
	listen               string
	port                 int
	dataDir              string
	logFile              string
	logLevel             string
	logFormat            string
	entraidAuthorizeURL  string
	entraidTokenURL      string
	entraidRedirectURI   string
	entraidClientID      string
	entraidClientSecret  string
	entraidUsersGroupID  string
	entraidAdminsGroupID string
	testAuthEnabled      bool
	testAuthUserID       string
	testAuthUserName     string
	testAuthPermitted    bool
}

func newRunOptions() *runOptions {
	return &runOptions{}
}

func (o *runOptions) bindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.configPath, "config", "./sithub.toml", "Path to config file")
	cmd.Flags().StringVar(&o.listen, "listen", "", "Server listen address")
	cmd.Flags().IntVar(&o.port, "port", 0, "Server listen port")
	cmd.Flags().StringVar(&o.dataDir, "data-dir", "", "Directory for the SQLite database and data files")
	cmd.Flags().StringVar(&o.logFile, "log-file", "", "Log file path, or '-' for stdout")
	cmd.Flags().StringVar(&o.logLevel, "log-level", "", "Log level (debug, info, warn, error)")
	cmd.Flags().StringVar(&o.logFormat, "log-format", "", "Log format (text or json)")
	cmd.Flags().StringVar(&o.entraidAuthorizeURL, "entraid-authorize-url", "", "Entra ID OAuth authorize URL")
	cmd.Flags().StringVar(&o.entraidTokenURL, "entraid-token-url", "", "Entra ID OAuth token URL")
	cmd.Flags().StringVar(&o.entraidRedirectURI, "entraid-redirect-uri", "", "Entra ID OAuth redirect URI")
	cmd.Flags().StringVar(&o.entraidClientID, "entraid-client-id", "", "Entra ID client ID")
	cmd.Flags().StringVar(&o.entraidClientSecret, "entraid-client-secret", "", "Entra ID client secret")
	cmd.Flags().StringVar(&o.entraidUsersGroupID, "entraid-users-group-id", "", "Entra ID users group ID")
	cmd.Flags().StringVar(&o.entraidAdminsGroupID, "entraid-admin-group-id", "", "Entra ID admins group ID")
	cmd.Flags().BoolVar(&o.testAuthEnabled, "test-auth-enabled", false, "Enable test auth (development only)")
	cmd.Flags().StringVar(&o.testAuthUserID, "test-auth-user-id", "", "Test auth user ID")
	cmd.Flags().StringVar(&o.testAuthUserName, "test-auth-user-name", "", "Test auth user name")
	cmd.Flags().BoolVar(&o.testAuthPermitted, "test-auth-permitted", false, "Whether test auth user is permitted")
}

func (o *runOptions) overrides(cmd *cobra.Command) map[string]interface{} {
	overrides := map[string]interface{}{}
	set := func(flag, key string, value interface{}) {
		if cmd.Flags().Changed(flag) {
			overrides[key] = value
		}
	}
	set("listen", "main.listen", o.listen)
	set("port", "main.port", o.port)
	set("data-dir", "main.data_dir", o.dataDir)
	set("log-file", "log.file", o.logFile)
	set("log-level", "log.level", o.logLevel)
	set("log-format", "log.format", o.logFormat)
	set("entraid-authorize-url", "entraid.authorize_url", o.entraidAuthorizeURL)
	set("entraid-token-url", "entraid.token_url", o.entraidTokenURL)
	set("entraid-redirect-uri", "entraid.redirect_uri", o.entraidRedirectURI)
	set("entraid-client-id", "entraid.client_id", o.entraidClientID)
	set("entraid-client-secret", "entraid.client_secret", o.entraidClientSecret)
	set("entraid-users-group-id", "entraid.users_group_id", o.entraidUsersGroupID)
	set("entraid-admin-group-id", "entraid.admins_group_id", o.entraidAdminsGroupID)
	set("test-auth-enabled", "test_auth.enabled", o.testAuthEnabled)
	set("test-auth-user-id", "test_auth.user_id", o.testAuthUserID)
	set("test-auth-user-name", "test_auth.user_name", o.testAuthUserName)
	set("test-auth-permitted", "test_auth.permitted", o.testAuthPermitted)
	return overrides
}
