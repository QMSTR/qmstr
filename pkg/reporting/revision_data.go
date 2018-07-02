package reporting

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/QMSTR/qmstr/pkg/service"
)

// RevisionData contains metadata about a specific revision.
type RevisionData struct {
	VersionIdentifier      string      // Usually a Git hash, but any string can be used
	VersionIdentifierShort string      // The short version of the version identifier
	ChangeDateTime         string      // The change timestamp
	Author                 string      // The author of the change
	Message                string      // The commit message
	Summary                string      // The short form of the commit message
	Package                PackageData // The package this version is associated with.
}

// CommitMessageSummary returns the summary of the commit message according to the usual guidelines
// (see https://chris.beams.io/posts/git-commit/, "Limit the subject line to 50 characters")
func CommitMessageSummary(message string) string {
	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return ""
	}
	summary := strings.TrimSpace(lines[0])
	if len(summary) > 50 {
		summary = fmt.Sprintf("%s...", summary[:47])
	}
	return summary
}

// ShortenedVersionIdentifier calculates shortened version of the version identifier, in upper-case characters
func ShortenedVersionIdentifier(message string) string {
	const cutoff = 8 // this should be configurable
	result := strings.ToUpper(message)
	if len(result) <= cutoff {
		return result
	}
	return result[:cutoff]
}

// GetRevisionData extracts the revision data from the package node
func GetRevisionData(packageNode *service.PackageNode, packageData PackageData) (RevisionData, error) {
	revisionData := RevisionData{"(SHA long)", "(SHA #8)", "(commit datetime)", "(author)", "(commit message)", "(commit summary)", packageData}

	ps := reflect.ValueOf(&packageData)
	s := ps.Elem()
	for _, inode := range packageNode.AdditionalInfo {
		if inode.Type == "metadata" {
			for _, dnode := range inode.DataNodes {
				f := s.FieldByName(dnode.Type)
				if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
					f.SetString(dnode.Data)
				}
			}
		}
		if inode.Type == "Revision" {
			for _, dnode := range inode.DataNodes {
				switch dnode.Type {
				case "AuthorName":
					revisionData.Author = dnode.Data
				case "CommitMessage":
					revisionData.Message = dnode.Data
				case "CommitID":
					log.Printf("WARN: using CommitID instead of description this can be misleading as it does not cover not commited changes")
					revisionData.VersionIdentifier = dnode.Data
				case "CommitterDate":
					revisionData.ChangeDateTime = dnode.Data
				}
			}
		}
	}
	revisionData.Summary = CommitMessageSummary(revisionData.Message)
	revisionData.VersionIdentifierShort = ShortenedVersionIdentifier(revisionData.VersionIdentifier)
	return revisionData, nil
}
