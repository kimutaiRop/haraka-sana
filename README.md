# haraka-sana

haraka sana is an opensource distribution system it just kicked of as of October 22, 2024
at 1915HRS

![project Architecture](https://github.com/kimutaiRop/haraka-sana/blob/main/architecture.png)

## Dependencies and Stack

- go-lang
  - Gorm
  - Gin
- databases
  - posgresql
  - valkey

## RUNNING DEVELOPMENT SERVER

  1. start docker to run databases
    - `docker compose up`
  2. Start the web server
    - `go run .`

## Features

* [ ] Impliment Ouath Authorization & Authentication [Implimentation details](https://aaronparecki.com/oauth-2-simplified/)

Oauth creds
```json
{
    "clien_secret": "PISnXzaspCkJSSsXJscznbLJmTgHqHrsddCebTQsuSFMYzSvehjZDKbCqDmTItzh",
    "client_id": "WFTyrSNoxJfLdlrkrPsWgYCQs"
}
```
http://127.0.0.1:8080/api/v1/oauth2/authorize?redirect_uri=http://localhost:3000/haraka-sana&client_id=WFTyrSNoxJfLdlrkrPsWgYCQs&grant_type=code
