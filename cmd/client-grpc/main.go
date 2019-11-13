package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	v1 "aniqma/aniqma/crudgrpc/pkg/api/v1"
)

//Users
type Users struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt time.Time
}

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {
	//parse json
	mux := http.NewServeMux()

	mux.HandleFunc("/echo", echoHandler)
	http.ListenAndServe(":5000", mux)

	// id := res1.Id

	// // Read
	// req2 := v1.ReadRequest{
	// 	Api: apiVersion,
	// 	Id:  id,
	// }
	// res2, err := c.Read(ctx, &req2)
	// if err != nil {
	// 	log.Fatalf("Read failed: %v", err)
	// }
	// log.Printf("Read result: <%+v>\n\n", res2)

	// // Update
	// req3 := v1.UpdateRequest{
	// 	Api: apiVersion,
	// 	ToDo: &v1.ToDo{
	// 		Id:          res2.ToDo.Id,
	// 		Title:       res2.ToDo.Title,
	// 		Description: res2.ToDo.Description + " + updated",
	// 		Reminder:    res2.ToDo.Reminder,
	// 	},
	// }
	// res3, err := c.Update(ctx, &req3)
	// if err != nil {
	// 	log.Fatalf("Update failed: %v", err)
	// }
	// log.Printf("Update result: <%+v>\n\n", res3)

	// // Call ReadAll
	// req4 := v1.ReadAllRequest{
	// 	Api: apiVersion,
	// }
	// res4, err := c.ReadAll(ctx, &req4)
	// if err != nil {
	// 	log.Fatalf("ReadAll failed: %v", err)
	// }
	// log.Printf("ReadAll result: <%+v>\n\n", res4)

	// // Delete
	// req5 := v1.DeleteRequest{
	// 	Api: apiVersion,
	// 	Id:  id,
	// }
	// res5, err := c.Delete(ctx, &req5)
	// if err != nil {
	// 	log.Fatalf("Delete failed: %v", err)
	// }
	// log.Printf("Delete result: <%+v>\n\n", res5)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	user := Users{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	user.CreatedAt = time.Now().Local()

	userJson, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(userJson)

	// get configuration
	address := flag.String("server", "localhost:9090", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewUsersServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	// pfx := t.Format(time.RFC3339Nano)

	// Call Create
	req1 := v1.CreateRequest{
		Api: apiVersion,
		USers: &v1.Users{
			Username: user.Username,
			Password: user.Password,
			Reminder: reminder,
		},
	}

	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

}
