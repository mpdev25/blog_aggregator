# Blog Aggregator

The blog aggregator is designed to let logged in users follow RSS blog feeds and store the data in PostgreSQL tables for later access.

## Requirements

In order to run the blog aggregator you will need the following programs installed.

Go (version 1.26 or higher)

PostgreSQL

The program uses goose migration to create the postgreSQL tables and **sqlc generate** to generate Go code from SQL.

## Installation and Setup

To install run:

go install ://github.com/mpdev25/blog_aggregator@latest


Manually create a config file in your home directory ~/.gatorconfig.json

It should contain:

{
  "db_url": "protocol://username:password@host:port/database?sslmode=disable",
  "current_user_name": "username_goes_here"
}

Run go build .

## Commands

The following CLI commands can be used, prefixed by gator :

    login <username> - Logs a user in.

	setdb - Sets the database url.

	register <username> - Registers a new user.

	reset - Deletes the user table.

	users - Lists all registered users.

	agg <time-between-feeds> - Starts collecting feeds at the specified interval, for example 1m30s for every 1 minute and 30 seconds.

	addfeed <name> <url> - Adds a new feed by url and gives it the name provided.

	feeds - Lists all avaialable feeds.

	following - Lists allthe feeds the current user is following.

	follow <url> - Sets the given url to be followed by the current user.

	unfollow <url> - Unfollow the feed with the given url for the current user.

	browse [limit] - Browse the latest posts with an optional limit. If no limit is specified, the default is 2.