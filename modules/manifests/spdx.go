package manifests

import (
	"fmt"
	"io"

	"github.com/QMSTR/qmstr/lib/go-qmstr/validation"
	"github.com/spdx/tools-golang/v0/tvloader"
)

type SPDXManifest struct {
	pi validation.PackageInfo
	fi []validation.FileInfo
}

func (mani *SPDXManifest) PackageInfo() validation.PackageInfo {
	return mani.pi
}

func (mani *SPDXManifest) FileInfo() []validation.FileInfo {
	return mani.fi
}

func NewSPDXManifest(r io.Reader) (validation.Manifest, error) {
	doc, err := tvloader.Load2_1(r)
	if err != nil {
		return nil, fmt.Errorf("SPDX mfst: loading tag-value: %v", err)
	}

	mani := &SPDXManifest{
		pi: validation.PackageInfo{
			Name: doc.Packages[0].PackageName,
		},
	}

	for _, f := range doc.Packages[0].Files {
		if f.LicenseConcluded == "NOASSERTION" || f.LicenseConcluded == "" {
			return nil, validation.MissingLicenseInfoError{
				Name: f.FileName,
			}
		}
		if f.FileCopyrightText == "NOASSERTION" || f.FileCopyrightText == "" {
			return nil, validation.MissingCopyrightInfoError{
				Name: f.FileName,
			}
		}
		mani.fi = append(mani.fi, validation.FileInfo{
			Name: f.FileName,
			SHA1: f.FileChecksumSHA1,
		})
	}

	return mani, nil
}
