package master

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseReport struct {
	genericServerPhase
}

func newReportPhase(session string, masterConfig *config.MasterConfig, db *database.DataBase, server *server) serverPhase {
	return &serverPhaseReport{
		genericServerPhase{Name: "Report", session: session, masterConfig: masterConfig, db: db, server: server},
	}
}

func (phase *serverPhaseReport) Activate() error {
	log.Println("Reporting activated")
	for idx, reporterConfig := range phase.masterConfig.Reporting {
		reporterName := reporterConfig.Reporter

		cmd := exec.Command(reporterName, "--rserv", phase.masterConfig.Server.RPCAddress, "--rid", fmt.Sprintf("%d", idx))
		out, err := cmd.CombinedOutput()
		if err != nil {
			logModuleError(reporterName, out)
			return err
		}
		log.Printf("Reporter %s finished successfully: %s\n", reporterName, out)
	}
	return nil
}

func (phase *serverPhaseReport) Shutdown() error {
	return nil
}

func (phase *serverPhaseReport) GetPhaseID() int32 {
	return PhaseIDReport
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

	// Set cachedir, if not overriden
	if _, ok := config.Config["cachedir"]; !ok {
		config.Config["cachedir"] = filepath.Join(phase.masterConfig.Server.CacheDir, config.Reporter, config.PosixName)
	}
	// Set output dir, if not overriden
	if _, ok := config.Config["outputdir"]; !ok {
		config.Config["outputdir"] = filepath.Join(phase.masterConfig.Server.OutputDir, config.Reporter, config.PosixName)
	}

	return &service.ReporterConfigResponse{ConfigMap: config.Config, Session: phase.session,
		Name: config.Name}, nil
}

func (phase *serverPhaseReport) GetInfoData(in *service.InfoDataRequest) (*service.InfoDataResponse, error) {
	db, err := phase.getDataBase()
	if err != nil {
		return nil, err
	}
	infos, err := db.GetInfoData(in.RootID, in.Infotype, in.Datatype)
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
