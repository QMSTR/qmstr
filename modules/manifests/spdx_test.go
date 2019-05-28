package manifests

import (
	"bytes"
	"testing"
)

func TestSPDX(t *testing.T) {
	mani, err := NewSPDXManifest(bytes.NewBufferString(curlSPDX))
	if err != nil {
		t.Errorf("Create new SPDX manifest fail: %v", err)
	}
	vals := map[string]string{
		"./usr/bin/curl":                           "b0b18f019df78d96d90160b8e226b54745a6a347",
		"./usr/share/doc/curl/NEWS.Debian.gz":      "79f0b8dec438bdd43920e20faed50c0f83d621b8",
		"./usr/share/doc/curl/changelog.Debian.gz": "66889b6a79756642692d0b22fdeb0cd13df6f08d",
		"./usr/share/doc/curl/changelog.gz":        "c1c4e88a537f9abc50facdf0aaa6772addfbfce9",
		"./usr/share/doc/curl/copyright":           "2bd36d29a0f4bff70886598f8e968d30d9aee53c",
		"./usr/share/man/man1/curl.1.gz":           "10878e12a84f177ea4f292b16b0edfb239478b42",
		"./usr/share/zsh/vendor-completions/_curl": "e1e0f8c3b5e4c081cda3e788488808abda187e24",
	}
	for _, f := range mani.FileInfo() {
		if vals[f.Name] != f.SHA1 {
			t.Errorf("hash does not match for file %q (%s<>%s)", f.Name, vals[f.Name], f.SHA1)
		}
	}
}

var curlSPDX = `DocumentName: curl-SPDX-manifest
SPDXID: SPDXRef-DOCUMENT
DataLicense: CC0-1.0
SPDXVersion: SPDX-2.1
Creator: Tool: QMSTR

PackageName: curl
SPDXID: SPDXRef-pkg-curl
PackageVersion: 7.52.1
PackageFileName: curl_7.52.1-5+deb9u9_amd64.deb
PackageDownloadLocation: http://ftp.us.debian.org/debian/pool/main/c/curl/curl_7.52.1-5+deb9u9_amd64.deb
PackageChecksum: MD5: e4a99108a544c5e33a8dc38cf64c3ad7
PackageChecksum: SHA256: d494ca641912ffdc6acf0c46664898ec1637cd252328a3f75f0b01cf7470a65d
PackageLicenseInfoFromFiles: NONE
PackageLicenseConcluded: curl
PackageLicenseDeclared: curl
PackageCopyrightText: <text>Daniel Stenberg <daniel@haxx.se></text>
PackageVerificationCode: faadc7f2d8c9dc53eb1eca2dfba34d29b9796a52

FileName: ./usr/share/man/man1/curl.1.gz
SPDXID: SPDXRef-File
FileChecksum: SHA1: 10878e12a84f177ea4f292b16b0edfb239478b42
LicenseConcluded: LGPL-2.0
FileCopyrightText: <text>Daniel Stenberg <daniel@haxx.se></text>
LicenseInfoInFile: NONE

FileName: ./usr/share/zsh/vendor-completions/_curl
SPDXID: SPDXRef-File
FileChecksum: SHA1: e1e0f8c3b5e4c081cda3e788488808abda187e24
LicenseConcluded: curl
FileCopyrightText:  NONE
LicenseInfoInFile: NONE

FileName: ./usr/share/doc/curl/changelog.gz
SPDXID: SPDXRef-File
FileChecksum: SHA1: c1c4e88a537f9abc50facdf0aaa6772addfbfce9
LicenseConcluded: curl
FileCopyrightText:  NONE
LicenseInfoInFile: NONE

FileName: ./usr/share/doc/curl/copyright
SPDXID: SPDXRef-File
FileChecksum: SHA1: 2bd36d29a0f4bff70886598f8e968d30d9aee53c
LicenseConcluded: curl
FileCopyrightText:  NONE
LicenseInfoInFile: NONE

FileName: ./usr/share/doc/curl/NEWS.Debian.gz
SPDXID: SPDXRef-File
FileChecksum: SHA1: 79f0b8dec438bdd43920e20faed50c0f83d621b8
LicenseConcluded: curl
FileCopyrightText:  NONE
LicenseInfoInFile: NONE

FileName: ./usr/share/doc/curl/changelog.Debian.gz
SPDXID: SPDXRef-File
FileChecksum: SHA1: 66889b6a79756642692d0b22fdeb0cd13df6f08d
LicenseConcluded: curl
FileCopyrightText:  NONE
LicenseInfoInFile: NONE

FileName: ./usr/bin/curl
SPDXID: SPDXRef-File
FileChecksum: SHA1: b0b18f019df78d96d90160b8e226b54745a6a347
LicenseConcluded: curl
FileCopyrightText:  NONE
LicenseInfoInFile: NONE`
