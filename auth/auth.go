package auth

import "github.com/gin-gonic/gin"
import "fmt"
import "io/ioutil"
import "net/http"
import "bytes"
import "encoding/xml"
import "text/template"

type AuthData struct {
	SourceOffice string
	Originator   string
	PassLength   string
	Pass         string
}

type authEnvelope struct {
	CreateHeader createHeader `xml:"Header"`
	CreateBody   createBody   `xml:"Body"`
}
type createHeader struct {
	Createsession createSession `xml:"Session"`
}
type createSession struct {
	SessionId      string `xml:"SessionId"`
	SequenceNumber string `xml:"SequenceNumber"`
	SecurityToken  string `xml:"SecurityToken"`
}
type createBody struct {
	Security_AuthenticateReply security_AuthenticateReply `xml:"Security_AuthenticateReply"`
}
type security_AuthenticateReply struct {
	ProcessStatus processStatus `xml:"processStatus"`
	ErrorSection  errorSection  `xml:"errorSection"`
}
type processStatus struct {
	StatusCode string `xml:"statusCode"`
}
type errorSection struct {
	ApplicationError applicationError `xml:"applicationError"`
}
type applicationError struct {
	ErrorDetails errorDetails `xml:"errorDetails"`
}
type errorDetails struct {
	ErrorCode      string `xml:"errorCode"`
	ErrorCategory  string `xml:"errorCategory"`
	ErrorCodeOwner string `xml:"errorCodeOwner"`
}

var fileloc string

func Run(c *gin.Context) {
	var doc bytes.Buffer
	var err error

	a := AuthData{}
	a.SourceOffice = "MDEVU28AA"
	a.Originator = "WSAVRABE"
	a.PassLength = "8"
	a.Pass = "WVZ5d2VRN1k="

	t := template.New("AuthData template")
	t, err = template.ParseFiles("xmls/auth1.xml")
	if err != nil {
		fmt.Println("error leyendo xml template: " + err.Error())
		return
	}
	err = t.Execute(&doc, a)
	if err != nil {
		fmt.Println("error haciendo merge: " + err.Error())
		return
	}

	url := "https://nodeD1.test.webservices.amadeus.com/"
	fmt.Println("URL :> ", url)
	s := doc.Bytes()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(s))
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", "http://webservices.amadeus.com/1ASIWABEAVR/VLSSLQ_06_1_1A")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error al hacer post: " + err.Error())
		resp.Body.Close()
		return
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var authEnv authEnvelope
	err = xml.Unmarshal(body, &authEnv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c.String(200, string(body)+"\n\n")
	c.String(200, "status code: "+authEnv.CreateBody.Security_AuthenticateReply.ProcessStatus.StatusCode+"\n")
	c.String(200, "token: "+authEnv.CreateHeader.Createsession.SecurityToken+"\n")
	c.String(200, "sequence: "+authEnv.CreateHeader.Createsession.SequenceNumber+"\n")
	c.String(200, "sessionId: "+authEnv.CreateHeader.Createsession.SessionId+"\n")
	c.String(200, "error code: "+authEnv.CreateBody.Security_AuthenticateReply.ErrorSection.ApplicationError.ErrorDetails.ErrorCode+"\n")
	c.String(200, "error cat: "+authEnv.CreateBody.Security_AuthenticateReply.ErrorSection.ApplicationError.ErrorDetails.ErrorCategory+"\n")
	c.String(200, "error own: "+authEnv.CreateBody.Security_AuthenticateReply.ErrorSection.ApplicationError.ErrorDetails.ErrorCodeOwner+"\n")
}
