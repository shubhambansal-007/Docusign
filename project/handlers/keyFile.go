package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

type KeyEC struct {
	publicKey  ed25519.PublicKey `json:"Public_key"`
	privateKey ed25519.PublicKey `json:"Private_key"`
}

var KEC KeyEC

func generateKeyEC() {
	//fb := "something"
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
	}

	KEC.privateKey = ed25519.PublicKey(private)
	KEC.publicKey = ed25519.PublicKey(public)
}

func createSignEC(data string) string {

	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	ciphr := ed25519.Sign(ed25519.PrivateKey(KEC.privateKey), msg)
	sign := hex.EncodeToString(ciphr)
	return sign
}

func verifySignatureEC(data string, signature string) bool {
	sign, err := hex.DecodeString(signature)
	if err != nil {
		fmt.Println(err)
	}
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	check := ed25519.Verify(KEC.publicKey, msg, sign)
	return check
}

func encryptAES(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key := []byte(keyString)
	plaintext := []byte(stringToEncrypt)

	fmt.Println("line 130:------------------ ")
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decryptAES(keyString string, encryptedString string) (decryptedString string) {

	key := []byte(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

//------------------
//---------------------------
//--------------RSA KEYS------generate key, encrypt, decrypt, create sign & verify sign
//---------------------------
//------------------

// type Key struct {
// 	publicKey  *rsa.PublicKey  `json:"Public_key"`
// 	privateKey *rsa.PrivateKey `json:"Private_key"`
// }

// var K Key
// func NewKey() {
// 	tempKey, err := rsa.GenerateKey(rand.Reader, 2048)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	K.publicKey = &tempKey.PublicKey
// 	K.privateKey = tempKey
// 	fmt.Println("Generated public and private key")
// }

// func encryptMessage(secretMessage string, key rsa.PublicKey) string {
// 	label := []byte("OAEP Encrypted")
// 	rng := rand.Reader
// 	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
// 	//ciphertext, err := ed25519.
// 	if err != nil {
// 		fmt.Println("error in encrypting :- ", err)
// 	}
// 	return base64.StdEncoding.EncodeToString(ciphertext)
// }

// func decryptMessage(cipherText string, privKey rsa.PrivateKey) string {

// 	ct, _ := base64.StdEncoding.DecodeString(cipherText)
// 	label := []byte("OAEP Encrypted")
// 	rng := rand.Reader
// 	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
// 	if err != nil {
// 		fmt.Println("error in decrypting :- ", err)
// 	}
// 	//fmt.Println("Decrypted Message is:- ", string(plaintext))
// 	return string(plaintext)
// }

// func createSignature(message string, privateKey *rsa.PrivateKey) []byte {

// 	msg := []byte(message)

// 	msgHash := sha256.New()
// 	_, err := msgHash.Write(msg)
// 	if err != nil {
// 		fmt.Println("error check createSignature line 42 :- ", err)
// 	}
// 	msgHashSum := msgHash.Sum(nil)

// 	signature, err2 := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
// 	if err2 != nil {
// 		fmt.Println("error check createSignature line 48 :- ", err2)
// 	}

// 	return signature
// }

// func checkSignature(message string, publicKey rsa.PublicKey, signature []byte) {
// 	msg := []byte(message)
// 	msgHash := sha256.New()
// 	_, err := msgHash.Write(msg)
// 	if err != nil {
// 		fmt.Println("error check createSignature line 46 :- ", err)
// 	}
// 	msgHashSum := msgHash.Sum(nil)

// 	err2 := rsa.VerifyPSS(&publicKey, crypto.SHA256, msgHashSum, signature, nil)
// 	if err2 != nil {
// 		fmt.Println("could not verify signature line 66:-   ", err2)
// 		return
// 	} else {
// 		fmt.Println("signature verified!!  /n")
// 	}
// }
