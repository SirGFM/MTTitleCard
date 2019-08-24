package srlprofile

import (
    "fmt"
    "github.com/SirGFM/MTTitleCard/getelementbyid"
    "golang.org/x/net/html"
    "net/http"
    "reflect"
    "strconv"
)

type User struct {
    // The username
    Name string `parent:"profile_name" id:"profile_playername" final:"true"`
    // URL to the user's twitch avatar
    SrlAvatar string `parent:"profile_name" id:"avatarHolder"`
    // Date of the first race ever
    FirstRace string `parent:"profile_races" id:"date" final:"true"`
    // How many races this user has taken part of
    NumRaces int `parent:"profile_races" id:"races" final:"true"`
    // Sum of race times
    TotalTimePlayed string `parent:"profile_races" id:"played" final:"true"`
    // Number of different games played
    NumGames int `parent:"profile_races" id:"games" final:"true"`
    // Number of times that the user has gotten first place
    NumFirst int `parent:"profile_races" id:"firsts" final:"true"`
    // Number of times that the user has gotten second place
    NumSecond int `parent:"profile_races" id:"seconds" final:"true"`
    // Number of times that the user has gotten third place
    NumThird int `parent:"profile_races" id:"thirds" final:"true"`
    // Number of times that the user has gotten forfeited
    NumForfeit int `parent:"profile_races" id:"quits" final:"true"`
    // User streaming channel
    Channel string

    // Nodes used to more quickly retrieve an item
    nodes map[string]*html.Node
}

// Get SRL profile page and parse it. THIS FUNCTION WORKS ONLY ON STATIC PAGES!
func Get(url string) User {
    // Download the user data
    resp, err := http.Get(url)
    if err != nil {
        panic(fmt.Sprintf("Failed to get racer page: %+v", err))
    }
    defer resp.Body.Close()
    doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Fail to parse: %+v", err))
	}

    // Parse it into a struct (using the tags/strings in the structure)
    var u User
    u.nodes = make(map[string]*html.Node)

    _type := reflect.TypeOf(u)
    _val := reflect.ValueOf(&u).Elem()
    for i := 0; i < _type.NumField(); i++ {
        // Retrieve each actual field
        field := _type.Field(i)
        if field.Name[0] >= 'a' && field.Name[0] <= 'z' {
            // Ignore private fields
            continue
        }

        // Retrieve the parent field ID in the HTML
        parent, ok := field.Tag.Lookup("parent")
        if !ok {
            panic(fmt.Sprintf("Failed to find parent of field %s", field.Name))
        }
        parentNode, ok := u.nodes[parent]
        if !ok {
            parentNode = getelementbyid.Find(doc, parent)
            u.nodes[parent] = parentNode
        }

        // Retrieve the field itself from the HTML (if it shouldn't have any children)
        id, ok := field.Tag.Lookup("id")
        if !ok {
            panic(fmt.Sprintf("Failed to get value of field %s", field.Name))
        }
        _, ok = field.Tag.Lookup("final")
        if ok {
            var data string
            outerNode := getelementbyid.Find(parentNode, id)
            if outerNode.FirstChild != nil {
                data = outerNode.FirstChild.Data
            } else {
                data = "Unknown"
            }

            // Set it back into the structure
            switch (field.Type.Kind()) {
            case reflect.String:
                _val.Field(i).SetString(data)
            case reflect.Int,
                reflect.Int8,
                reflect.Int16,
                reflect.Int32,
                reflect.Int64:

                iData, err := strconv.ParseInt(data, 10, 64)
                if err != nil {
                    iData = 0
                }
                _val.Field(i).SetInt(iData)
            }
        } else {
            // TODO
        }
    }

    return u
}
