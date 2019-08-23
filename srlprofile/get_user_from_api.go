package srlprofile

import (
    "fmt"
    "encoding/json"
    "net/http"
    "time"
)

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

type SrlPlayer struct {
    Id int
    Name string
    Channel string
    Api string
    Twitter string
    Youtube string
    Country string
}

type SrlGame struct {
     Name string
     Abbrev string
}

type SrlApiProfile struct {
    Stats SrlStats
    Player SrlPlayer
    Game SrlGame
}

func GetFromApi(url string) User {
    // Download the user data
    resp, err := http.Get(url)
    if err != nil {
        panic(fmt.Sprintf("Failed to get user from API: %+v", err))
    }
    defer resp.Body.Close()

    dec := json.NewDecoder(resp.Body)

    var api SrlApiProfile
    err = dec.Decode(&api)
    if err != nil {
        panic(fmt.Sprintf("Failed to decode the JSON: %+v", err))
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

    return u
}

func GetFromUsername(username string) User {
    url := fmt.Sprintf("http://api.speedrunslive.com/stat?player=%s", username)
    return GetFromApi(url)
}
