package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const PORT = ":1234"

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := "Thanks for visiting!\n"
	fmt.Fprintf(w, "%s", Body)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get telephone
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)
	telephone := paramStr[2]
	err := deleteEntry(telephone)
	if err != nil {
		fmt.Println(err)
		Body := err.Error() + "\n"
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", Body)
		return
	}

	Body := telephone + " deleted!\n"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", Body)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := list()
	fmt.Fprintf(w, "%s", Body)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := fmt.Sprintf("Total entries: %d\n", len(data))
	fmt.Fprintf(w, "%s", Body)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 5 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not enough arguments: "+r.URL.Path)
		return
	}

	name := paramStr[2]
	surname := paramStr[3]
	tel := paramStr[4]

	t := strings.ReplaceAll(tel, "-", "")
	if !matchTel(t) {
		fmt.Println("Not a valid telephone number:", tel)
		return
	}

	temp := &Entry{Name: name, Surname: surname, Tel: t}
	err := insert(temp)
	if err != nil {
		Body := "Failed to add record\n"
		w.WriteHeader(http.StatusNotModified) // With this status the ResponseWriter does not output message
		// w.WriteHeader(http.StatusOK)  // With this status the output works
		fmt.Fprintf(w, "%s", Body)
	} else {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		Body := "New record added successfully\n"
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", Body)
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Get telephone
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	var Body string
	telephone := paramStr[2]
	t := search(telephone)
	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		Body = "Could not be found: " + telephone + "\n"
	} else {
		w.WriteHeader(http.StatusOK)
		Body = t.Name + " " + t.Surname + " " + t.Tel + "\n"
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)
	fmt.Fprintf(w, "%s", Body)
}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	var tmpFileName string

	//create temporary file name
	f, err := os.CreateTemp("", "data*.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	tmpFileName = f.Name()

	defer os.Remove(tmpFileName)

	// Save data to file
	err = saveCSVFile(tmpFileName)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Cannot create: " + tmpFileName)
		return
	}

	fmt.Println("Serving: ", tmpFileName)

	http.ServeFile(w, r, tmpFileName)

	// Sleep few seconds to get the file
	time.Sleep(10 * time.Second)
}
