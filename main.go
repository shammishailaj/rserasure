package main

import (
	"log"
	"net/http"
	"strconv"

	Pkg "./pkg"
)

// EncodeRS - HTTP Handler for reed-solomon encoder
func EncodeRS(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		source := r.FormValue("source")
		log.Printf("source = %s", source)
		target := r.FormValue("targetdir")
		log.Printf("target = %s", target)
		dataShards, _ := strconv.Atoi(r.FormValue("datashards"))
		log.Printf("data shards = %s", dataShards)
		parShards, _ := strconv.Atoi(r.FormValue("parityshards"))
		log.Printf("parity shards  = %s", parShards)
		Pkg.Encode(source, target, dataShards, parShards)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		w.Write([]byte("POSTed data sent to Encoder"))
	}
	Pkg.Encode()
}

func main() {
	http.HandleFunc("/encode", EncodeRS)
	http.ListenAndServe(":3000", nil)
}
