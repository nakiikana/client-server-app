package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"tools/models"
)

func main() {

	certData, _ := ioutil.ReadFile("./minica/minica.pem")
	cp, _ := x509.SystemCertPool() // current pool of certificate authorities that we trust
	cp.AppendCertsFromPEM(certData)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: cp,
			},
		},
	}
	kind := flag.String("k", "", "Parsed kind of the candy")
	count := flag.Int("c", -1, "Parsed candy count")
	cash := flag.Int("m", -1, "Parsed candy amount")
	flag.Parse()
	if *kind == "" || *count == -1 || *cash == -1 {
		fmt.Println("Usage: -k [the candy kind] -c [the amount of candy] -m [the amount of money]")
		return
	}
	order := models.Order{
		Money:      *cash,
		CandyType:  *kind,
		CandyCount: *count,
	}
	orderJson, _ := json.Marshal(&order)
	resp, err := client.Post("https://localhost:8200", "application/json", bytes.NewBuffer(orderJson))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
}
