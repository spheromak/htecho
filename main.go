package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/kelseyhightower/envconfig"
)

type options struct {
	Port int    `short:"p" long:"port" description:"Port to listen on" default:"9999"`
	File string `short:"c" long:"config" description:"Configfile to load" default:".proxypass.json"`
	Bind string `short:"b" long:"bind" description:"Ip address to bind too" default:"127.0.0.1"`
}

type response map[string]interface{}

func (r response) String() (s string) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func main() {
	opts := options{}
	err := envconfig.Process("htecho", &opts)
	if err != nil {
		log.Fatalf("Error parsing ENV vars %s", err)
	}
	if _, err := flags.Parse(&opts); err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Println(err.Error())
			os.Exit(1)
		}
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response{
			"headers":     r.Header,
			"method":      r.Method,
			"host":        r.Host,
			"proto":       r.Proto,
			"request_uri": r.RequestURI,
			"remote_addr": r.RemoteAddr,
		})
	}

	http.HandleFunc("/", hf)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", opts.Bind, opts.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
