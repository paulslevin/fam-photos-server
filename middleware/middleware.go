package middleware

import (
	"encoding/json"
	"fam-photos-server/models"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
)

func init() {
	context = AWSContext{}
	context.openSession()
}

// GetFirstName connects to AWS Dynamo DB to get the user's first name.
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

// GetFamilies connects to AWS Dynamo DB to get the families to which the user belongs.
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

// GetImageDataByFamily connects to AWS S3 to get the images for the given family.
func GetImageDataByFamily(w http.ResponseWriter, r *http.Request) {

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

	numberOfImages := len(response.Contents)
	// Here we make the buffer size numberOfImages - 1 since we don't include
	// the enclosing family folder
	imageChan := make(chan models.ImageData, numberOfImages-1)

	var waitGroup sync.WaitGroup

	for _, item := range response.Contents {
		if strings.HasSuffix(*item.Key, "/") {
			continue
		}
		waitGroup.Add(1)
		go getImageMetadata(&waitGroup, item, svc, imageChan)
	}

	waitGroup.Wait()
	close(imageChan)

	x := []map[string]string{}
	for imageData := range imageChan {
		d := map[string]string{}
		d["url"] = imageData.URL
		d["caption"] = imageData.Caption
		x = append(x, d)
	}

	sort.Slice(x, func(i, j int) bool { return x[i]["url"] < x[j]["url"] })

	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(x)

}
