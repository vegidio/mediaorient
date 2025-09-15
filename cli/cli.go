package main

import (
	"cli/internal/charm"
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	. "github.com/vegidio/mediaorient"
)

func buildCliCommands() *cli.Command {
	var media []Media
	var files []string
	var directory string
	var output string
	var recursive bool
	var mediaType string
	var ignoreErrors bool
	var includeZero bool
	var autoFix bool
	var err error

	return &cli.Command{
		Name:            "mediaorient",
		Usage:           "a tool to calculate the orientation of images & videos",
		UsageText:       "mediaorient <command>",
		Version:         Version,
		HideHelpCommand: true,
		Commands: []*cli.Command{
			{
				Name:      "files",
				Usage:     "calculate the orientation of one or more files",
				UsageText: "mediaorient files <file1> [<file2> ...] ",
				Flags:     []cli.Flag{},
				Action: func(ctx context.Context, command *cli.Command) error {
					files = command.Args().Slice()
					amount := len(files)
					includeZero = true

					if amount < 1 {
						return fmt.Errorf("at least one files must be specified")
					}

					if output == "report" {
						charm.PrintCalculateFiles(amount)
					}

					files = lo.Map(files, func(file string, _ int) string {
						fullFile, _ := expandPath(file)
						return fullFile
					})

					result := CalculateFilesOrientation(files)
					media, err = getMedia(result, len(files), output)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:      "dir",
				Usage:     "calculate the orientation of files in a directory",
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
						Usage:       "type of media to calculate; image | video | all",
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
					includeZero = false

					if output == "report" {
						charm.PrintCalculateDirectory(directory)
					}

					directory, err = expandPath(directory)
					if err != nil {
						return nil
					}

					result, total := CalculateDirectoryOrientation(directory, mediaType, recursive)
					media, err = getMedia(result, total, output)
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "format how orientation is reported; report | json | csv",
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
				Name:        "auto-fix",
				Aliases:     []string{"af"},
				Usage:       "automatically fix orientation of files",
				Value:       false,
				DefaultText: "false",
				Destination: &autoFix,
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			return fmt.Errorf("command missing; try 'mediaorient --help' for more information")
		},
		After: func(ctx context.Context, command *cli.Command) error {
			charm.PrintReport(media, includeZero)
			return nil
		},
	}
}

// region - Private functions

func getMedia(
	result <-chan Result[Media],
	total int,
	output string,
) ([]Media, error) {
	var media []Media
	var err error

	if output == "report" {
		media, err = charm.StartProgress(result, total)
		if err != nil {
			return nil, err
		}
	} else {
		results := lo.ChannelToSlice(result)
		media = lo.FilterMap(results, func(r Result[Media], _ int) (Media, bool) {
			return r.Data, r.IsSuccess()
		})
	}

	return media, nil
}

// endregion
