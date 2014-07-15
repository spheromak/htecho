package main

import (
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"os"
)

type Options struct {
	Port int    `short:"p" long:"port" description:"Port to listen on" default:"9999"`
	File string `short:"c" long:"config" description:"Configfile to load" default:".proxypass.json"`
	Bind string `short:"b" long:"bind" description:"Ip address to bind too" default:"127.0.0.1"`
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func main() {
	opts := Options{}
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
		fmt.Fprint(w, Response{
			"headers":     r.Header,
			"method":      r.Method,
			"host":        r.Host,
			"proto":       r.Proto,
			"request_uri": r.RequestURI,
		})
	}

	http.HandleFunc("/", hf)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", opts.Bind, opts.Port), nil)
	if err != nil {
		log.Fatal(err)
	}

}
