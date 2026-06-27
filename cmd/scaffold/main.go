package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"scaffold/internal/engine"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "scaffold",
		Short:         "Generate a Go + React Router monorepo from a feature catalog",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newProjectCmd())
	return root
}

func newProjectCmd() *cobra.Command {
	var with, catalog string
	cmd := &cobra.Command{
		Use:   "new <target-dir>",
		Short: "Create a new project in <target-dir>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNew(args[0], catalog, with, cmd.Flags().Changed("with"))
		},
	}
	cmd.Flags().StringVar(&with, "with", "", "comma-separated optional features to enable (skips the interactive picker)")
	cmd.Flags().StringVar(&catalog, "catalog", "features", "path to the feature catalog directory")
	return cmd
}

func runNew(target, catalog, with string, withSet bool) error {
	cat, err := engine.LoadCatalog(catalog)
	if err != nil {
		return fmt.Errorf("loading catalog: %w", err)
	}

	var selected []string
	if withSet {
		for _, s := range strings.Split(with, ",") {
			if s = strings.TrimSpace(s); s != "" {
				selected = append(selected, s)
			}
		}
	} else {
		selected, err = chooseOptional(cat)
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Fprintln(os.Stderr, "cancelled")
			os.Exit(0)
		}
		if err != nil {
			return err
		}
	}

	plan, err := engine.Resolve(cat, selected)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", target, err)
	}
	if err := engine.Write(plan, target); err != nil {
		return fmt.Errorf("writing plan to %s: %w", target, err)
	}
	for _, step := range plan.PostSteps {
		fmt.Println("  >", step)
		cmd := exec.Command("bash", "-uc", step)
		cmd.Dir = target
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("post-step %q: %w", step, err)
		}
	}
	fmt.Printf("Created %s with features: %s\n", target, strings.Join(plan.Features, ", "))
	return nil
}

func chooseOptional(cat []engine.Feature) ([]string, error) {
	var opts []huh.Option[string]
	for _, f := range cat {
		if !f.Always {
			opts = append(opts, huh.NewOption(fmt.Sprintf("%s — %s", f.Name, f.Description), f.Name))
		}
	}
	if len(opts) == 0 {
		return nil, nil
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].Value < opts[j].Value })
	var chosen []string
	err := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title("Optional features").
			Options(opts...).
			Value(&chosen),
	)).Run()
	return chosen, err
}
