package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/config"
	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/handlers"
	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/logging"
	s3service "bitbucket.tylertech.com/spy/scm/tcp-auditor/server/pkg/data"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ddgin "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// MaxWorker number of goroutines
	MaxWorker = os.Getenv("MAX_WORKERS")
	// MaxQueue is max number of jobs in queue
	MaxQueue = os.Getenv("MAX_QUEUE")
)

func main() {
	// Setup datadog tracing
	addr := net.JoinHostPort(
		os.Getenv("DD_AGENT_HOST"),
		os.Getenv("DD_TRACE_AGENT_PORT"),
	)
	tracer.Start(tracer.WithAgentAddr(addr))
	defer tracer.Stop()

	// Configure log formatting
	logging.InitLogger()

	// Read service configuration
	config.InitConfig()

	// create shared s3 session
	s3service.CreateS3Session()

	r := gin.Default()

	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// TODO: implement ccf token handoff from sdk
	// r.Use(tokenvalidationmiddleware.JwtAuthenticationMiddleware())

	//Setup Middlewares
	r.Use(ddgin.Middleware("tcp-auditor-go"))

	// Setup Handlers
	r.GET("/audit/", handlers.GetAuditRecordHandler)
	r.POST("/audit/:environment/:product/:service", handlers.PostAuditRecordHandler)
	// Register pprof handlers
	pprof.Register(r, "dev/pprof")

	s.ListenAndServe()
}
