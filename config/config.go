package config

import (
    "encoding/json"
    "github.com/pkg/errors"
    "io"
    "io/ioutil"
    "log"
    "os"
)

// SheetRange defines a range within a spreadsheet
type SheetRange struct {
    SheetName string
    FirstColumn string
    LastColumn string
    FirstRow int
}

type Config struct {
    // Port for serving the webpages
    Port int
    // Client ID used when retrieving values from twitch's API
    TwitchClientID string
    // Path to a JSON file with the credentials to access Google's API
    CredentialFile string
    // Path to a JSON file with the token to access Google's API
    TokenFile string
    // ID of the MT Career spreadsheet
    MtCareerSpreasheet string
    // Range used to extract the number of entrants
    TourneyInfo SheetRange
    // Range used to extract the a tournament entrant
    UserInfo SheetRange
    // Range used to extract the entrants standings
    StandingsInfo SheetRange
    // Index, in the spreadsheet, of the number of tournaments entered by the user
    JoinedMtIdx int
    // Index, in the spreadsheet, of the user's name
    NameIdx int
    // Index, in the spreadsheet, of the username
    TorneyCountIdx int
    // Index, in the spreadsheet, of the user's number of victories
    WinIdx int
    // Index, in the spreadsheet, of the user's number of losses
    LoseIdx int
    // Index, in the spreadsheet, of the user's value(?) in draft points
    DraftIdx int
    // Path to a CSS file used to override the default CSS
    CssFile string
    // Style sheet for the player page
    cssData []byte
    // Path to a HTML-template file used to override the default page template
    TemplateFile string
    // Template for the player page
    templateData []byte
    // URI of the service within the server. Mostly used to set the path to the CSS file.
    ServiceUri string
}

// Store the loaded configuration
var config Config

// Get a copy of the configuration
func Get() Config {
    return config
}

// WriteCss, if CssFile was supplied, or the fallback if not.
func (c Config) WriteCss(w io.Writer, fallback []byte) error {
    data := fallback
    if len(c.cssData) != 0 {
        data = c.cssData
    }
    // TODO: Check 'n' here
    _, err := w.Write(data)
    // XXX: if err == nil, errors.Wrap returns nil as well!
    return errors.Wrap(err, "Failed to write the CSS data")
}

// PageTemplate returns a the page template or the supplied fallback
func (c Config) PageTemplate(fallback string) string {
    if len(c.templateData) != 0 {
        return string(c.templateData)
    }
    return fallback
}

// Retrieve the default configurations
func GetDefault() Config {
    return Config {
        Port: 8080,
        TwitchClientID: "",
        CredentialFile: "credentials.json",
        TokenFile: "token.json",
        MtCareerSpreasheet: "1DWYq3T1w8u1N0CWWJ72tqQRv67c1eY098u0wyuiMEmA",
        TourneyInfo: SheetRange {
            SheetName: "STATS",
            FirstColumn: "C",
            LastColumn: "F",
            FirstRow: 16,
        },
        UserInfo: SheetRange {
            SheetName: "MT Career",
            FirstColumn: "A",
            LastColumn: "H",
            FirstRow: 2,
        },
        StandingsInfo: SheetRange {
            SheetName: "MT Career Standings",
            FirstColumn: "A",
            LastColumn: "S",
            FirstRow: 1,
        },
        JoinedMtIdx: 0,
        NameIdx: 1,
        TorneyCountIdx: 2,
        WinIdx: 3,
        LoseIdx: 4,
        DraftIdx: 7,
    }
}

// Load the supplied configuration
func LoadConfig(customConfig Config) error {
    var err error

    if customConfig.CssFile != "" {
        customConfig.cssData, err = ioutil.ReadFile(customConfig.CssFile)
        if err != nil {
            return errors.Wrap(err, "Failed to read the custom CSS file")
        }
    }

    if customConfig.TemplateFile != "" {
        customConfig.templateData, err = ioutil.ReadFile(customConfig.TemplateFile)
        if err != nil {
            return errors.Wrap(err, "Failed to read the custom template file")
        }
    }

    config = customConfig
    return nil
}

// Load the supplied configuration on file, or the default configuration
func Load(path string) error {
    if path == "" {
        log.Print("Using the default configuration...")
        config = GetDefault()
        return nil
    }

    log.Printf("Loading the configuration from '%s'...", path)
    f, err := os.Open(path)
    if err != nil {
        return errors.Wrap(err, "Failed to open config file")
    }
    defer f.Close()
    dec := json.NewDecoder(f)

    var _config Config
    err = dec.Decode(&_config)
    if err != nil {
        return errors.Wrap(err, "Failed to decode config JSON")
    }

    return LoadConfig(_config)
}
