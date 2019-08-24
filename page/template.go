package page

const style = `
body {
    background: #202225;
    color: #ffffff;
    font-weight: bold;
    font-family: sans-serif;
}
.channel {
    font-size: small;
    margin: 1em;
    margin-bottom: 2em;
    margin-top: 0.5em;
    display: block;
    width: auto;
    text-align: right;
}
.user {
    margin: 1.5em;
}
.username {
    font-size: xx-large;
    display: block;
}
.stats {
    margin: 1.5em;
    width: 90%;
}
.stats_label {
    width: 80%;
}
.stats_field {
}
`

const pageTemplate = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <title> {{.Username}}'s MT Title Card </title>
        <link rel="stylesheet" href="/style.css">
        <meta charset="UTF-8">
    </head>
    <body>
        <label class="channel" id="channel">
            twitch.tv/{{.Channel}}
        </label>
        <div class="user" id="user">
            <!-- Avatar -->
            <label class="username" id="username">
                {{.Username}}
            </label>
        </div>
        <table class="stats" id="stats"><tbody>
            <tr>
                <td class="stats_label" id="stats_label">Joined</td>
                <td class="stats_field" id="stats_field">{{.Joined}}</td>
            </tr>
            <tr>
                <td class="stats_label" id="stats_label">MT Count</td>
                <td class="stats_field" id="stats_field">{{.MtCount}}</td>
            </tr>
            <tr>
                <td class="stats_label" id="stats_label">Wins</td>
                <td class="stats_field" id="stats_field">{{.Wins}}</td>
            </tr>
            <tr>
                <td class="stats_label" id="stats_label">Losses</td>
                <td class="stats_field" id="stats_field">{{.Losses}}</td>
            </tr>
            <tr>
                <td class="stats_label" id="stats_label">Win Rate</td>
                <td class="stats_field" id="stats_field">{{.WinRate}}</td>
            </tr>
            <tr>
                <td class="stats_label" id="stats_label">Draft Points</td>
                <td class="stats_field" id="stats_field">{{.DraftPoints}}</td>
            </tr>
            <tr>
                <td class="stats_label" id="stats_label">Highest Placement</td>
                <td class="stats_field" id="stats_field">{{.HighestPlacement}}</td>
            </tr>
        </tbody></table>
    </body>
</html>
`