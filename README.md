# Votefin

This is a free and open source server that I made so my Jellyfin users could help me decide what movies I should put on my server next!

It integrates with your jellyfin server and tmdb so that you can easily have your users vote and make your decisions for you.

## Installation

Build the server from source. Thankfully golang makes this process extremely simple.

First ensure you have golang 1.25.1 minimum installed on your system. If it isn't either update your tool chain or reference the [go docs](https://go.dev/doc/install) on how to get it installed on your computer.

Next ensure you have git installed. This will allow you to clone down the repo by running `git clone https://github.com/jkboyo/votefin`.

To build the program run `go build .` inside of the cloned directory.

## Setup

### Environment

Once you have everything cloned down and the binary built you need to fill out the .example.env and save it as your .env.

```bash
TMDB_API_KEY="tmbdapikeyhere" # This code will be something you need to sign up for an account with TMDB.org for
JELLYFIN_URL="yourdomainhere.com"
SRVER_ID="votefin-1"
VOTE_LIMIT="5" # limits how many votes each user can use.
POSTER_IM_WIDTH="154" # This can be either 92, 154, 185, 342, 500, 780, or original
SECURE_ONLY_COOKIE="true" # If you want this to be a local only server set to false. This enforces TLS which doesn't work with IP addresses
```

The main hurdle will be getting the api key from tmdb.

Simply setup an account [here](https://www.themoviedb.org/signup) and then request an api-key [here](https://www.themoviedb.org/settings/api).

The api is free for personal use so you should be able to get a key and use it here.

### Database

The next step is to get the database running.

Votefin uses [goose](https://github.com/pressly/goose) to do it's db migrations. To get it correctly setup you will need

to install goose. The simplest way is by running `go install github.com/pressly/goose/v3/cmd/goose@latest` in the cli.

If you don't want to install it by that method there are instructions on how to get it appropriately setup for your system.

Once it is installed you simply have to run the db creation command `make dev-db-up` at the root of the project and the db will initialize.

## Usage

From there you should be able to login with your jellyfin admin login. As admin you can add any movies that you want people to vote on using the top search bar.

Once the movies are added any jellyfin user you have can login and vote for their picks! 

Once you have the movies added on the server you can go to votefin and mark on server and it will clear all the votes and remove the movie from voting.
