package main

import(
	"fmt"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
}

type AdvertEXternal struct {
	Type string
	Text string
	Title string
	Href string
	CampaignId int
	//AdGroupId int `json:"AdGroupId"`
}

func main(){
	var router App
	router.Initialize("root","ksCnhtkjr_97","AdStatistics")
	log.Fatal(http.ListenAndServe(":3003", router.Router))
}

func (a *App) Initialize(user, password, dbname string) {
	a.Router = mux.NewRouter().StrictSlash(true)
	a.InitializeRoutes()

	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) InitializeRoutes() {
	a.Router = a.Router.PathPrefix("/AdStatisticsApp/adverts").Subrouter()
	a.Router.HandleFunc("/",a.StartPage).Methods("POST")
	a.Router.HandleFunc("/add",a.Add).Methods("POST")
	a.Router.HandleFunc("/get",a.Get).Methods("POST")
	a.Router.HandleFunc("/update",a.Update).Methods("POST")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) StartPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "By this service you can add, update or get info about adverts!")
	return
}

func (a *App) Add(w http.ResponseWriter, r *http.Request) {

}

func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	
}

