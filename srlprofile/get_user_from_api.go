package srlprofile

import (
    "encoding/json"
    "fmt"
    "github.com/pkg/errors"
    "io"
    "net/http"
    "time"
)

// Mapping for the 'stats' field in SRL's API
type SrlStats struct {
    Rank int
    TotalRaces int
    TotalGames int
    FirstRace int
    FirstRaceDate int
    TotalTimePlayed int
    TotalFirstPlace int
    TotalSecondPlace int
    TotalThirdPlace int
    TotalQuits int
    TotalDisqualifications int
}

// Mapping for the 'player' field in SRL's API
type SrlPlayer struct {
    Id int
    Name string
    Channel string
    Api string
    Twitter string
    Youtube string
    Country string
}

// Mapping for the 'game' field in SRL's API
type SrlGame struct {
     Name string
     Abbrev string
}

// Object retrieved when doing a get for 'stat' in SRL's API
type SrlApiProfile struct {
    Stats SrlStats
    Player SrlPlayer
    Game SrlGame
}

// XXX: This was retrieved from SRL's profile page... sowwy D:
var client_id string

func getUserAvatar(channel string) (string, error) {
    url := fmt.Sprintf("https://api.twitch.tv/kraken/channels/%s?client_id=%s",
            channel, client_id)
    resp, err := http.Get(url)
    if err != nil {
        return "", errors.Wrap(err, "Failed to get twitch info")
    }
    defer resp.Body.Close()

    dec := json.NewDecoder(resp.Body)
    getNext := false
    for {
        t, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
            return "", errors.Wrap(err, "Failed to parse twitch info")
        }
        switch val := t.(type) {
        case string:
            if getNext {
                return val, nil
            }
            getNext = (val == "logo")
        default:
            continue
        /* Other valid cases are the following, but we don't really care abou those...
        case json.Delim:
        case bool:
        case float64:
        case json.Number:
        case nil:
        */
        }
    }

    return "", errors.New("Failed to get avatar from twitch info")
}

// GetFromApi retrieves and parses the user info retrieved from url, which must be a:
//   http://api.speedrunslive.com/stat?player=<username>
func GetFromApi(url string) (User, error) {
    // Download the user data
    resp, err := http.Get(url)
    if err != nil {
        return User{}, errors.Wrap(err, "Failed to get user from API")
    }
    defer resp.Body.Close()

    dec := json.NewDecoder(resp.Body)

    var api SrlApiProfile
    err = dec.Decode(&api)
    if err != nil {
        return User{}, errors.Wrap(err, "Failed to decode the JSON")
    }

    var u User

    u.Name = api.Player.Name
    u.Channel = api.Player.Channel
    u.FirstRace = time.Unix(int64(api.Stats.FirstRaceDate), 0).Format("Jan 2, 2006")
    u.NumRaces = api.Stats.TotalRaces
    dur := time.Duration(api.Stats.TotalTimePlayed)
    u.TotalTimePlayed = (dur * time.Second).String()
    u.NumGames = api.Stats.TotalGames
    u.NumFirst = api.Stats.TotalFirstPlace
    u.NumSecond = api.Stats.TotalSecondPlace
    u.NumThird = api.Stats.TotalThirdPlace
    u.NumForfeit = api.Stats.TotalQuits

    u.SrlAvatar, err = getUserAvatar(api.Player.Channel)
    if err != nil {
        // XXX: Failing to get the avatar isn't (imo) a critical error...
        fmt.Printf("Failed to get the player's avatar:\n\n%+v\n", err)
    }
    return u, nil
}

// GetFromUsername retrieves a user from SRL's API.
func GetFromUsername(username string) (User, error) {
    url := fmt.Sprintf("http://api.speedrunslive.com/stat?player=%s", username)
    return GetFromApi(url)
}
