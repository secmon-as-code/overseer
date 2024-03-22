package cli

import (
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/overseer/pkg/cli/config"
	"github.com/m-mizutani/overseer/pkg/domain/model"
	"github.com/m-mizutani/overseer/pkg/domain/types"
	"github.com/m-mizutani/overseer/pkg/infra"
	"github.com/m-mizutani/overseer/pkg/usecase"
	"github.com/m-mizutani/overseer/pkg/utils"
	"github.com/urfave/cli/v2"
)

func runCommand() *cli.Command {
	var (
		queryDir string
		bq       config.BigQuery
		pubsub   config.PubSub
		sentry   config.Sentry
		tags     cli.StringSlice
		ids      cli.StringSlice
	)
	return &cli.Command{
		Name:    "run",
		Usage:   `Run overseer task`,
		Aliases: []string{"r"},
		Flags: mergeFlags([]cli.Flag{
			&cli.StringFlag{
				Name:        "query-dir",
				Usage:       "Directory path of query files",
				Category:    "query",
				Destination: &queryDir,
				Aliases:     []string{"d"},
				EnvVars:     []string{"OVERSEER_QUERY_DIR"},
				Required:    true,
			},

			&cli.StringSliceFlag{
				Name:        "tag",
				Usage:       "Filter tasks by tag",
				Category:    "task",
				Destination: &tags,
				Aliases:     []string{"t"},
				EnvVars:     []string{"OVERSEER_TASK_TAG"},
			},
			&cli.StringSliceFlag{
				Name:        "id",
				Usage:       "Filter tasks by ID",
				Category:    "task",
				Destination: &ids,
				Aliases:     []string{"i"},
				EnvVars:     []string{"OVERSEER_TASK_ID"},
			},
		}, bq.Flags(), pubsub.Flags(), sentry.Flags()),
		Action: func(ctx *cli.Context) error {
			utils.Logger().Info("Run overseer task",
				"queryDir", queryDir,
				"tags", tags.Value(),
				"ids", ids.Value(),
				"bq", bq,
				"pubsub", pubsub,
				"sentry", sentry,
			)

			queryFiles, err := listQueryFiles(queryDir)
			if err != nil {
				return err
			}

			if len(queryFiles) == 0 {
				return goerr.Wrap(types.ErrInvalidOption, "No query files")
			}

			bqClient, err := bq.Configure(ctx.Context)
			if err != nil {
				return err
			}
			pubsubClient, err := pubsub.Configure(ctx.Context)
			if err != nil {
				return err
			}

			target := &model.Target{
				Tags: tags.Value(),
				IDs:  ids.Value(),
			}
			if err := target.Validate(); err != nil {
				return err
			}

			clients := infra.New(bqClient, pubsubClient)

			var tasks model.Tasks
			for _, queryFile := range queryFiles {
				fd, err := os.Open(queryFile)
				if err != nil {
					return goerr.Wrap(err, "Fail to open query file").With("file", queryFile)
				}
				defer utils.SafeClose(fd)

				task, err := model.NewTask(fd)
				if err != nil {
					return err
				}

				tasks = append(tasks, task)
			}

			if err := tasks.Validate(); err != nil {
				return err
			}

			if err := usecase.RunTasks(ctx.Context, clients, tasks, target); err != nil {
				return err
			}

			return nil
		},
	}
}

func listQueryFiles(queryDir string) ([]string, error) {
	entries, err := os.ReadDir(queryDir)
	if err != nil {
		return nil, goerr.Wrap(err, "Fail to read query directory")
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			subFiles, err := listQueryFiles(filepath.Join(queryDir, e.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else if filepath.Ext(e.Name()) == ".sql" {
			files = append(files, filepath.Join(queryDir, e.Name()))
		}
	}

	return files, nil
}
