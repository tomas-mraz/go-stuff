package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

const defaultRepoURL = "https://github.com/tomas-mraz/go-stuff.git"
const defaultFromTag = "1.0.0"
const defaultToTag = "1.0.1"

type FileChange struct {
	Action  string `json:"action"`
	Path    string `json:"path"`
	OldPath string `json:"old_path,omitempty"`
}

func main() {
	fromTag := defaultFromTag
	toTag := defaultToTag

	files, err := fetchRepositoryChanges(defaultRepoURL, fromTag, toTag)
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(files); err != nil {
		log.Fatalf("encode result: %v", err)
	}
}

func fetchRepositoryChanges(repoURL, fromTag, toTag string) ([]FileChange, error) {
	repo, err := fetchTagsRepository(repoURL, fromTag, toTag)
	if err != nil {
		return nil, err
	}

	fromCommit, err := resolveTagCommit(repo, fromTag)
	if err != nil {
		return nil, err
	}

	toCommit, err := resolveTagCommit(repo, toTag)
	if err != nil {
		return nil, err
	}

	fromTree, err := fromCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("load tree for tag %q: %w", fromTag, err)
	}

	toTree, err := toCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("load tree for tag %q: %w", toTag, err)
	}

	changes, err := object.DiffTreeWithOptions(
		context.Background(),
		fromTree,
		toTree,
		object.DefaultDiffTreeOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("diff trees between tags %q and %q: %w", fromTag, toTag, err)
	}

	return toFileChanges(changes)
}

func fetchTagsRepository(repoURL string, tags ...string) (*git.Repository, error) {
	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return nil, fmt.Errorf("init in-memory repository: %w", err)
	}

	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{repoURL},
	}); err != nil {
		return nil, fmt.Errorf("create remote: %w", err)
	}

	refSpecs := make([]config.RefSpec, 0, len(tags))
	for _, tag := range tags {
		refSpecs = append(refSpecs, config.RefSpec(tagRefSpec(tag)))
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: git.DefaultRemoteName,
		Depth:      1,
		Tags:       git.NoTags,
		RefSpecs:   refSpecs,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, fmt.Errorf("fetch tags %v: %w", tags, err)
	}

	return repo, nil
}

func resolveTagCommit(repo *git.Repository, tagName string) (*object.Commit, error) {
	revision := plumbing.Revision(plumbing.NewTagReferenceName(tagName).String())
	hash, err := repo.ResolveRevision(revision)
	if err != nil {
		return nil, fmt.Errorf("resolve tag %q: %w", tagName, err)
	}

	commit, err := repo.CommitObject(*hash)
	if err != nil {
		return nil, fmt.Errorf("load commit for tag %q: %w", tagName, err)
	}

	return commit, nil
}

func tagRefSpec(tag string) string {
	ref := plumbing.NewTagReferenceName(tag).String()
	return "+" + ref + ":" + ref
}

func toFileChanges(changes object.Changes) ([]FileChange, error) {
	result := make([]FileChange, 0, len(changes))

	for _, change := range changes {
		if !isFileChange(change) {
			continue
		}

		action, err := change.Action()
		if err != nil {
			return nil, fmt.Errorf("resolve change action: %w", err)
		}

		fileChange := FileChange{}
		switch action {
		case merkletrie.Insert:
			fileChange.Action = "add"
			fileChange.Path = change.To.Name
		case merkletrie.Delete:
			fileChange.Action = "remove"
			fileChange.Path = change.From.Name
		case merkletrie.Modify:
			if change.From.Name != change.To.Name {
				fileChange.Action = "rename"
				fileChange.OldPath = change.From.Name
				fileChange.Path = change.To.Name
			} else {
				fileChange.Action = "modify"
				fileChange.Path = change.To.Name
			}
		default:
			return nil, fmt.Errorf("unsupported change action %v", action)
		}

		result = append(result, fileChange)
	}

	return result, nil
}

func isFileChange(change *object.Change) bool {
	return change.From.TreeEntry.Mode.IsFile() || change.To.TreeEntry.Mode.IsFile()
}
