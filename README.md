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

- [x] Impliment Ouath Authorization & Authentication [Implimentation details](https://aaronparecki.com/oauth-2-simplified/)
- [x] Third patry integration and order creating
- [ ] Agent Integrations
- [ ] Alert Third patry Application progress
- [ ] Location Tracking System
- [ ] Local Agent Ditsribution Sytem (Same direction Same person)
- [ ] Admin/Staff managemnt system
- [ ] Agent Distibution App / web app

  Oauth creds

```json
{
  "clien_secret": "HVvNqUcEoCEMuJUGxKXBehzCqTvTgcXZfzAYjLIsTUqSmnGbZlojjJqhjJZUguEo",
  "client_id": "nukfagiKxpgaTGnCvZfJbaQNf"
}
```

http://127.0.0.1:8080/api/v1/oauth2/authorize?redirect_uri=http://localhost:3000/haraka-sana&client_id=nukfagiKxpgaTGnCvZfJbaQNf&grant_type=code
