package analysis

import (
	"encoding/csv"
	"log"
	"os"
	"os/exec"

	"github.com/QMSTR/qmstr/pkg/database"
)

type NinkaAnalyzer struct {
	cmd     string
	cmdargs []string
	Config  map[string]string
	db      *database.DataBase
}

func NewNinkaAnalyzer(config map[string]string, db *database.DataBase) *NinkaAnalyzer {
	return &NinkaAnalyzer{"ninka", []string{"-i"}, config, db}
}

func (ninka *NinkaAnalyzer) Analyze(node *AnalysisNode) error {
	licenses := ninka.detectLicenses(node.GetPath())
	// TODO merge with currently set licenses
	for _, license := range licenses {
		node.SetLicense(license)
	}
	return nil
}

func (ninka *NinkaAnalyzer) detectLicenses(srcFilePath string) []*database.License {
	var licenses []*database.License
	cmd := exec.Command(ninka.cmd, append(ninka.cmdargs, srcFilePath)...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("License analysis failed for %s", srcFilePath)
	}

	licenseFile, err := os.Open(srcFilePath + ".license")
	if err != nil {
		log.Println(err)
	}
	r := csv.NewReader(licenseFile)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		log.Println(err)
	}

	for _, fields := range records {
		if len(fields) > 0 {
			licenses = append(licenses, &database.License{Key: fields[0]})
		}
	}
	if len(licenses) == 0 {
		licenses = []*database.License{database.UnknownLicense}
	}
	return licenses
}
