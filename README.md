<h1 align="center"><strong>VAT HMRC MTD Commandline Tool to submit VAT returns</strong></h1>

<br />

- The goal is create a simple CLI tool with Golang to submit VAT tax returns and get historical data. This tool is a bridging tool and does not store data on a server. It will probably never graduate to a fully fledged accounting softare.
- The project is open-source so feel free to use the code for your own projects and suggest improvements and raise issues.
- Please note that I have a busy job and this is a side project so I may not be on here for weeks or months at a time.
- Every application needs to be approved by HMRC, more below.
- I'm thinking about setting up a free website with React, Netlify and a server so that people can submit VAT tax returns online. But I haven't gotten round to it yet.

## Current status

**Testing stage** with a view to moving into production soon.

## Useful links

 - [VAT (MTD) API docs](https://developer.service.hmrc.gov.uk/api-documentation/docs/api/service/vat-api/1.0)
 - [HMRC VAT API GitHub project](https://github.com/hmrc/vat-api)
 - [Tutorials on HMRC Developer Hub](https://developer.service.hmrc.gov.uk/api-documentation/docs/tutorials)
 - [HMRC API Auth docs](https://developer.service.hmrc.gov.uk/api-documentation/docs/authorisation)
 - [HMRC Developer Hub Login](https://developer.service.hmrc.gov.uk/developer/login)
 - [HMRC API Nodejs Example Client GitHub](https://github.com/hmrc/api-example-nodejs-client)
 - [Making Tax Digital (MTD)-VAT End-to-End (E2E) Customer Journeys v5.0 PDF](https://developer.service.hmrc.gov.uk/api-documentation/assets/content/documentation/f66c79c2c4fc2f0cf27c158b2411a1b2-MTD-VAT%20End-to-End%20(E2E)%20Customer%20Journeys.pdf)

 - [Golang Oauth2 docs](https://godoc.org/golang.org/x/oauth2)

## Getting Started

**ASSUMPTION** I'm assume users will have decent computer knowledge and can install Golang and run a CLI tool.

1. Install Golang
2. Install `oauth2` via `go get golang.org/x/oauth2`
2. Install `uuid` via `go get github.com/google/uuid`
3. Register your app with HMRC (this is a time consuming step, see **Approval Process**)
4. Save your settings to the env.json file in the root directory
```json
{
  "API_URL": "https://test-api.service.hmrc.gov.uk",
  "CLIENT_ID": "YOUR-CLIENT-ID",
  "CLIENT_SECRET": "YOUR-CLIENT-SECRET",
  "SERVER_TOKEN": "YOUR-SERVER-TOKEN",
  "STATE_CHECK": "RANDOM-STRING-YOU-GENERATE-VIA-UUID-OR-MD5",
  "VRN": "YOUR-COMPANY-VRN-NUMBER",
}
```
5. You are ready to run the app using the CLI via `go run *.go help`
6. Run `go run *.go auth` to launch a local server (at [http://localhost:3000/](http://localhost:3000/)) in order to login and get your token which is saved to the same directory as the code.
7. You can then start to retrieve and submit VAT returns.

## Approval Process

- Firstly create an account on the Developer Hub. This allows you to register a client for the test API. Save the config to your `env.json` and off you go, you can start testing the app. You'll need to create a test login in order to test the server.
- Once you have tested your app, you will need to apply for approval for a product app. HMRC sends you an email and you need to reply in order to setup a demo.
- Before setting up the demo, they ask you a bunch of standard questions:

1. What is your Company Name and Registration Number (if you have one)?
2. Please confirm your business address and telephone number.
3. Do you have a website, if so please provide the URL
4. Is your business on LinkedIn or Twitter, if they are please provide the links
5. Is it a full digital record keeping solution, is it file-only (bridging product) or is it both?
6. Are you targeting a particular customer demographic or business sector?
7. Will your product be used as an in-house solution only or do you plan to sell/licence the product?
8. Is it built to UK standards and conventions (e.g. UK dates, £ not $ etc.)?
9. Are you developing more than one product? If you are what are they named?
10. Is your product a white label? (In general, white label branding is a practice in which a product or service is produced by one company and then rebranded by another company to make it appear to be their own)
11. Is your product GDPR compliant?
12. Does your product meet the [Web Content Accessibility Guidelines](https://www.w3.org/TR/WCAG21/) (minimum level AA)?
13. Does your product capture and send the [fraud prevention headers](https://developer.service.hmrc.gov.uk/api-documentation/docs/fraud-prevention)? 
14. What VAT schemes does it support?
    - Cash accounting
    - Annual accounting
    - Flat rate
    - Retail     
    - Margin      
    - Exemption
 
I have not yet completed a demo so no insight in that at the moment.