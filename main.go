package main

import (
    "fmt"
    "github.com/gocolly/colly"
    "github.com/lalamove/konfig"
    "github.com/lalamove/konfig/loader/klfile"
    "github.com/lalamove/konfig/parser/kpjson"
    "log"
)

var configFiles = []klfile.File{
    {
        Path:   "./config.json",
        Parser: kpjson.Parser,
    },
}

func init() {
    konfig.Init(konfig.DefaultConfig())
    konfig.RegisterLoader(
        klfile.New(&klfile.Config{
            Files: configFiles,
        }),
    )
    konfig.Load()
}

func main() {
    loginURL := konfig.String("host") + "/auth/login"
    checkInURL := konfig.String("host") + "/user/checkin"

    loginCollector := colly.NewCollector(
        //colly.Async(true),
        colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:80.0) Gecko/20100101 Firefox/80.0"),
    )

    loginCollector.Limit(&colly.LimitRule{
        Parallelism: 5,
    })

    loginCollector.OnRequest(func(r *colly.Request) {
        log.Println("Visiting", r.URL)
    })

    loginCollector.OnError(func(_ *colly.Response, err error) {
        log.Println("Something went wrong:", err)
    })

    checkInCollector := loginCollector.Clone()


    checkInCollector.Limit(&colly.LimitRule{
        Parallelism: 5,
    })

    checkInCollector.OnRequest(func(r *colly.Request) {
        log.Println("Visiting", r.URL)
    })

    checkInCollector.OnError(func(_ *colly.Response, err error) {
        log.Println("Something went wrong:", err)
    })

    checkInCollector.OnResponse(func(response *colly.Response) {
        fmt.Println(string(response.Body))
    })

    loginCollector.Visit(loginURL)

    // attach callback after visit the URL
    loginCollector.OnResponse(func(response *colly.Response) {
        fmt.Println(string(response.Body))
        // check in
        checkInCollector.Post(checkInURL, nil)
    })

    // login
    loginData := map[string]string{}
    loginData["email"] = konfig.String("username")
    loginData["passwd"] = konfig.String("password")
    loginCollector.Post(loginURL, loginData)

    loginCollector.Wait()
}
