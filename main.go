package main

import (
	"Pretest-OCA/app"
	"encoding/json"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	en_translations "gopkg.in/go-playground/validator.v10/translations/en"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

var Cars = []app.Car{}

func GetByType (w http.ResponseWriter, r *http.Request) {
	mapp := map[string][]app.Car{} //create mapp for divide mpv and suv
	SUV := []app.Car{} //create array/slice suv
	MPV := []app.Car{} //create slice mpv

	for index, item := range Cars {
		if item.Tipe == "SUV" { //append to suv
			SUV = append(SUV, Cars[index])
		} else if item.Tipe == "MPV" { //append to mpv
			MPV = append(MPV, Cars[index])
		}
	}
	mapp["SUV"] = SUV
	mapp["MPV"] = MPV

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(mapp)
}

func GetCarsByType (w http.ResponseWriter, r *http.Request) {
	Total := 0

	//get form input
	tipe := r.FormValue("tipe")

	for _, item := range Cars {
		if item.Tipe == tipe {
			Total += 1
		}
	}

	//make map for display result
	mapp := map[string]int{}
	mapp["jumlah_kendaraan"] = Total

	w.Header().Set("Content-Type", "application/json")
	_ =json.NewEncoder(w).Encode(mapp)
}

func GetCarsByColors (w http.ResponseWriter, r *http.Request) {
	// make map for fetch the result
	mapp := map[string][]string{}

	//hold plat that requested by client
	var plat []string

	//get client input
	inputs := r.FormValue("warna")

	for _, item := range Cars {
		if item.Warna == inputs {
			plat = append(plat, item.Plat_nomor)
		}
	}
	mapp["plat_nomor"] = plat

	w.Header().Set("Content-Type", "application/json")
	_ =json.NewEncoder(w).Encode(mapp)
}

func Create (w http.ResponseWriter, r *http.Request) {
	//get client input
	plat_nomor := r.FormValue("plat_nomor")
	warna := r.FormValue("warna")
	tipe := r.FormValue("tipe")

	//translator for validator input
	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	check := app.Car{
		Plat_nomor:    plat_nomor,
		Warna:         warna,
		Tipe:          tipe,
		Tanggal_masuk: time.Now(),
	}

	err := v.Struct(check)

	if err != nil {
		for _, errors := range err.(validator.ValidationErrors) {
			w.WriteHeader(400)
			_ = json.NewEncoder(w).Encode(errors.Translate(trans))
		}
	} else {
		//add to array Cars
		Cars = append(Cars, app.Car{
			Plat_nomor:    plat_nomor,
			Warna:         warna,
			Tipe:          tipe,
			Tanggal_masuk: time.Now(),
		})

		//fetch result to display
		mapp := map[string]string{}

		//parse time to string
		t := time.Now()
		newT := t.Format("2006-01-02 15:04")

		mapp["plat_nomor"] = plat_nomor
		mapp["parking_lot"] = "A1"
		mapp["tanggal_masuk"] = newT

		w.Header().Set("Content-Type", "application/json")
		_ =json.NewEncoder(w).Encode(mapp)
	}
}

func Out (w http.ResponseWriter, r *http.Request) {
	//get deleted item
	mapp := map[string]string{}

	//get biaya parkir
	biaya := 0

	//get client request
	plat := r.FormValue("plat_nomor")

	for index, item := range Cars {
		if item.Plat_nomor == plat {
			charge := 0
			if item.Tipe == "SUV" { //get biaya clasification
				charge = 25000
			} else if item.Tipe == "MPV" {
				charge = 35000
			}

			duration := time.Since(item.Tanggal_masuk)
			if int(duration.Hours()) < 2 {
				biaya = charge
			} else if int(duration.Hours()) >= 2 {
				multiply := (20 * charge)/100
				biaya = charge + (multiply * int(duration.Hours()))
			}

			//parse time to string
			masuk := item.Tanggal_masuk.Format("2006-01-02 15:04")
			keluar := time.Now().Format("2006-01-02 15:04")


			mapp["plat_nomor"] = item.Plat_nomor
			mapp["tanggal_masuk"] = masuk
			mapp["tanggal_keluar"] = keluar
			mapp["jumlah_bayar"] = strconv.Itoa(biaya)
			Cars = append(Cars[:index], Cars[index+1:]...) //delete array
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ =json.NewEncoder(w).Encode(mapp)
}

func main () {
	// dummy data
	Cars = append(Cars, app.Car{
		Plat_nomor:    "B 123 45",
		Warna:         "Biru",
		Tipe:          "SUV",
		Tanggal_masuk: time.Date(2020, 04, 15, 17, 0,0, 0, time.Local),
	})
	Cars = append(Cars, app.Car{
		Plat_nomor:    "B 123 44",
		Warna:         "Hitam",
		Tipe:          "SUV",
		Tanggal_masuk: time.Date(2020, 04, 15, 18, 0,0, 0, time.Local),
	})
	Cars = append(Cars, app.Car{
		Plat_nomor:    "B 123 33",
		Warna:         "Biru",
		Tipe:          "MPV",
		Tanggal_masuk: time.Date(2020, 04, 15, 20, 0,0, 0, time.Local),
	})
	Cars = append(Cars, app.Car{
		Plat_nomor:    "B 123 31",
		Warna:         "Putih",
		Tipe:          "MPV",
		Tanggal_masuk: time.Date(2020, 04, 15, 21, 0,0, 0, time.Local),
	})

	// init new router
	r := mux.NewRouter()

	// router list
	//get all cars in parking lot by SUV or MPV
	r.HandleFunc("/cars", GetByType).Methods("GET")
	//get all total cars by its type
	r.HandleFunc("/carTotalType", GetCarsByType).Methods("GET")
	//get all cars by colors
	r.HandleFunc("/carPlat", GetCarsByColors).Methods("GET")
	//regist new car in parking lot
	r.HandleFunc("/create", Create).Methods("POST")
	//cars out from parking lots
	r.HandleFunc("/out", Out).Methods("POST")

	//serve a server
	_ = http.ListenAndServe(":8000", r)
}
