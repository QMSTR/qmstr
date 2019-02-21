package reporting

import (
	"fmt"
	"strings"

	"github.com/QMSTR/qmstr/pkg/service"
)

// RevisionData contains metadata about a specific revision.
type RevisionData struct {
	VersionIdentifier string       // Usually a Git hash, but any string can be used
	ChangeDateTime    string       // The change timestamp
	Author            string       // The author of the change
	Message           string       // The commit message
	Package           *PackageData // The package this version is associated with.
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
	// escape double quotes
	summary = strings.Replace(summary, `\`, `\\`, -1)
	summary = strings.Replace(summary, `"`, `\"`, -1)
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

// GetRevisionData extracts the revision data from the BOM
func GetRevisionData(bom *service.BOM, packageData *PackageData) *RevisionData {
	return &RevisionData{
		bom.VersionInfo.Id,
		bom.VersionInfo.CommitDate,
		bom.VersionInfo.Author.Name,
		bom.VersionInfo.Message,
		packageData,
	}
}
