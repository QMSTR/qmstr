package spdxreporter

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
	"github.com/spdx/tools-golang/v0/spdx"
	"github.com/spdx/tools-golang/v0/tvsaver"
)

const (
	ModuleName = "reporter-spdx"
)

type SPDXReporter struct {
	enableWarnings bool
	enableErrors   bool
	outputdir      string
}

func (r *SPDXReporter) Configure(config map[string]string) error {
	if outDir, ok := config["outputdir"]; ok {
		r.outputdir = outDir
	} else {
		return fmt.Errorf("no output directory configured")
	}

	return nil
}

func (r *SPDXReporter) Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient, session string) error {

	bom, err := rserv.GetBOM(context.Background(), &service.BOMRequest{Session: session, Warnings: r.enableWarnings, Errors: r.enableErrors})
	if err != nil {
		return err
	}
	log.Printf("%v", bom)

	files := []*spdx.File2_1{}

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
	}

	pkg := &spdx.Package2_1{
		PackageName: bom.PackageInfo.Name,
		// this should be unique
		PackageSPDXIdentifier:   "SPDXRef-pkg-" + bom.PackageInfo.Name,
		PackageDownloadLocation: "NOASSERTION",

		// this is not correct. we do analyze files but if this is set
		// to true we need to implement the verification algorithm like
		// here: https://spdx.org/spdx-specification-21-web-version#h.2p2csry
		FilesAnalyzed: false,
		// this is the license we detect
		PackageLicenseConcluded: "NOASSERTION",
		PackageLicenseDeclared:  bom.PackageInfo.LicenseDeclared,
		PackageCopyrightText:    "NOASSERTION",
		Files:                   files,
	}

	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:       "SPDX-2.1",
			DataLicense:       "CC0-1.0",
			SPDXIdentifier:    "SPDXRef-DOCUMENT",
			DocumentName:      bom.PackageInfo.Name,
			DocumentNamespace: "https://qmstr.org",
			Created:           time.Now().Format(time.RFC3339),
			CreatorTools: []string{
				"QMSTR",
			},
		},
		Packages: []*spdx.Package2_1{pkg},
	}

	fName := filepath.Join(r.outputdir, "report.spdx")
	out, err := os.Create(fName)
	if err != nil {
		return fmt.Errorf("failed to create out file %q: %v", fName, err)
	}
	defer out.Close()

	tvsaver.Save2_1(doc, out)

	return nil
}

func (r *SPDXReporter) PostReport() error {
	return nil
}
