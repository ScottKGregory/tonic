
# 🍸 Tonic

[![Build Status](https://travis-ci.com/ScottKGregory/tonic.svg?branch=main)](https://travis-ci.com/ScottKGregory/tonic)
[![Go Report Card](https://goreportcard.com/badge/github.com/ScottKGregory/tonic)](https://goreportcard.com/report/github.com/ScottKGregory/tonic)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/scottkgregory/tonic)
![GitHub](https://img.shields.io/github/license/scottkgregory/tonic)

## What is Tonic?
Tonic is a set of helpers and pre-built endpoints helpful for setting up a Gin based website/API.
It's primarily set up to enable quick setup of a web API with cookie and bearer token auth using an OIDC provider.
It may become more configurable but for now it's just something to speed up my creation of new sites!

Technologies in use:
  - [Gin](https://github.com/gin-gonic/gin)
  - [Zerolog](https://github.com/rs/zerolog): Each request will add a logger to the context which contains pre-populated values
  - [Mabma](https://github.com/scottkgregory/mamba): The options are all annotated ready to be supplied directly to mamba

## Quick start

1. Install Tonic
```
go get github.com/scottkgregory/mamba
```

2. Call `tonic.Init`

```go
router, authedRouter, logger, err := tonic.Init(cfg.Tonic)
if err != nil {
  panic(err)
}
```

3. Use the Gin engine and router groups as normal
```go
router.GET("/tonic", func(c *gin.Context) {
  // Calling GetLogger with context will populate the logger with request specific values
  logger := helpers.GetLogger(c)
  logger.Info().Msg("Cheers!")
  c.JSON(http.StatusOK, gin.H{"message": "Hello!"})
})
```
4. Start the listener
```go
// Get a pre-configured logger
logger := helpers.GetLogger()
logger.Info().Msg("Starting listener")
r.Run(fmt.Sprintf(":%d", cfg.Port))
```

5. Visit the site, the homepage should show a Tonic default with a log in button. Logging in using your configured provider
should insert the user details in to the provided backend and present you with a token.

## What's not here?
There are a few things that aren't currently set up how I'd like and may change going forward:
  1. More configurable backends, currently defaults to storing user details in mongodb
  2. Configuration (or optional removal) of default pages