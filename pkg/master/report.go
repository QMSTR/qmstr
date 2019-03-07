package master

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/reporting"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseReport struct {
	genericServerPhase
}

func newReportPhase(masterConfig *config.MasterConfig, db *database.DataBase, server *server) serverPhase {
	return &serverPhaseReport{
		genericServerPhase{Name: "Report", masterConfig: masterConfig, db: db, server: server},
	}
}

func (phase *serverPhaseReport) Activate() error {
	log.Println("Reporting activated")
	for idx, reporterConfig := range phase.masterConfig.Reporting {
		reporterName := reporterConfig.Reporter

		log.Printf("Running reporter %s ...\n", reporterName)
		phase.server.publishEvent(&service.Event{Class: service.EventClass_MODULE, Message: fmt.Sprintf("Running reporter %s", reporterName)})
		cmd := exec.Command(reporterName, "--rserv", phase.masterConfig.Server.RPCAddress, "--rid", fmt.Sprintf("%d", idx))
		out, err := cmd.CombinedOutput()
		if err != nil {
			logModuleError(reporterName, out)
			errMsg := fmt.Sprintf("Reporter %s failed: %v", reporterName, err)
			phase.server.publishEvent(&service.Event{Class: service.EventClass_MODULE, Message: errMsg})
			return errors.New(errMsg)
		}
		phase.server.publishEvent(&service.Event{Class: service.EventClass_MODULE, Message: fmt.Sprintf("Reporter %s successfully finished", reporterName)})
		log.Printf("Reporter %s finished successfully:\n%s\n", reporterName, out)
	}
	return nil
}

func (phase *serverPhaseReport) Shutdown() error {
	return nil
}

func (phase *serverPhaseReport) GetPhaseID() service.Phase {
	return service.Phase_REPORT
}

func (phase *serverPhaseReport) GetReporterConfig(in *service.ReporterConfigRequest) (*service.ReporterConfigResponse, error) {
	idx := in.ReporterID
	if idx < 0 || idx >= int32(len(phase.masterConfig.Reporting)) {
		return nil, fmt.Errorf("Invalid reporter id %d", idx)
	}
	config := phase.masterConfig.Reporting[idx]

	if config.Config == nil {
		config.Config = make(map[string]string)
	}

	config.Config["name"] = config.Name
	config.Config["cachedir"] = filepath.Join(ServerCacheDir, config.Reporter, config.PosixName)
	config.Config["outputdir"] = filepath.Join(ServerOutputDir, config.Reporter, config.PosixName)

	return &service.ReporterConfigResponse{ConfigMap: config.Config, Name: config.Name}, nil
}

func (phase *serverPhaseReport) GetBOM(in *service.BOMRequest) (*service.BOM, error) {
	db, err := phase.getDataBase()
	if err != nil {
		return nil, err
	}
	pkgNode, err := db.GetPackageNode()
	if err != nil {
		return nil, err
	}
	projNode, err := db.GetProjectNode()
	if err != nil {
		return nil, err
	}
	bom, err := reporting.GetBOM(projNode, pkgNode, in.Warnings, in.Errors)
	if err != nil {
		return nil, err
	}
	return bom, nil
}

func (phase *serverPhaseReport) GetInfoData(in *service.InfoDataRequest) (*service.InfoDataResponse, error) {
	db, err := phase.getDataBase()
	if err != nil {
		return nil, err
	}
	var infos []string

	if in.Datatype == "" {
		infos, err = db.GetAllInfoData(in.Infotype)
		if err != nil {
			return nil, err
		}
		return &service.InfoDataResponse{Data: infos}, nil
	} else {
		infos, err = db.GetInfoData(in.RootID, in.Infotype, in.Datatype)
	}
	if err != nil {
		return nil, err
	}

	infoSet := map[string]struct{}{}
	for _, data := range infos {
		infoSet[data] = struct{}{}
	}

	uniqInfos := []string{}
	for key := range infoSet {
		uniqInfos = append(uniqInfos, key)
	}

	return &service.InfoDataResponse{Data: uniqInfos}, nil
}
