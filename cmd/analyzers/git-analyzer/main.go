package main

import (
	"log"
	"os"
	"time"

	git "gopkg.in/libgit2/git2go.v26"
)

type Revision struct {
	Description    string
	CommitID       string
	CommitMessage  string
	CommitterName  string
	CommitterEmail string
	CommitterDate  time.Time
	AuthorName     string
	AuthorEmail    string
	AuthorDate     time.Time
}

type GitAnalyzer struct {
	repo     *git.Repository
	revision Revision
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Oh no %v", err)
	}
	ga := NewGitAnalyzer(cwd)
	ga.fillRevision()
	log.Printf("You git %s", ga.revision)
}

func NewGitAnalyzer(workdir string) *GitAnalyzer {
	repo, err := git.OpenRepository(workdir)
	if err != nil {
		log.Fatalf("Oh no %v", err)
	}
	return &GitAnalyzer{repo: repo, revision: Revision{}}
}

func (ga *GitAnalyzer) obtainDesc() {

	descRes, err := ga.repo.DescribeWorkdir(&git.DescribeOptions{Strategy: git.DescribeAll})
	if err != nil {
		log.Fatalf("oh no %v", err)
	}

	desc, err := descRes.Format(&git.DescribeFormatOptions{AlwaysUseLongFormat: true, AbbreviatedSize: 7, DirtySuffix: "dirty"})
	if err != nil {
		log.Fatalf("oh no %v", err)
	}

	ga.revision.Description = desc
}

func (ga *GitAnalyzer) fillRevision() error {
	ref, err := ga.repo.Head()
	if err != nil {
		return err
	}
	oid := ref.Target()
	commit, err := ga.repo.LookupCommit(oid)
	if err != nil {
		return err
	}

	ga.revision.CommitMessage = commit.Message()
	ga.revision.CommitID = commit.Id().String()

	ga.revision.CommitterName = commit.Committer().Name
	ga.revision.CommitterEmail = commit.Committer().Email
	ga.revision.CommitterDate = commit.Committer().When

	ga.revision.AuthorName = commit.Author().Name
	ga.revision.AuthorEmail = commit.Author().Email
	ga.revision.AuthorDate = commit.Author().When

	ga.obtainDesc()

	return nil
}
