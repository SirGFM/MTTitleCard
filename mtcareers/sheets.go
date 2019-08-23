package mtcareers

import (
    "encoding/json"
    "fmt"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/sheets/v4"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
)

type Sheet struct {
    srv *sheets.Service
    id string
    TotalEntrants int
    LatestEntrants int
}

type Arg struct {
    ApiToken string
    // Path to JSON downloaded after enabling the API
    CredentialToken string
    SpreadsheetId string
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, arg *Arg) *http.Client {
    // The file token.json stores the user's access and refresh tokens, and is
    // created automatically when the authorization flow completes for the first
    // time.
    if arg.ApiToken == "" {
        arg.ApiToken = "token.json"
    }
    tok, err := tokenFromFile(arg.ApiToken)
    if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(arg.ApiToken, tok)
    }
    return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the "+
            "authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
        log.Fatalf("Unable to read authorization code: %v", err)
    }

    tok, err := config.Exchange(context.TODO(), authCode)
    if err != nil {
        log.Fatalf("Unable to retrieve token from web: %v", err)
    }
    return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
}

func cellToInt(cell interface{}) (int, error) {
    switch c := cell.(type) {
    case string:
        val, err := strconv.ParseInt(c, 10, 32)
        return int(val), err
    case int:
        return c, nil
    case float32:
        val := int(c)
        return int(val), nil
    case float64:
        val := int(c)
        return int(val), nil
    case bool:
        return 0, nil
    }
    /* TODO Error */
    return -1, nil
}

func cellToFloat(cell interface{}) (float32, error) {
    switch c := cell.(type) {
    case string:
        val, err := strconv.ParseFloat(c, 32)
        return float32(val), err
    case int:
        return float32(c), nil
    case float32:
        return c, nil
    case float64:
        return float32(c), nil
    case bool:
        return 0, nil
    }
    /* TODO Error */
    return -1, nil
}

func cellToStr(cell interface{}) (string, error) {
    switch c := cell.(type) {
    case string:
        return c, nil
    case int:
        return strconv.Itoa(c), nil
    case float32:
        return strconv.Itoa(int(c)), nil
    case float64:
        return strconv.Itoa(int(c)), nil
    case bool:
        return strconv.FormatBool(c), nil
    }
    /* TODO Error */
    return "", nil
}

func GetSheet(arg *Arg) (*Sheet, error) {
    b, err := ioutil.ReadFile(arg.CredentialToken)
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
        return nil, err
    }

    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
    if err != nil {
        log.Fatalf("Unable to parse client secret file to config: %v", err)
        return nil, err
    }
    client := getClient(config, arg)

    var sheet Sheet
    sheet.srv, err = sheets.New(client)
    if err != nil {
        log.Fatalf("Unable to retrieve Sheets client: %v", err)
        return nil, err
    }
    sheet.id = arg.SpreadsheetId

    return &sheet, nil
}
