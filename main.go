package main

import (
	"archive/tar"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

//go:embed static/*
var content embed.FS

func runshell(cmd string) bool {
	if err := exec.Command("/bin/bash", "-c", cmd).Run(); err != nil {
		log.Printf("\"%s\" %s", cmd, err)
		return false
	}
	return true
}

type certServer struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
}

type Response struct {
	ServerCRT string
	ServerKEY string
	RootCRT   string
}

func readFile(filename string) []byte {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("read file error %s", err)
	}
	return body
}
func (cert certServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if err := json.NewDecoder(r.Body).Decode(&cert); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("请求参数错误"))
		log.Printf("解析请求参数错误 %s", err)
		return
	}

	log.Printf("make cert for \"%s\" \"%s\"", cert.Domain, cert.IP)
	var cmd = fmt.Sprintf("perl mkcert --root-crt=ca.crt --root-key=ca.key --domain='%s' --ip='%s'", cert.Domain, cert.IP)
	if runshell(cmd) {
		log.Printf("make cert success, download beginning")
		w.Header().Set("Content-Disposition", "attachment; filename=cert.tar")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		var writer = tar.NewWriter(w)
		var files = []struct {
			Name string
			Body []byte
		}{
			{"ca.crt", readFile("ca.crt")},
			{"server.crt", readFile("server.crt")},
			{"server.key", readFile("server.key")},
		}
		for _, file := range files {
			hdr := &tar.Header{
				Name: file.Name,
				Mode: 0600,
				Size: int64(len(file.Body)),
			}
			if err := writer.WriteHeader(hdr); err != nil {
				log.Fatal(err)
			}
			if _, err := writer.Write(file.Body); err != nil {
				log.Fatal(err)
			}
		}
		writer.Close()
	} else {
		w.WriteHeader(500)
		log.Printf("make cert for \"%s\" \"%s\" failed", cert.Domain, cert.IP)
	}
}
func main() {
	var address string
	flag.StringVar(&address, "address", "127.0.0.1:8080", "server binding address")
	flag.Parse()
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	var s certServer
	http.Handle("/static/", http.StripPrefix("/", http.FileServer(http.FS(content))))
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || r.URL.Path == "/" {
			http.Redirect(rw, r, "/static/", 302)
		}
	})
	http.Handle("/cert", s)
	log.Printf("run server on %s", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err)
	}
}
