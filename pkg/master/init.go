package master

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseInit struct {
	genericServerPhase
}

func newInitServerPhase(masterConfig *config.MasterConfig) *serverPhaseInit {
	return &serverPhaseInit{genericServerPhase{Name: "Init", masterConfig: masterConfig}}
}

func (phase *serverPhaseInit) Activate() error {
	// Connect to database (dgraph)
	db, err := database.Setup(phase.masterConfig.Server.DBAddress, phase.masterConfig.Server.DBWorkers)
	if err != nil {
		return fmt.Errorf("Could not setup database: %v", err)
	}
	phase.db = db

	if !snapshotAvailable() {
		phase.db.OpenInsertQueue()
		if phase.masterConfig.Name != "" {
			phase.initProject()
		}
		phase.db.CloseInsertQueue()
		return nil
	}

	if err := importSnapshot(); err != nil {
		return fmt.Errorf("Failed to import snapshot: %v", err)
	}

	qmstrState, err := phase.db.GetQmstrStateNode()
	if err != nil {
		return fmt.Errorf("Failed to reconstruct qmstr state after snapshot import: %v", err)
	}
	phase.postInitPhase = &qmstrState.Phase
	return nil
}

func snapshotAvailable() bool {
	_, err := os.Stat(common.ContainerGraphImportPath)
	// TODO fix exist check
	return !os.IsNotExist(err)
}

func importSnapshot() error {
	snapshotFile, err := os.Open(common.ContainerGraphImportPath)
	if err != nil {
		return fmt.Errorf("Failed opening snapshot: %v", snapshotFile)
	}

	rdfFile, err := ioutil.TempFile("", "rdf")
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(rdfFile.Name() + ".gz")
	}()

	var rdf, schema bool

	tr := tar.NewReader(snapshotFile)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return fmt.Errorf("Failed reading snapshot: %v", err)
		}
		if hdr.Typeflag == tar.TypeReg && strings.Contains(hdr.Name, "schema.gz") {

			// check scheme
			r, err := gzip.NewReader(tr)
			if err != nil {
				return err
			}
			checkSchema, err := ioutil.ReadAll(r)
			schema = database.CheckSchema(string(checkSchema))
			r.Close()
			continue
		}
		if hdr.Typeflag == tar.TypeReg && strings.Contains(hdr.Name, "rdf.gz") {
			_, err := io.Copy(rdfFile, tr)
			if err != nil {
				return err
			}
			rdfFile.Close()
			os.Rename(rdfFile.Name(), rdfFile.Name()+".gz")
			rdf = true
			continue
		}
		if strings.Contains(hdr.Name, common.ContainerPushFilesDirName) {
			targetDir := filepath.Join(common.ContainerBuildDir, common.ContainerPushFilesDirName)
			switch hdr.Typeflag {
			case tar.TypeDir:
				if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
					return err
				}
			case tar.TypeReg:
				dest, err := os.Create(filepath.Join(targetDir, filepath.Base(hdr.Name)))
				if err != nil {
					return err
				}
				_, err = io.Copy(dest, tr)
				if err != nil {
					return err
				}
				dest.Close()
			}
		}
	}

	if !schema || !rdf {
		return errors.New("invalid snapshot supplied")
	}

	importCmd := exec.Command("dgraph", "live", "-r", rdfFile.Name()+".gz")
	out, err := importCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("live replay failed: %v\n%s", err, out)
	}
	return nil
}

func (phase *serverPhaseInit) initProject() {
	projectNode := &service.ProjectNode{Name: phase.masterConfig.Name}
	tmpInfoNode := &service.InfoNode{Type: "metadata"}
	for key, val := range phase.masterConfig.MetaData {
		tmpInfoNode.DataNodes = append(tmpInfoNode.DataNodes, &service.InfoNode_DataNode{Type: key, Data: val})
	}

	if len(tmpInfoNode.DataNodes) > 0 {
		projectNode.AdditionalInfo = []*service.InfoNode{tmpInfoNode}
	}

	phase.db.AddProjectNode(projectNode)
}

func (phase *serverPhaseInit) Shutdown() error {
	return nil
}

func (phase *serverPhaseInit) GetPhaseID() service.Phase {
	return service.Phase_INIT
}
