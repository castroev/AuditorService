# TCP-Auditor-Go
<img src="https://eks-tcpci.s3-us-west-2.amazonaws.com/auditgopher.png" width="300" />

## Features

* [GO 1.12](https://golang.org/)
* [Amazon AWS SDK for S3](https://docs.aws.amazon.com/sdk-for-go/api/service/s3/)
* [Docker](https://www.docker.com/)

## Golang Project Prerequisites


  * Install the Go language
  * Make sure you have a $GOPATH set to /home/go
  * Create the following directory structure in your GoPATH: `/src/bitbucket.tylertech.com/spy/scm/`
  * Clone this repo into the `/scm/` folder
  * Make sure the env variable: `GO111MODULE=on` is set
  * On a command line, navigate to the project directory and run `go get`. This will install all dependencies needed to run the project.
  * Test by running the docker-compose.yaml file which spins up a MongoDB container, and then run the service with `f5`

## Project Description & Usage

This is the TCP-Auditor project. It is a service that inserts audit records into a specified AWS S3 bucket. The records are sent from the Platform Service API when POST, PUT, and DELETE actions are taken on the API. This is a service that maintains audit records of all data modifications that occur on our core services. 

POST endpoint example URL (fill in the blanks) `http://localhost:8080/audit/{env}/{product}/{service}`
With the S3 configuration set to use the `auditrecorde2etesting` S3 bucket, the above url will create a record in S3 in a subdirectory of the `auditrecorde2etesting` bucket with the following path `{env}/{product}/{service}/{guid}.json`

* Use this sample data with postman: 
```JSON

{
    "firstName":   "wowsa",
	"lastName":    "boogie",
    "email":       "billy.tester@tylertech.com",
	"payload":     "{fake: virus}",
	"actionType":  "tcp-user",
	"requestUrl":  "https://fakingmys3service.com/testing/",
	"identitySub": "12ufk38383nfn!34tu",
	"action":      "DeleteAllTheThingsVirusRunner()",
	"timeStamp": "2019-09-27T07:00:00.000Z"
}
```

### local development env

You can utilize the included docker-compose project, which will connect itself to the platform-dev-compose project. NOTE: For this to work, you will need to set the following environment variables on the tcp-platformservice appsettings:

```json
{
    "AuditServiceUri": "http://tcp-auditor:8080/audit/localdev/CloudPlatform/testing/",
    "EnablePlatformAuditing": true
}
```


### Unit testing
### TO RUN THIS TEST, YOU MUST SET THE `AWS_ACCESS_KEY_ID` AND `AWS_SECRET_ACCESS_KEY` IN YOUR MACHINES ENV. 

Right now there is one test that ultimately tests the entire service. Open VSCode and navigate to `auditing_test.go`. Ensure you have installed the VSCode Golang extensions. You should see the "run test" and "debug test" buttons just above the following code:
```go
   func Test_AuditServiceCreate(t *testing.T) {
	   t.Run("CreateAuditRecord", createAuditRecord_should_insert_content_into_s3)
   }
```
NOTE: You will need consul running in the background as well, since the auditor service reads its configuration from consul. 



