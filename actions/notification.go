package actions

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
	"text/template"
	"time"

	libs "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/k8s"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"gopkg.in/gomail.v2"
	"k8s.io/client-go/pkg/api/v1"
)

var lastHipchat = time.Time{}
var lastHipchatMutex = &sync.Mutex{}

const messageTemplate = `

{{if .LastWarning}}
<img src="https://raw.githubusercontent.com/mpritter76/images/master/stop.png" style="width:20px;height:20px;">
{{else}}
<img src="https://raw.githubusercontent.com/mpritter76/images/master/warning.png" style="width:20px;height:20px;">
{{end}}

Violation in namespace <b>{{.Namespace}}</b> in <b>{{.Cluster}}</b>:
<p>
<ul>
<li>{{.EntityType}}: {{.EntitySource}}</li>
<li>Violation: {{.ViolationType}}</li>
<li>Source: {{.ViolationSource}}</li>
<li>Warning Count: {{.WarningCount}}</li>
</ul>

{{if .LastWarning}}
<b>This is the last warning before taking action!</b>
{{end}}
</p>
`

type actionMessage struct {
	Namespace       string
	Cluster         string
	EntityType      string
	EntitySource    string
	ViolationType   string
	ViolationSource string
	WarningCount    int
	LastWarning     bool
}

func NotifyOfViolation(actionMessage actionMessage) {

	tmpl, err := template.New("actionMessage").Parse(messageTemplate)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, actionMessage)
	if err != nil {
		panic(err)
	}

	clientset, err := k8s.LoadClientset()
	if err != nil {
		panic(err)
	}

	ns, err := clientset.CoreV1().Namespaces().Get(actionMessage.Namespace)
	if err != nil {
		panic(err)
	}

	go notifyHipChat(tpl.String(), ns, actionMessage.LastWarning)
	notifyEmail(tpl.String(), ns, actionMessage.LastWarning)

}

func notifyHipChat(message string, namespace *v1.Namespace, lastWarning bool) {

	lastHipchatMutex.Lock()
	canChat := time.Now().Sub(lastHipchat) < libs.Cfg.DurationBetweenChatNotifications
	if canChat {
		time.Sleep(libs.Cfg.DurationBetweenChatNotifications)
	}
	lastHipchat = time.Now()
	lastHipchatMutex.Unlock()

	if len(libs.Cfg.HipchatBaseURL) == 0 && len(libs.Cfg.HipchatRoomID) == 0 {
		libs.Log.Debug("Skipping Hipchat due to empty hipchat base url or empty hipichat room id")
		return
	}

	c := hipchat.NewClient(libs.Cfg.HipchatToken)
	hipchatUrl, err := url.Parse("https://tgtbullseye.hipchat.com/v2/")
	c.BaseURL = hipchatUrl

	color := hipchat.ColorYellow

	if lastWarning {
		color = hipchat.ColorRed
	}

	if teamHipchatIds, ok := namespace.Annotations["team/hipchat-ids"]; ok && libs.Cfg.HipchatTagNamespaceOwner {
		tags := []string{}
		for _, hipchatId := range strings.Split(teamHipchatIds, ",") {
			hId := strings.TrimSpace(hipchatId)
			strings.Replace(hId, "@", "", 1)
			tags = append(tags, "@"+hId)
		}

		notifRq := &hipchat.NotificationRequest{Message: strings.Join(tags, " ") + " (downvote)", MessageFormat: "text", Color: color}
		resp, err := c.Room.Notification(libs.Cfg.HipchatRoomID, notifRq)
		if err != nil {
			if resp != nil {
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				bodyString := string(bodyBytes)
				libs.Log.Error(bodyString)
			}
			libs.Log.Error(err)
		}
	}

	notifRq := &hipchat.NotificationRequest{Message: message, MessageFormat: "html", Color: color}
	resp, err := c.Room.Notification(libs.Cfg.HipchatRoomID, notifRq)
	if err != nil {
		if resp != nil {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			libs.Log.Error(bodyString)
		}
		libs.Log.Error(err)
	}

}

func notifyEmail(actionMessage string, namespace *v1.Namespace, lastWarning bool) {

	if len(libs.Cfg.SmtpServer) == 0 {
		libs.Log.Debug("Skipping emailing due to empty smtp server config")
		return
	}

	teamEmails := []string{}

	if libs.Cfg.SmtpSendToNamespaceOwner == false {
		if len(libs.Cfg.SmtpFallbackSendTo) == 0 {
			return
		}
		libs.Log.Warn("Not emailing namespace owner, sending to fallback instead.")
		teamEmails = strings.Split(libs.Cfg.SmtpFallbackSendTo, ",")
	} else if teamEmailString, ok := namespace.Annotations["team/email-ids"]; ok {
		teamEmails = strings.Split(teamEmailString, ",")
	} else {
		if len(libs.Cfg.SmtpFallbackSendTo) == 0 {
			return
		}
		libs.Log.Warn("The namespace " + namespace.Name + " does not have a team/email-ids annotation using fallback!")
		teamEmails = strings.Split(libs.Cfg.SmtpFallbackSendTo, ",")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", libs.Cfg.SmtpSendFrom)
	m.SetHeader("To", teamEmails...)
	if lastWarning {
		m.SetHeader("Subject", "LAST WARNING: Kubernetes Violation!")
	} else {
		m.SetHeader("Subject", "Kubernetes Violation!")
	}
	m.SetBody("text/html", actionMessage)
	d := gomail.NewDialer(libs.Cfg.SmtpServer, libs.Cfg.SmtpPort, libs.Cfg.SmtpUsername, libs.Cfg.SmtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)

	}
}
