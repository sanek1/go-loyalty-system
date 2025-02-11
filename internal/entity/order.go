package entity

import "github.com/Azure/go-autorest/autorest/date"

type Order struct {
	ID           uint      `json:"ID"`
	UserID       int       `json:"USER_ID"`
	StatusID     string    `json:"Status"`
	CreationDate string    `json:"CreationDate"`
	Number       string    `json:"Number"`
	UploadedAt   date.Time `json:"Uploaded"`
}
