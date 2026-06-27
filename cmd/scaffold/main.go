package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"

	"scaffold/internal/engine"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "new" {
		fmt.Fprintln(os.Stderr, "usage: scaffold new <target-dir> [--with a,b] [--catalog DIR]")
		os.Exit(2)
	}
	if err := runNew(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func runNew(args []string) error {
	var target, with, catalog string
	catalog = "features"
	withSet := false
	rest := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--with":
			i++
			if i >= len(args) {
				return fmt.Errorf("--with requires a value")
			}
			with = args[i]
			withSet = true
		case "--catalog":
			i++
			if i >= len(args) {
				return fmt.Errorf("--catalog requires a value")
			}
			catalog = args[i]
		default:
			rest = append(rest, args[i])
		}
	}
	if len(rest) != 1 {
		return fmt.Errorf("expected exactly one target directory")
	}
	target = rest[0]

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
