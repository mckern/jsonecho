package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"

	"github.com/mckern/pflag"
)

var versionNumber string
var whoami = path.Base(os.Args[0])
var version = fmt.Sprintf("%s %s", whoami, versionNumber)
var host string
var port string

func init() {
	var versionFlag bool
	var helpFlag bool
	var bind string

	pflag.StringVarP(&bind, "bind", "b", ":9090", "`address` and port to listen on")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "show this help")
	pflag.BoolVarP(&versionFlag, "version", "v", false, "print version number")

	pflag.Usage = func() {
		fmt.Fprintf(pflag.CommandLine.Output(), "%s: run a JSON prettifying service\n\n", whoami)
		fmt.Fprintf(pflag.CommandLine.Output(), "usage: %s [-hv] [-b address]\n", whoami)
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if helpFlag {
		pflag.CommandLine.SetOutput(os.Stdout)
		pflag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if isFlagPassed("bind") {
		host, port, _ = net.SplitHostPort(bind)
	} else {
		host, port, _ = net.SplitHostPort(getEnv("BIND", ":9090"))
	}
}

func isFlagPassed(name string) bool {
	found := false
	pflag.Visit(func(f *pflag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func prettyJSON(body io.Reader) []byte {
	var blob map[string]interface{}
	err := json.NewDecoder(body).Decode(&blob)

	if err != nil {
		log.Printf("cannot decode JSON '%+v'", body)
	}

	resp, _ := json.MarshalIndent(blob, "", "  ")
	return resp
}

func methodDispatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", r)

	if r.Method == "POST" {
		w.Write(prettyJSON(r.Body))
		return
	}

	if r.Method == "GET" {
		fmt.Fprint(w, "POST some JSON to this URL instead of GET")
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping server (pid %v)", sig, os.Getpid())
			os.Exit(1)
		}
	}()

	http.HandleFunc("/", methodDispatch)

	listen := fmt.Sprintf("%v:%v", host, port)
	log.Printf("starting server (pid %v) on %v", os.Getpid(), listen)

	log.Fatal(http.ListenAndServe(listen, nil))
}
