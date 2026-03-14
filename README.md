# Votefin

This is a free and open source server that I made so my Jellyfin users could help me decide what movies I would work on getting on the server next!
It integrates with your jellyfin server and tmdb so that you can easily have users voting on movies that you're unsure which one you wanna work on next and let others make your decision for you.

## Installation

Currently the only method for installing is to build the server from source. Thankfully golang makes this process extremely simple.

First ensure you have golang 1.25.1 minimum installed on your system. If it isn't either update your tool chain or reference the [go docs](https://go.dev/doc/install) on how to get it installed on your computer.

Next ensure you have git installed. This will allow you to clone down the repo running by running `git clone https://github.com/jkboyo/votefin`.

From there you should just have to type in `go build .` inside of the cloned down directory.

## Setup

Once you have everything cloned down you need to fill out the .example.env as shown below.

```bash
TMDB_API_KEY="tmbdapikeyhere" # This code will be something you need to sign up for an account with TMDB.org for
TMDB_DATA="/path/to/tmdb/data/here"
DB_PATH=../../votefin.db
JELLYFIN_URL="yourdomainhere.com"
SRVER_ID="votefin-1"
VOTE_LIMIT="5" # limits how many votes each user can use.
POSTER_IM_WIDTH="154" # This can be either 92, 154, 185, 342, 500, 780, or original
SECURE_ONLY_COOKIE="true" # If you want this to be a local only server set to false. This enforces TLS which doesn't work with IP addresses
```

The main hurdle will be getting the api key from tmdb.

Simply setup an account [here](https://www.themoviedb.org/signup) and then request an api-key [here](https://www.themoviedb.org/settings/api).

The api is free for personal use so you should be able to get a key and use it here.


