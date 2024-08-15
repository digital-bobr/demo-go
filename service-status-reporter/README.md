# SERVICE-STATUS-REPORTER -- DEMO PROJECT
## add/edit/obtain health check and status info on your microservices

## Contribution
Everyone is welcome to contribute and collaborate with other team members on:
- adding new features
- adding tests
- improving overall performance
- create Jira requests to enrich tool's functionality
- fix issues and bugs

We ask you to be respectful:
- Review PRs carefully
- Follow GoLang best practices
- Add your commitments on time
- Test changes prior to creating PRs

## Prerequisites:
This repository contains source code of the tool and an existing database nested under:
'/service-status-reporter/database/' folder. No executables.

- install GO: https://golang.org/dl/ | sudo apt install golang-go
- install Sqlite: sudo apt install golang-go | sudo yum install sqlite
- clone this project repo
- navigate to the project's repo
- install dependencies: `go mod tidy`
- build the service: `go build -o reporter`
- run the application: `./reporter`
