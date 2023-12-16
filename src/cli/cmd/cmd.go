package cmd

import (
	"os"

	"github.com/aicirt2012/fileintegrity"
	"github.com/aicirt2012/fileintegrity/doc/license"
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     `fileintegrity`,
		Short:   `Creates, updates or verifies file integrity`,
		Version: fileintegrity.Version,
	}
	cmd.AddCommand(upsert())
	cmd.AddCommand(verify())
	cmd.AddCommand(check())
	cmd.AddCommand(licenseTxt())
	return cmd
}

func upsert() *cobra.Command {
	var quiet bool
	var cmd = &cobra.Command{
		Use:   `upsert <dir>`,
		Short: `Upsert integrity`,
		Long:  `Creates or updated integrity file if needed`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fileintegrity.Upsert(args[0], options(&quiet))
		},
	}
	addQuietFlag(cmd, &quiet)
	return cmd
}

func verify() *cobra.Command {
	var quiet bool
	var cmd = &cobra.Command{
		Use:   `verify <dir>`,
		Short: `Verify integrity`,
		Long:  `Verify integrity file if exist`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fileintegrity.Verify(args[0], options(&quiet))
		},
	}
	addQuietFlag(cmd, &quiet)
	return cmd
}

func check() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   `check`,
		Short: `Check integrity directory`,
		Long:  `Check integrity directory regarding several aspects`,
	}
	cmd.AddCommand(checkDuplicates())
	cmd.AddCommand(checkContained())
	cmd.AddCommand(checkStyleIssue())
	return cmd
}

func checkDuplicates() *cobra.Command {
	var quiet bool
	var cmd = &cobra.Command{
		Use:   `duplicates <dir>`,
		Short: `Check duplicates`,
		Long:  `Check duplicates within integrity file`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fileintegrity.CheckDuplicates(args[0], options(&quiet))
		},
	}
	addQuietFlag(cmd, &quiet)
	return cmd
}

func checkContained() *cobra.Command {
	var quiet, fix bool
	var cmd = &cobra.Command{
		Use:   `contained <dir> <externalDir>`,
		Short: `Check contained`,
		Long:  `Check contained within integrity file`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fileintegrity.CheckContained(args[0], args[1], fix, options(&quiet))
		},
	}
	cmd.Flags().BoolVarP(&fix, "fix", "f", false, "delete contained and duplicate files within the external directory")
	addQuietFlag(cmd, &quiet)
	return cmd
}

func checkStyleIssue() *cobra.Command {
	var quiet bool
	var cmd = &cobra.Command{
		Use:   `style <dir>`,
		Short: `Check style issues`,
		Long:  `Check style issues within integrity file`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fileintegrity.CheckStyleIssues(args[0], options(&quiet))
		},
	}
	addQuietFlag(cmd, &quiet)
	return cmd
}

func licenseTxt() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   `license`,
		Short: `license text`,
		Long:  `Shows the full license text`,
		Run: func(cmd *cobra.Command, args []string) {
			content, err := license.Text()
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}
			cmd.Println(content)
		},
	}
	return cmd
}

func addQuietFlag(cmd *cobra.Command, p *bool) {
	cmd.Flags().BoolVarP(p, "quiet", "q", false, "enable quiet mode")
}

func options(quiet *bool) fileintegrity.Options {
	return fileintegrity.LogOptions(quiet)
}
