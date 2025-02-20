package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net"
	"os"
	"runtime"
	"strconv"

	"github.com/ibm/opentalaria/config"

	"github.com/ibm/opentalaria/protocol"

	"github.com/ibm/opentalaria/api"
	"golang.org/x/sync/semaphore"
)

type Server struct {
	host   string
	port   string
	config *config.Config
}

type Client struct {
	conn   net.Conn
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	var host, port string
	if len(config.Broker.Listeners) > 0 {
		listener := config.Broker.Listeners[0]
		host = listener.Host
		port = strconv.Itoa(int(listener.Port))
	}

	return &Server{
		host:   host,
		port:   port,
		config: config,
	}
}

func (server *Server) Run() {
	ctx := context.TODO()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		slog.Error("error creating tcp listener", "err", err)
		return
	}
	defer listener.Close()

	slog.Info(fmt.Sprintf("tcp server listening on %s:%s", server.host, server.port))

	cpu := os.Getenv("GOMAXPROCS")
	if cpu == "" {
		cpu = "0"
	}
	numberOfCpu, err := strconv.Atoi(cpu)
	if err != nil {
		slog.Error("error creating connection", "error", err)
		return
	}
	// Adding more CPU's only helps up to number of available Go routines
	// For example GOMAXPROCS(8) and semaphore.NewWeighted(8) means each Go routine will be executed on different CPU
	// However if we set GOMAXPROCS(4) and semaphore.NewWeighted(8) we will have only 4 CPU's to handle 8 Go routines
	runtime.GOMAXPROCS(numberOfCpu)
	slog.Debug("number of available CPU's ", "GOMAXPROCS", numberOfCpu)

	var conCapacity int64
	conPoolStr := server.config.Env.GetString("max.connections")
	if conPoolStr == "" {
		//If env variable max.connections was not set we use default val of MaxInt64
		conCapacity = math.MaxInt64
	} else {
		//If env variable is set, we need to convert it to int64
		c, err := strconv.ParseInt(conPoolStr, 10, 64)
		if err != nil {
			slog.Error("error setting max.connections", "error", err)
			return
		}
		conCapacity = c
	}

	slog.Debug("max.connections set to ", "max.connections", conCapacity)

	//semaphore package mimics a typical “worker pool” pattern,
	//but without the need to explicitly shut down idle workers when the work is done
	sem := semaphore.NewWeighted(int64(conCapacity))

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("error accepting tcp connections", "err", err)
		}

		client := &Client{
			conn:   conn,
			config: server.config,
		}

		if err := sem.Acquire(ctx, 1); err != nil {
			slog.Error("Failed to acquire semaphore: %v", "err", err)
			break
		}
		go func() {
			defer sem.Release(1)
			client.handleRequest()
		}()
	}
	// Acquire all of the tokens to wait for any remaining workers to finish
	if err := sem.Acquire(ctx, int64(conCapacity)); err != nil {
		slog.Error("Failed to acquire semaphore: %v", "err", err)
	}
}

func (client *Client) handleRequest() {
	defer client.conn.Close()

Exit:
	// read from socket until there are no more bytes left.
	for {
		// first 4 bytes contain the message size
		sizeBytes := make([]byte, 4)
		_, err := io.ReadFull(client.conn, sizeBytes[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("tcp read error", "err", err)
			break
		}
		size := binary.BigEndian.Uint32(sizeBytes)

		// read the rest of the message into the buffer.
		messageBytes := make([]byte, size)

		if _, err := io.ReadFull(client.conn, messageBytes[:]); err != nil {
			slog.Error("error decoding message", "err", err)
			break
		}

		// save the message to a file to use for testing later.
		// encoded := hex.EncodeToString(messageBytes)
		// fmt.Println(encoded)

		// We parse the header twice, first time parse only API key and API version, from which we can
		// infer the correct header version and then parse that again in the API code to get the full header.
		header := &protocol.RequestHeader{}
		protocol.VersionedDecode(messageBytes, header, 1)

		slog.Debug(header.String())

		var apiHandler api.API
		switch header.RequestApiKey {
		case (&protocol.ApiVersionsRequest{}).GetKey():
			req, err := makeRequest(messageBytes,
				client.conn,
				(&protocol.ApiVersionsRequest{Version: header.RequestApiVersion}).GetHeaderVersion(),
				client.config)
			if err != nil {
				slog.Error("error creating request", "err", err)
				// This break exits the outer for loop and closes the socket connection.
				// If there is an error in the metadata exchange for example, we don't want to continue consuming the rest of the APIs.
				break Exit
			}
			apiHandler = api.APIVersionsAPI{Request: req, Config: client.config}
		case (&protocol.MetadataRequest{}).GetKey():
			req, err := makeRequest(messageBytes,
				client.conn,
				(&protocol.MetadataRequest{Version: header.RequestApiVersion}).GetHeaderVersion(),
				client.config)
			if err != nil {
				slog.Error("error creating request", "err", err)
				break Exit
			}
			apiHandler = api.MetadataAPI{Request: req, Config: client.config}
		case (&protocol.ProduceRequest{}).GetKey():
			req, err := makeRequest(messageBytes,
				client.conn,
				(&protocol.ProduceRequest{Version: header.RequestApiVersion}).GetHeaderVersion(),
				client.config)
			if err != nil {
				slog.Error("error creating request", "err", err)
				break Exit
			}
			apiHandler = api.ProduceAPI{Request: req, Config: client.config}
		case (&protocol.CreateTopicsRequest{}).GetKey():
			req, err := makeRequest(messageBytes,
				client.conn,
				(&protocol.CreateTopicsRequest{Version: header.RequestApiVersion}).GetHeaderVersion(),
				client.config)
			if err != nil {
				slog.Error("error creating request", "err", err)
				break Exit
			}
			apiHandler = api.CreateTopicsAPI{Request: req, Config: client.config}
		default:
			slog.Error("Unknown API key", "key", header.RequestApiKey)
		}

		err = api.HandleResponse(apiHandler)
		if err != nil {
			slog.Error("error handling response", "err", err)
			break
		}
	}
}

func makeRequest(msg []byte, conn net.Conn, headerVersion int16, config *config.Config) (api.Request, error) {
	// parse the full header, based on API key and version
	header := &protocol.RequestHeader{}
	headerSize, err := protocol.VersionedDecode(msg, header, headerVersion)
	if err != nil {
		return api.Request{}, err
	}

	return api.Request{
		Header:  *header,
		Message: msg[headerSize:],
		Conn:    conn,
	}, nil
}
