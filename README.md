# token-budget-insight

This small tool aims to automate the creation of the long-term token for a Budget Insight application; the process is explain in the [Budget Insight API docs](https://docs.budget-insight.com/guides/add-first-user-connection). It's planned to be run locally and exposed with a tunnel system, like ngork.

## Usage

First, an account and application are required on Budget Insight. Also, your callback URI need to be defined inside your app client configuration. More information [here](https://docs.budget-insight.com/guides/quick-start#create-your-account).

```bash
go build
./token-budget-insight -domain <domainName> -clientid <clientID> -yourcallbackuri <hostname>/callback -clientsecret <clientSecret>
```

After opening your web browser on [localhost:8080](http://localhost:8080), click on the WebView link and select your account inside the Budget Insight WebView then the token is displayed inside your page or your terminal. This access token is a long term token, more detail [here](https://docs.budget-insight.com/reference/authentication#exchange-a-temporary-code-for-an-access-token).

### Example (with fake values)

```bash
go build
./token-budget-insight -domain test45-sandbox  -clientid 12548 -yourcallbackuri https://website.com/callback -clientsecret fsfefR4fRghop5
```