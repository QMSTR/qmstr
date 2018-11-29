package spdxreporter

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
	"github.com/spdx/tools-golang/v0/spdx"
	"github.com/spdx/tools-golang/v0/tvsaver"
)

const (
	ModuleName  = "reporter-spdx"
	outFileName = "%s.spdx"
)

type SPDXReporter struct {
	enableWarnings bool
	enableErrors   bool
	outputdir      string
	nsURI          string
}

func (r *SPDXReporter) Configure(config map[string]string) error {
	if outDir, ok := config["outputdir"]; ok {
		r.outputdir = outDir
	} else {
		return fmt.Errorf("no output directory configured")
	}

	if nsURI, ok := config["namespaceURI"]; ok {
		r.nsURI = nsURI
	} else {
		// see: https://spdx.org/spdx-specification-21-web-version#h.1gdfkutofa90
		r.nsURI = "http://spdx.org/spdxdocs/"
	}

	return nil
}

func (r *SPDXReporter) Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient, session string) error {

	bom, err := rserv.GetBOM(context.Background(), &service.BOMRequest{Session: session, Warnings: r.enableWarnings, Errors: r.enableErrors})
	if err != nil {
		return err
	}

	files := []*spdx.File2_1{}
	hashes := []string{}

	for _, trgt := range bom.Targets {
		fl := &spdx.File2_1{
			FileName: trgt.Name,
			// this should be unique
			FileSPDXIdentifier: "SPDXRef-file-" + trgt.Name,
			FileChecksumSHA1:   trgt.Sha1,
			LicenseConcluded:   "NOASSERTION",
			LicenseInfoInFile:  []string{"NOASSERTION"},
			FileCopyrightText:  "NOASSERTION",
		}
		files = append(files, fl)
		hashes = append(hashes, trgt.Sha1)
	}

	dnldLocation := bom.PackageInfo.SourceURL
	if dnldLocation == "" {
		dnldLocation = "NOASSERTION"
	}
	pkg := &spdx.Package2_1{
		PackageName: bom.PackageInfo.Name,
		// this should be unique
		PackageSPDXIdentifier:   "SPDXRef-pkg-" + bom.PackageInfo.Name,
		PackageDownloadLocation: dnldLocation,
		FilesAnalyzed:           true,
		PackageVerificationCode: calcSHA1Hash(hashes),
		// this is the license we detect
		PackageLicenseConcluded:     "NOASSERTION",
		PackageLicenseInfoFromFiles: []string{"NOASSERTION"},
		PackageLicenseDeclared:      bom.PackageInfo.LicenseDeclared,
		PackageCopyrightText:        "NOASSERTION",
		Files:                       files,
	}

	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:       "SPDX-2.1",
			DataLicense:       "CC0-1.0",
			SPDXIdentifier:    "SPDXRef-DOCUMENT",
			DocumentName:      bom.PackageInfo.Name,
			DocumentNamespace: r.nsURI + url.PathEscape(bom.PackageInfo.Name),
			Created:           time.Now().Format(time.RFC3339),
			CreatorTools: []string{
				"QMSTR",
			},
		},
		Packages: []*spdx.Package2_1{pkg},
	}

	fName := filepath.Join(r.outputdir, fmt.Sprintf(outFileName, common.GetPosixFullyPortableFilename(bom.PackageInfo.Name)))
	out, err := os.Create(fName)
	if err != nil {
		return fmt.Errorf("failed to create out file %q: %v", fName, err)
	}
	defer out.Close()

	// export to tag-value format
	err = tvsaver.Save2_1(doc, out)
	if err != nil {
		return fmt.Errorf("failed to export SPDX doc: %v", err)
	}
	return nil
}

func (r *SPDXReporter) PostReport() error {
	return nil
}

func calcSHA1Hash(l []string) string {
	// sort the list
	sort.Strings(l)

	h := sha1.New()
	for _, s := range l {
		h.Write([]byte(s))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
