# Gator
## ABOUT THIS REPOSITORY
This is a repository for a cli application named "gator". The application is part of a guided project in relation to [boot.dev](https://www.boot.dev/) course for building a blog aggregator using Go language.

## Requirements
This program requires installation of Go and PostgreSQL for the basic minimum in regards to running the application.

To install the application, run the following:
````
go install gator
````

## Development
We use `sqlc` to generate the necessary go code for our application, and using goose to migrate our schemas.
