package openapi

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/extension-marketplace-service/shared"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func cloneRepo(repoPath string) (*git.Repository, *ssh.PublicKeys, error) {
	// Clone the repository

	var sshKey []byte
	var err error
	if shared.SetupEnv() == shared.DEV {
		sshKey, err = os.ReadFile(DEPLOY_KEY_PATH)
		if err != nil {
			log.Fatalf("Failed to read deploy key file: %v", err)
		}
	} else {
		sshKey = getPemFromEnv()
	}

	publicKeys, err := ssh.NewPublicKeys("git", sshKey, "")

	if err != nil {
		log.Fatalf("Failed to create public keys from private key: %v", err)
	}

	env := shared.SetupEnv()
	files := ""
	if env == shared.STAGE || env == shared.PROD {
		files = "root/.ssh/known_hosts"
	} else {
		home := os.Getenv("HOME")
		files = fmt.Sprintf("%s/.ssh/known_hosts", home)
	}

	publicKeys.HostKeyCallback, err = ssh.NewKnownHostsCallback(files)
	if err != nil {
		log.Fatalf("Failed to create public keys from private key: %v", err)
	}

	repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      "git@github.com:serranolabs-io/bookera-extension-hub.git",
		Auth:     publicKeys,
		Progress: os.Stdout,
	})

	return repo, publicKeys, err
}

func openRepo(cloneErr error, repoPath string) *git.Repository {
	var err error
	var repo *git.Repository
	if cloneErr == git.ErrRepositoryAlreadyExists {
		log.Info("Repository already exists, proceeding with existing repository")
		repo, err = git.PlainOpen(repoPath)
		if err != nil {
			log.Fatalf("Failed to open existing repository: %v", err)
		}
	} else {
		log.Fatalf("Failed to clone repository: %v", cloneErr)
	}

	return repo
}

const REPO_PATH = "./temp-repo"

func initGitRepo() (*git.Repository, *git.Worktree, *ssh.PublicKeys, string) {
	repo, publicKeys, err := cloneRepo(REPO_PATH)

	if err != nil {
		repo = openRepo(err, REPO_PATH)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Failed to get worktree: %v", err)
	}

	return repo, worktree, publicKeys, REPO_PATH
}

func delRepo() {
	err := os.RemoveAll(REPO_PATH)
	if err != nil {
		log.Fatalf("Failed to delete repository folder: %v", err)
	}
	log.Info("Repository folder deleted successfully")
}

func commitChanges(commitMessage string, worktree *git.Worktree, repo *git.Repository, publicKeys *ssh.PublicKeys) {
	// Commit the changes
	_, err := worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Extension Marketplace Bot",
			Email: "daviddserranodev@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Fatalf("Failed to commit changes: %v", err)
	}

	// Push the changes to the remote repository
	err = repo.Push(&git.PushOptions{
		Auth: publicKeys,
	})
	if err != nil {
		if err.Error() == "non-fast-forward update: refs/heads/main" {
			log.Warn("Non-fast-forward update detected, changes not pushed")
			err = repo.Fetch(&git.FetchOptions{
				Auth: publicKeys,
			})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				log.Fatalf("Failed to fetch changes: %v", err)
			}
			err = worktree.Pull(&git.PullOptions{
				Auth:       publicKeys,
				RemoteName: "origin",
			})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				log.Fatalf("Failed to rebase changes: %v", err)
			}

			err = repo.Push(&git.PushOptions{
				Auth: publicKeys,
			})

			if err != nil {
				log.Fatalf("Failed to push changes after rebase: %v", err)
			}

		} else {
			log.Fatalf("Failed to push changes: %v", err)
		}
	}

	log.Info(commitMessage + " and pushed to repository")
}

func getPemFromEnv() []byte {
	b64 := os.Getenv("PEM_KEY")

	str, _ := base64.StdEncoding.DecodeString(b64)

	return str
}
