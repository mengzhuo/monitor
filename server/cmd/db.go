package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.ntppool.org/monitor/ntpdb"
)

func (cli *CLI) dbCmd() *cobra.Command {

	var serverCmd = &cobra.Command{
		Use:   "db",
		Short: "db utility functions",
		// DisableFlagParsing: true,
		// Args:  cobra.ExactArgs(1),
	}

	serverCmd.PersistentFlags().AddGoFlagSet(cli.Config.Flags())

	serverCmd.AddCommand(
		&cobra.Command{
			Use:   "mon",
			Short: "monitor config debug",
			RunE:  cli.Run(cli.dbMonitorConfig),
		})

	return serverCmd
}

func (cli *CLI) dbMonitorConfig(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("db mon [monitername]")
	}

	name := args[0]

	ctx := context.Background()

	dbconn, err := cli.OpenDB()
	if err != nil {
		return err
	}
	db := ntpdb.New(dbconn)

	mon, err := db.GetMonitorTLSName(ctx, sql.NullString{String: name, Valid: true})
	if err != nil {
		return err
	}

	smon, err := db.GetSystemMonitor(ctx, "settings", mon.IpVersion)
	if err == nil {
		log.Printf("system defaults: %s", smon.Config)
		mconf, err := mon.GetConfigWithDefaults([]byte(smon.Config))
		if err != nil {
			return err
		}
		fmt.Printf("mconf: %+v", mconf)
	}

	return nil
}
