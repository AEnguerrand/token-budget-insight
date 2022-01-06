package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	DefaultBudgetInsightURITemplateWebview   = "https://{{.domain}}.biapi.pro/2.0/auth/webview/connect?client_id={{.clientId}}&redirect_uri={{.yourCallbackUri}}"
	DefaultBudgetInsightURITemplateAuthToken = "https://{domain}.biapi.pro/2.0/auth/token/access"
)

type AppConfiguration struct {
	budgetInsight struct {
		URITemplate struct {
			webview   string
			authToken string
		}
		domain          string
		clientId        string
		yourCallbackUri string
	}
}

var appConfig AppConfiguration

func webviewRedirect(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start: Redirect to the Budget Insight Webview")

	tmpl, err := template.New("webview").Parse(appConfig.budgetInsight.URITemplate.webview)
	if err != nil {
		panic(err)
	}

	parms := map[string]string{
		"domain":          appConfig.budgetInsight.domain,
		"clientId":        appConfig.budgetInsight.clientId,
		"yourCallbackUri": appConfig.budgetInsight.yourCallbackUri,
	}

	var tplredirectURI bytes.Buffer
	if err := tmpl.Execute(&tplredirectURI, parms); err != nil {
		panic(err)
	}

	log.Printf("Redirect to the Budget Insight Webview: %s", tplredirectURI.String())
	http.Redirect(w, r, tplredirectURI.String(), http.StatusTemporaryRedirect)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<a href=/webview>webview</a>")
}

func webviewCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Callback URI")
	log.Printf("Full URI: %s\n", r.RequestURI)
}

func appInit() {
	log.SetFlags(2 | 3)
	log.Print("Starting ...")

	flag.StringVar(&appConfig.budgetInsight.URITemplate.webview, "URITemplate.webview", DefaultBudgetInsightURITemplateWebview, "Change default budgetInsight.URITemplate.webview")
	flag.StringVar(&appConfig.budgetInsight.URITemplate.authToken, "URITemplate.authToken", DefaultBudgetInsightURITemplateAuthToken, "Change default budgetInsight.URITemplate.authToken")

	flag.StringVar(&appConfig.budgetInsight.domain, "domain", "none", "Domain for Budget Insight")
	flag.StringVar(&appConfig.budgetInsight.clientId, "clientid", "none", "ClientID for Budget Insight")
	flag.StringVar(&appConfig.budgetInsight.yourCallbackUri, "yourcallbackuri", "none", "CallbackUri call after the webview of Budget Insight")
	flag.Parse()

	log.Printf("%+v\n", appConfig)

	if appConfig.budgetInsight.domain == "none" || appConfig.budgetInsight.clientId == "none" || appConfig.budgetInsight.yourCallbackUri == "none" {
		log.Print("Mandatory flags is not set (domain, clientId, yourCallbackUri)")
		log.Print("Usage: ")
		flag.PrintDefaults()
		os.Exit(-1)
	}
}

func main() {

	appInit()

	http.HandleFunc("/", homepage)
	http.HandleFunc("/webview", webviewRedirect)
	http.HandleFunc("/callback", webviewCallback)

	http.ListenAndServe(":8080", nil)
}
