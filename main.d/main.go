/*
 * @author    Emmanuel Kofi Bessah
 * @email     ekbessah@uew.edu.gh
 * @created   Sat Jun 30 2018 11:41:21
 * @copyright © 2018 University of Education, Winneba
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	juggler "bekinsoft.com/ds-juggler"
	"bekinsoft.com/eritars/iam-ms/datasource"
	"bekinsoft.com/eritars/srm-ms/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/employee", testHandler)
	r.HandleFunc("/employee/{id}", testHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)

	// str := `{"where": {"username": "ekbessah"}}`

	filter, err := juggler.GetFilterParamMap(r)
	// filter, err := juggler.GetFilterParamMapFromJSONString(str)
	filter.Valid = true
	fmt.Println(filter)
	if err != nil {
		panic(err)
	}
	var mod model.StudMain

	if filter.Valid {

		tx := database().Begin()
		tx.LogMode(true)

		tx, err := juggler.FilterQuery(filter, tx)
		if err != nil {
			panic(err)
		}

		tx.Find(&mod)
		fmt.Println(mod)
	}

	// fmt.Println(r.Method)
	// fmt.Println(r.Body)
	// fmt.Println(r.RemoteAddr)
	// fmt.Println(r.URL)
	// fmt.Println(r.Response)

	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	// io.WriteString(w, `{"alive": true}`)
	json.NewEncoder(w).Encode(mod)
}

func database() *gorm.DB {
	var DBClient datasource.IGormClient
	DBClient = &datasource.GormClient{}
	DBClient.SetupDBForTest("root@192.168.99.100:26257/eritars_srm?sslmode=disable&charset=utf8&parseTime=True")
	// DBClient.SetupDBForTest("root@41.74.82.219:26257/eritars_iam?sslmode=disable&charset=utf8&parseTime=True")
	db := DBClient.GetInstance()

	return db
}
