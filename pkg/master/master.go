package master

import (
	"fmt"
	"log"
	"net"
	"path/filepath"

	"golang.org/x/net/context"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/QMSTR/qmstr/pkg/database"
	"google.golang.org/grpc"
)

var quitServer chan interface{}

type server struct {
	db *database.DataBase
}

func (s *server) Build(ctx context.Context, in *pb.BuildMessage) (*pb.BuildResponse, error) {
	// Compiliation
	for _, compile := range in.GetCompilations() {
		log.Printf("Compiled %v", compile)

		uidTrgt, err := s.db.HasNode(compile.Target.GetHash())
		if err != nil {
			return &pb.BuildResponse{Success: false}, err
		}

		// no such node exist
		if uidTrgt == "" {
			uidSrc, err := s.db.HasNode(compile.Source.GetHash())
			if err != nil {
				return &pb.BuildResponse{Success: false}, err
			}

			var src database.Node

			// no such node exist
			if uidSrc == "" {
				src = database.NewNode(compile.Source.GetPath(), compile.Source.GetHash())
				src.Type = database.ArtifactTypeSrc
			} else {
				src = database.Node{
					Uid: uidSrc,
				}
			}
			trgt := database.NewNode(compile.Target.GetPath(), compile.Target.GetHash())
			trgt.DerivedFrom = []database.Node{src}
			trgt.Type = database.ArtifactTypeObj

			uidTrgt, err := s.db.AddNode(&trgt)
			if err != nil {
				return &pb.BuildResponse{Success: false}, err
			}
			log.Printf("Target node with UID: %s added\n", uidTrgt)
		}
	}

	// Linking
	for _, bin := range in.GetBinary() {
		log.Printf("Linked target: %v", bin)

		uidTrgt, err := s.db.HasNode(bin.Target.GetHash())
		if err != nil {
			return &pb.BuildResponse{Success: false}, err
		}

		deps := []database.Node{}
		// no such node exist
		if uidTrgt == "" {
			for _, dep := range bin.GetInput() {
				uidDep, err := s.db.HasNode(dep.GetHash())
				if err != nil {
					return &pb.BuildResponse{Success: false}, err
				}

				depNode := database.Node{}

				// dep not in db
				if uidDep == "" {
					depNode = database.NewNode(dep.GetPath(), dep.GetHash())
					depNode.Name = filepath.Base(dep.GetPath())
				} else {
					depNode = database.Node{
						Uid: uidDep,
					}
				}
				deps = append(deps, depNode)
			}
			trgt := database.NewNode(bin.Target.GetPath(), bin.Target.GetHash())
			trgt.DerivedFrom = deps
			trgt.Type = database.ArtifactTypeLink

			uidTrgt, err := s.db.AddNode(&trgt)
			if err != nil {
				return &pb.BuildResponse{Success: false}, err
			}
			log.Printf("Target node with UID: %s added\n", uidTrgt)
		}
	}

	return &pb.BuildResponse{Success: true}, nil
}

func (s *server) Log(ctx context.Context, in *pb.LogMessage) (*pb.LogResponse, error) {
	log.Printf("REMOTE: %s", string(in.Msg))
	return &pb.LogResponse{Success: true}, nil
}

func (s *server) Quit(ctx context.Context, in *pb.QuitMessage) (*pb.QuitResponse, error) {
	if in.Kill {
		log.Fatalf("qmstr was killed hard by client")
	}

	// Wait for pending tasks to complete e.g. synchronize channels

	// Schedule shutdown
	quitServer <- nil

	return &pb.QuitResponse{Success: true}, nil
}

func ListenAndServe(rpcAddr string, dbAddr string) error {

	// Connect to backend database (dgraph)
	db, err := database.Setup(dbAddr)
	if err != nil {
		return fmt.Errorf("Could not setup database: %v", err)
	}

	// Setup buildservice
	lis, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return fmt.Errorf("Failed to setup socket and listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBuildServiceServer(s, &server{
		db: db,
	})

	quitServer = make(chan interface{})
	go func() {
		<-quitServer
		log.Println("qmstr terminated by client")
		s.GracefulStop()
		close(quitServer)
		quitServer = nil
	}()

	log.Print("qmstr master running")
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("Failed to start rpc service %v", err)
	}
	return nil
}
