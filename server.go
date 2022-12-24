package main

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pterm/pterm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func CreateAndStartServer(ipAddress string, port int, serverInfo *pterm.SpinnerPrinter) {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/{type:[a-z]+}/v{version:[.,0-9]+}/{update_name}", handleRequest)
	http.Handle("/", rtr)
	serverAddr := fmt.Sprintf("%v:%d", ipAddress, port)
	ScreenDisplayHelp(serverAddr)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: rtr,
	}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverInfo.Fail("HTTP server error: ", err)
		}
	}()

	serverInfo.Success("Server Started")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		serverInfo.Fail("HTTP shutdown error")
		log.Fatalf("Err : %v", err)
	}
	serverInfo.Warning("Server Stopped")
	os.Exit(1)
}
func handleRequest(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	updateName := params["update_name"]
	filetype := params["type"]
	version := params["version"]
	if strings.HasSuffix(updateName, ".md5") {
		filename, _ := strings.CutSuffix(updateName, ".md5")

		filename = filetype + version + ".tar.gz"

		if Exists("update/" + filename) {
			res.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(res, generate_md5("update/"+filename))
			return
		}
	} else {
		filename := updateName
		filename = filetype + version + ".tar.gz"
		if Exists("update/" + filename) {
			fileBytes, err := ioutil.ReadFile("update/" + filename)
			if err != nil {
				panic(err)
			}
			res.WriteHeader(http.StatusOK)
			res.Header().Set("Content-Type", "application/octet-stream")
			res.Write(fileBytes)
			ScreenNewLog(fmt.Sprintf("Request Time: %s : Remote Address: %s ", time.Now().Format("2006-01-02 15:04:05"), req.RemoteAddr), Info)
			return
		}
	}
	fmt.Fprint(res, "Not Found")

}

func Exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return false
}

func generate_md5(filename string) string {
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		panic(err)
	}

	hashfile := fmt.Sprintf("%x", hash.Sum(nil))
	return hashfile
}
