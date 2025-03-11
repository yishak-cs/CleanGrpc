package main

import (
	"fmt"
	"log"
	"net"

	"github.com/yishak-cs/CleanGrpc/Internal/db"
	interfaces "github.com/yishak-cs/CleanGrpc/pkg/v1"
	repository "github.com/yishak-cs/CleanGrpc/pkg/v1/Repository"
	usecase "github.com/yishak-cs/CleanGrpc/pkg/v1/UseCase"
	handler "github.com/yishak-cs/CleanGrpc/pkg/v1/handler/grpc"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func main() {

	// connect to a database
	db := db.DBconn()

	//grpc server listen tcp connection on address string
	listener, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		fmt.Println("unable to get Listener")
	}
	server := grpc.NewServer()

	// get a type that implements UseCaseInterface
	uc := initUserServer(db)

	//register the UserService handler on the server
	handler.NewUserServer(server, uc)

	// start serving to the address
	log.Fatal(server.Serve(listener))
}

func initUserServer(db *gorm.DB) interfaces.UseCaseInterface {
	//create a type that implements RepoInterface
	repo := repository.NewRepo(db)
	//return the UseCaseInterface instance to the called
	return usecase.NewUseCase(repo)
}
