package server

import (
	"fmt"
	"time"
	"strings"
	"net/http"
	"net/url"

	"github.com/eug48/fhir/auth"
)

// Config is used to hold information about the configuration of the FHIR server.
type Config struct {
	// ServerURL is the full URL for the root of the server. This may be used
	// by other middleware to compute redirect URLs
	ServerURL string

	// Auth determines what, if any authentication and authorization will be used
	// by the FHIR server
	Auth auth.Config

	// IndexConfigPath is the path to an indexes.conf configuration file, specifying
	// what mongo indexes the server should create (or verify) on startup
	IndexConfigPath string

	// DatabaseURI is the url of the mongo replica set to use for the FHIR database.
	// A replica set is required for transactions support
	// e.g. mongodb://db1:27017,db2:27017/?replicaSet=rs1
	DatabaseURI string

	// DatabaseName is the name of the mongo database used for the fhir database.
	// Typically this will be the default DatabaseName "fhir".
	DatabaseName string

	// DatabaseSocketTimeout is the amount of time the mgo driver will wait for a response
	// from mongo before timing out.
	DatabaseSocketTimeout time.Duration

	// DatabaseOpTimeout is the amount of time GoFHIR will wait before killing a long-running
	// database process. This defaults to a reasonable upper bound for slow, pipelined queries: 30s.
	DatabaseOpTimeout time.Duration

	// DatabaseKillOpPeriod is the length of time between scans of the database to kill long-running ops.
	DatabaseKillOpPeriod time.Duration

	// CountTotalResults toggles whether the searcher should also get a total
	// count of the total results of a search. In practice this is a performance hit
	// for large datasets.
	CountTotalResults bool

	// EnableCISearches toggles whether the mongo searches uses regexes to maintain
	// case-insesitivity when performing searches on string fields, codes, etc.
	EnableCISearches bool

	// Whether to support storing previous versions of each resource
	EnableHistory bool

	// Whether to allow retrieving resources with no meta component,
	// meaning Last-Modified & ETag headers can't be generated (breaking spec compliance)
	// May be needed to support previous databases
	AllowResourcesWithoutMeta bool

	// ValidatorURL is an endpoint to which validation requests will be sent
	ValidatorURL string

	// ReadOnly toggles whether the server is in read-only mode. In read-only
	// mode any HTTP verb other than GET, HEAD or OPTIONS is rejected.
	ReadOnly bool

	// Enables requests and responses using FHIR XML MIME-types
	EnableXML bool

	// Debug toggles debug-level logging.
	Debug bool
}

// DefaultConfig is the default server configuration
var DefaultConfig = Config{
	ServerURL:             "",
	IndexConfigPath:       "config/indexes.conf",
	DatabaseURI:           "mongodb://localhost:27017/?replicaSet=rs0",
	DatabaseName:          "fhir",
	DatabaseSocketTimeout: 2 * time.Minute,
	DatabaseOpTimeout:     90 * time.Second,
	DatabaseKillOpPeriod:  10 * time.Second,
	Auth:                  auth.None(),
	EnableCISearches:      true,
	EnableHistory:         true,
	EnableXML:             true,
	CountTotalResults:     true,
	ReadOnly:              false,
	Debug:                 false,
}

func (config *Config) responseURL(r *http.Request, paths ...string) *url.URL {

	if config.ServerURL != "" {
		theURL := fmt.Sprintf("%s/%s", strings.TrimSuffix(config.ServerURL, "/"), strings.Join(paths, "/"))
		responseURL, err := url.Parse(theURL)

		if err == nil {
			return responseURL
		}
	}

	responseURL := url.URL{}

	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		responseURL.Scheme = "https"
	} else {
		responseURL.Scheme = "http"
	}
	responseURL.Host = r.Host
	responseURL.Path = fmt.Sprintf("/%s", strings.Join(paths, "/"))

	return &responseURL
}