package mtcareers

import (
    "fmt"
    "github.com/pkg/errors"
    "log"
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
        log.Fatalf("Unable to retrieve data from sheet: %v", err)
        return err
    }

    if len(resp.Values) == 0 {
        fmt.Println("No data found.")
        // TODO return err
    } else {
        row := resp.Values[0]
        s.TotalEntrants, err = cellToInt(row[0])
        if err != nil {
            return err
        }
        s.LatestEntrants, err = cellToInt(row[len(row)-1])
        if err != nil {
            return err
        }
    }

    return nil
}

func rowToUser(row []interface{}) (u User, err error) {
    u.Username, err = cellToStr(row[nameIdx])
    if err != nil {
        return
    }
    u.FirstMT, err = cellToStr(row[joinedMtIdx])
    if err != nil {
        return
    }
    if u.FirstMT[0] == '.' {
        u.FirstMT = u.FirstMT[1:]
    }
    u.TourneyCount, err = cellToInt(row[torneyCountIdx])
    if err != nil {
        return
    }
    u.WinCount, err = cellToInt(row[winIdx])
    if err != nil {
        return
    }
    u.LoseCount, err = cellToInt(row[loseIdx])
    if err != nil {
        return
    }
    u.DraftPoints, err = cellToFloat(row[draftIdx])
    return
}

func min(a,b int) int {
    if a < b {
        return a
    }
    return b
}

func (u *User) setHighestPosition(row []interface{}) {
    u.HighestPosition = 999

    for i, v := range row {
        if i < 2 || i >= len(_idxToPlace) {
            continue
        }
        cell, err := cellToInt(v)
        if err == nil && cell > 0 {
            u.HighestPosition = min(u.HighestPosition, _idxToPlace[i])
        }
    }
}

func (s *Sheet) GetUserInfo(username string) (User, error) {
    if _tourneyCache == nil {
        _range := fmt.Sprintf(baseRange,
            userInfoSheet,
            userFirstCol,
            userFirstRow,
            userLastCol,
            s.TotalEntrants)
        resp, err := s.srv.Spreadsheets.Values.Get(s.id, _range).Do()
        if err != nil {
            log.Fatalf("Unable to retrieve data from sheet: %v", err)
            return User{}, err
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
        resp, err := s.srv.Spreadsheets.Values.Get(s.id, _range).Do()
        if err != nil {
            log.Fatalf("Unable to retrieve data from sheet: %v", err)
            return User{}, err
        }

        _standingsCache = resp.Values
        for _, row := range _standingsCache[0] {
            st, ok := row.(string)
            if !ok || len(st) < 2 {
                continue
            }
            val, err := strconv.ParseInt(st[:len(st)-2], 10, 64)
            if err != nil {
                continue
            }
            _idxToPlace = append(_idxToPlace, int(val))
        }
    }

    var u User

    err := errors.New("User not found")
    for _, row := range _tourneyCache {
        if row[1] == username {
            u, err = rowToUser(row)
            break
        }
    }
    if err == nil {
        for _, row := range _standingsCache {
            if row[0] == username {
                u.setHighestPosition(row)
                break
            }
        }
    }

    return u, err
}
