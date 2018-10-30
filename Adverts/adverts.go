package main

import(
	"fmt"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"os"
	"strconv"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
}

type AdvertExternal struct {
	AdvertId int
	Type string
	Text string
	Title string
	Href string
	CampaignId int
	//AdGroupId int `json:"AdGroupId"`
}

type AdvertInternal struct {
	AdvertId int `json:"AdvertId"`
	Type int `json:"Type"`
	Text string `json:"Text"`
	Title string `json:"Title"`
	Href string `json:"Href"`
	CampaignId int `json:"CampaignId"`
	//AdGroupId int `json:"AdGroupId"`
}

type Response struct {
	Result Result `json:"Result,omitempty"`
	Error string `json:"Error,omitempty"`
}

type Result struct {
	AddResult AdvertInternal `json:"AddResult,omitempty"`
	UpdateResult AdvertInternal `json:"UpdateResult,omitempty"`
	Advert AdvertInternal `json:"User,omitempty"`
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
	a.Router.HandleFunc("/getAllAdvertsForCampaignById",a.GetAllAdvertsForCampaignById).Methods("POST")
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

func (a *App) GetAllAdvertsForCampaignById(w http.ResponseWriter, r *http.Request) {
	
}

func (a *App) AdvertParseForm(w http.ResponseWriter,r *http.Request) (bool,AdvertInternal){
	var res AdvertInternal
	err := r.ParseForm()
	if CheckError(err){
		fmt.Fprint(w, err.Error())
		return false,res
	}
	res = SetAdvert(ParseAdvert(r))
	return true,res
}

func ParseAdvert(r *http.Request) AdvertExternal {
	var result AdvertExternal
	id,err := strconv.Atoi(r.Form.Get("AdvertId"))
	CheckError(err)
	result.AdvertId = id
	result.Type = r.Form.Get("Type")
	result.Text = r.Form.Get("Text")
	result.Title = r.Form.Get("Title")
	result.Href = r.Form.Get("Href")
	id, err = strconv.Atoi(r.Form.Get("CampaignId"))
	CheckError(err)
	result.CampaignId = id
	return result
}

func SetAdvert(r AdvertExternal) AdvertInternal {
	var result AdvertInternal
	result.AdvertId = r.AdvertId
	result.Type = SetType(r.Type)
	result.Text = r.Text
	result.Title = r.Title
	result.Href = r.Href
	result.CampaignId = r.CampaignId
	return result
}

func SetType(TypeName string) int{
	var result int
	switch TypeName{
		case "TEXT_AD":
			result = 1
		case "MOBILE_APP_AD":
			result = 2
		case "DYNAMIC_TEXT_AD":
			result = 3
		case "IMAGE_AD":
			result = 4
		case "CPM_BANNER_AD":
			result = 5
	}
	return result
}

func CheckError(err error) bool{
	if err!=nil{
		return true
	}else{
		return false
	}
}

func CreateResponseError(message string) []byte {
	var result Response
	result.Error = message
	b, err := json.Marshal(result)
	CheckError(err)
	return b
}

func LogMessage(message string) {
	file := "logs.txt"
	current_time := time.Now().Local()
	_, err := os.Stat(file)
    if err != nil { 
    	_, err := os.Create(file)
    	if err != nil {
			fmt.Println(("Error: "+err.Error()))
		}
    }
    f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
		fmt.Println(("Error: "+err.Error()))
	}
	defer f.Close()

	_, err = f.WriteString("\n["+current_time.Format("Mon Jan 2 15:04:05 2006")+"] "+message)
	CheckError(err)
	if err != nil {
		fmt.Println(("Error: "+err.Error()))
	}
}