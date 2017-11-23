package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type chronoInfo struct {
	Chrono int
	Point  int
}

type resultInfo struct {
	Number   int
	Laps     int
	Dur      string
	Dur1     string
	Dur2     string
	Dur3     string
	Dur4     string
	Dur5     string
	Name     string
	Category string
	Team     string
}

var (
	glMtx       *sync.Mutex
	glDbs       map[string]*sql.DB
	glChronoTpl *template.Template
	glResultTpl *template.Template
)

func openDb(name string) *sql.DB {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatal("Cannot open DB:", err)
	}
	db.Exec(sqlCreateDb)
	return db
}

func parseIntParam(m url.Values, k string) (int, error) {
	v, ok := m[k]
	if !ok || len(v) != 1 {
		return 0, errors.New("wrong parameter")
	}
	r, err := strconv.Atoi(v[0])
	if err != nil {
		return 0, errors.New("wrong parameter")
	}
	return r, nil
}

func copyHandler(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()
	c, err := parseIntParam(m, "c")
	if err != nil {
		log.Println(err)
		return
	}
	name := fmt.Sprintf("chrono-%d.db", c)
	glMtx.Lock()
	db, ok := glDbs[name]
	if ok && db != nil {
		db.Close()
		glDbs[name] = nil
	}
	buff, err := ioutil.ReadFile(name)
	if err != nil {
		log.Println(err)
		return
	}
	glMtx.Unlock()
	src := bytes.NewReader(buff)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	io.Copy(w, src)
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()
	c, err := parseIntParam(m, "c")
	if err != nil {
		log.Println(err)
		return
	}
	p, err := parseIntParam(m, "p")
	if err != nil {
		log.Println(err)
		return
	}
	name := fmt.Sprintf("chrono-%d.db", c)
	glMtx.Lock()
	db, ok := glDbs[name]
	if !ok || db == nil {
		db = openDb(name)
		glDbs[name] = db
	}
	glMtx.Unlock()
	if r.Body == nil {
		log.Println("requst body is null")
		return
	}
	body := csv.NewReader(r.Body)
	sql := &bytes.Buffer{}
	num := 0
	io.WriteString(sql, "INSERT INTO Passes(number, point, pass) VALUES ")
	for {
		record, err := body.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}
		if num > 0 {
			io.WriteString(sql, ",")
		}
		io.WriteString(sql, "(")
		io.WriteString(sql, record[0])
		io.WriteString(sql, ",")
		io.WriteString(sql, strconv.Itoa(p))
		io.WriteString(sql, ",")
		io.WriteString(sql, record[1])
		io.WriteString(sql, ")")
		num++
	}
	db.Exec(sql.String())
}

func chronoHandler(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()
	c, err := parseIntParam(m, "c")
	if err != nil {
		log.Println(err)
		return
	}
	p, err := parseIntParam(m, "p")
	if err != nil {
		log.Println(err)
		return
	}
	info := &chronoInfo{Chrono: c, Point: p}
	glChronoTpl.Execute(w, info)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()
	c, err := parseIntParam(m, "c")
	if err != nil {
		log.Println(err)
		return
	}
	name := fmt.Sprintf("chrono-%d.db", c)
	glMtx.Lock()
	db, ok := glDbs[name]
	if !ok || db == nil {
		db = openDb(name)
		glDbs[name] = db
	}
	glMtx.Unlock()
	rows, err := db.Query(sqlResults)
	if err != nil {
		log.Println(err)
		return
	}
	info := make([]resultInfo, 0)
	for rows.Next() {
		var number, laps, dur, dur1, dur2, dur3, dur4, dur5 int
		var name, category, team string
		err = rows.Scan(&number, &laps, &dur, &dur1, &dur2, &dur3, &dur4, &dur5, &name, &category, &team)
		if err != nil {
			log.Println(err)
			break
		}
		info = append(info, resultInfo{Number: number, Laps: laps, Dur: millisToTime(dur), Dur1: millisToTime(dur1),
			Dur2: millisToTime(dur2), Dur3: millisToTime(dur3), Dur4: millisToTime(dur4), Dur5: millisToTime(dur5),
			Name: name, Category: category, Team: team})
	}
	rows.Close()
	glResultTpl.Execute(w, info)
}

func millisToTime(millis int) string {
	seconds := (millis / 1000) % 60
	minutes := ((millis / (1000 * 60)) % 60)
	hours := ((millis / (1000 * 60 * 60)) % 24)
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, (millis % 1000))
}

func init() {
	var err error
	glMtx = &sync.Mutex{}
	glDbs = make(map[string]*sql.DB)
	glChronoTpl, err = template.ParseFiles("chrono.tpl")
	if err != nil {
		log.Fatal("Cannot read chrono.tpl!")
	}
	glResultTpl, err = template.ParseFiles("result.tpl")
	if err != nil {
		log.Fatal("Cannot read result.tpl!")
	}
}

func main() {
	http.HandleFunc("/chrono", chronoHandler)
	http.HandleFunc("/sync", syncHandler)
	http.HandleFunc("/copy", copyHandler)
	http.HandleFunc("/result", resultHandler)
	http.ListenAndServe(":8080", nil)
}
