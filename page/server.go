package page

import (
    "fmt"
    "github.com/pkg/errors"
    "github.com/SirGFM/MTTitleCard/config"
    "github.com/SirGFM/MTTitleCard/mtcareers"
    "html/template"
    "log"
    "net/http"
)

// pageServer wraps the request handler, and everything used by it
type pageServer struct {
    // userPage is a template used to fill a user's info
    userPage *template.Template
    // renewPage is a template used to renew the server's token
    renewPage *template.Template
    // httpServer handling requests from the client
    httpServer *http.Server
}
// srv is the currently executing server, mainly used to top the server
var srv pageServer

// request wraps parameters common to handlers to simplify calling sub-functions
type request struct {
    // p is the running server (should be the same as srv)
    p *pageServer
    // w is this request's response writer
    w http.ResponseWriter
    // req is the actual received request
    req *http.Request
    // path is the received URL (e.g., "index", "favicon.ico
    path string
}

// getCss retrieve the CSS used by the server
func (r *request) getCss() {
    r.w.Header().Set("Content-Type", "text/css")
    r.w.WriteHeader(http.StatusOK)
    config.Get().WriteCss(r.w, []byte(style))
}

// getUserData parse username info, fit it into the template and return the
// resulting page. In case of error, the stack trace is returned back to the
// client
func (r *request) getUserData(username string) {
    err := GenerateData(username, username)
    if err == nil {
        r.w.Header().Set("Content-Type", "text/html")
        r.w.WriteHeader(http.StatusOK)
        r.p.userPage.Execute(r.w, _cache[username])
    } else {
        serr := fmt.Sprintf("%+v", err)
        http.Error(r.w, serr, http.StatusNotFound)
        log.Print(serr)
    }
}

// Data supplied to the renew token page
type RenewData struct {
    Url string
}

// getRenewToken and display it to the client, with a form to post the
// renewed token
func (r *request) getRenewToken() {
    url, err := mtcareers.CheckToken()
    if err != nil {
        serr := fmt.Sprintf("%+v", err)
        http.Error(r.w, serr, http.StatusNotFound)
        log.Print(serr)
        return
    }
    d := RenewData {
        Url: url,
    }
    r.w.Header().Set("Content-Type", "text/html")
    r.w.WriteHeader(http.StatusOK)
    r.p.renewPage.Execute(r.w, d)
}

// get handles GET requests
func (r *request) get() {
    switch r.path {
    case "style.css":
        r.getCss()
    case "",
        "index",
        "index.html":

        r.getRenewToken()
    case "favicon.ico":
        http.Error(r.w, "Missing a favicon...", http.StatusNotFound)
    default:
        r.getUserData(r.path)
    }
}

// post handles POST requests
func (r *request) post() {
    // XXX: This was really rushed... D:
    err := r.req.ParseForm()
    if err != nil {
        serr := fmt.Sprintf("%+v", err)
        http.Error(r.w, serr, http.StatusNotFound)
        log.Print(serr)
        return
    }

    tok := r.req.PostFormValue("token")
    err = mtcareers.SaveAuthentication(tok)
    if err != nil {
        serr := fmt.Sprintf("%+v", err)
        http.Error(r.w, serr, http.StatusNotFound)
        log.Print(serr)
        return
    }

    r.w.Header().Set("Content-Type", "text/html")
    r.w.WriteHeader(http.StatusOK)
    r.w.Write([]byte("<h1>Token saved successfully!</h1>"))
}

// ServeHTTP is called by Go's http package whenever a new HTTP request arrives
func (p *pageServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    log.Printf("New request from %+v: %s %+v", req.RemoteAddr, req.Method, req.URL.Path)
    var r request = request {
        p: p,
        w: w,
        req: req,
        path: req.URL.Path,
    }

    if r.path[0] == '/' {
        r.path = r.path[1:]
    }

    switch (req.Method) {
    case "GET":
        r.get()
    case "POST":
        r.post()
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

// StartServer starts a new server in the requested port
func StartServer(port int) error {
    var err error

    if srv.httpServer != nil {
        return errors.New("Server is already running!")
    }

    srv.httpServer = &http.Server {
        Addr: fmt.Sprintf(":%d", port),
        Handler: &srv,
    }

    srv.userPage = template.New("")
    _, err = srv.userPage.Parse(config.Get().PageTemplate(pageTemplate))
    if err != nil {
        return errors.Wrap(err, "Failed to parse template page")
    }

    srv.renewPage = template.New("")
    _, err = srv.renewPage.Parse(renewTemplate)
    if err != nil {
        return errors.Wrap(err, "Failed to parse renew server template page")
    }

    go func() {
        log.Print("Waiting...")
        srv.httpServer.ListenAndServe()
    } ()

    return nil
}

// StopServer stops the currently executing server
func StopServer() {
    if srv.httpServer != nil {
        srv.httpServer.Close()
        srv.httpServer = nil
    }
}
