package http

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"
	"strings"
	"path/filepath"
	"strconv"
	"slices"
	"fmt"
	"github.com/maheshkumaarbalaji/proteus/lib/fs"
)

// Returns the file media type for the given file path.
func getContentType(CompleteFilePath string) (string, bool) {
	pathType, err := fs.GetPathType(CompleteFilePath)
	if err == nil {
		if pathType == fs.FILE_TYPE_PATH {
			fileExtension := filepath.Ext(CompleteFilePath)
			fileExtension = strings.TrimSpace(fileExtension)
			fileExtension = strings.ToLower(fileExtension)
			contentType, exists := AllowedContentTypes[fileExtension]
			if exists {
				return contentType, exists
			} else {
				return strings.TrimSpace(ServerDefaults["content_type"]), true
			}
		}
	}
	return "", false
}

// Returns the default port number from the list of default configuration values.
func getDefaultPort() int {
	portNumberValue := ServerDefaults["port"]
	portNumber, _ := strconv.Atoi(portNumberValue)
	return portNumber
}

// Returns the value for the given key from server default configuration values.
func getServerDefaults(key string) string {
	value := ServerDefaults[strings.TrimSpace(key)]
	value = strings.TrimSpace(value)
	return value
}

// Gets the highest version of HTTP supported by the web server.
func getHighestVersion() string {
	var maxVersion float64 = 0.0
	for versionNo := range Versions {
		currentVersion, err := strconv.ParseFloat(versionNo, 64)
		if err == nil {
			if currentVersion > maxVersion {
				maxVersion = currentVersion
			}
		}
	}

	return fmt.Sprintf("%.1f", maxVersion)
}

// Gets an array of all the versions of HTTP supported by the web server.
func getAllVersions() []string {
	vers := make([]string, 0)
	for versionNo := range Versions {
		tempVer := strings.TrimSpace(versionNo)
		vers = append(vers, tempVer)
	}

	return vers
}

// Gets the list of allowed HTTP methods supported by the web server for the given HTTP version.
func getAllowedMethods(version string) string {
	for versionNo, AllowedMethods := range Versions {
		if strings.EqualFold(versionNo, version) {
			return strings.Join(AllowedMethods, ", ")
		}
	}

	return ""
}

// Checks if the given HTTP method is supported by the web server for the given version.
func isMethodAllowed(version string, requestMethod string) bool {
	for versionNo, AllowedMethods := range Versions {
		if strings.EqualFold(versionNo, version) && slices.Contains(AllowedMethods, requestMethod) {
			return true
		}
	}

	return false
}

// Returns the HTTP response version for the given request version value.
func getResponseVersion(requestVersion string) string {
	isCompatible := false

	for _, version := range getAllVersions() {
		if strings.EqualFold(version, requestVersion) {
			isCompatible = true
			break
		}
	}

	if isCompatible {
		return requestVersion
	} else {
		return getHighestVersion()
	}
}

// Creates and returns pointer to a new instance of HTTP request.
func newRequest(Connection net.Conn) *HttpRequest {
	var httpRequest HttpRequest
	httpRequest.initialize()
	reader := bufio.NewReader(Connection)
	httpRequest.setReader(reader)
	return &httpRequest
}

// Creates and returns pointer to a new instance of HTTP response.
func newResponse(Connection net.Conn, request *HttpRequest) *HttpResponse {
	var httpResponse HttpResponse
	httpResponse.initialize(getResponseVersion(request.Version))
	writer := bufio.NewWriter(Connection)
	httpResponse.setWriter(writer)
	return &httpResponse
}

// Creates and returns pointer to a new instance of Router.
func newRouter() *Router {
	router := new(Router)
	router.Routes = make([]Route, 0)
	router.RouteTree = createTree()
	return router
}

// Returns the current UTC time in RFC 1123 format.
func getRfc1123Time() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC1123)
}

// Returns an instance of HTTP web server.
func NewServer() *HttpServer {
	if SrvLogger == nil {
		SrvLogger = log.New(os.Stdout, "", log.Ldate | log.Ltime)
	}

	if ServerInstance == nil {
		var server HttpServer
		server.HostAddress = "";
		server.PortNumber = 0
		server.innerRouter = newRouter()
		ServerInstance = &server
		return &server
	}

	return ServerInstance
}