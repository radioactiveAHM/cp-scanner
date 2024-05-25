package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	probing "github.com/prometheus-community/pro-bing"
)

type Conf struct {
	Hostname   string   `json:"Hostname"`
	SNI        string   `json:"SNI"`
	MaxPing    int      `json:"MaxPing"`
	Goroutines int      `json:"Goroutines"`
	Scans      int      `json:"Scans"`
	Maxletency int64    `json:"Maxletency"`
	Scheme     string   `json:"Scheme"`
	Alpn       []string `json:"Alpn"`
}

func main() {
	// load config file
	cfile, cfile_err := os.ReadFile("conf.json")
	if cfile_err != nil {
		log.Fatalln(cfile_err.Error())
	}

	conf := Conf{}
	conf_err := json.Unmarshal(cfile, &conf)
	if conf_err != nil {
		log.Fatalln(conf_err.Error())
	}

	log.Println("start of app")
	// input
	hostname := conf.Hostname
	sni := conf.SNI
	maxping := conf.MaxPing
	goroutines := conf.Goroutines
	scans := conf.Scans
	var maxletency int64 = conf.Maxletency
	scheme := conf.Scheme
	alpn := conf.Alpn

	ch := make(chan string)
	for range goroutines {
		go func() {
			for range scans {
				// pick an ip
				file, _ := os.ReadFile("ipv4.txt")
				ranges := strings.Split(string(file), "\r\n")
				n4 := strconv.Itoa(rand.Intn(255))
				selected := ranges[rand.Intn(len(ranges))]
				ip := selected + n4

				// ping ip
				pinger, ping_err := probing.NewPinger(ip)
				pinger.SetPrivileged(true)
				pinger.Timeout = time.Duration(maxping) * time.Millisecond
				if ping_err != nil {
					log.Println(ping_err.Error())
					continue
				}
				pinger.Count = 1
				pinging_err := pinger.Run()
				if pinging_err != nil {
					log.Println(pinging_err.Error())
					continue
				}

				fmt.Printf("%s\t%s\n", ip, pinger.Statistics().MinRtt)

				if pinger.Statistics().PacketLoss > 0 || pinger.Statistics().MinRtt > (time.Duration(maxping)*time.Millisecond) {
					continue
				}

				// generate http req
				req := http.Request{Method: "GET", URL: &url.URL{Scheme: scheme, Host: ip, Path: "/"}, Host: hostname}
				req.Header = map[string][]string{
					"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:125.0)"},
					"Accept":          {"*/*"},
					"Accept-Language": {"en-US,en;q=0.5"},
					"Accept-Encoding": {"gzip", "deflate", "br", "zstd"},
				}

				var client *http.Client

				if conf.Scheme == "https" {
					tr := &http.Transport{
						TLSClientConfig: &tls.Config{ServerName: sni, NextProtos: alpn, MinVersion: tls.VersionTLS13},
						WriteBufferSize: 16384,
						ReadBufferSize:  32768,
					}
					client = &http.Client{Transport: tr}
				} else {
					client = http.DefaultClient
				}

				client.Timeout = time.Millisecond * time.Duration(maxletency)
				s := time.Now()
				// send request
				respone, http_err := client.Do(&req)
				e := time.Now()
				latency := e.UnixMilli() - s.UnixMilli()
				if http_err != nil {
					color.Red("%s", http_err.Error())
					continue
				}
				if respone.Header.Get("Server") != "cloudflare" {
					continue
				}

				println(respone.StatusCode)
				if respone.StatusCode == 200 {
					rep := fmt.Sprintf("%s\t%s\t%d\n", ip, pinger.Statistics().MinRtt, latency)
					color.Cyan("%s", rep)
					ch <- rep
				}
			}
			ch <- "end"
		}()
	}

	file, _ := os.OpenFile("result.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	deadgoroutines := 0
	for {
		if deadgoroutines == goroutines {
			break
		}
		v, ok := <-ch
		if !ok {
			break
		}
		if v == "end" {
			deadgoroutines += 1
			color.Green("end of goroutine")
			continue
		}
		file.Write([]byte(v))
	}
}
