package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

func getkey_Signature(phone string, ctx context.Context, clientDataStore *datastore.Client) {

	q := new(Request)
	nameKey1 := datastore.NameKey("phone", phone, nil)

	//finding keys in datastore
	err := clientDataStore.Get(ctx, nameKey1, q)
	if err != nil {
		fmt.Println("Keys not available:------------ ", err)
		return
	} else {
		fmt.Println("-----------------------Keys are available-------------- ")

		// decoding keys find from datastore and unmarshal them.
		keyBytesPrivate, _ := base64.StdEncoding.DecodeString(q.Private_Key)
		keyBytesPublic, _ := base64.StdEncoding.DecodeString(q.Public_Key)
		err := json.Unmarshal(keyBytesPrivate, &KEC.privateKey)
		if err != nil {
			fmt.Println(err)
			return
		}
		err2 := json.Unmarshal(keyBytesPublic, &KEC.publicKey)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		fmt.Println("privateKey is:- ", q.Private_Key)
		fmt.Println("publicKey is:- ", q.Public_Key)
	}
}

func ReadJSONHandleFunc(w http.ResponseWriter, r *http.Request) {

	fmt.Println("------------------------we are in readJSON.go==================")

	// getting phone number from user
	phone := r.URL.Query().Get("phone")
	fmt.Println(phone)

	fileNameSavedInBucket := phone + ".qf"

	ctx := context.Background()

	//creating client
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println("problem in creating storage client:-------", err)
	}

	//reading file from bucket
	rc, err := client.Bucket("filestorage1_bucket").Object(fileNameSavedInBucket).NewReader(ctx)
	if err != nil {
		fmt.Println("problem in reading file from bucket:--------", err)
	}

	if err := rc.Close(); err != nil {
		fmt.Println("error in creation of client:--------- ", err.Error())
	}

	tempData, err2 := ioutil.ReadAll(rc)
	if err2 != nil {
		fmt.Println(err)
	}

	var readable_data writeInFile

	err3 := json.Unmarshal(tempData, &readable_data)
	if err3 != nil {
		fmt.Println("problem in unmarshal of data:--------", err3)
	}

	keyForEncryptionDecryption := "thisis32bitlongpassphraseimusing"

	//---------decrypting data
	data := decryptAES(keyForEncryptionDecryption, readable_data.Payload)

	//------------------checking signature
	clientDataStore, error := datastore.NewClient(ctx, "urlmanager-386505")
	if error != nil {
		fmt.Println("error in creation of client: ", error.Error())
	}
	getkey_Signature(phone, ctx, clientDataStore)
	checkingSignature := verifySignatureEC(data, readable_data.Signature)
	if checkingSignature {
		fmt.Println("-----------signature is valid!!!!!!!------------")
	} else {
		fmt.Println("-----------signature is invalid!!!!!!!------------")
		return
	}

	//creating file in local
	f, err4 := os.Create(readable_data.Header.FileName)
	if err4 != nil {
		fmt.Println("problem in creating file:-----", err4)
	}
	defer f.Close()
	f.Write([]byte(data))

	//sending data to client
	w.Header().Add("Content-Type", readable_data.Header.FileType)
	w.Write([]byte(data))
	fmt.Println("---------send file to user-----------")
}
