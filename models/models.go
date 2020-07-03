package models

// User represents a user object stored in AWS Dynamo DB.
type User struct {
	Username  string
	FirstName string
	Surname   string
	Families  []string
}

// ImageData represents image information stored in AWS S3.
type ImageData struct {
	URL     string
	Caption string
}
