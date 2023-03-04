package main

import (
	"api_gateway/bootstrap"
	"log"
)

func main() {

	//conn, err := grpc.Dial("localhost:8086", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//
	//in := &api.ListRequest{
	//	Offset:  0,
	//	Limit:   1,
	//	OrderBy: "id",
	//}
	//out := new(api.ListResponse)
	//
	//err = conn.Invoke(context.Background(), "/pb.AuthService/List", in, out)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("%v", out)

	err := bootstrap.Run()
	if err != nil {
		log.Fatal(err)
	}
}
