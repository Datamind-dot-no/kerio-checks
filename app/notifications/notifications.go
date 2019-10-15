package notifications

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/datamind-dot-no/kerio-checks/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Notifications struct {
	kerioChkConf *config.Config
}

func New(conf *config.Config) *Notifications {
	return &Notifications{
		kerioChkConf: conf,
	}
}

// SendNotification function is used to alert the support crew is an issue is
// flagged by sending an email.
func (n *Notifications) SendNotification(QueueLength int) (*ses.SendEmailOutput, error) {

	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Replace the SES e-mail template data key placeholders with current values
	theSubject := strings.Replace(n.kerioChkConf.SubjectT, "{{servername}}", n.kerioChkConf.ServerName, -1)
	theHTMLBody := strings.Replace(n.kerioChkConf.HTMLBodyT, "{{servername}}", n.kerioChkConf.ServerName, -1)
	theHTMLBody = strings.Replace(theHTMLBody, "{{queuelength}}", strconv.Itoa(QueueLength), -1)
	theTextBody := strings.Replace(n.kerioChkConf.TextBodyT, "{{servername}}", n.kerioChkConf.ServerName, -1)
	theTextBody = strings.Replace(theTextBody, "{{queuelength}}", strconv.Itoa(QueueLength), -1)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(n.kerioChkConf.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(n.kerioChkConf.CharSet),
					Data:    aws.String(theHTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(n.kerioChkConf.CharSet),
					Data:    aws.String(theTextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(n.kerioChkConf.CharSet),
				Data:    aws.String(theSubject),
			},
		},
		Source: aws.String(n.kerioChkConf.Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return nil, fmt.Errorf(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return nil, fmt.Errorf(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return nil, fmt.Errorf(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				return nil, fmt.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, fmt.Errorf(err.Error())
		}

		return nil, err
	}

	fmt.Println("Email Sent to address: " + n.kerioChkConf.Recipient)
	fmt.Println(result)
	return result, nil
}
