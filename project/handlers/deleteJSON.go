package handlers

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
)

func DeleteJSONHandleFunc(w http.ResponseWriter, r *http.Request) {

	// getting phone number from user
	phone := r.URL.Query().Get("phone")
	fmt.Println(phone)

	fileNameSavedInBucket := phone + ".qf"

	ctx := context.Background()

	// creating client
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println("problem in creating storage client:-------", err)
	}

	//deleting user details from bucket
	err2 := client.Bucket("filestorage1_bucket").Object(fileNameSavedInBucket).Delete(ctx)
	if err2 != nil {
		fmt.Println("problem in reading file from bucket:--------", err)
	} else {
		fmt.Println("----------------file deleted from bucket--------")
	}
}
