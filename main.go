package main

import (
	"log"
	"net/http"
	"net/mail"
	"strings"
	"text/template"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/uploadEmail", uploadEmail)
	log.Fatal(http.ListenAndServe(":18080", nil))
}

func uploadEmail(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Must be 'POST' request", 400)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Maximum file size exceeded", 400)
		return
	}

	file, _, err := r.FormFile("filename")
	if err != nil {
		http.Error(w, "Error parsing form", 500)
		return
	}

	defer file.Close()

	mail, err := mail.ReadMessage(file)
	if err != nil {
		http.Error(w, "File type not allowed", 400)
		return
	}

	var infoToReturn emailContent

	infoToReturn.To = sanitizeEmail(mail.Header["To"][0])
	infoToReturn.From = sanitizeEmail(mail.Header["From"][0])
	infoToReturn.Date = mail.Header["Date"][0]
	infoToReturn.Subject = mail.Header["Subject"][0]
	infoToReturn.MessageID = mail.Header["Message-Id"][0]

	w.Header().Add("Content-Type", "text/html")
	templates := template.New("template")
	templates.New("doc").Parse(doc)
	templates.Lookup("doc").Execute(w, infoToReturn)

}

type emailContent struct {
	To        string `json:"to"`
	From      string `json:"from"`
	Date      string `json:"date"`
	Subject   string `json:"subject"`
	MessageID string `json:"message-id"`
}

const doc = `
<!DOCTYPE html>
<html>
    <head>
		<h2>Upload successful</h2>
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js"></script>
    </head>
	<body>
		<p>To - {{.To}}</p>
		<p>From - {{.From}}</p>
		<p>Date - {{.Date}}</p>
		<p>Subject - {{.Subject}}</p>
		<br/>
		<p>Upload another file?</p>
		<form enctype="multipart/form-data" action="http://localhost:18080/uploadEmail" method="POST">
		<input type="file" name="filename">
		<input type="submit" value="upload" disabled>
	</form>
	<script src="emailUpload.js"></script>
    </body>
</html>
`

func sanitizeEmail(candidate string) string {
	candidate = strings.ReplaceAll(candidate, "<", "")
	return strings.ReplaceAll(candidate, ">", "")
}
