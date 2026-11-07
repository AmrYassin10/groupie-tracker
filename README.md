# Groupie Tracker

A web application written in Go that fetches data from an external REST API and displays music artists and their concerts, with dynamic search, filtering, and server-side rendering.

## Features

- Fetches artist, location, date, and concert-relation data from an external JSON REST API
- Renders HTML pages on the server using Go templates
- Dynamic search and filtering of artists
- Detail page for each artist <check: confirm you have one>

## Tech stack

- Go (standard library: `net/http`, `html/template`, `encoding/json`) <check: remove anything you didn't use>
- HTML, CSS
- External REST API (JSON)

## How it works

1. The Go server starts and registers its HTTP routes.
2. When a page is requested, the handler calls the external API and decodes the JSON response into Go structs.
3. The data is filtered or searched according to the request.
4. A Go template renders the result as HTML and sends it to the browser.

## Run it locally

Requirements: Go (see `go.mod` for the version)

```bash
git clone https://github.com/AmrYassin10/groupie-tracker.git
cd groupie-tracker
go run .
```

Then open http://localhost:<port> in your browser.

## Project structure

```
<paste the output of `ls` or `tree` here>
```

## What I learned

- Consuming a REST API and decoding JSON into structs
- Building an HTTP server with routing and server-side templates
- Implementing search and filtering logic
- Handling errors and invalid requests

## Author

Amr Yassin: [GitHub](https://github.com/AmrYassin10) | [LinkedIn](https://www.linkedin.com/in/amr-yassin-7a4441156)

Built as part of the Software Engineering programme at 01 School.
