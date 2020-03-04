// NOTE: You must have docker-compose running as well as set the AWS credentials in your env to run this test.
package handlers_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/config"
	s3service "bitbucket.tylertech.com/spy/scm/tcp-auditor/server/pkg/data"
	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/pkg/models"
	"github.com/rs/xid"
)

func Test_AuditServiceCreate(t *testing.T) {
	t.Run("CreateAuditRecord", createAuditRecord_should_insert_content_into_s3)
}

func createAuditRecord_should_insert_content_into_s3(t *testing.T) {

	// TO RUN THIS TEST, YOU MUST SET THE AWS_ACCESS_KEY_ID AND AWS_SECRET_ACCESS_KEY IN YOUR ENV.
	os.Setenv("AWS_ACCESS_KEY_ID", "")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "")

	config.InitConfig()
	auditRecord := models.AuditRecord{
		FirstName:   "wowsa",
		LastName:    "boogie",
		Email:       "billy.tester@tylertech.com",
		Payload:     "{fake: virus}",
		ActionType:  "tcp-user",
		RequestURL:  "https://fakingmys3service.com/testing/",
		IdentitySub: "12ufk38383nfn!34tu",
		Action:      "DeleteAllTheThingsVirusRunner()",
		TimeStamp:   time.Now(),
	}

	product := "TylerCloudPlatform"
	service := "AuditServiceE2E"
	guid := xid.New()
	environment := "e2eTestLocal"

	t.Logf("Sending Audit Record to S3 for Product: %s, and Service: %s.", product, service)

	stringData, _ := json.Marshal(auditRecord)
	key := fmt.Sprintf("%s/%s-%s/%s.json", environment, product, service, guid)
	s3service.SaveToS3Bucket(key, stringData)
	return
}
