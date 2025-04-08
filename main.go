package main

import (
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "os"
    "github.com/fatih/color"
    "time"
)

type IPAPI struct {
    Status     string  `json:"status"`
    Country    string  `json:"country"`
    CountryCode string `json:"countryCode"`
    Region     string  `json:"regionName"`
    RegionCode string  `json:"region"`
    City       string  `json:"city"`
    Zip        string  `json:"zip"`
    Lat        float64 `json:"lat"`
    Lon        float64 `json:"lon"`
    Timezone   string  `json:"timezone"`
    ISP        string  `json:"isp"`
    Org        string  `json:"org"`
    AS         string  `json:"as"`
    Mobile     bool    `json:"mobile"`
    Proxy      bool    `json:"proxy"`
    Hosting    bool    `json:"hosting"`
    Query      string  `json:"query"`
}

type IPWhois struct {
    Continent      string `json:"continent"`
    ContinentCode  string `json:"continent_code"`
    Currency       string `json:"currency"`
    ASN            string `json:"asn"`
    TimezoneGMT    string `json:"timezone_gmt"`
}

type IPInfo struct {
    Hostname string `json:"hostname"`
}

func banner() {
    red := color.New(color.FgRed).SprintFunc()
    white := color.New(color.FgWhite).SprintFunc()

    fmt.Println("\n" + red(`
    ________     ______                __
   /  _/ __ \   /_  __/________ ______/ /_____  _____
   / // /_/ /    / / / ___/ __ '/ ___/ //_/ _ \/ ___/
 _/ // ____/    / / / /  / /_/ / /__/ ,< /  __/ /
/___/_/        /_/ /_/   \__,_/\___/_/|_|\___/_/
`))
    fmt.Println(red("           [ ") + white("www.kirovgroup.org") + red(" ]\n"))
}



func resolveDomain(domain string) string {
    ip, err := net.LookupIP(domain)
    if err != nil || len(ip) == 0 {
        fmt.Println("❌ Не удалось резолвить домен!")
        os.Exit(1)
    }
    return ip[0].String()
}

func fetchJSON(url string, target interface{}) {
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        fmt.Println("❌ Ошибка при запросе:", err)
        os.Exit(1)
    }
    defer resp.Body.Close()
    json.NewDecoder(resp.Body).Decode(target)
}


func main() {
    banner()

    white := color.New(color.FgWhite).SprintFunc()
    green := color.New(color.FgGreen).SprintFunc()

    if len(os.Args) < 2 {
        fmt.Println("\n" + white("Использование: go run main.go <IP/Домен>") + "\n")
        os.Exit(1)
    }

    ip := os.Args[1]
    if net.ParseIP(ip) == nil {
        fmt.Println(white("🟡 Резолвинг домена: ") + green(ip))
        ip = resolveDomain(ip)
    }

    var ipapi IPAPI
    var ipwhois IPWhois
    var ipinfo IPInfo

    fetchJSON("http://ip-api.com/json/"+ip+"?fields=66846719", &ipapi)
    fetchJSON("https://ipwhois.app/json/"+ip, &ipwhois)
    fetchJSON("https://ipinfo.io/"+ip+"/json", &ipinfo)

    fmt.Println()
    fmt.Printf(white("IP               : ") + green("%s\n\n"), ip)

    fmt.Printf(white("Статус           : ") + green("%s\n"), ipapi.Status)
    fmt.Printf(white("Континент        : ") + green("%s (%s)\n"), ipwhois.Continent, ipwhois.ContinentCode)
    fmt.Printf(white("Страна           : ") + green("%s (%s)\n"), ipapi.Country, ipapi.CountryCode)
    fmt.Printf(white("Регион           : ") + green("%s (%s)\n"), ipapi.Region, ipapi.RegionCode)
    fmt.Printf(white("Город            : ") + green("%s\n"), ipapi.City)
    fmt.Printf(white("ZIP-Код          : ") + green("%s\n"), ipapi.Zip)
    fmt.Printf(white("Широта           : ") + green("%f\n"), ipapi.Lat)
    fmt.Printf(white("Долгота          : ") + green("%f\n"), ipapi.Lon)
    fmt.Printf(white("Часовой пояс     : ") + green("%s (GMT %s)\n"), ipapi.Timezone, ipwhois.TimezoneGMT)
    fmt.Printf(white("Валюта           : ") + green("%s\n"), ipwhois.Currency)
    fmt.Printf(white("ISP              : ") + green("%s\n"), ipapi.ISP)
    fmt.Printf(white("Организация      : ") + green("%s\n"), ipapi.Org)
    fmt.Printf(white("AS               : ") + green("%s\n"), ipapi.AS)
    fmt.Printf(white("ASN              : ") + green("%s\n"), ipwhois.ASN)
    fmt.Printf(white("Реверс DNS       : ") + green("%s\n"), ipinfo.Hostname)
    fmt.Printf(white("Мобильная сеть   : ") + green("%v\n"), ipapi.Mobile)
    fmt.Printf(white("Прокси           : ") + green("%v\n"), ipapi.Proxy)
    fmt.Printf(white("Хостинг          : ") + green("%v\n"), ipapi.Hosting)
    fmt.Printf(white("Запрос           : ") + green("%s\n"), ipapi.Query)
}


