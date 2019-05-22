package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/org/unhookd/helm"
	"github.com/org/unhookd/lib"
)

func instastageServerHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize)

	jsonRequestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println("Error reading request body", err)
		return
	}

	values := string(jsonRequestBody)

	release := r.URL.Query().Get("release")
	if !lib.VerifyQueryStringParams(release) {
		http.Error(w, "Error reading query string param release", 500)
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if !lib.VerifyQueryStringParams(namespace) {
		namespace = release
	}

	chart := r.URL.Query().Get("chart")
	if !lib.VerifyQueryStringParams(chart) {
		http.Error(w, "Error reading query string param chart", 500)
		return
	}

	async := r.URL.Query().Get("async")
	shouldWait := !lib.VerifyQueryStringParams(async)

	dryRun := r.URL.Query().Get("dryRun")
	shouldDryRun := lib.VerifyQueryStringParams(dryRun)

	helmOut, _, err := helm.HelmBinUpgradeInstall(release, namespace, chart, "instastage", values, shouldWait, shouldDryRun)

	log.Println("--------------------------------")
	log.Println("Values:")
	log.Println(values)
	log.Println(string(helmOut))
	log.Println("++++++++++++++++++++++++++++++++")

	if err != nil {
		http.Error(w, fmt.Sprintf("Error running helm upgrade: %v %v", string(helmOut), err.Error()), 500)
		return
	}

	w.WriteHeader(201)
	w.Write(helmOut)
}

func instastageServer(runHelm bool) {
	helm.Setup(runHelm)

	secretPath := lib.GetEnv("SECRET_PATH", "/")
	http.HandleFunc(secretPath, instastageServerHandler)
	http.HandleFunc("/status", healthCheck)
	listen := "0.0.0.0:8080"
	fmt.Println("Starting unhookD ... Please standby", listen)

	srv := &http.Server{
		Addr:         listen,
		ReadTimeout:  601 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	srv.ListenAndServe()
}
