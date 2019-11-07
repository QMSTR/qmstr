package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/QMSTR/qmstr/lib/go-qmstr/module"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/spdx/tools-golang/v0/spdx"
	"github.com/spdx/tools-golang/v0/tvsaver"
)

const (
	ModuleName  = "reporter-package-manifest"
	outFileName = "%s.spdx"
)

type PkgManifestReporter struct {
	enableWarnings bool
	enableErrors   bool
	outputdir      string
	nsURI          string
	pathRegexp     string
	pathReplace    string
}

func (r *PkgManifestReporter) Configure(config map[string]string) error {
	if outDir, ok := config["outputdir"]; ok {
		r.outputdir = outDir
	} else {
		return fmt.Errorf("no output directory configured")
	}
	if str, ok := config["pathSubst"]; ok {
		arr := strings.Split(str, "||")
		if len(arr) != 2 {
			return fmt.Errorf("invalid pathSubst: %s", str)
		}
		r.pathRegexp = arr[0]
		r.pathReplace = arr[1]
	}

	if nsURI, ok := config["namespaceURI"]; ok {
		r.nsURI = nsURI
	} else {
		// see: https://spdx.org/spdx-specification-21-web-version#h.1gdfkutofa90
		r.nsURI = "http://spdx.org/spdxdocs/"
	}

	return nil
}

func (r *PkgManifestReporter) Report(masterClient *module.MasterClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pnps, err := masterClient.GetPackageNodes()
	if err != nil {
		return fmt.Errorf("Couldn't get package node, %v", err)
	}

	for _, pnp := range pnps {
		if err = r.generateSPDX(pnp, masterClient.RptSvcClient, ctx); err != nil {
			return fmt.Errorf("Error while generating SPDX, %v", err)
		}
	}
	return nil
}

func (r *PkgManifestReporter) generateSPDX(pkgNode *module.PackageNodeProxy, rserv service.ReportServiceClient, ctx context.Context) error {
	files := []*spdx.File2_1{}
	hashes := []string{}
	targets, err := pkgNode.GetTargets()
	if err != nil {
		return err
	}
	for _, trgt := range targets {
		licenses, err := rserv.GetInfoData(ctx, &service.InfoDataRequest{RootID: trgt.Uid, Infotype: "license", Datatype: "name"})
		if err != nil {
			return fmt.Errorf("Couldn't get license node, %v", err)
		}
		copyrights, err := rserv.GetInfoData(ctx, &service.InfoDataRequest{RootID: trgt.Uid, Infotype: "copyright", Datatype: "author"})
		if err != nil {
			return fmt.Errorf("Couldn't get copyright node, %v", err)
		}
		if r.pathRegexp != "" {
			re, err := regexp.Compile(r.pathRegexp)
			if err != nil {
				return fmt.Errorf("failed to compile regexp: %v, %s", err, r.pathRegexp)
			}
			trgt.Path = re.ReplaceAllString(trgt.Path, r.pathReplace)
		}

		if err != nil {
			return fmt.Errorf("Couldn't get copyright node, %v", err)
		}
		fl := &spdx.File2_1{
			FileName: trgt.Path,
			// this should be unique
			FileSPDXIdentifier: "SPDXRef-file-" + trgt.Name,
			FileChecksumSHA1:   trgt.FileData.GetHash(),
			LicenseConcluded:   "NOASSERTION",
			LicenseInfoInFile:  []string{"NOASSERTION"},
			FileCopyrightText:  "NOASSERTION",
		}
		if len(licenses.Data) != 0 {
			fl.LicenseInfoInFile = licenses.Data
			lic := strings.Join(licenses.Data, " AND ")
			fl.LicenseConcluded = lic
		}
		if len(copyrights.Data) != 0 {
			cprights := strings.Join(copyrights.Data, " AND ")
			fl.FileCopyrightText = cprights
		}
		files = append(files, fl)
		hashes = append(hashes, trgt.FileData.GetHash())
	}

	dnldLocation := pkgNode.GetMetaData("SourceURL", "NOASSERTION")
	pkg := &spdx.Package2_1{
		PackageName: pkgNode.Name,
		// this should be unique
		PackageSPDXIdentifier:   "SPDXRef-pkg-" + pkgNode.Name,
		PackageDownloadLocation: dnldLocation,
		FilesAnalyzed:           true,
		PackageVerificationCode: calcSHA1Hash(hashes),
		// this is the license we detect
		PackageLicenseConcluded:     "NOASSERTION",
		PackageLicenseInfoFromFiles: []string{"NOASSERTION"},
		PackageLicenseDeclared:      pkgNode.GetMetaData("license_declared", "NOASSERTION"),
		PackageCopyrightText:        pkgNode.GetMetaData("cr_text", "NOASSERTION"),
		Files:                       files,
	}

	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:       "SPDX-2.1",
			DataLicense:       "CC0-1.0",
			SPDXIdentifier:    "SPDXRef-DOCUMENT",
			DocumentName:      pkgNode.Name,
			DocumentNamespace: r.nsURI + url.PathEscape(pkgNode.Name),
			Created:           time.Now().Format(time.RFC3339),
			CreatorTools: []string{
				"QMSTR",
			},
		},
		Packages: []*spdx.Package2_1{pkg},
	}

	fName := filepath.Join(r.outputdir, fmt.Sprintf(outFileName, common.GetPosixFullyPortableFilename(pkgNode.Name)))
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

func (r *PkgManifestReporter) PostReport() error {
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
