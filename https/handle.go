package https

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/qq2383/tunnel/config"
	"github.com/qq2383/tunnel/process"
	"github.com/qq2383/tunnel/server"
	"github.com/qq2383/tunnel/user"
)

var (
	Err400 = fmt.Errorf("400")
	Err404 = fmt.Errorf("404")
)

type Handler struct {
	root embed.FS
	r    *http.Request
	w    http.ResponseWriter
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r = r
	h.w = w

	path := r.URL.Path
	log.Printf("handle path: %s\n", path)

	defer func() {
		r.Body.Close()
	}()

	var buf []byte
	var err error
	extp := strings.LastIndex(path, ".")
	if extp != -1 {
		ext := path[extp:]
		if strings.ToLower(ext) != ".html" {
			ct := mime.TypeByExtension(ext)
			w.Header().Set("Content-Type", ct)
			buf, err = h.load(path[1:])
			if err != nil {
				err = Err404
			}
		}
	} else {
		s := NerSession(w, r)
		switch r.Method {
		case http.MethodPost:
			buf, err = h.post(path, s)
		case http.MethodGet:
			buf, err = h.get(path, s)
			if err != nil {
				err = Err404
			}
		default:
			err = Err400
		}
		w.Header().Set("Content-Type", "charset=utf-8")
	}

	switch err {
	case Err404:
		w.WriteHeader(http.StatusNotFound)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(buf)

}

func (h *Handler) get(path string, s *Session) ([]byte, error) {
	if path == "/" {
		return h.load("def.html")
	}

	if logon, ok := s.data["logon"]; ok && logon.(bool) {
		switch path {
		case "/mng":
			return h.load("mng.html")
		}
	} else {
		return h.load("def.html")
	}
	return nil, Err404
}

func (h *Handler) load(path string) ([]byte, error) {
	//root, _ := os.Getwd()
	//path = fmt.Sprintf("%s/html/%s", root, path)
	//bs, _ := os.ReadFile(path)

	return h.root.ReadFile(path)
}

func (h *Handler) login() *ResultData {
	r := h.r
	rd := &ResultData{}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err == nil {
			pwd := r.FormValue("pwd")
			cnf := config.Get("Http").(*Config)
			if pwd == cnf.Passwd {
				rd.SetOk([]any{true})
			}
		}
	}
	return rd
}

func (h *Handler) post(path string, s *Session) ([]byte, error) {
	var rd *ResultData

	switch path {
	case "/login":
		rd = h.login()
		if rd.Rows != 0 {
			s.SetAttr("logon", true)
		}
	case "/get/tunnel":
		rd = h.getTunnel()
	case "/post/tunnel":
		rd = h.postTunnel()
	case "/get/https":
		rd = h.getHttps()
	case "/post/https":
		rd = h.postHttps()
	case "/post/passwd":
		rd = h.postPasswd()
	case "/get/user":
		rd = h.getUser()
	case "/post/user":
		rd = h.postUser()
	case "/post/play":
		rd = h.postPlay()
	}

	return rd.ToJson()
}

func (h *Handler) getTunnel() *ResultData {
	var rd ResultData
	data := make([]any, 0)
	b := make(map[string]any)

	cnf := config.Get("Server").(*server.Config)
	b["tport"] = cnf.Port
	data = append(data, b)
	rd.SetOk(data)
	return &rd
}

func (h *Handler) postTunnel() *ResultData {
	var rd ResultData
	r := h.r
	err := r.ParseForm()
	if err == nil {
		_port := r.FormValue("tport")
		port, err := strconv.Atoi(_port)

		cnf := config.Get("Server").(*server.Config)
		if err == nil && port != cnf.Port {
			cnf.Port = port
			err = config.Write()
			if err == nil {
				rd.Rows = 1
			}
		}
	}
	return &rd
}

func (h *Handler) getHttps() *ResultData {
	var rd ResultData
	data := make([]any, 0)
	b := make(map[string]any)

	cnf := config.Get("Http").(*Config)
	b["hport"] = cnf.Port
	data = append(data, b)
	rd.SetOk(data)
	return &rd
}

func (h *Handler) postHttps() *ResultData {
	var rd ResultData
	r := h.r
	err := r.ParseForm()
	if err == nil {
		_port := r.FormValue("hport")
		port, err := strconv.Atoi(_port)

		cnf := config.Get("Http").(*Config)
		if err == nil && port != cnf.Port {
			cnf.Port = port
			err = config.Write()
			if err == nil {
				rd.Rows = 1
			}
		}
	}
	return &rd
}

func (h *Handler) postPasswd() *ResultData {
	var rd ResultData
	r := h.r
	err := r.ParseForm()
	if err == nil {
		npwd := r.FormValue("npwd")
		cpwd := r.FormValue("cpwd")
		if cpwd == "" || npwd == "" {
			rd.Rows = 1
		} else if npwd == cpwd {
			cnf := config.Get("Http").(*Config)
			cnf.Passwd = cpwd
			err := config.Write()
			if err != nil {
				rd.Rows = 3
			}
		} else {
			rd.Rows = 2
		}
	}
	return &rd
}

func (h *Handler) getUser() *ResultData {
	users := user.Get();
	var rd ResultData
	data := make([]any, len(users))
	for i, u := range users {
		data[i] = u
	}
	rd.SetOk(data)
	return &rd
}

func (h *Handler) postUser() *ResultData {
	r := h.r
	err := r.ParseForm()
	if err == nil {
		t := r.FormValue("t")
		o := r.FormValue("o")
		if t == "r" {
			user.Remove(o)
		} else if t == "s" {
			u := r.FormValue("u")
			p := r.FormValue("p")
			user.Modify(o, u, p)
		}
	}
	return h.getUser()
}

func (h *Handler) postPlay() *ResultData {
	var rd ResultData
	rd.Rows = 1

	r := h.r
	err := r.ParseForm()
	if err == nil {
		t := r.FormValue("type")
		switch t {
		case "hstart":
			process.Start("https")
		case "hstop":
			process.Stop("https")
		case "hrestart":
			process.Restart("https", time.Second*5)
		case "hstat":
			process.Status("https")
		case "tstart":
			process.Start("tunnel")
		case "tstop":
			process.Stop("tunnel")
		case "trestart":
			process.Restart("tunnel", time.Second*5)
		case "tstat":
			process.Status("tunnel")
		case "astat":
			re := process.Statusall()
			rd.SetOk([]any{re})
		default:
			rd.Rows = 0
		}
	}
	return &rd
}

type ResultData struct {
	Rows int
	Data []any
}

func (rd *ResultData) SetError() {
	rd.Rows = 0
	rd.Data = nil
}

func (rd *ResultData) SetOk(data []any) {
	rd.Rows = len(data)
	rd.Data = data
}

func (rd *ResultData) ToJson() ([]byte, error) {
	return json.Marshal(rd)
}
