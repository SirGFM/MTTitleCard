package mtcareers

import (
    goErrors "errors"
    "fmt"
    "github.com/pkg/errors"
    "strconv"
)

const baseRange = "%s!%s%d:%s%d"

// TODO retrieve exact row with stats
const torneyInfoSheet = "STATS"
const tourneyTotalCol = "C"
const currentTotalCol = "F"
const totalsRow = 16

const userInfoSheet = "MT Career"
const userFirstCol = "A"
const userLastCol = "H"
const userFirstRow = 2

const joinedMtIdx = 0
const nameIdx = 1
const torneyCountIdx = 2
const winIdx = 3
const loseIdx = 4
const draftIdx = 7

const standingSheet = "MT Career Standings"
const standingsFirstCol = "A"
const standingsLastCol = "S"
const standingsFirstRow = 1

var _tourneyCache [][]interface{} = nil
var _standingsCache [][]interface{} = nil
// Initialize the places skipping the name and MT count columns
var _idxToPlace []int = []int{0, 0}

var errNoPlacing error = goErrors.New("User haven't played in any tournament yet")

type User struct {
    Username string
    FirstMT string
    TourneyCount int
    WinCount int
    LoseCount int
    HighestPosition int
    DraftPoints float32
}

func (s *Sheet) GetTourneyInfo() error {
    if s.TotalEntrants != 0 && s.LatestEntrants != 0 {
        // Info already cached, no need to do anything
        return nil
    }

    _range := fmt.Sprintf(baseRange,
        torneyInfoSheet,
        tourneyTotalCol,
        totalsRow,
        currentTotalCol,
        totalsRow+1)
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

func rowToUser(row []interface{}) (u User, err error) {
    u.Username, err = cellToStr(row[nameIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse Username from sheet")
        return
    }
    u.FirstMT, err = cellToStr(row[joinedMtIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse FirstMT from sheet")
        return
    }
    if u.FirstMT[0] == '.' {
        u.FirstMT = u.FirstMT[1:]
    }
    u.TourneyCount, err = cellToInt(row[torneyCountIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse TourneyCount from sheet")
        return
    }
    u.WinCount, err = cellToInt(row[winIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse WinCount from sheet")
        return
    }
    u.LoseCount, err = cellToInt(row[loseIdx])
    if err != nil {
        err = errors.Wrap(err, "Failed to parse LoseCount from sheet")
        return
    }
    u.DraftPoints, err = cellToFloat(row[draftIdx])
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

func (s *Sheet) GetUserInfo(username string) (u User, err error) {
    if _tourneyCache == nil {
        _range := fmt.Sprintf(baseRange,
            userInfoSheet,
            userFirstCol,
            userFirstRow,
            userLastCol,
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
            standingSheet,
            standingsFirstCol,
            standingsFirstRow,
            standingsLastCol,
            s.TotalEntrants)
        resp, gerr := s.srv.Spreadsheets.Values.Get(s.id, _range).Do()
        if gerr != nil {
            err = errors.Wrap(gerr, "Unable to retrieve standings from sheet")
            return
        }

        _standingsCache = resp.Values
        for i, row := range _standingsCache[0] {
            st, ok := row.(string)
            if !ok || len(st) < 2 || i < 2 {
                continue
            }
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

    err = errors.New(fmt.Sprintf("User not found: '%s'", username))
    for _, row := range _tourneyCache {
        if row[1] == username {
            u, err = rowToUser(row)
            break
        }
    }
    if err == nil {
        for _, row := range _standingsCache {
            if row[0] == username {
                err = u.setHighestPosition(row)
                if errors.Cause(err) == errNoPlacing {
                    err = nil
                }
                break
            }
        }
    }

    return
}
