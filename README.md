# SRL's Mystery Tournament Title Card

# Quick guide

Install Golang and then run:

```
go get github.com/SirGFM/MTTitleCard
go get golang.org/x/net/html
go get -u google.golang.org/api/sheets/v4
go get -u golang.org/x/oauth2/google
cd ${GOPATH}/src/github.com/SirGFM/MTTitleCard
go build .
```

To cross-compile from linux, run:

```
GOOS=windows go build .
```

# Using Google Sheets API

Go to https://developers.google.com/sheets/api/quickstart/go and click `ENABLE THE GOOGLE SHEETS API`
