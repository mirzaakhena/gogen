package remoting

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"reflect"
)

type RemoteListener struct {
	handler any
	port    int
}

type RemoteCaller struct {
	port int
}

func NewRemoteListener(port int) *RemoteListener {
	return &RemoteListener{
		port: port,
	}
}

func NewRemoteCaller(port int) *RemoteCaller {
	return &RemoteCaller{port: port}
}

func (r *RemoteListener) SetHandler(Handler any) {
	r.handler = Handler
}

func (r *RemoteListener) Run() {

	err := rpc.Register(r.handler)
	if err != nil {
		log.Fatal(err.Error())
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", r.port))
	if err != nil {
		log.Fatal("listen error:", err)
	}

	fmt.Printf("server %v is running and waiting for request from client...\n", reflect.TypeOf(r.handler))

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("serve error:", err)
		return
	}

}

func (r *RemoteCaller) Call(methodName string, args any, reply any) error {

	var err error

	client, err := rpc.DialHTTP("tcp", fmt.Sprintf(":%d", r.port))
	if err != nil {
		return err
	}

	defer func(client *rpc.Client) {
		tempErr := client.Close()
		if err == nil {
			err = tempErr
		}
	}(client)

	err = client.Call(methodName, args, &reply)
	if err != nil {
		return err
	}

	return nil

}
