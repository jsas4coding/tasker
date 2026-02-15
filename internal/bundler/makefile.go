package bundler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tasker.jsas.dev/internal/constants"
	"tasker.jsas.dev/internal/resolver"
)

// makefileShellEscape escapes a string for safe use inside double-quoted
// shell strings in Makefile recipes. This prevents injection via name and
// description fields that are interpolated into @echo commands.
func makefileShellEscape(s string) string {
	r := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		`$`, `$$`,
		"`", "\\`",
		"\n", " ",
		"\r", "",
	)
	return r.Replace(s)
}

// makefileCommentSafe sanitizes a string for use in Makefile ## comments.
// Newlines would start a new Makefile rule, so they must be removed.
func makefileCommentSafe(s string) string {
	r := strings.NewReplacer(
		"\n", " ",
		"\r", "",
	)
	return r.Replace(s)
}

// GenerateMakefile creates Makefile content from a resolved project.
func GenerateMakefile(project *resolver.ResolvedProject) []byte {
	var b strings.Builder

	b.WriteString(constants.HeaderGenerated + "\n")
	b.WriteString(constants.HeaderSource + "\n")
	b.WriteString(constants.HeaderRegenerate + "\n\n")

	b.WriteString("-include .env\n")
	b.WriteString("export\n\n")

	// Collect all targets for .PHONY
	var targets []string
	targets = append(targets, "help")

	groupKeys := resolver.SortedGroupKeys(project.Groups)
	for _, groupKey := range groupKeys {
		tasks := project.Groups[groupKey]
		for _, rt := range tasks {
			targets = append(targets, TaskColonToDash(rt.FullKey))
		}
	}

	b.WriteString(".PHONY: ")
	b.WriteString(strings.Join(targets, " "))
	b.WriteString("\n\n")

	b.WriteString(".DEFAULT_GOAL := help\n\n")

	// Help target
	b.WriteString("help: ## Show available targets\n")
	b.WriteString(fmt.Sprintf("\t@echo \"%s - %s\"\n", makefileShellEscape(project.Config.Name), makefileShellEscape(project.Config.Description)))
	b.WriteString("\t@echo \"\"\n")

	for _, groupKey := range groupKeys {
		group := project.Config.Groups[groupKey]
		tasks := project.Groups[groupKey]

		b.WriteString(fmt.Sprintf("\t@echo \"%s  %s\"\n", makefileShellEscape(group.Name), makefileShellEscape(group.Description)))

		for _, rt := range tasks {
			target := TaskColonToDash(rt.FullKey)
			b.WriteString(fmt.Sprintf("\t@echo \"  %-30s %s\"\n", target, makefileShellEscape(rt.Description)))
		}
		b.WriteString("\t@echo \"\"\n")
	}
	b.WriteString("\n")

	// Task targets
	for _, groupKey := range groupKeys {
		tasks := project.Groups[groupKey]
		for _, rt := range tasks {
			target := TaskColonToDash(rt.FullKey)
			b.WriteString(fmt.Sprintf("%s: ## %s\n", target, makefileCommentSafe(rt.Description)))

			if rt.Environment != "" {
				b.WriteString(fmt.Sprintf("\t@if [ -n \"$$ENV\" ] && [ \"$$ENV\" != \"%s\" ]; then echo \"Task %s requires ENV=%s (current: $$ENV)\"; exit 1; fi\n",
					rt.Environment, rt.FullKey, rt.Environment))
			}

			for _, cmd := range rt.Cmds {
				b.WriteString(fmt.Sprintf("\t%s\n", cmd))
			}
			b.WriteString("\n")
		}
	}

	return []byte(b.String())
}

// WriteMakefile writes the generated Makefile to disk.
func WriteMakefile(project *resolver.ResolvedProject, dir string) error {
	data := GenerateMakefile(project)
	return os.WriteFile(filepath.Join(dir, constants.MakefileOutput), data, constants.FilePermissions)
}
