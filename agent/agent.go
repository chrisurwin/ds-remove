package agent

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Agent struct {
	sync.WaitGroup

	probePeriod time.Duration
	httpClient  http.Client
	awsConfig   *aws.Config
	awsSession  *session.Session
}

const (
	AWSMetaData = "test.aws-cleanup.demo-cattle.route53.dns-route.ml:7777"
)

func NewAgent(probePeriod time.Duration, arn string) *Agent {
	region := getAWSInfo("/latest/meta-data/placement/availability-zone")
	region = region[:len(region)-1]
	fmt.Print(region)
	sess, _ := session.NewSession()

	var c = &aws.Config{Region: aws.String(region)}

	if arn != "" {
		creds := stscreds.NewCredentials(sess, arn)
		c = &aws.Config{
			Region:      aws.String(region),
			Credentials: creds,
		}
	}
	return &Agent{
		probePeriod: probePeriod,
		httpClient: http.Client{
			Timeout: time.Duration(2 * time.Second),
		},
		awsSession: sess,
		awsConfig:  c,
		//log:         log.WithField("pkg", "agent"),
	}
}

func (a *Agent) Start() {

	t := time.NewTicker(a.probePeriod)
	for _ = range t.C {
		go a.checkDiskSpace()

		a.Wait()
	}
}

func (a *Agent) checkDiskSpace() {
	log.Info("Checking diskspace... ")
	var stat syscall.Statfs_t

	wd, _ := os.Getwd()

	syscall.Statfs(wd, &stat)

	// Available blocks * size per block = available space in bytes
	var free = stat.Bfree * uint64(stat.Bsize)
	if free > 50000000 {
		log.Info("Disk Space check passed..")
	} else {
		log.Error("disk full")
		/*var i = getAWSInfo("/latest/meta-data/instance-id")
		var r bool = awsfunctions.TerminateInstance(i, a.awsSession, a.awsConfig)
		if r == true {
			log.Info("Termination initiated")
		} else {
			log.Error("Terminateion failed")
		}*/
	}

}

func getAWSInfo(path string) string {

	resp, err := http.Get("http://" + AWSMetaData + path)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	return string(body)

}
