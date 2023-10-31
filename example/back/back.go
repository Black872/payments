package back

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Records struct {
	RecordsList     []string
	NumberOfRecords int
}

func Start() {
	http.HandleFunc("/example", viewHandler)
	http.HandleFunc("/example/new", newHandler)
	http.HandleFunc("/example/create", createHandler)
	err := http.ListenAndServe(":8090", nil)
	log.Fatal(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("front/index.html")
	check(err)
	input := readFile("data/records.txt")
	records := Records{
		RecordsList:     input,
		NumberOfRecords: len(input),
	}
	err = html.Execute(w, records)
	check(err)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("front/new.html")
	check(err)
	err = html.Execute(w, nil)
	check(err)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	record := r.FormValue("record")
	opts := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile("data/records.txt", opts, 0664)
	check(err)
	defer func() {
		check(file.Close())
	}()

	_, err = fmt.Fprintln(file, record)
	check(err)
	http.Redirect(w, r, "/example", http.StatusFound)
}

func readFile(fileName string) (list []string) {
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return list
	}
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	return list
}
