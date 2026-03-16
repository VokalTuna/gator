# Gator
## ABOUT THIS REPOSITORY
This is a repository for a cli application named "gator". "gator" application is a multi user tool to aggregate rss feeds and able to view the posts. The application is part of a guided project in relation to [boot.dev](https://www.boot.dev/) course for building a blog aggregator using Go language.

## Installation
This program requires installation of Go and PostgreSQL for the basic minimum in regards to running the application.
Ensure that it is the lates available go version.

### Installing the Go application
To install the application, run the following:
````
go install gator
````

### Making the config file
We access a config file in the home folder. Create a `.gatorconfig.json` file in the home directory with the following content:
```json
{
  "db_url": "postgres://<your username>:@<localhost name>:5432/database?sslmode=disable"
}
```

Remember to replace the contents including the greater than signs `<...>` with your database connection information.

## Usage

A couple of functions needs you to create a new user. Use the following commands for that:
```bash
gator register <name>
```
Add a feed to your user:
```bash
gator addfeed <name> <url>
```

Start the aggregator with the desired duration between each fetch:
```bash
gator agg 30s
```

View the posts available:

```bash
gator browse [limit]
```

Other available commands are:
- `gator login` - Log in as a user that already exists
- `gator users` - List all users
- `gator feeds` - List all feeds
- `gator follow` - Follow an already existing feed.
- `gator following` - List all feeds that you are following.
- `gator unfollowing <url>` - Unfollow feeds you are currently following.
- `gator reset` - Delete all contents from the database.

## Development
We use `sqlc` to generate the necessary go code for our application, and using `goose` to migrate our schemas.
