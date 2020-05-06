package master

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/QMSTR/qmstr/lib/go-qmstr/config"
	"github.com/QMSTR/qmstr/lib/go-qmstr/database"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

type serverPhaseReport struct {
	genericServerPhase
}

func newReportPhase(masterConfig *config.MasterConfig, db *database.DataBase, server *server, done bool) serverPhase {
	return &serverPhaseReport{
		genericServerPhase{Name: "Report", masterConfig: masterConfig, db: db, server: server, done: done},
	}
}

func (phase *serverPhaseReport) Activate() error {
	log.Println("Reporting activated")
	phase.server.publishEvent(&service.Event{Class: service.EventClass_PHASE, Message: "Activating reporting phase"})
	return nil
}

func (phase *serverPhaseReport) InitModule(in *service.InitModuleRequest) (*service.InitModuleResponse, error) {
	log.Printf("Running reporter %s ...\n", in.ModuleName)
	phase.server.publishEvent(&service.Event{Class: service.EventClass_MODULE, Message: fmt.Sprintf("Running reporter %s", in.ModuleName)})
	return &service.InitModuleResponse{}, nil
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

func (phase *serverPhaseReport) GetInfoData(in *service.InfoDataRequest) (*service.InfoDataResponse, error) {
	db, err := phase.getDataBase()
	if err != nil {
		return nil, err
	}
	var infos []string

	if in.Infotype == "" {
		infos, err = db.GetAllInfoData(in.Infotype)
		if err != nil {
			return nil, err
		}
	} else {
		infos, err = db.GetInfoDataByTrustLevel(in.RootID, in.Infotype, in.Datatype)
	}
	// remove duplicate data
	datamap := map[string]struct{}{}
	for _, info := range infos {
		info = strings.TrimSpace(info)
		if info != "" {
			datamap[info] = struct{}{}
		}
	}
	var data []string
	for value := range datamap {
		data = append(data, value)
	}
	if err != nil {
		return nil, err
	}
	return &service.InfoDataResponse{Data: data}, nil
}
