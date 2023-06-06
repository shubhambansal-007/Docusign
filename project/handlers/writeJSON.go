package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

type Request struct {
	Public_Key  string `json:"Public_Key"`
	Private_Key string `json:"Private_Key"`
}

type header struct {
	Algo     string `json:"algo"`
	Api      string `json:"api"`
	Ext      string `json:"ext"`
	FileType string `json:"fileType"`
	FileName string `json:"fileName"`
}
type writeInFile struct {
	Header    header `json:"header"`
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

type Data struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	ProjectID    string `json:"projectID"`
	Valid_kind   string `json:"valid_kind"`
	Invalid_kind string `json:"invalid_kind"`
	Text         string `json:"text"`
}

var configFile Data

func getkey_phone(phone string, ctx context.Context, clientDataStore *datastore.Client) {

	q := new(Request)
	nameKey1 := datastore.NameKey("phone", phone, nil)

	//finding keys in datastore
	err := clientDataStore.Get(ctx, nameKey1, q)
	if err != nil {
		fmt.Println("Keys not available, so we have to create new keypair and error is:------------ ", err)
		createNewKey(phone, ctx, clientDataStore)
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

func createNewKey(phone string, ctx context.Context, clientDataStore *datastore.Client) {

	//generating new keys
	generateKeyEC()

	q := new(Request)
	keyBytesPrivate, _ := json.Marshal(KEC.privateKey)
	keyBytesPublic, _ := json.Marshal(KEC.publicKey)
	q.Private_Key = base64.StdEncoding.EncodeToString(keyBytesPrivate)
	q.Public_Key = base64.StdEncoding.EncodeToString(keyBytesPublic)

	//storing keys to datastore
	key_generated := datastore.NameKey("phone", phone, nil)
	fmt.Println("Value sent to the datastore is: ", q)
	_, err := clientDataStore.Put(ctx, key_generated, q)
	if err != nil {
		fmt.Println("error in adding Keys to database:", err.Error())
	}
}

// taking request of client and adding in respective json file
func (data *Data) WriteJSONHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request is: ", r.Body)
	//connect with config.json file
	configFile = *data

	fmt.Println("------------------------we are in writeJSON.go==================")
	// setting client to access datastore
	ctx := context.Background()
	clientDataStore, error := datastore.NewClient(ctx, configFile.ProjectID)
	if error != nil {
		fmt.Println("error in creation of client: ", error.Error())
	}

	//doing multipartform to get data
	r.ParseMultipartForm(1024)
	phones := r.MultipartForm.Value["phone"]
	files := r.MultipartForm.File["file"]

	// fetching file and phone number from request
	fileType := files[0].Header.Get("Content-Type")
	contentType := fileType
	fileName := files[0].Filename
	phone := phones[0]
	f, _ := files[0].Open()
	fileType = fileType[:5]
	fmt.Println("phone number received:-    ", phone)
	fmt.Println("type of file received:-   ", fileType)
	if fileType == "audio" {
		fmt.Println("Audio Type Not-Supported")
		return
	}
	if fileType == "video" {
		fmt.Println("Video Type Not-Supported")
		return
	}
	if fileType == "image" {
		fmt.Println("Image Type Not-Supported")
		return
	}

	//reading from file
	fileData, error2 := ioutil.ReadAll(f)
	if error2 != nil {
		fmt.Println(error2.Error())
	}
	dataReceivedFromUser := string(fileData)
	fmt.Println("Data inside file:-    ", dataReceivedFromUser)

	//getting keys from user
	getkey_phone(phone, ctx, clientDataStore)

	fmt.Println("Here we are 1:----------------------------------------")

	//encoding , decoding and signature
	var writtable_data writeInFile
	keyForEncryptionDecryption := "thisis32bitlongpassphraseimusing"
	//fmt.Println("dataReceivedFromUser:--------------", dataReceivedFromUser)
	encoded_text := encryptAES(dataReceivedFromUser, keyForEncryptionDecryption)
	//fmt.Println("Encrypted text is:--------------", encoded_text)
	decryptAES(keyForEncryptionDecryption, encoded_text)
	//fmt.Println("Decrypted text is:--------------", str)
	signature := createSignEC(dataReceivedFromUser)

	fmt.Println("Here we are 2:----------------------------------------")
	//writting data in qf type of file
	writtable_data.Header.Algo = "sha-256"
	writtable_data.Header.Api = "AES-GCM"
	writtable_data.Header.Ext = "qf"
	writtable_data.Header.FileType = contentType
	writtable_data.Header.FileName = fileName
	writtable_data.Payload = encoded_text
	writtable_data.Signature = signature
	fmt.Println("Signature of data:- ", signature)

	// saving data in dataOfUser.qf file
	tempData, err6 := json.Marshal(writtable_data)
	if err6 != nil {
		fmt.Println("error in marshal of writtable_data:-------------", err6.Error())
	}
	err7 := ioutil.WriteFile("dataOfUser.qf", tempData, 0644)
	if err7 != nil {
		fmt.Println("error to write writtable_data in dataOfUser.json file:------------", err7.Error())
	}

	// creating google bucket path and data in bucket
	client, error := storage.NewClient(ctx)
	if error != nil {
		fmt.Println("error in creation of client: ", error.Error())
	}

	fileNameSavedInBucket := phone + ".qf"
	wc := client.Bucket("filestorage1_bucket").Object(fileNameSavedInBucket).NewWriter(ctx)
	wc.ContentType = "application/json"

	if _, err := wc.Write(tempData); err != nil {
		fmt.Println("error in creation of client2: ", err.Error())
	}
	if err := wc.Close(); err != nil {
		fmt.Println("error in creation of client4: ", err.Error())
	}
}
