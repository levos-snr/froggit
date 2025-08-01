package git

import (
	"fmt"
	"strings"
)

// Add stages a file
func Add(filename string) error {
	return NewGitClient("").Add(filename)
}

func (g *GitClient) Add(filename string) error {
	_, err := g.runGitCommandCombinedOutput("add", filename)
	return err
}

// Reset unstages a file
func Reset(filename string) error {
	return NewGitClient("").Reset(filename)
}

func (g *GitClient) Reset(filename string) error {
	_, err := g.runGitCommandCombinedOutput("reset", "HEAD", filename)
	return err
}

// Commit creates a commit
func Commit(message string) error {
	return NewGitClient("").Commit(message)
}

func (g *GitClient) Commit(message string) error {
	_, err := g.runGitCommandCombinedOutput("commit", "-m", message)
	return err
}

func Merge(branch string) error {
	return NewGitClient("").Merge(branch)
}

func (g *GitClient) Merge(branch string) error {
	_, err := g.runGitCommandCombinedOutput("merge", branch)
	return err
}

func Rebase(branch string) error {
	return NewGitClient("").Rebase(branch)
}

func (g *GitClient) Rebase(branch string) error {
	_, err := g.runGitCommandCombinedOutput("rebase", branch)
	return err
}

func SaveStash(message string) error {
	return NewGitClient("").SaveStash(message)
}

func (g *GitClient) SaveStash(message string) error {
	_, err := g.runGitCommandCombinedOutput("stash", "push", "-m", message)
	return err
}

func StashPop() error {
	return NewGitClient("").StashPop()
}

func (g *GitClient) StashPop() error {
	_, err := g.runGitCommandCombinedOutput("stash", "pop")
	return err
}

func StashList() (string, error) {
	return NewGitClient("").StashList()
}

func (g *GitClient) StashList() (string, error) {
	output, err := g.runGitCommand("stash", "list")
	if err != nil {
		return "", fmt.Errorf("failed to list stashes: %w", err)
	}
	return string(output), nil
}

// HasRemoteChanges checks if local branch is behind remote
func HasRemoteChanges(branch string) (bool, error) {
	return NewGitClient("").HasRemoteChanges(branch)
}

func (g *GitClient) HasRemoteChanges(branch string) (bool, error) {
	if err := g.Fetch(); err != nil {
		return false, err
	}

	output, err := g.runGitCommand("rev-list", "--count", fmt.Sprintf("HEAD..origin/%s", branch))
	if err != nil {
		return false, err
	}

	count := strings.TrimSpace(string(output))
	return count != "0", nil
}

// LogsGraph returns a graph view of logs
func LogsGraph() (string, error) {
	return NewGitClient("").LogsGraph()
}

func (g *GitClient) LogsGraph() (string, error) {
	output, err := g.runGitCommand("log", "--graph", "--oneline", "--all")
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}
	return string(output), nil
}

// GetConflictFiles devuelve los archivos en conflicto usando git diff --name-only --diff-filter=U
func GetConflictFiles() ([]string, error) {
	output, err := NewGitClient("").runGitCommand("diff", "--name-only", "--diff-filter=U")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, l := range lines {
		if l != "" {
			files = append(files, l)
		}
	}
	return files, nil
}

// MergeContinue ejecuta git merge --continue
func MergeContinue() error {
	_, err := NewGitClient("").runGitCommandCombinedOutput("merge", "--continue")
	return err
}

// MergeAbort ejecuta git merge --abort
func MergeAbort() error {
	_, err := NewGitClient("").runGitCommandCombinedOutput("merge", "--abort")
	return err
}

// RebaseContinue runs git rebase --continue
func RebaseContinue() error {
	_, err := NewGitClient("").runGitCommandCombinedOutput("rebase", "--continue")
	return err
}

// RebaseAbort runs git rebase --abort
func RebaseAbort() error {
	_, err := NewGitClient("").runGitCommandCombinedOutput("rebase", "--abort")
	return err
}

func StashApply(stashRef string) error {
	return NewGitClient("").StashApply(stashRef)
}

func (g *GitClient) StashApply(stashRef string) error {
	_, err := g.runGitCommandCombinedOutput("stash", "apply", stashRef)
	return err
}

func StashDrop(stashRef string) error {
	return NewGitClient("").StashDrop(stashRef)
}

func (g *GitClient) StashDrop(stashRef string) error {
	_, err := g.runGitCommandCombinedOutput("stash", "drop", stashRef)
	return err
}

func StashShow(stashRef string) (string, error) {
	return NewGitClient("").StashShow(stashRef)
}

func (g *GitClient) StashShow(stashRef string) (string, error) {
	output, err := g.runGitCommand("stash", "show", "-p", stashRef)
	if err != nil {
		return "", fmt.Errorf("failed to show stash: %w", err)
	}
	return string(output), nil
}

// GetStashRef extracts stash reference from stash list line
func GetStashRef(stashLine string) string {
	parts := strings.Split(stashLine, ":")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return "stash@{0}"
}

func HasCommitsToush() (bool, error) {
	return NewGitClient("").HasCommitsToush()
}

func (g *GitClient) HasCommitsToush() (bool, error) {
	output, err := g.runGitCommand("rev-list", "--count", "HEAD")
	if err != nil {
		return false, nil
	}

	localCommits := strings.TrimSpace(string(output))
	if localCommits == "0" {
		return false, nil
	}

	output, err = g.runGitCommand("rev-list", "--count", "HEAD", "^@{u}")
	if err != nil {
		return localCommits != "0", nil
	}

	ahead := strings.TrimSpace(string(output))
	return ahead != "0", nil
}
