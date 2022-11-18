# Teletekstbot
There are actually two bots: a 101 bot and a reply bot. The 101 bot fetches pages 104 til 150, checks for each page if it is updated and then posts a toot with a screenshot of the page. The reply bot checks if the bot is mentioned with a request like 'Pagina 101' and then replies with a screenshot of that page.

## Prerequisites
* Install the latest version of Go
* Install Docker
* Run `docker pull leonjza/gowitness:latest`, we need it to take screenshots from web pages
* Create a `.env` file from `.env_sample` using your Mastodon credentials. Go to `https://<mastodon_instance>/settings/applications` and create a new application to get those credentials.

## Run
The 101 bot: `go run cmd/101/bot_101.go`  
The reply bot: `go run cmd/reply/bot_reply.go`
