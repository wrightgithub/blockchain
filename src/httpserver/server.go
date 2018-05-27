package httpserver

import (
	"os"
	"net/http"
	"log"
	"time"
	"github.com/gorilla/mux"
	"encoding/json"
	"io"
	"model"
	"github.com/davecgh/go-spew/spew"
	"utils"
)

type Message struct {
	BPM int
}

func Run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(model.Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		utils.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	resultBlockchain, err := model.Run(m.BPM);
	if err != nil {
		utils.RespondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	spew.Dump(resultBlockchain)
	utils.RespondWithJSON(w, r, http.StatusCreated, model.GetLatestBlock())

}
