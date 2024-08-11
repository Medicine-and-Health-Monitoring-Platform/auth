package main

import (
	"Auth/api"
	"Auth/config"
	pbu "Auth/genproto/users"
	"Auth/service"
	"Auth/storage"
	"Auth/storage/postgres"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatalf("error while connecting to database: %v", err)
	}
	defer db.Close()
	cfg := config.Load()

	var wg sync.WaitGroup
	wg.Add(2)

	go RunService(&wg, db, cfg)
	go RunRouter(&wg, db, cfg)

	wg.Wait()
}

func RunService(wg *sync.WaitGroup, s storage.IStorage, cfg *config.Config) {
	defer wg.Done()

	lis, err := net.Listen("tcp", cfg.AUTH_SERVICE_PORT)
	if err != nil {
		log.Fatalf("error while listening: %v", err)
	}
	defer lis.Close()

	u := service.NewUserService(s)
	a := service.NewAdminService(s)
	server := grpc.NewServer()
	pbu.RegisterAuthServiceServer(server, u)
	pbu.RegisterAdminServer(server, a)

	log.Printf("Service is listening on port %s...\n", cfg.AUTH_SERVICE_PORT)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("error while serving auth service: %s", err)
	}
}

func RunRouter(wg *sync.WaitGroup, s storage.IStorage, cfg *config.Config) {
	defer wg.Done()

	r := api.NewRouter(s)
	log.Printf("Router is running on port %s...\n", cfg.AUTH_ROUTER_PORT)
	r.Run(config.Load().AUTH_ROUTER_PORT)
}
