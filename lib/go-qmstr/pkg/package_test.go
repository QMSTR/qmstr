package pkg

import (
	"testing"
)

type testManifest struct {
	pi PackageInfo
	fi []FileInfo
}

func (mani testManifest) PackageInfo() PackageInfo {
	return mani.pi
}

func (mani testManifest) FileInfo() []FileInfo {
	return mani.fi
}

type testPackage struct {
	pi PackageInfo
	fi []FileInfo
}

func (p testPackage) Validate(mani Manifest) error {

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
		return FileNotInManifestError{Name: name}
	}
	for name := range mIdx {
		return FileNotInPackageError{Name: name}
	}
	return nil
}

func validateFileInfos(a, b FileInfo) error {
	if a.SHA1 != b.SHA1 {
		return HashMissmatchError{
			Algo: "SHA1",
			A:    a.SHA1,
			B:    b.SHA1,
		}
	}
	return nil
}

func indexByName(fis []FileInfo) map[string]FileInfo {
	m := map[string]FileInfo{}
	for _, fi := range fis {
		m[fi.Name] = fi
	}
	return m
}

func TestValidation(t *testing.T) {
	mani := testManifest{
		pi: PackageInfo{
			Name: "TestPackage",
		},
		fi: []FileInfo{
			FileInfo{
				Name: "path/to/a",
				SHA1: "aabb",
			},
			FileInfo{
				Name: "/path/to/b",
				SHA1: "bbcc",
			},
		},
	}

	pkgValid := testPackage{
		pi: PackageInfo{
			Name: "TestPackage",
		},
		fi: []FileInfo{
			FileInfo{
				Name: "path/to/a",
				SHA1: "aabb",
			},
			FileInfo{
				Name: "/path/to/b",
				SHA1: "bbcc",
			},
		},
	}

	pkgWrongSHA1 := testPackage{
		pi: PackageInfo{
			Name: "TestPackage",
		},
		fi: []FileInfo{
			FileInfo{
				Name: "path/to/a",
				SHA1: "aabb",
			},
			FileInfo{
				Name: "/path/to/b",
				SHA1: "bbccdd",
			},
		},
	}

	pkgFileNotInPackage := testPackage{
		pi: PackageInfo{
			Name: "TestPackage",
		},
		fi: []FileInfo{
			FileInfo{
				Name: "path/to/a",
				SHA1: "aabb",
			},
		},
	}

	pkgFileNotInManifest := testPackage{
		pi: PackageInfo{
			Name: "TestPackage",
		},
		fi: []FileInfo{
			FileInfo{
				Name: "path/to/a",
				SHA1: "aabb",
			},
			FileInfo{
				Name: "/path/to/b",
				SHA1: "bbcc",
			},
			FileInfo{
				Name: "/path/to/c",
				SHA1: "ddfr",
			},
		},
	}

	testCases := []struct {
		desc string
		pkg  Package
		err  string
	}{
		{
			desc: "Valid Package",
			pkg:  pkgValid,
			err:  "",
		},
		{
			desc: "SHA1 Missmatch",
			pkg:  pkgWrongSHA1,
			err:  "hash SHA1 missmatch: bbcc <> bbccdd",
		},
		{
			desc: "File not in Package",
			pkg:  pkgFileNotInPackage,
			err:  "file /path/to/b not present in package",
		},
		{
			desc: "File not in Manifest",
			pkg:  pkgFileNotInManifest,
			err:  "file /path/to/c not documented in manifest",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.pkg.Validate(mani)
			if err != nil {
				if err.Error() != tC.err {
					t.Fatal(err)
				}
				return
			}
			if tC.err != "" {
				t.Fatal("error should be nil")
			}
		})
	}
}
