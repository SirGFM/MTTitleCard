# SRL's Mystery Tournament Title Card

## Quick guide

Install Golang and then run:

```
go get github.com/SirGFM/MTTitleCard
go get github.com/pkg/errors
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

## Running the application

### Setting up

Before running the application for the first time, be sure to access
https://developers.google.com/sheets/api/quickstart/go and click `ENABLE THE GOOGLE SHEETS API`
to generate a project and a API(?) credential.

Place the generated `credentials.json` in the same directory as you will
execute the application from.

Before generating title cards, you'll need to generate an OAuth2 token.

Start the server:

```
./MTTitleCard
```

Then open a browser, access `http://localhost:8080` and follow the
instructions in the page.

### Accessing title cards

To access the title card for a given SRL user, start the server and use the
username as the page's path. For example:

```
./MTTitleCard
curl http://localhost:8080/GFM
```

### Configuring

You may customize the server by specifying a JSON file:

```
./MTTitleCard -config config.json
```

Sample JSON object:

```
{
    port: 8080,
    twitchClientID: "",
    credentialFile: "credentials.json",
    tokenFile: "token.json",
    mtCareerSpreasheet: "1DWYq3T1w8u1N0CWWJ72tqQRv67c1eY098u0wyuiMEmA",
    tourneyInfo: {
        sheetName: "STATS",
        firstColumn: "C",
        lastColumn: "F",
        firstRow: 16
    },
    userInfo: {
        sheetName: "MT Career",
        firstColumn: "A",
        lastColumn: "H",
        firstRow: 2
    },
    standingsInfo: {
        sheetName: "MT Career Standings",
        firstColumn: "A",
        lastColumn: "S",
        firstRow: 1
    },
    joinedMtIdx: 0,
    nameIdx: 1,
    torneyCountIdx: 2,
    winIdx: 3,
    loseIdx: 4,
    draftIdx: 7,
    cssFile: "style.css",
    templateFile: "template.html"
}
```

Explanation:

* port: Port for serving the webpages
* twitchClientID: Client ID used when retrieving values from twitch's API
* credentialFile: Path to a JSON file with the credentials to access Google's API
* tokenFile: Path to a JSON file with the token to access Google's API
* mtCareerSpreasheet: ID of the MT Career spreadsheet
* tourneyInfo: "Range" used to extract the number of entrants
* userInfo: "Range" used to extract the a tournament entrant
* standingsInfo: "Range" used to extract the entrants standings
* joinedMtIdx: Index, in the spreadsheet, of the number of tournaments entered by the user
* nameIdx: Index, in the spreadsheet, of the user's name
* torneyCountIdx: Index, in the spreadsheet, of the username
* winIdx: Index, in the spreadsheet, of the user's number of victories
* loseIdx: Index, in the spreadsheet, of the user's number of losses
* draftIdx: Index, in the spreadsheet, of the user's value(?) in draft points
* cssFile: Path to a CSS file used to override the default CSS
* templateFile: Path to a HTML-template file used to override the default page template

The range is an object with the following fields:

* sheetName: Name of the specific sheet within the spreadsheet
* firstColumn: First column to be downloaded
* lastColumn: Last column to be downloaded
* firstRow: First row to be downloaded

The number of rows is automatically calculated.
