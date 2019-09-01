package page

import (
    "fmt"
    "github.com/pkg/errors"
    "github.com/SirGFM/MTTitleCard/mtcareers"
    "github.com/SirGFM/MTTitleCard/srlprofile"
    "strconv"
)

// Data maps every API-retrieved user information into an structure understood
// by the template page
type Data struct {
    Channel string
    Username string
    Avatar string
    Joined string
    MtCount int
    Wins int
    Losses int
    WinRate string
    DraftPoints int
    HighestPlacement string
}

// _cache of already downloaded and parsed users
var _cache map[string]Data = map[string]Data{}
// _fmtNumber maps the unit of a position to its suffix
var _fmtNumber map[int]string = map[int]string {
    1: "%dst",
    2: "%dnd",
    3: "%drd",
}

// GenerateData downloads, parses and caches data for a given username.
// srlUsername and username should be the same.
func GenerateData(srlUsername, username string) error {
    if _, ok := _cache[username]; ok {
        // User already parsed and cached
        return nil
    }

    srlUser, err := srlprofile.GetFromUsername(srlUsername)
    if err != nil {
        return errors.Wrap(err, "Failed to get SRL Profile to generate user data")
    }

    sh, err := mtcareers.GetSheet()
    if err != nil {
        return errors.Wrap(err, "Failed to retrieve MT Career spreadsheet to generate user data")
    }
    err = sh.GetTourneyInfo()
    if err != nil {
        return errors.Wrap(err, "Failed to get tourney info to generate user data")
    }
    mtUser, err := sh.GetUserInfo(username)
    if err != nil {
        return errors.Wrap(err, "Failed to get MT Career user info to generate user data")
    }
    _cache[username] = generateDataFromUser(srlUser, mtUser)

    return nil
}

// generateDataFromUser merges the SRL User and the MT Career User in a single
// structure accepted by the template
func generateDataFromUser(srlUser srlprofile.User, mtUser mtcareers.User) Data {
    // For most numbers (and for 11, 12 and 13), just add a "th". In every other
    // case, simply add the correct suffix (e.g., 21st).
    unit := mtUser.HighestPosition % 10
    f, ok := _fmtNumber[unit]
    if !ok || (mtUser.HighestPosition % 100) - unit == 10 {
        f = "%dth"
    }
    pos := fmt.Sprintf(f, mtUser.HighestPosition)

    var rateStr string
    if mtUser.WinCount + mtUser.LoseCount != 0 {
        rate := float32(mtUser.WinCount)
        rate /= float32(mtUser.WinCount + mtUser.LoseCount)
        rate *= 100
        rateStr = strconv.Itoa(int(rate))
    } else {
        rateStr = "N/A"
    }

    return Data {
        Channel: srlUser.Channel,
        Username: mtUser.Username,
        Avatar: srlUser.SrlAvatar,
        Joined: mtUser.FirstMT,
        MtCount: mtUser.TourneyCount,
        Wins: mtUser.WinCount,
        Losses: mtUser.LoseCount,
        WinRate: rateStr,
        DraftPoints: int(mtUser.DraftPoints),
        HighestPlacement: pos,
    }
}
