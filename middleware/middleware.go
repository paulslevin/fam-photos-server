package middleware

import (
	"encoding/json"
	"fam-photos-server/models"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	s3URL  = "https://fam-photos-photos.s3.eu-west-2.amazonaws.com/"
	bucket = "fam-photos-photos"
)

type AWSContext struct {
	session session.Session
}

func (a *AWSContext) OpenSession() {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	})
	a.session = *sess
}

var context AWSContext

func init() {
	context = AWSContext{}
	context.OpenSession()
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func GetFirstName(w http.ResponseWriter, r *http.Request) {

	var username interface{}

	json.NewDecoder(r.Body).Decode(&username)

	svc := dynamodb.New(&context.session)
	tableName := "users"
	result, _ := svc.GetItem(
		&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"username": {
					S: aws.String(fmt.Sprint(username)),
				},
			},
		},
	)

	var user models.User
	dynamodbattribute.Unmarshal(result.Item["firstName"], &user.FirstName)

	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(user.FirstName)
}

func GetFamilies(w http.ResponseWriter, r *http.Request) {

	var username interface{}

	json.NewDecoder(r.Body).Decode(&username)

	svc := dynamodb.New(&context.session)
	tableName := "users"
	result, _ := svc.GetItem(
		&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"username": {
					S: aws.String(fmt.Sprint(username)),
				},
			},
		},
	)

	var user models.User
	dynamodbattribute.Unmarshal(result.Item["families"], &user.Families)

	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(user.Families)

}

func GetImageURLsByFamily(w http.ResponseWriter, r *http.Request) {

	svc := s3.New(&context.session)

	var family interface{}
	json.NewDecoder(r.Body).Decode(&family)

	response, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(fmt.Sprint(family)),
	})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	urls := []string{}

	for _, item := range response.Contents {

		// Exclude the folder object from results
		if strings.HasSuffix(*item.Key, "/") {
			continue
		}

		imageURL := fmt.Sprintf("%s%s", s3URL, *item.Key)
		imageURL = strings.ReplaceAll(imageURL, " ", "%20")
		fmt.Println("URL", imageURL)
		urls = append(urls, imageURL)
	}

	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(urls)

}
