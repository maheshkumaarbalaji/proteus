package http

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Structure to create an instance of a web server.
type HttpServer struct {
	// Hostname of the web server instance.
	HostAddress string
	// Port number where web server instance is listening for incoming requests.
	PortNumber int
	// Server socket created and bound to the port number.
	Socket net.Listener
	// Router instance that contains all the routes and their associated handlers.
	innerRouter *Router
}

// Define a static route and map to a static file or folder in the file system.
func (srv *HttpServer) Static(Route string, TargetPath string) error {
	err := srv.innerRouter.addStaticRoute("GET", Route, TargetPath)
	if err != nil {
		return err
	}

	err = srv.innerRouter.addStaticRoute("HEAD", Route, TargetPath)
	if err != nil {
		return err
	}

	return nil
}

// Setup the web server instance to listen for incoming HTTP requests at the given hostname and port number.
func (srv * HttpServer) Listen(PortNumber int, HostAddress string) {
	if PortNumber == 0 {
		srv.PortNumber = getDefaultPort()
	} else {
		srv.PortNumber = PortNumber
	}

	if HostAddress == "" {
		srv.HostAddress = getServerDefaults("hostname")
	} else {
		srv.HostAddress = strings.TrimSpace(HostAddress)
	}

	serverAddress := srv.HostAddress + ":" + strconv.Itoa(srv.PortNumber)
	server, err := net.Listen("tcp", serverAddress)
	if err != nil {
		LogError(fmt.Sprintf("Error occurred while setting up listener socket: %s", err.Error()))
		return
	}

	srv.Socket = server
	defer srv.Socket.Close()
	LogInfo(fmt.Sprintf("Web server is listening at http://%s", serverAddress))

	for {
		clientConnection, err := srv.Socket.Accept()
		if err != nil {
			LogError(fmt.Sprintf("Error occurred while accepting a new client: %s", err.Error()))
			continue
		}

		LogInfo(fmt.Sprintf("A new client - %s has connected to the server", clientConnection.RemoteAddr().String()))
		go srv.handleClient(clientConnection)
	}
}

// Handles incoming HTTP requests sent from each individual client trying to connect to the web server instance.
func (srv *HttpServer) handleClient(ClientConnection net.Conn) {
	defer ClientConnection.Close()
	httpRequest := newRequest(ClientConnection)
	err := httpRequest.read()
	if err != nil {
		LogError(err.Error())
		return
	}

	httpResponse := newResponse(ClientConnection, httpRequest)

	if !isMethodAllowed(httpResponse.Version, strings.ToUpper(strings.TrimSpace(httpRequest.Method))) {
		httpResponse.Status(StatusMethodNotAllowed)
		ErrorHandler(httpRequest, httpResponse)
	} else {
		routeHandler, err := srv.innerRouter.matchRoute(httpRequest)
		if err != nil {
			LogError(err.Error())
			httpResponse.Status(StatusNotFound)
			ErrorHandler(httpRequest, httpResponse)
		} else {
			routeHandler(httpRequest, httpResponse)
		}
	}
}

// Creates a new GET endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Get(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("GET", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new HEAD endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Head(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("HEAD", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new POST endpoint at the given route path and sets the handler function to be invoked when the route is requested by the user.
func (srv *HttpServer) Post(routePath string, handlerFunc Handler) error {
	routePath = strings.TrimSpace(routePath)
	err := srv.innerRouter.addDynamicRoute("POST", routePath, handlerFunc)
	if err != nil {
		return err
	}

	return nil
}