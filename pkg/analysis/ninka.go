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
	licenses, err := ninka.detectLicenses(node.GetPath())
	if err != nil {
		node.SetLicense("UNKNOWN")
		return err
	}

	// TODO merge with currently set licenses
	for _, license := range licenses {
		node.SetLicense(license)
	}
	return nil
}

func (ninka *NinkaAnalyzer) detectLicenses(srcFilePath string) ([]string, error) {
	licenses := []string{}
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
		return []string{}, err
	}
	r := csv.NewReader(licenseFile)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return []string{}, err
	}

	for _, fields := range records {
		if len(fields) > 0 {
			licenses = append(licenses, fields[0])
		}
	}
	return licenses, nil
}
