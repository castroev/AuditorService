package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/xid"

	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/logging"
	s3service "bitbucket.tylertech.com/spy/scm/tcp-auditor/server/pkg/data"
	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/pkg/models"
	"github.com/gin-gonic/gin"
)

var (
	// MaxRoutines number of goroutines
	MaxRoutines = os.Getenv("MAX_ROUTINES")
)

// PostAuditRecordHandler is a gin handler used for creating an audit record
func PostAuditRecordHandler(c *gin.Context) {
	// Respond to caller before processing the request to ensure no blockages
	c.JSON(http.StatusOK, "A new audit record is being processed.")

	auditRecord := &models.AuditRecord{}
	err := json.NewDecoder(c.Request.Body).Decode(&auditRecord)

	if err != nil {
		logging.Logger.Errorf("No audit record was provided, and the service was unable to create the record. Error: %s", err.Error())
		c.JSON(http.StatusBadRequest, "Error: No audit record was provided, and the service was unable to create the record. "+err.Error())
	}

	environment := c.Param("environment")
	product := c.Param("product")
	service := c.Param("service")
	guid := xid.New()

	if product != "" && service != "" {
		stringData, _ := json.Marshal(*auditRecord)
		key := fmt.Sprintf("%s/%s-%s/%s.json", environment, product, service, guid)

		// Create a channel and waitgroup so that we can ensure the goroutine gets released properly after execution
		uploaderChannel := make(chan bool)
		// Create a goroutine
		go func() {
			defer close(uploaderChannel)
			auditFileUploaded := s3service.SaveToS3Bucket(key, stringData)
			uploaderChannel <- auditFileUploaded
		}()

	} else {
		c.JSON(http.StatusBadRequest, "Error: Product and Service names must be specified.")
	}
}

// GetAuditRecordHandler is a gin handler used for getting an audit record
func GetAuditRecordHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get audit handler not implemented yet",
	})
}
