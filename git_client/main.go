package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

const defaultRepoURL = "https://github.com/tomas-mraz/go-stuff.git"

type FileChange struct {
	Action  string `json:"action"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	OldPath string `json:"old_path,omitempty"`
	OldName string `json:"old_name,omitempty"`
}

func main() {
	refName := flag.String("ref", "", "branch or tag to clone; default is remote HEAD")
	flag.Parse()

	files, err := fetchRepositoryChanges(defaultRepoURL, *refName)
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(files); err != nil {
		log.Fatalf("encode result: %v", err)
	}
}

func fetchRepositoryChanges(repoURL, refName string) ([]FileChange, error) {
	repo, err := cloneRepository(repoURL, refName)
	if err != nil {
		return nil, err
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("resolve HEAD: %w", err)
	}

	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("load HEAD commit: %w", err)
	}

	headTree, err := headCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("load HEAD tree: %w", err)
	}

	var parentTree *object.Tree
	if headCommit.NumParents() > 0 {
		parentCommit, err := headCommit.Parent(0)
		if err != nil {
			return nil, fmt.Errorf("load parent commit: %w", err)
		}

		parentTree, err = parentCommit.Tree()
		if err != nil {
			return nil, fmt.Errorf("load parent tree: %w", err)
		}
	}

	changes, err := object.DiffTreeWithOptions(
		context.Background(),
		parentTree,
		headTree,
		object.DefaultDiffTreeOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("diff commit trees: %w", err)
	}

	return toFileChanges(changes)
}

func cloneRepository(repoURL, refName string) (*git.Repository, error) {
	if refName == "" {
		return git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL:           repoURL,
			Depth:         2,
			SingleBranch:  true,
			ReferenceName: plumbing.HEAD,
			Tags:          git.NoTags,
		})
	}

	branchRepo, branchErr := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           repoURL,
		Depth:         2,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(refName),
		Tags:          git.NoTags,
	})
	if branchErr == nil {
		return branchRepo, nil
	}

	tagRepo, tagErr := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           repoURL,
		Depth:         2,
		SingleBranch:  true,
		ReferenceName: plumbing.NewTagReferenceName(refName),
	})
	if tagErr == nil {
		return tagRepo, nil
	}

	return nil, fmt.Errorf("clone ref %q as branch failed: %v; as tag failed: %w", refName, branchErr, tagErr)
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
			fileChange.Name = path.Base(change.To.Name)
		case merkletrie.Delete:
			fileChange.Action = "remove"
			fileChange.Path = change.From.Name
			fileChange.Name = path.Base(change.From.Name)
		case merkletrie.Modify:
			if change.From.Name != change.To.Name {
				fileChange.Action = "rename"
				fileChange.OldPath = change.From.Name
				fileChange.OldName = path.Base(change.From.Name)
				fileChange.Path = change.To.Name
				fileChange.Name = path.Base(change.To.Name)
			} else {
				fileChange.Action = "modify"
				fileChange.Path = change.To.Name
				fileChange.Name = path.Base(change.To.Name)
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
