package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
)

// SSL/TLS Email Example

func main() {
	key := []byte("passphrasewhichneedstobe32bytes!")
	ciphertext, err := ioutil.ReadFile("pass")

	if err != nil {
		fmt.Println(err)
	}

	cip, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(cip)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}

	myscanner := bufio.NewScanner(os.Stdin)

	fmt.Print("To: ")
	myscanner.Scan()
	to_address := myscanner.Text()

	fmt.Print("Subject: ")
	myscanner.Scan()
	subject_text := myscanner.Text()

	fmt.Print("Body:\n ")
	myscanner.Scan()
	body_text:= myscanner.Text()


	from := mail.Address{"", "test4lab@yandex.ru"}
	to   := mail.Address{"", to_address}
	subj := subject_text
	body := body_text

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k,v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.yandex.ru:465"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("","test4lab@yandex.ru", string(plaintext), host)

	// TLS config
	tlsconfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

}