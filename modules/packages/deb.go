package packages

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/QMSTR/qmstr/lib/go-qmstr/validation"
	"github.com/blakesmith/ar"
	"github.com/ulikunitz/xz"
)

type DebPackage struct {
	pi validation.PackageInfo
	fi []validation.FileInfo
}

func (p *DebPackage) Validate(mani validation.Manifest) error {
	pIdx := indexByName(p.fi)
	mIdx := indexByName(mani.FileInfo())

	for name, pfi := range pIdx {
		if mfi, ok := mIdx[name]; ok {
			if err := validateFileInfos(mfi, pfi); err != nil {
				return err
			}
			delete(mIdx, name)
			continue
		}
		return validation.FileNotInManifestError{Name: name}
	}
	for name := range mIdx {
		return validation.FileNotInPackageError{Name: name}
	}
	return nil
}

func NewDebPackage(in io.Reader) (*DebPackage, error) {
	r := ar.NewReader(in)
	dataFound := false
	hdr := &ar.Header{}
	var err error

	for {
		hdr, err = r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(hdr.Name, "data.") {
			dataFound = true
			break
		}
	}
	if !dataFound {
		return nil, errors.New("deb pkg: no data archive found")
	}

	parts := strings.Split(hdr.Name, ".")
	cmp := ""

	if len(parts) == 3 {
		cmp = parts[2]
	} else {
		if len(parts) != 2 {
			return nil, fmt.Errorf("deb pkg: wrong data archive name: %s", hdr.Name)
		}
	}

	rc, err := decompress(cmp, r)
	if err != nil {
		return nil, fmt.Errorf("deb pkg: %v", err)
	}

	tarr := tar.NewReader(rc)
	p := &DebPackage{}

	for {
		thdr, err := tarr.Next()
		// done
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("deb pkg: data archive: %v", err)
		}

		if !thdr.FileInfo().IsDir() {
			h, err := common.Hash(tarr)
			if err != nil {
				return nil, fmt.Errorf("deb pkg: hashing file %s: %v", thdr.Name, err)
			}
			p.fi = append(p.fi, validation.FileInfo{
				Name: thdr.Name,
				SHA1: h,
			})
		}
	}
	return p, nil
}

func decompress(cmp string, r io.Reader) (io.ReadCloser, error) {
	switch cmp {
	case "":
		return ioutil.NopCloser(r), nil
	case "gz", "gzip":
		return gzip.NewReader(r)
	case "bz", "bz2", "bzip2":
		return ioutil.NopCloser(bzip2.NewReader(r)), nil
	case "xz":
		rdr, err := xz.NewReader(r)
		return ioutil.NopCloser(rdr), err
	default:
		return nil, fmt.Errorf("unsuported compression: %s", cmp)
	}
}
