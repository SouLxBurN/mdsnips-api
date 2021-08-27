# SouLxSnips

Backend Restful API for storing and retrieving Markdown for markdown snippets.

## Overview

### Running the server
1. Copy of rename `soulxsnips.env` to `.env` and update the missing environment variables.
	- SOULXSNIPS_USER: Basic API Authentication user name.
	- SOULXSNIPS_PASS: Basic API Authentication password.
	- SOULXSNIPS_MONGO_CONN: MongoDB Connection String.
2. To run the server, simply execute one fo the following:
	```
	go run main.go
	```
	or
	```
	fresh
	```

### Regenerating Swagger Documentation

Run the following in the project root:
```
swag init
```

Swagger documentation is hosted at `/swagger/index.html`
