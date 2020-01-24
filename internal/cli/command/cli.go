package command

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

// State represents the state of the Adhesive workflow.
type State struct {
	workflow string
}

// Adhesive represents a running instance of the Adhesive application.
type AdhesiveCli struct {
	State  *State
	Config *Config

	FoundConfigFile bool

	cfn *cloudformation.CloudFormation
	s3  *s3.S3
}

func NewAdhesiveCli() (*AdhesiveCli, error) {
	var (
		config          = NewConfig()
		foundConfigFile bool
	)

	// Try reading configuration from adhesive.toml
	err := LoadConfigFileInto(config, "adhesive.toml")
	if pathErr, ok := err.(*os.PathError); ok && os.IsNotExist(pathErr) {
		foundConfigFile = true
	} else if err != nil {
		return nil, err
	}

	return &AdhesiveCli{
		Config:          config,
		FoundConfigFile: foundConfigFile,
	}, nil
}

func (cli *AdhesiveCli) InitializeClients() error {
	sess, err := session.NewSession(&aws.Config{
		Logger: aws.LoggerFunc(func(args ...interface{}) {
			log.Debug(args...)
		}),
		LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
		Region:   aws.String(cli.Config.Region),
	})
	if err != nil {
		return err
	}

	cli.cfn = cloudformation.New(sess)
	cli.s3 = s3.New(sess)
	return nil
}

func (cli *AdhesiveCli) S3() *s3.S3 {
	return cli.s3
}

func (cli *AdhesiveCli) CloudFormation() *cloudformation.CloudFormation {
	return cli.cfn
}
