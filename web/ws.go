package web

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"umutsevdi/pengui/sys"

	"log"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type wsRequest struct {
	Type    string `json:"req"`
	Message string `json:"msg"`
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	log.Println("ws::connected#", conn.LocalAddr().String())
	if err != nil {
		log.Println("error:", err)
		return
	}
	go func() {
		defer conn.Close()
		_, op, _ := wsutil.ReadClientData(conn)
		err = wsutil.WriteServerMessage(conn, op, []byte("{\"ready\":true}"))
		started := false
		_, err := sys.NewShell()

		if err != nil {
			log.Println(conn.LocalAddr().String(), err)
			return
		}
		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				log.Println(conn.LocalAddr().String(), err)
				break
			}
			var request wsRequest
			err = json.Unmarshal(msg, &request)
			if err != nil {
				log.Println(conn.LocalAddr().String(), err)
				continue
			}
			if request.Type == "fetch" {

			} else if request.Type == "in" {

			}

			if started && string(msg) != "stream#next" {
				log.Println(conn, "error: Invalid text")
				break
			} else if !started && string(msg) == "stream#begin" {
				started = true
			}
			data, err := json.Marshal(sys.StreamResource())
			if err != nil {
				log.Println(conn.LocalAddr().String(), err)
			}
			err = wsutil.WriteServerMessage(conn, op, data)
			if err != nil {
				log.Println(conn.LocalAddr().String(), err)
			}
		}
		log.Println(conn.LocalAddr().String(), "client_exit#")
	}()
}

func ServeTerm(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	log.Println("ws/term::connected#", conn.LocalAddr().String())
	if err != nil {
		log.Println("error:", err)
		return
	}
	go func() {
		defer conn.Close()
		msg, op, _ := wsutil.ReadClientData(conn)
		log.Println(conn.LocalAddr().String(), string(msg))
		if err != nil {
			log.Println(conn.LocalAddr().String(), err)
		}
		err = wsutil.WriteServerMessage(conn, op, "{}")
		if err != nil {
			log.Println(conn.LocalAddr().String(), err)
		}
		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				log.Println(conn.LocalAddr().String(), err)
				break
			}
			log.Println(conn.LocalAddr().String(), string(msg))
			if len(string(msg)) == 0 {
				continue
			}
			log.Println("before bash")
			stdout, err := exec.Command().S.Exec(string(msg))
			log.Println("after bash")
			if err != nil {
				log.Println(conn.LocalAddr().String(), err.Error())
				err = wsutil.WriteServerMessage(conn, op, []byte("Error while interpreting command."))
				continue
			}
			log.Println(conn.LocalAddr().String(), ".", string(msg), "->", stdout)
			err = wsutil.WriteServerMessage(conn, op, []byte("> "+stdout))
			if err != nil {
				log.Println(conn.LocalAddr().String(), err)
			}
		}
		log.Println(conn.LocalAddr().String(), "client_exit#")
		bash.Exit()

	}()

}
