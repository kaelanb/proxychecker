package main

import (
    "io/ioutil"
    "log"
    "fmt"
    "os"
    "bufio"
    "strings"
    "net/http"
    "net/url"
    "time"
)

type Data struct{
	origin string
}

func main() {

    //creating the proxy URL list
    list, err := readLines("proxylist.txt")
    if err != nil {
    	log.Fatal(err)
    }
	
	for _, proxy := range list {
	
    proxyURL, err := url.Parse(strings.TrimSpace("http://" + proxy))
    if err != nil {
        log.Fatal(err)
    }

    //creating the URL to be loaded through the proxy
    urlStr := "http://httpbin.org/ip"
    url, err := url.Parse(urlStr)
    if err != nil {
        log.Fatal(err)
    }

    //adding the proxy settings to the Transport object
    transport := &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    }

    //adding the Transport object to the http Client
    client := &http.Client{
        Transport: transport,
        Timeout: 20 * time.Second,
    }

    //generating the HTTP GET request
    request, err := http.NewRequest("GET", url.String(), nil)
    if err != nil {
        log.Println(err)
    }
	
	fmt.Println("Sending request through " + proxyURL.String())	
    //calling the response
    response, err := client.Do(request)
    if err != nil {
        log.Fatal(err)
    }

    //getting the response content
    data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println(err)
    }
	defer response.Body.Close()

    if response.StatusCode != 200 {
    	fmt.Println("Result: " + string(data) + " - DOWN")
    }else {
    	fmt.Println("Result: " + string(data) + " - UP")
    }
    }
}

func readLines(path string) ([]string, error) {
        file, err := os.Open(path)
        if err != nil {
                return nil, err
        }
        defer file.Close()

        var lines []string
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                lines = append(lines, scanner.Text())
        }
        return lines, scanner.Err()
}
