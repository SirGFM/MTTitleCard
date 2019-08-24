package page

import (
    "fmt"
    "github.com/SirGFM/MTTitleCard/mtcareers"
    "github.com/SirGFM/MTTitleCard/srlprofile"
)

type Data struct {
    Channel string
    Username string
    Joined string
    MtCount int
    Wins int
    Losses int
    WinRate int
    DraftPoints int
    HighestPlacement string
}

var _cache map[string]Data = map[string]Data{}

var arg mtcareers.Arg = mtcareers.Arg {
    CredentialToken: "credentials.json",
    SpreadsheetId: "1LE6z_xRRxtIcCKYDzH9ag_1Iry6iHlqhpc09mqTZfiU",
}

func GenerateData(srlUsername, username string) error {
    if _, ok := _cache[username]; ok {
        return nil
    }

    srlUser, err := srlprofile.GetFromUsername(srlUsername)
    if err != nil {
        return err
    }

    sh, err := mtcareers.GetSheet(&arg)
    if err != nil {
        return err
    }
    err = sh.GetTourneyInfo()
    if err != nil {
        //panic(fmt.Sprintf("Failed to get tourney info: %+v", err))
        return err
    }
    mtUser, err := sh.GetUserInfo(username)
    if err != nil {
        //panic(fmt.Sprintf("Failed to user info: %+v", err))
        return err
    }

    var pos string
    switch p := mtUser.HighestPosition; p % 10 {
    case 1:
        pos = fmt.Sprintf("%dst", p)
    case 2:
        pos = fmt.Sprintf("%dnd", p)
    case 3:
        pos = fmt.Sprintf("%drd", p)
    default:
        pos = fmt.Sprintf("%dth", p)
    }

    rate := float32(mtUser.WinCount)
    rate /= float32(mtUser.WinCount + mtUser.LoseCount)
    var d Data = Data {
        Channel: srlUser.Channel,
        Username: mtUser.Username,
        Joined: mtUser.FirstMT,
        MtCount: mtUser.TourneyCount,
        Wins: mtUser.WinCount,
        Losses: mtUser.LoseCount,
        WinRate: int(rate * 100),
        DraftPoints: int(mtUser.DraftPoints),
        HighestPlacement: pos,
    }
    _cache[username] = d

    return nil
}
