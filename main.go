package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	DefaultBudgetInsightURITemplateWebview   = "https://{{.domain}}.biapi.pro/2.0/auth/webview/connect?client_id={{.clientID}}&redirect_uri={{.yourCallbackURI}}"
	DefaultBudgetInsightURITemplateAuthToken = "https://{{.domain}}.biapi.pro/2.0/auth/token/access"
)

type AppConfiguration struct {
	budgetInsight struct {
		URITemplate struct {
			webview   string
			authToken string
		}
		domain          string
		clientID        string
		clientSecret    string
		yourCallbackURI string
	}
}

type BudgetInsightAuth struct {
	code         string
	result       string
	connectionID string
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
		"clientID":        appConfig.budgetInsight.clientID,
		"yourCallbackURI": appConfig.budgetInsight.yourCallbackURI,
	}

	var tplURI bytes.Buffer
	if err := tmpl.Execute(&tplURI, parms); err != nil {
		panic(err)
	}

	log.Printf("Redirect to the Budget Insight Webview: %s", tplURI.String())
	http.Redirect(w, r, tplURI.String(), http.StatusTemporaryRedirect)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<a href=/webview>webview</a>")
}

func webviewCallback(w http.ResponseWriter, r *http.Request) {
	log.Print("Callback URI trigger")

	queryValues := r.URL.Query()

	if err := queryValues.Get("error"); err != "" {
		log.Printf("Callback error: %s", err)
		fmt.Fprintf(w, "Callback error: %s", err)

		return
	}

	userAuth := BudgetInsightAuth{
		code:         queryValues.Get("code"),
		connectionID: queryValues.Get("connection_id"),
		result:       "none",
	}

	getAuthToken(&userAuth)

	log.Printf("Result Auth (client secret):\n %s", userAuth.result)
	fmt.Fprintf(w, "Result Auth (client secret): %s", userAuth.result)
}

func getAuthToken(userAuth *BudgetInsightAuth) {
	tmpl, err := template.New("webview").Parse(appConfig.budgetInsight.URITemplate.authToken)
	if err != nil {
		panic(err)
	}

	parms := map[string]string{
		"domain": appConfig.budgetInsight.domain,
	}

	var tplURI bytes.Buffer
	if err := tmpl.Execute(&tplURI, parms); err != nil {
		panic(err)
	}

	payload := map[string]interface{}{
		"code":          userAuth.code,
		"client_id":     appConfig.budgetInsight.clientID,
		"client_secret": appConfig.budgetInsight.clientSecret,
	}

	byts, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodPost, tplURI.String(), bytes.NewBuffer(byts))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	log.Printf("Send POST to the Budget Insight API: %s", tplURI.String())
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		log.Printf("Error API token: %s", string(body))
	}

	userAuth.result = string(body)
}

func appInit() {
	log.SetFlags(2 | 3)
	log.Print("Starting ...")

	flag.StringVar(&appConfig.budgetInsight.URITemplate.webview, "URITemplate.webview", DefaultBudgetInsightURITemplateWebview, "Change default budgetInsight.URITemplate.webview")
	flag.StringVar(&appConfig.budgetInsight.URITemplate.authToken, "URITemplate.authToken", DefaultBudgetInsightURITemplateAuthToken, "Change default budgetInsight.URITemplate.authToken")

	flag.StringVar(&appConfig.budgetInsight.domain, "domain", "none", "Domain for Budget Insight")
	flag.StringVar(&appConfig.budgetInsight.clientID, "clientid", "none", "ClientID for App Budget Insight")
	flag.StringVar(&appConfig.budgetInsight.clientSecret, "clientsecret", "none", "ClientSecret for App Budget Insight")
	flag.StringVar(&appConfig.budgetInsight.yourCallbackURI, "yourcallbackuri", "none", "CallbackUri call after the webview of Budget Insight")
	flag.Parse()

	log.Printf("%+v\n", appConfig)

	if appConfig.budgetInsight.domain == "none" || appConfig.budgetInsight.clientID == "none" || appConfig.budgetInsight.clientSecret == "none" || appConfig.budgetInsight.yourCallbackURI == "none" {
		log.Print("Mandatory flags is not set (domain, clientId, yourCallbackURI, clientSecret)")
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

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
