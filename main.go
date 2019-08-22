package main

import (
    "github.com/SirGFM/MTTitleCard/srltitlecard"
    "flag"
    "fmt"
)

func main() {
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
