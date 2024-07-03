package main

import (
	"fmt"
	"strings"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

func ChangesetExample() {
	var packages, majorPackages, minorPackages, patchPackages []string

	prompts.Intro(picocolors.BgCyan(picocolors.Black(" changesets ")))

	packages, err := prompts.GroupMultiSelect(prompts.GroupMultiSelectParams[string]{
		Message: "Which packages would you like to include?",
		Options: map[string][]prompts.MultiSelectOption[string]{
			"changed packages": {
				{Label: "@scope/a"},
				{Label: "@scope/b"},
				{Label: "@scope/c"},
			},
			"unchanged packages": {
				{Label: "@scope/x"},
				{Label: "@scope/y"},
				{Label: "@scope/z"},
			},
		},
	})
	HandleCancel(err)

	packageOptions := make([]*prompts.MultiSelectOption[string], len(packages))
	for i, packageOption := range packages {
		packageOptions[i] = &prompts.MultiSelectOption[string]{
			Label: packageOption,
		}
	}
	majorPackages, err = prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: fmt.Sprintf("Which packages should have a %s bump?", picocolors.Red("major")),
		Options: packageOptions,
	})
	HandleCancel(err)

	minorOptions := []*prompts.MultiSelectOption[string]{}
packagesLoop:
	for _, packageOption := range packages {
		for _, majorPackage := range majorPackages {
			if majorPackage == packageOption {
				continue packagesLoop
			}
		}
		minorOptions = append(minorOptions, &prompts.MultiSelectOption[string]{
			Label: packageOption,
		})
	}

	if len(minorOptions) > 0 {
		minorPackages, err = prompts.MultiSelect(prompts.MultiSelectParams[string]{
			Message: fmt.Sprintf("Which packages should have a %s bump?", picocolors.Yellow("minor")),
			Options: minorOptions,
		})
		HandleCancel(err)
	}

packagesLoop2:
	for _, _package := range packages {
		for _, majorPackage := range majorPackages {
			if majorPackage == _package {
				continue packagesLoop2
			}
		}
		for _, minorPackage := range minorPackages {
			if minorPackage == _package {
				continue packagesLoop2
			}
		}
		patchPackages = append(patchPackages, _package)
	}

	if len(patchPackages) > 0 {
		note := strings.Join(patchPackages, picocolors.Dim(", "))
		prompts.Step(fmt.Sprintf("These packages will have a %s bump: %s", picocolors.Green("patch"), picocolors.Dim(note)))
	}

	_, err = prompts.Text(prompts.TextParams{
		Message:     "Please enter a summary for this change",
		Placeholder: "Summary",
	})
	HandleCancel(err)

	prompts.Outro(fmt.Sprintf("Changeset added! %s", picocolors.Underline(picocolors.Cyan(".changeset/orange-crabs-sing.md"))))
}
