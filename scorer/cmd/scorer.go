package cmd

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/spf13/cobra"
	"go.ntppool.org/monitor/ntpdb"
	"go.ntppool.org/monitor/scorer"
)

func (cli *CLI) scorerCmd() *cobra.Command {

	var scorerCmd = &cobra.Command{
		Use:   "scorer",
		Short: "scorer execution",
	}

	scorerCmd.PersistentFlags().AddGoFlagSet(cli.Config.Flags())

	scorerCmd.AddCommand(
		&cobra.Command{
			Use:   "run",
			Short: "scorer run once",
			RunE:  cli.Run(cli.scorer),
		})

	scorerCmd.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "run continously",
			RunE:  cli.Run(cli.scorerServer),
		})

	scorerCmd.AddCommand(
		&cobra.Command{
			Use:   "setup",
			Short: "setup scorers",
			RunE:  cli.Run(cli.scorerSetup),
		})

	return scorerCmd
}

func (cli *CLI) scorerServer(cmd *cobra.Command, args []string) error {

	for {
		count, err := cli.scorerRun(cmd, args)
		if err != nil {
			return err
		}

		if count == 0 {
			time.Sleep(20 * time.Second)
		}

	}

}

func (cli *CLI) scorer(cmd *cobra.Command, args []string) error {
	_, err := cli.scorerRun(cmd, args)
	return err
}

func (cli *CLI) scorerRun(cmd *cobra.Command, args []string) (int, error) {

	ctx := context.Background()

	dbconn, err := ntpdb.OpenDB(cli.Config.Database)
	if err != nil {
		return 0, err
	}

	sc, err := scorer.New(ctx, dbconn)
	if err != nil {
		return 0, nil
	}
	count, err := sc.Run()
	if err != nil {
		return count, err
	}
	log.Printf("Processed %d log scores", count)
	return count, nil
}

func (cli *CLI) scorerSetup(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	dbconn, err := ntpdb.OpenDB(cli.Config.Database)
	if err != nil {
		return err
	}

	tx, err := dbconn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	db := ntpdb.New(dbconn).WithTx(tx)

	scr, err := scorer.New(ctx, dbconn)
	if err != nil {
		return err
	}

	dbScorers, err := db.GetScorers(ctx)
	if err != nil {
		return err
	}
	existingScorers := map[string]bool{}

	for _, dbS := range dbScorers {
		existingScorers[dbS.Name] = true
	}

	log.Printf("dbScorers: %+v", dbScorers)

	codeScorers := scr.Scorers()

	minLogScoreID, err := db.GetMinLogScoreID(ctx)
	if err != nil {
		return err
	}

	for name := range codeScorers {
		if _, ok := existingScorers[name]; ok {
			log.Printf("%s already configured", name)
			continue
		}
		log.Printf("setting up %s", name)

		insert, err := db.InsertScorer(ctx, ntpdb.InsertScorerParams{
			Name:    name,
			TlsName: sql.NullString{String: name + ".scores.ntp.dev", Valid: true},
		})
		if err != nil {
			return err
		}
		scorerID, err := insert.LastInsertId()
		if err != nil {
			return err
		}
		db.InsertScorerStatus(ctx, ntpdb.InsertScorerStatusParams{
			ScorerID:   int32(scorerID),
			LogScoreID: sql.NullInt64{Int64: minLogScoreID, Valid: true},
		})

	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}
