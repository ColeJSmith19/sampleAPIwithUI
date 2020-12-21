package main

import (
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"strings"
)

const tenMegabyteMax = 10 << 20

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/uploadEmail", uploadEmail)
	log.Fatal(http.ListenAndServe(":18080", nil))
}

func uploadEmail(w http.ResponseWriter, r *http.Request) {

	// Although my frontend exclusively uses POST, I figured this was still worth including
	// Since the endpoint would be exposed in the request, someone could use something like postman
	// to maliciously make requests of other types against this. This prevents those requests from getting too far.
	// If the '/uploadEmail' endpoint was intended to handle more than one type of request, I would have put this into
	// a switch statement that switched over the r.Method, where the default case would have been a message saying
	// that the request method was invalid. That said, an implementation using that methodology would be more loosely coupled.
	// If you ask me to add 'GET' functionality to this API for the '/uploadEmail' endpoint, I would need to do a lot of
	// restructuring. If I had implemented this to be more loosely coupled from the start than this would not be an issue.
	if r.Method != "POST" {
		http.Error(w, "Must be 'POST' request", 400)
		return
	}

	err := r.ParseMultipartForm(tenMegabyteMax)
	if err != nil {
		// This error is misleading. After some research I learned that this is the max file size stored in memory.
		// The rest is written to a temp directory. So a really big email would have the capacity to crash my API.
		// Before parsing the forms, I should have checked the file size.
		http.Error(w, "Maximum file size exceeded", 400)
		return
	}

	// Grab form by key
	file, _, err := r.FormFile("uploadedEmail")
	if err != nil {
		http.Error(w, "Error parsing request", 500)
		return
	}

	// Still close file if any of the following code fails. This is better than manually closing the form when I am done.
	defer file.Close()

	// Use go's net/mail package to parse the file
	mail, err := mail.ReadMessage(file)
	if err != nil {
		http.Error(w, "File type not allowed", 400)
		return
	}

	// var infoToReturn emailContent

	// infoToReturn.To = sanitizeEmail(mail.Header["To"][0])
	// infoToReturn.From = sanitizeEmail(mail.Header["From"][0])
	// infoToReturn.Date = mail.Header["Date"][0]
	// infoToReturn.Subject = mail.Header["Subject"][0]
	// infoToReturn.MessageID = mail.Header["Message-Id"][0]

	// Which index should I use in the case of multiple headers for a single key?
	// This change is more idiomatic
	infoToReturn := emailContent{
		To:        sanitizeEmail(mail.Header["To"][0]),
		From:      sanitizeEmail(mail.Header["From"][0]),
		Date:      mail.Header["Date"][0],
		Subject:   mail.Header["Subject"][0],
		MessageID: mail.Header["Message-Id"][0],
	}

	// Set response to type html
	w.Header().Add("Content-Type", "text/html")

	// templates := template.New("template")
	// templates.New("doc").Parse(doc)
	// templates.Lookup("doc").Execute(w, infoToReturn)

	// This change gets rid of the need to create and rename a template called 'template'
	templates, _ := template.New("doc").Parse(doc)
	// Fetch the template named "doc" and apply the parsed template to the variable 'infoToReturn'
	// This is where the data from the struct is dynamically being applied to the template.
	templates.Lookup("doc").Execute(w, infoToReturn)
}

type emailContent struct {
	To        string `json:"to"`
	From      string `json:"from"`
	Date      string `json:"date"`
	Subject   string `json:"subject"`
	MessageID string `json:"message-id"`
}

// https://adlerhsieh.com/blog/rendering-dynamic-data-in-go-http-template
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
		<input type="file" name="uploadedEmail">
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
