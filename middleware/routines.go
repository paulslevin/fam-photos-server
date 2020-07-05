package middleware

import (
	"fam-photos-server/models"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getImageMetadata(
	waitGroup *sync.WaitGroup,
	item *s3.Object,
	svc *s3.S3,
	imageChan chan models.ImageData,
) {
	defer waitGroup.Done()

	imageData := models.ImageData{}

	response, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(*item.Key),
	})

	if err != nil {
		exitErrorf("Unable to retrieve item %q from bucket %q, %v", *item.Key, bucket, err)
	}

	if caption, ok := response.Metadata["Caption"]; ok {
		imageData.Caption = *caption
	}

	imageURL := fmt.Sprintf("%s%s", s3URL, *item.Key)
	imageURL = strings.ReplaceAll(imageURL, " ", "%20")

	imageData.URL = imageURL

	imageChan <- imageData

}
