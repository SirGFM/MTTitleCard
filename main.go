package main

import (
    "github.com/SirGFM/MTTitleCard/mtcareers"
    "github.com/SirGFM/MTTitleCard/srlprofile"
    "github.com/SirGFM/MTTitleCard/page"
    "flag"
    "fmt"
    "os"
    "os/signal"
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
        u := srlprofile.Get(srlURL);
        fmt.Printf("%+v\n", u)
    } else if srlUser != "" {
        u, err := srlprofile.GetFromUsername(srlUser);
        fmt.Printf("%+v\nerr: %+v\n", u, err)
    }
}

func testPage() {
    err := page.StartServer(8080)
    if err != nil {
        panic(fmt.Sprintf("Failed to start server: %+v", err))
    }

    signalTrap := make(chan os.Signal, 1)
    wait := make(chan struct{}, 1)
    go func (c chan os.Signal) {
        _ = <-c
        page.StopServer()
        wait <- struct{}{}
    } (signalTrap)
    signal.Notify(signalTrap, os.Interrupt)

    <-wait
}

func main() {
    //testMtCareers()
    //testSrlTitleCard()
    testPage()
}
