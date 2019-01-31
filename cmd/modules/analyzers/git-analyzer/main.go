package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/QMSTR/go-qmstr/analysis"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/go-qmstr/service"
	git "gopkg.in/libgit2/git2go.v26"
)

type Revision struct {
	Description    string
	CommitID       string
	CommitMessage  string
	CommitterName  string
	CommitterEmail string
	CommitterDate  string
	AuthorName     string
	AuthorEmail    string
	AuthorDate     string
}

type GitAnalyzer struct {
	repo     *git.Repository
	revision Revision
}

func main() {
	analyzer := analysis.NewAnalyzer(&GitAnalyzer{})
	if err := analyzer.RunAnalyzerModule(); err != nil {
		log.Printf("%v failed: %v\n", analyzer.GetModuleName(), err)
		os.Exit(master.ReturnAnalyzerFailed)
	}
}

func (ga *GitAnalyzer) Configure(configMap map[string]string) error {
	if workdir, ok := configMap["workdir"]; ok {
		repo, err := git.OpenRepository(workdir)
		if err != nil {
			log.Fatalf("Failed to open git repository: %v", err)
		}
		ga.repo = repo
		ga.revision = Revision{}
		return nil
	}
	return fmt.Errorf("Misconfigured git analyzer")
}

func (ga *GitAnalyzer) Analyze(controlService service.ControlServiceClient, analysisService service.AnalysisServiceClient, token int64) error {
	ga.fillRevision()
	log.Printf("Found git revision %s", ga.revision)
	pkgNode, err := controlService.GetPackageNode(context.Background(), &service.PackageRequest{})
	tempDataNodes := []*service.InfoNode_DataNode{}

	v := reflect.ValueOf(ga.revision)

	for i := 0; i < v.NumField(); i++ {
		tempDataNodes = append(tempDataNodes, &service.InfoNode_DataNode{
			Type: v.Type().Field(i).Name,
			Data: v.Field(i).Interface().(string),
		})
	}

	gitNodes := &service.InfoNode{
		Type:      "Revision",
		DataNodes: tempDataNodes,
	}

	infoNodeMsg := &service.InfoNodeMessage{Token: token, Infonode: gitNodes, Uid: pkgNode.Uid}

	send_stream, err := analysisService.SendInfoNodes(context.Background())
	if err != nil {
		return err
	}

	send_stream.Send(infoNodeMsg)

	reply, err := send_stream.CloseAndRecv()
	if err != nil {
		return err
	}
	if reply.Success {
		log.Println("Git analyzer sent InfoNodes")
	}

	return nil
}

func (ga *GitAnalyzer) obtainDesc() {

	descRes, err := ga.repo.DescribeWorkdir(&git.DescribeOptions{Strategy: git.DescribeAll, ShowCommitOidAsFallback: true})
	if err != nil {
		log.Printf("Failed to describe workdir: %v", err)
		return
	}

	desc, err := descRes.Format(&git.DescribeFormatOptions{AlwaysUseLongFormat: true, AbbreviatedSize: 7, DirtySuffix: "dirty"})
	if err != nil {
		log.Printf("Failed to format description %v", err)
		return
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
	ga.revision.CommitterDate = commit.Committer().When.String()

	ga.revision.AuthorName = commit.Author().Name
	ga.revision.AuthorEmail = commit.Author().Email
	ga.revision.AuthorDate = commit.Author().When.String()

	ga.obtainDesc()

	return nil
}

func (ga *GitAnalyzer) PostAnalyze() error {
	return nil
}
