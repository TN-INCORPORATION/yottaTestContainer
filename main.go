package main

import (
	"cErr"
	"context"
	"dataObject"
	"dbUtil"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"lang.yottadb.com/go/yottadb"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type A struct {
	S string
}

func main() {
	defer func() {
		fmt.Println("Main defer, before call yottadb.Exit()")
		yottadb.Exit()
	}()

	router := mux.NewRouter()
	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/cifs", CreateCIF).Methods("POST")
	router.HandleFunc("/pics", CreatePic).Methods("POST")
	// router.HandleFunc("/cifs/{id}", updateCIF).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8010", router))
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API server for YT test is ready!")
	fmt.Println("Endpoint Hit: home")
}

func CreateCIF(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	start := time.Now()
	var cif dataObject.CIF
	err := json.NewDecoder(r.Body).Decode(&cif)
	if err != nil {
		fmt.Println("error : ", err)
		roundTrip := time.Now().Sub(start).Nanoseconds() / 1000000
		w.Header().Set("x-roundtrip", strconv.Itoa(int(roundTrip)))
		fmt.Fprintf(w, "format is wrong")
		return
	}

	// init database connection
	dbCon := dataObject.GetCon(64, 1,64, 256, 64,yottadb.YDB_MAX_ERRORMSG)
	defer dbCon.Free()

	// under TP

	cusErr := dbUtil.WrapTransaction(ctx, func(c context.Context, u uint64) cErr.Error {

		var tptoken uint64
		//if tptoken, ok := c.Value("YOTTATOKEN").(uint64); ok {
		//	tptoken = tptoken
		//} else {
		//	tptoken = yottadb.NOTTP
		//}
		tptoken = u

		id, err := dataObject.GetNextSeq(tptoken, dbCon)
		if err != nil {
			return err
		}

		// this logic is for test application panic, but still continue
		if cif.Firstname == "XXX" && cif.Lastname == "XXX" {
			fmt.Println("App panic - nil pointer reference = not panic")
			var a *A
			fmt.Println(a.S)
		}

		// this logic is for test application panic
		// todo : it still make server hang, so comment it out for now. we need to detect it and exit application
		if cif.Firstname == "YYY" && cif.Lastname == "YYY" {
			//fmt.Println("App panic - index out of range = panic")
			//exceed := [2]int{1, 2}
			//fmt.Println(exceed[len(cif.Firstname)])
		}

		// insert
		cif.ID = id
		err = cif.Save(tptoken, dbCon)
		if err != nil {
			return err
		}

		return nil
	})

	if cusErr != nil {
		fmt.Println(cusErr)
		roundTrip := time.Now().Sub(start).Nanoseconds() / 1000000
		w.Header().Set("x-roundtrip", strconv.Itoa(int(roundTrip)))
		fmt.Fprintf(w, "CIF not Created")
		return
	}

	roundTrip := time.Now().Sub(start).Nanoseconds() / 1000000
	w.Header().Set("x-roundtrip", strconv.Itoa(int(roundTrip)))
	// fmt.Fprintf(w, "CIF was Created : " + strconv.Itoa(cif.ID))
	json.NewEncoder(w).Encode(cif)
}

func CreatePic(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	picKey := rand.Intn(1000)
	pic := dataObject.Pic{ID: picKey, Data: dataObject.PicData}
	ctx := context.Background()

	start := time.Now()

	// init database connection
	dbCon := dataObject.GetCon(64, 1,64, 500000, 64,yottadb.YDB_MAX_ERRORMSG)
	defer dbCon.Free()

	// under TP
	cusErr := dbUtil.WrapTransaction(ctx, func(c context.Context, u uint64) cErr.Error {

		var tptoken uint64
		//if tptoken, ok := c.Value("YOTTATOKEN").(uint64); ok {
		//	tptoken = tptoken
		//} else {
		//	tptoken = yottadb.NOTTP
		//}
		tptoken = u

		// insert
		err := pic.Save(tptoken, dbCon)
		if err != nil {
			return err
		}

		return nil
	})

	if cusErr != nil {
		fmt.Println(cusErr)
		roundTrip := time.Now().Sub(start).Nanoseconds() / 1000000
		w.Header().Set("x-roundtrip", strconv.Itoa(int(roundTrip)))
		fmt.Fprintf(w, "Pic not Created")
		return
	}

	roundTrip := time.Now().Sub(start).Nanoseconds() / 1000000
	w.Header().Set("x-roundtrip", strconv.Itoa(int(roundTrip)))
	fmt.Fprintf(w, "Pic was Created : " + strconv.Itoa(pic.ID))
	// json.NewEncoder(w).Encode(pic)
}

//func updateCIF(w http.ResponseWriter, r *http.Request) {
//
//	start := time.Now()
//	params := mux.Vars(r)
//	var cif dataObject.CIF
//	err := json.NewDecoder(r.Body).Decode(&cif)
//	if err != nil {
//		fmt.Println("error : ", err)
//		roundTrip := time.Now().Sub(start).Nanoseconds() / 1000000
//		w.Header().Set("x-roundtrip", strconv.Itoa(int(roundTrip)))
//		fmt.Fprintf(w, "format is wrong")
//		return
//	}
//	cif.ID,_ = strconv.Atoi(params["id"])
//
//	// init database connection
//	var dbCon dataObject.YottaCon
//	dbCon.Alloc(64, 1,64, 256, 0,128)
//	defer dbCon.Free()
//
//	// under TP
//	tptoken := yottadb.NOTTP
//	err = yottadb.TpE(tptoken, &dbCon.ErrStr, func(tptoken uint64, errstrp *yottadb.BufferT) int32 {
//		// insert
//		err = cif.Save(tptoken, dbCon)
//		if handleError(err) { return int32(yottadb.ErrorCode(err)) }
//
//		return 0
//	}, "CS", []string{})
//
//	if err != nil {
//		roundTrip := time.Now().Sub(start).Nanoseconds() / 1000000
//		w.Header().Set("x-roundtrip", strconv.Itoa(int(roundTrip)))
//		fmt.Fprintf(w, "CIF not Updated")
//		return
//	}
//
//	roundTrip := time.Now().Sub(start).Nanoseconds() / 1000000
//	w.Header().Set("x-roundtrip", strconv.Itoa(int(roundTrip)))
//	fmt.Fprintf(w, "CIF was Updated")
//}



