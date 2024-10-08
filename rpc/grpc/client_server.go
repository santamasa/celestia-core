package coregrpc

import (
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cmtnet "github.com/tendermint/tendermint/libs/net"
)

// Config is an gRPC server configuration.
type Config struct {
	MaxOpenConnections int
}

// StartGRPCServer starts a new gRPC BroadcastAPIServer using the given
// net.Listener.
// NOTE: This function blocks - you may want to call it in a go-routine.
func StartGRPCServer(ln net.Listener) error {
	grpcServer := grpc.NewServer()
	RegisterBroadcastAPIServer(grpcServer, &broadcastAPI{})
	api := newBlockAPI()
	go api.listenForHeights()
	RegisterBlockAPIServer(grpcServer, api)
	return grpcServer.Serve(ln)
}

// StartGRPCClient dials the gRPC server using protoAddr and returns a new
// BroadcastAPIClient.
func StartGRPCClient(protoAddr string) BroadcastAPIClient {
	conn, err := grpc.Dial(protoAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialerFunc))
	if err != nil {
		panic(err)
	}
	return NewBroadcastAPIClient(conn)
}

// StartBlockAPIGRPCClient dials the gRPC server using protoAddr and returns a new
// BlockAPIClient.
func StartBlockAPIGRPCClient(protoAddr string) BlockAPIClient {
	conn, err := grpc.Dial(protoAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialerFunc))
	if err != nil {
		panic(err)
	}
	return NewBlockAPIClient(conn)
}

func dialerFunc(ctx context.Context, addr string) (net.Conn, error) {
	return cmtnet.Connect(addr)
}
