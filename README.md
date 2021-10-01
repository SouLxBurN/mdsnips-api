# MDSnips

Backend Restful API for storing and retrieving Markdown for markdown snippets.

## Overview

### Running the server
1. Copy of rename `mdsnips.env` to `.env` and update the missing environment variables.
	- MDSNIPS_USER: Basic API Authentication user name.
	- MDSNIPS_PASS: Basic API Authentication password.
	- MDSNIPS_MONGO_CONN: MongoDB Connection String.
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
