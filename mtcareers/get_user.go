package mtcareers

import (
    "fmt"
    "log"
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

const nameIdx = 1
const torneyCountIdx = 2
const winIdx = 3
const loseIdx = 4

var _tourneyCache [][]interface{} = nil

type User struct {
    Username string
    TourneyCount int
    WinCount int
    LoseCount int
    HighestPosition int
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
    u.TourneyCount, err = cellToInt(row[torneyCountIdx])
    if err != nil {
        return
    }
    u.WinCount, err = cellToInt(row[winIdx])
    if err != nil {
        return
    }
    u.LoseCount, err = cellToInt(row[loseIdx])
    return
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

    for _, row := range _tourneyCache {
        if row[1] == username {
            return rowToUser(row)
        }
    }

    return User{}, nil
}
