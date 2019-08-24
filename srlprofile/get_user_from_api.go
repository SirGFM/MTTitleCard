package srlprofile

import (
    "fmt"
    "github.com/pkg/errors"
    "encoding/json"
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

    return u, nil
}

// GetFromUsername retrieves a user from SRL's API.
func GetFromUsername(username string) (User, error) {
    url := fmt.Sprintf("http://api.speedrunslive.com/stat?player=%s", username)
    return GetFromApi(url)
}
