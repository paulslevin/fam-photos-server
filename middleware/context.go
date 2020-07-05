package middleware

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	s3URL  = "https://fam-photos-photos.s3.eu-west-2.amazonaws.com/"
	bucket = "fam-photos-photos"
)

// AWSContext represents AWS session.
type AWSContext struct {
	session session.Session
}

// openSession creates a session with AWS.
func (a *AWSContext) openSession() {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	})
	a.session = *sess
}

var context AWSContext
