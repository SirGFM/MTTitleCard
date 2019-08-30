package mtcareers

import (
    "encoding/json"
    "fmt"
    "github.com/pkg/errors"
    mttcConfig "github.com/SirGFM/MTTitleCard/config"
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
    // Object used to access the spreadsheet
    srv *sheets.Service
    // ID of the spreadsheet being accessed
    id string
    // Number of entrants through every tournament
    TotalEntrants int
    // Number of entrants on the latest tournament
    LatestEntrants int
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
    authToken, err := CheckToken()
    if err != nil {
        return nil, errors.Wrap(err, "Failed to check OAuth token")
    } else if authToken != "" {
        serr := fmt.Sprintf("OAuth token not found!\n\nAccess http://localhost:%d\n", mttcConfig.Get().Port)
        return nil, errors.New(serr)
    }

    // This shouldn't ever fail, since the token has already been checked
    tok, err := tokenFromFile(mttcConfig.Get().TokenFile)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to retrieve the OAuth token")
    }
    return config.Client(context.Background(), tok), nil
}

// Get the OAuth2 config for the given credential file
func getConfig() (*oauth2.Config, error) {
    b, err := ioutil.ReadFile(mttcConfig.Get().CredentialFile)
    if err != nil {
        return nil, errors.Wrap(err, "Unable to read client secret file")
    }

    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
    // XXX: if err == nil, errors.Wrap returns nil as well!
    return config, errors.Wrap(err, "Unable to parse client secret file to config")
}

// CheckToken. If it's not valid (mainly because it doesn't exist), generate
// and return a auth URL.
func CheckToken() (string, error) {
    _, err := tokenFromFile(mttcConfig.Get().TokenFile)
    if err != nil {
        config, err := getConfig()
        if err != nil {
            return "", errors.Wrap(err, "Unable to parse client secret file to config")
        }
        return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), nil
    }
    return "", nil
}

// SaveAuthentication finishes authenticating with OAuth2 and saves the token
func SaveAuthentication(authCode string) error {
    config, err := getConfig()
    if err != nil {
        return errors.Wrap(err, "Unable to parse client secret file to config")
    }
    tok, err := config.Exchange(context.TODO(), authCode)
    if err != nil {
        return errors.Wrap(err, "Unable to retrieve token from web")
    }
    err = saveToken(mttcConfig.Get().TokenFile, tok)
    // XXX: if err == nil, errors.Wrap returns nil as well!
    return errors.Wrap(err, "Failed to save the OAuth token")
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
    log.Printf("Saving credential file to: %s", path)
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
func GetSheet() (*Sheet, error) {
    // If modifying these scopes, delete your previously saved token.json.
    config, err := getConfig()
    if err != nil {
        return nil, errors.Wrap(err, "Unable to parse client secret file to config")
    }
    client, err := getClient(config)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to initialize the Google API client")
    }

    var sheet Sheet
    sheet.srv, err = sheets.New(client)
    if err != nil {
        return nil, errors.Wrap(err, "Unable to retrieve Sheets accessor")
    }
    sheet.id = mttcConfig.Get().MtCareerSpreasheet

    return &sheet, nil
}
