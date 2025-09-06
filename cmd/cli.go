package main

import (
	"cli/internal/charm"
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"github.com/vegidio/mediaorient"
)

func buildCliCommands() *cli.Command {
	var media []mediaorient.Media
	var files []string
	var directory string
	var output string
	var recursive bool
	var mediaType string
	var ignoreErrors bool
	var dryRun bool
	var err error

	return &cli.Command{
		Name:            "mediaorient",
		Usage:           "a tool to determine/fix the orientation of images & videos",
		UsageText:       "mediaorient <command>",
		Version:         mediaorient.Version,
		HideHelpCommand: true,
		Commands: []*cli.Command{
			{
				Name:      "files",
				Usage:     "determine/fix the orientation of one or more files",
				UsageText: "mediaorient files <file1> [<file2> ...] ",
				Flags:     []cli.Flag{},
				Action: func(ctx context.Context, command *cli.Command) error {
					files = command.Args().Slice()

					if len(files) < 1 {
						return fmt.Errorf("at least one files must be specified")
					}

					files = lo.Map(files, func(file string, _ int) string {
						fullFile, _ := expandPath(file)
						return fullFile
					})

					if output == "report" {
						charm.PrintCalculateFiles(len(files))
					}

					media, err = mediaorient.CalculateFilesOrientation(files)

					if err != nil {
						return err
					}
					if len(media) == 0 {
						return nil
					}

					fmt.Printf("%+v\n", media)

					//groups := groupAndReport(media, threshold, output)
					//
					//switch output {
					//case "report":
					//	charm.PrintGroupReport(groups)
					//case "json":
					//	charm.PrintGroupJson(groups)
					//case "csv":
					//	charm.PrintGroupCsv(groups)
					//}

					return nil
				},
			},
			{
				Name:      "dir",
				Usage:     "determine/fix the orientation of files in a directory",
				UsageText: "mediaorient dir <directory> [-r] [--mt <media-type>]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "recursive",
						Aliases:     []string{"r"},
						Usage:       "recursively search for files in the directory",
						Value:       false,
						DefaultText: "false",
						Destination: &recursive,
					},
					&cli.StringFlag{
						Name:        "media-type",
						Aliases:     []string{"mt"},
						Usage:       "type of media to compare; image | video | all",
						Value:       "all",
						DefaultText: "all",
						Destination: &mediaType,
						Validator: func(s string) error {
							if s != "image" && s != "video" && s != "all" {
								return fmt.Errorf("invalid media type")
							}

							return nil
						},
					},
				},
				Action: func(ctx context.Context, command *cli.Command) error {
					directory = command.Args().First()
					directory, err = expandPath(directory)
					if err != nil {
						return nil
					}

					if output == "report" {
						charm.PrintCalculateDirectory(directory)
					}

					media, err = mediaorient.CalculateDirectoryOrientation(directory, mediaType, recursive)
					if err != nil {
						return err
					}
					if len(media) == 0 {
						return nil
					}

					fmt.Printf("%+v\n", media)

					//groups := groupAndReport(media, threshold, output)
					//
					//switch output {
					//case "report":
					//	charm.PrintGroupReport(groups)
					//case "json":
					//	charm.PrintGroupJson(groups)
					//case "csv":
					//	charm.PrintGroupCsv(groups)
					//}

					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "format how similarity is reported; report | json | csv",
				Value:       "report",
				DefaultText: "report",
				Destination: &output,
				Validator: func(s string) error {
					if s != "report" && s != "json" && s != "csv" {
						return fmt.Errorf("invalid output format")
					}

					return nil
				},
			},
			&cli.BoolFlag{
				Name:        "ignore-errors",
				Aliases:     []string{"ie"},
				Usage:       "continues processing files even if an error occurs",
				Value:       false,
				DefaultText: "false",
				Destination: &ignoreErrors,
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Aliases:     []string{"dr"},
				Usage:       "do not rotate any media; only prints the ones that would be modified",
				Value:       false,
				DefaultText: "false",
				Destination: &dryRun,
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			return fmt.Errorf("command missing; try 'mediaorient --help' for more information")
		},
	}
}
