package main

import (
    "github.com/SirGFM/MTTitleCard/srltitlecard"
    "flag"
    "fmt"
)

func main() {
    var srlURL string

    flag.StringVar(&srlURL, "srlURL", "", "URL of the player's SRL page")
    flag.Parse()

    if srlURL == "" {
        panic("Missing -srlURL")
    }

    u := srltitlecard.Get(srlURL);
    fmt.Printf("%+v\n", u)
}
