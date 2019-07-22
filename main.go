package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/klauspost/reedsolomon"
	// Pkg "github.com/shammishailaj/rserasure/pkg"
)

// Encode - Actual function for reed-solomon encoder
func Encode(fileName, outDir string, dataShards, parShards int) {
	log.Printf("Encode() called with fileName = %s, outDir = %s, dataShards = %d, parShards = %d", fileName, outDir, dataShards, parShards)

	if (dataShards + parShards) > 256 {
		fmt.Fprintf(os.Stderr, "Error: sum of data and parity shards cannot exceed 256\n")
		os.Exit(1)
	}
	fname := fileName

	// Create encoding matrix.
	enc, err := reedsolomon.New(dataShards, parShards)
	checkErr(err)

	log.Printf("Opening file %s...", fname)
	b, err := ioutil.ReadFile(fname)
	checkErr(err)

	// Split the file into equally sized shards.
	shards, err := enc.Split(b)
	checkErr(err)
	log.Printf("File split into %d data+parity shards with %d bytes/shard.\n", len(shards), len(shards[0]))

	// Encode parity
	err = enc.Encode(shards)
	checkErr(err)

	// Write out the resulting files.
	dir, file := filepath.Split(fname)
	if outDir != "" {
		dir = outDir
	} else {
		log.Printf("Output directory no specified")
	}

	for i, shard := range shards {
		outfn := fmt.Sprintf("%s.%d", file, i)

		log.Println("Writing to", outfn)
		err = ioutil.WriteFile(filepath.Join(dir, outfn), shard, 0644)
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Printf("Error: %s", err.Error())
		os.Exit(2)
	}
}

// Decode - Actual function for reed-solomon decoder
func Decode(baseFilePath, targetDir string, dataShards, parShards int) {
	fname := baseFilePath

	// Create matrix
	enc, err := reedsolomon.New(dataShards, parShards)
	checkErr(err)

	// Create shards and load the data.
	shards := make([][]byte, dataShards+parShards)
	for i := range shards {
		infn := fmt.Sprintf("%s.%d", fname, i)
		log.Printf("Opening file %s...", infn)
		shards[i], err = ioutil.ReadFile(infn)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
		}
	}

	// Verify the shards
	ok, err := enc.Verify(shards)
	if ok {
		log.Println("No reconstruction needed")
	} else {
		log.Println("Verification failed. Reconstructing data")
		err = enc.Reconstruct(shards)
		if err != nil {
			log.Printf("Reconstruct failed - %#v", err)
			os.Exit(1)
		}
		ok, err = enc.Verify(shards)
		if !ok {
			log.Println("Verification failed after reconstruction, data likely corrupted.")
			os.Exit(1)
		}
		checkErr(err)
	}

	// Join the shards and write them
	outfn := targetDir
	if outfn == "" {
		outfn = fname
	}

	fmt.Println("Writing data to", outfn)
	f, err := os.Create(outfn)
	checkErr(err)

	// We don't know the exact filesize.
	err = enc.Join(f, shards, len(shards[0])*dataShards)
	checkErr(err)
}

// EncodeRS - HTTP Handler for reed-solomon encoder
func EncodeRS(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		source := r.FormValue("source")
		target := r.FormValue("targetdir")
		dataShards, _ := strconv.Atoi(r.FormValue("datashards"))
		parShards, _ := strconv.Atoi(r.FormValue("parityshards"))
		log.Printf("source = %s, target = %s, data shards = %d, parity shards  = %d", source, target, dataShards, parShards)
		if source == "" || target == "" || dataShards < 1 || parShards < 1 {
			http.Error(w, "{error:\"Must provide source, target, datashards and parityshards values\"}", http.StatusBadRequest)
		} else {
			Encode(source, target, dataShards, parShards)
			w.Write([]byte("{\"message\":\"POSTed data sent to Encoder\",\"status\":\"success\""))
		}
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
	}
}

// DecodeRS - HTTP Handler for reed-solomon decoder
func DecodeRS(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		source := r.FormValue("source")
		target := r.FormValue("targetdir")
		dataShards, _ := strconv.Atoi(r.FormValue("datashards"))
		parShards, _ := strconv.Atoi(r.FormValue("parityshards"))
		log.Printf("source = %s, target = %s, data shards = %d, parity shards  = %d", source, target, dataShards, parShards)
		if source == "" || target == "" || dataShards < 1 || parShards < 1 {
			http.Error(w, "{error:\"Must provide source, target, datashards and parityshards values\"}", http.StatusBadRequest)
		} else {
			Decode(source, target, dataShards, parShards)
			w.Write([]byte("{\"message\":\"POSTed data sent to Decoder\",\"status\":\"success\""))
		}
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
	}
}

func main() {
	http.HandleFunc("/encode", EncodeRS)
	http.HandleFunc("/decode", DecodeRS)
	http.ListenAndServe(":3000", nil)
}
