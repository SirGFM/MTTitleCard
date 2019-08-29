package mtcareers

import (
    goErrors "errors"
    "fmt"
    "github.com/pkg/errors"
    "github.com/SirGFM/MTTitleCard/config"
    "strconv"
    "strings"
)

// baseRange formats the lookup string for retrieving a range from a
// spreadsheet. For example, to retrieve the columns "A" to "D" and rows "5" to
// "20" from a spreadsheet "TEST", one would request the following range:
//     "TEST!A5:D20"
// This can be more easily (or in a more organized fashion) by doing:
//     fmt.Sprintf(baseRange, "TEST", "A", 5, "D", 20)
const baseRange = "%s!%s%d:%s%d"

// _tourneyCache stores the downloaded tournament info spreadsheet
var _tourneyCache [][]interface{} = nil
// _standingsCache stores the downloaded standings spreadsheet
var _standingsCache [][]interface{} = nil

// _idxToPlace converts an index in the standings row into a tournament
// placement. The first two values are initialized to zero to skip the name and
// MT count columns.
var _idxToPlace []int = []int{0, 0}

// errNoPlacing indicates that the user still doesn't have a best placement, as
// this is most likely their first MT
var errNoPlacing error = goErrors.New("User haven't played in any tournament yet")

// User stores every info retrieved from the spreadsheet
type User struct {
    // The username (the same for SRL and MT Tournament)
    Username string
    // Initials of the first MT that user joined
    FirstMT string
    // Number of MTs that the user has joined
    TourneyCount int
    // How many victories the user has obtained through all their tournaments
    WinCount int
    // How many losses the user has obtained through all their tournaments
    LoseCount int
    // Highest position achieved by the user in a MT
    HighestPosition int
    // How many draft points the user is worth (?)
    DraftPoints float32
}

// GetTourneyInfo retrieve the total number of entrants and the number of
// entrants in the latest tournament.
func (s *Sheet) GetTourneyInfo() error {
    if s.TotalEntrants != 0 && s.LatestEntrants != 0 {
        // Info already cached, no need to do anything
        return nil
    }

    _range := fmt.Sprintf(baseRange,
        config.Get().TourneyInfo.SheetName,
        config.Get().TourneyInfo.FirstColumn,
        config.Get().TourneyInfo.FirstRow,
        config.Get().TourneyInfo.LastColumn,
        config.Get().TourneyInfo.FirstRow+1)
    resp, err := s.srv.Spreadsheets.Values.Get(s.id, _range).Do()
    if err != nil {
        return errors.Wrap(err, "Failed to get the number of entrants: Unable to retrieve data from sheet")
    }

    if len(resp.Values) == 0 {
        return errors.New("Failed to get the number of entrants: No data found")
    } else {
        row := resp.Values[0]
        s.TotalEntrants, err = cellToInt(row[0])
        if err != nil {
            return errors.Wrap(err, "Failed to parse TotalEntrants from sheet")
        }
        s.LatestEntrants, err = cellToInt(row[len(row)-1])
        if err != nil {
            return errors.Wrap(err, "Failed to parse LatestEntrants from sheet")
        }
    }

    return nil
}

// rowToUser convert a row, retrieved from the spreadsheet, into a User
func rowToUser(row []interface{}) (u User, err error) {
    u.Username, err = cellToStr(row[config.Get().NameIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse Username from sheet")
        return
    }
    u.FirstMT, err = cellToStr(row[config.Get().JoinedMtIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse FirstMT from sheet")
        return
    }
    if u.FirstMT[0] == '.' {
        u.FirstMT = u.FirstMT[1:]
    }
    u.TourneyCount, err = cellToInt(row[config.Get().TorneyCountIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse TourneyCount from sheet")
        return
    }
    u.WinCount, err = cellToInt(row[config.Get().WinIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse WinCount from sheet")
        return
    }
    u.LoseCount, err = cellToInt(row[config.Get().LoseIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse LoseCount from sheet")
        return
    }
    u.DraftPoints, err = cellToFloat(row[config.Get().DraftIdx])
    // XXX: if err == nil, errors.Wrap returns nil as well!
    err = errors.Wrap(err, "Failed to parse DraftPoints from sheet")
    return
}

// min return the smallest of two values
func min(a,b int) int {
    if a < b {
        return a
    }
    return b
}

// setHighestPosition from a player's row in the spreadsheet
func (u *User) setHighestPosition(row []interface{}) (err error) {
    u.HighestPosition = 9999
    found := false

    for i, v := range row {
        if i < 2 || i >= len(_idxToPlace) {
            continue
        }
        cell, gerr := cellToInt(v)
        if gerr != nil {
            return errors.Wrap(gerr, "Failed to parse user's highest position")
        } else if cell > 0 {
            u.HighestPosition = min(u.HighestPosition, _idxToPlace[i])
            found = true
        }
    }

    if !found {
        err = errors.Wrap(errNoPlacing, "")
    }
    return
}

// colToStr try to retrieve a string from an interface, return "" if it fails
func colToStr(c interface{}) string {
    s, ok := c.(string)
    if !ok {
        return ""
    }
    return s
}

// getUserRow, doing a case insensitive look up
func getUserRow(username string, nameIdx int, rows[][]interface{}) []interface{} {
    name := strings.ToLower(username)
    // First assume that the list is sorted and try a binary search
    for l, r := 0, len(rows) - 1; l <= r; {
        m := (l + r) / 2
        row := rows[m]
        switch strings.Compare(name, strings.ToLower(colToStr(row[nameIdx]))) {
        case 0:
            return row
        case 1:
            l = m + 1
        case -1:
            r = m - 1
        }
    }
    // If not found, look sequentially
    for _, row := range rows {
        if strings.ToLower(colToStr(row[nameIdx])) == name {
            return row
        }
    }

    return nil
}

// GetUserInfo from the MT Career spreadsheet
func (s *Sheet) GetUserInfo(username string) (u User, err error) {
    // Download and cache the participants info and standings through every
    // tournament
    if _tourneyCache == nil {
        _range := fmt.Sprintf(baseRange,
            config.Get().UserInfo.SheetName,
            config.Get().UserInfo.FirstColumn,
            config.Get().UserInfo.FirstRow,
            config.Get().UserInfo.LastColumn,
            s.TotalEntrants)
        resp, gerr := s.srv.Spreadsheets.Values.Get(s.id, _range).Do()
        if gerr != nil {
            err = errors.Wrap(gerr, "Unable to retrieve tourney data from sheet")
            return
        }

        _tourneyCache = resp.Values
    }
    if _standingsCache == nil {
        _range := fmt.Sprintf(baseRange,
            config.Get().StandingsInfo.SheetName,
            config.Get().StandingsInfo.FirstColumn,
            config.Get().StandingsInfo.FirstRow,
            config.Get().StandingsInfo.LastColumn,
            s.TotalEntrants)
        resp, gerr := s.srv.Spreadsheets.Values.Get(s.id, _range).Do()
        if gerr != nil {
            err = errors.Wrap(gerr, "Unable to retrieve standings from sheet")
            return
        }

        // Convert an index in the standings cache to a tournament placement
        _standingsCache = resp.Values
        for i, row := range _standingsCache[0] {
            st, ok := row.(string)
            if !ok || len(st) < 2 || i < 2 {
                continue
            }
            // Remove the position suffix (e.g., 1st, 2nd etc)
            val, gerr := strconv.ParseInt(st[:len(st)-2], 10, 64)
            if gerr != nil {
                err = errors.Wrap(gerr, "Unable to map standing index in sheet to tournament placement")
                return
            } else if val == 0 {
                continue
            }
            _idxToPlace = append(_idxToPlace, int(val))
        }
    }

    // Retrieve the user info from the previously downloaded data
    row := getUserRow(username, config.Get().NameIdx, _tourneyCache)
    if row == nil {
        err = errors.New(fmt.Sprintf("User not found: '%s'", username))
        return
    }
    u, err = rowToUser(row)
    if err == nil {
        posRow := getUserRow(username, 0, _standingsCache)
        err = u.setHighestPosition(posRow)
        if errors.Cause(err) == errNoPlacing {
            err = nil
        }
    }

    return
}
