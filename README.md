# haraka-sana

haraka sana is an opensource distribution system it just kicked of as of October 22, 2024
at 1915HRS

![project arhitecture](https://github.com/kimutaiRop/haraka-sana/blob/main/architecture.png)

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
