package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"git.neds.sh/matty/entain/api/proto/racing"
)

var (
	apiEndpoint  = flag.String("api-endpoint", "localhost:8000", "API endpoint")
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9000", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Printf("failed running api server: %s\n", err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Create a new gRPC client connection
	conn, err := grpc.Dial(*grpcEndpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create a new RacingClient
	client := racing.NewRacingClient(conn)

	// Define the routes for each method
	listRacesRoute := "/v1/list-races"
	getRaceRoute := "/v1/get-race/{id}"

	// Register the ListRaces route handler
	r.HandleFunc(listRacesRoute, func(w http.ResponseWriter, r *http.Request) {
		req := &racing.ListRacesRequest{}

		resp, err := client.ListRaces(r.Context(), req)
		if err != nil {
			log.Printf("failed to handle ListRaces request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}).Methods(http.MethodGet)

	// Register the GetRace route handler
	r.HandleFunc(getRaceRoute, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

				// Convert id to int64
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Printf("failed to parse id: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Create a new GetRaceRequest with the ID provided in the URL path
		req := &racing.GetRaceRequest{Id: idInt}


		resp, err := client.GetRace(r.Context(), req)
		if err != nil {
			log.Printf("failed to handle GetRace request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}).Methods(http.MethodGet)

	log.Printf("API server listening on: %s\n", *apiEndpoint)

	return http.ListenAndServe(*apiEndpoint, r)
}
