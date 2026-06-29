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
	root.AddCommand(addCmd())
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
	if err := runPostSteps(plan.PostSteps, target); err != nil {
		return err
	}
	fmt.Printf("Created %s with features: %s\n", target, strings.Join(plan.Features, ", "))
	return nil
}

func runPostSteps(steps []string, dir string) error {
	for _, step := range steps {
		fmt.Println("  >", step)
		cmd := exec.Command("bash", "-uc", step)
		cmd.Dir = dir
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("post-step %q: %w", step, err)
		}
	}
	return nil
}

func addCmd() *cobra.Command {
	var dir, catalog string
	cmd := &cobra.Command{
		Use:   "add <feature>...",
		Short: "Add one or more optional features to an existing project",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(dir, catalog, args)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "path to the existing project")
	cmd.Flags().StringVar(&catalog, "catalog", "features", "path to the feature catalog directory")
	return cmd
}

func runAdd(dir, catalog string, requested []string) error {
	cat, err := engine.LoadCatalog(catalog)
	if err != nil {
		return fmt.Errorf("loading catalog: %w", err)
	}
	applied, err := engine.ReadManifest(dir)
	if err != nil {
		return fmt.Errorf("reading %s/scaffold.toml (is this a scaffold project?): %w", dir, err)
	}

	known := map[string]engine.Feature{}
	for _, f := range cat {
		known[f.Name] = f
	}
	appliedSet := map[string]bool{}
	for _, n := range applied {
		appliedSet[n] = true
	}

	var toAdd []string
	for _, name := range requested {
		f, ok := known[name]
		if !ok {
			return fmt.Errorf("unknown feature %q", name)
		}
		if f.Always {
			return fmt.Errorf("feature %q is always-on and is part of every project", name)
		}
		if appliedSet[name] {
			fmt.Printf("  - %s is already applied; skipping\n", name)
			continue
		}
		toAdd = append(toAdd, name)
	}
	if len(toAdd) == 0 {
		fmt.Println("Nothing to add.")
		return nil
	}

	// Resolve the full optional set (already-applied optionals plus the new
	// ones) so conflict/requirement checks see the whole picture.
	var selectedOptional []string
	for _, name := range applied {
		if f, ok := known[name]; ok && !f.Always {
			selectedOptional = append(selectedOptional, name)
		}
	}
	selectedOptional = append(selectedOptional, toAdd...)

	plan, err := engine.Resolve(cat, selectedOptional)
	if err != nil {
		return err
	}

	addSet := map[string]bool{}
	for _, name := range toAdd {
		addSet[name] = true
	}
	if err := engine.AddFeatures(plan, dir, addSet); err != nil {
		return fmt.Errorf("applying features to %s: %w", dir, err)
	}

	// Run the newly-added features' post-steps, plus the always-on features'
	// post-steps as finalizers — `go mod tidy` / `bun install` must re-resolve
	// the dependencies the added code introduced. They are idempotent. Ordered
	// the same as `new`: always-on first, then the added optional features.
	var steps []string
	for _, name := range plan.Features {
		if known[name].Always || addSet[name] {
			steps = append(steps, known[name].PostSteps...)
		}
	}
	if err := runPostSteps(steps, dir); err != nil {
		return err
	}

	fmt.Printf("Added %s to %s. Project features: %s\n",
		strings.Join(toAdd, ", "), dir, strings.Join(plan.Features, ", "))
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
