package main

import (
    "github.com/SirGFM/MTTitleCard/mtcareers"
    "github.com/SirGFM/MTTitleCard/srltitlecard"
    "flag"
    "fmt"
)

func testMtCareers() {
    arg := mtcareers.Arg {
        CredentialToken: "credentials.json",
        SpreadsheetId: "1LE6z_xRRxtIcCKYDzH9ag_1Iry6iHlqhpc09mqTZfiU",
    }
    sh, err := mtcareers.GetSheet(&arg)
    if err != nil {
        panic(fmt.Sprintf("Failed to load sheet: %+v", err))
    }

    err = sh.GetTourneyInfo()
    if err != nil {
        panic(fmt.Sprintf("Failed to get tourney info: %+v", err))
    }
    usr, err := sh.GetUserInfo("GFM")
    if err != nil {
        panic(fmt.Sprintf("Failed to get user info: %+v", err))
    }
    fmt.Printf("%+v\n", usr)
}


func testSrlTitleCard() {
    var srlURL string
    var srlUser string

    flag.StringVar(&srlURL, "srlURL", "", "URL of the player's SRL page")
    flag.StringVar(&srlUser, "srlUser", "", "Username of player on SRL")
    flag.Parse()

    if srlURL == "" && srlUser == "" {
        panic("Missing either -srlURL or -srlUser")
    }

    if srlURL != "" {
        u := srltitlecard.Get(srlURL);
        fmt.Printf("%+v\n", u)
    } else if srlUser != "" {
        u := srltitlecard.GetFromUsername(srlUser);
        fmt.Printf("%+v\n", u)
    }
}

func main() {
    //testMtCareers()
    testSrlTitleCard()
}
