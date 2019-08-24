package page

import (
    "fmt"
    "io"
    "html/template"
    "net/http"
)

type pageServer struct {
    userPage *template.Template
    body *template.Template
    httpServer http.Server
}
var srv pageServer

func (p *pageServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    username := req.URL.Path
    if username[0] == '/' {
        username = username[1:]
    }
    if username == "style.css" {
        w.Header().Set("Content-Type", "text/css")
        w.WriteHeader(http.StatusOK)
        io.WriteString(w, style)
    } else if len(username) > 0 {
        err := GenerateData(username, username)
        if err == nil {
            w.Header().Set("Content-Type", "text/html")
            w.WriteHeader(http.StatusOK)
            p.userPage.Execute(w, _cache[username])
        } else {
            w.WriteHeader(http.StatusNotFound)
        }
    } else {
        w.WriteHeader(http.StatusNotFound)
    }
}

func StarServer(port int) error {
    var err error

    srv.httpServer = http.Server {
        Addr: fmt.Sprintf(":%d", port),
        Handler: &srv,
    }

    srv.userPage = template.New("")
    _, err = srv.userPage.Parse(pageTemplate)
    if err != nil {
        return err
    }

    go func() {
        fmt.Println("waiting...")
        srv.httpServer.ListenAndServe()
    } ()

    return nil
}

func StopServer() {
    srv.httpServer.Close()
}
