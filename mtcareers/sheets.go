package mtcareers

import (
    "encoding/json"
    "fmt"
    "github.com/pkg/errors"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/sheets/v4"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
)

type Sheet struct {
    // Object used to access the spreadsheet
    srv *sheets.Service
    // ID of the spreadsheet being accessed
    id string
    // Number of entrants through every tournament
    TotalEntrants int
    // Number of entrants on the latest tournament
    LatestEntrants int
}

type Arg struct {
    // Path to the API token file, if not the default path
    ApiToken string
    // Path to JSON downloaded after enabling the API
    CredentialToken string
    // ID of the spreadsheet being accessed
    SpreadsheetId string
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, arg *Arg) (*http.Client, error) {
    // The file token.json stores the user's access and refresh tokens, and is
    // created automatically when the authorization flow completes for the first
    // time.
    if arg.ApiToken == "" {
        arg.ApiToken = "token.json"
    }
    tok, err := tokenFromFile(arg.ApiToken)
    if err != nil {
        tok, err = getTokenFromWeb(config)
        if err != nil {
            return nil, errors.Wrap(err, "Failed to retrieve the OAuth token")
        }
        err := saveToken(arg.ApiToken, tok)
        if err != nil {
            return nil, errors.Wrap(err, "Failed to save the OAuth token")
        }
    }
    return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the "+
            "authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
        return nil, errors.Wrap(err, "Unable to read authorization code")
    }

    // XXX: if err == nil, errors.Wrap returns nil as well!
    tok, err := config.Exchange(context.TODO(), authCode)
    return tok, errors.Wrap(err, "Unable to retrieve token from web")
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to open token file")
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    // XXX: if err == nil, errors.Wrap returns nil as well!
    return tok, errors.Wrap(err, "Failed to decode token JSON")
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return errors.Wrap(err, "Unable to cache oauth token")
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
    return nil
}

// cellToInt converts a cell to an integer
func cellToInt(cell interface{}) (i int, err error) {
    switch c := cell.(type) {
    case string:
        val, gerr := strconv.ParseInt(c, 10, 32)
        i = int(val)
        err = errors.Wrap(gerr, "Failed to parse string cell into integer")
    case int:
        i = c
    case float32:
        i = int(c)
    case float64:
        i = int(c)
    case bool:
        i = 0
    default:
        err = errors.New("Failed to understand and convert cell to integer")
    }
    return
}

// cellToFloat converts a cell to a float
func cellToFloat(cell interface{}) (f32 float32, err error) {
    switch c := cell.(type) {
    case string:
        val, gerr := strconv.ParseFloat(c, 32)
        f32 = float32(val)
        err = errors.Wrap(gerr, "Failed to parse string cell into float")
    case int:
        f32 = float32(c)
    case float32:
        f32 = c
    case float64:
        f32 = float32(c)
    case bool:
        f32 = 0
    default:
        err = errors.New("Failed to understand and convert cell to float")
    }
    return
}

// cellToString converts a cell to a string
func cellToStr(cell interface{}) (str string, err error) {
    switch c := cell.(type) {
    case string:
        str = c
    case int:
        str = strconv.Itoa(c)
    case float32:
        str = strconv.Itoa(int(c))
    case float64:
        str = strconv.Itoa(int(c))
    case bool:
        str = strconv.FormatBool(c)
    default:
        err = errors.New("Failed to understand and convert cell to string")
    }
    return
}

// GetSheet retrieves an object for accessing an spreadsheet
func GetSheet(arg *Arg) (*Sheet, error) {
    b, err := ioutil.ReadFile(arg.CredentialToken)
    if err != nil {
        return nil, errors.Wrap(err, "Unable to read client secret file")
    }

    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
    if err != nil {
        return nil, errors.Wrap(err, "Unable to parse client secret file to config")
    }
    client, err := getClient(config, arg)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to initialize the Google API client")
    }

    var sheet Sheet
    sheet.srv, err = sheets.New(client)
    if err != nil {
        return nil, errors.Wrap(err, "Unable to retrieve Sheets accessor")
    }
    sheet.id = arg.SpreadsheetId

    return &sheet, nil
}
