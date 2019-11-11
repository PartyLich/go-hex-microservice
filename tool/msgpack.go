package main

import (
	"bytes"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PartyLich/hex-microservice/shortUrl"
)

func httpPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	return fmt.Sprintf(":%s", port)
}

func main() {
	address := fmt.Sprintf("http://localhost%s", httpPort())
	redirect := shortUrl.Redirect{}
	redirect.URL = "https://google.com"

	body, err := msgpack.Marshal(&redirect)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.Post(address, "application/x-msgpack", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// deserialize and print
	msgpack.Unmarshal(body, &redirect)
	log.Printf("%v\n", redirect)
}
