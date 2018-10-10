package main

import(
	"fmt"
	"github.com/gorilla/mux"
	"github.com/fatih/structs"
	"net/http"
	"log"
	"encoding/json"
	"net/url"
	//"bytes"
	"io/ioutil"
)

type App struct {
	Router *mux.Router
}

type UserExternal struct {
	Name string
	Login string
	Description string
}

type UserInternal struct {
	Name string `json:"Name,omitempty"`
	Login string `json:"Login,omitempty"`
	Description string `json:"Description,omitempty"`
}

type AdvertExternal struct {
	Type string
	Text string
	Title string
	Href string
	//AdGroupName string
	CampaignName string
}

type AdvertInternal struct {
	Type string `json:"Type"`
	Text string `json:"Text"`
	Title string `json:"Title"`
	Href string `json:"Href"`
	//AdGroupId int `json:"AdGroupId"`
}

type CampaignExternal struct {
	Login string
	Email string
	Name string
	Type string
	StartDate string
	EndDate string
}

type CampaignInternal struct {
	ClientId int `json:"ClientId"`
	Email string `json:"Email"`
	Name string `json:"Name"`
	Type string `json:"Type"`
	StartDate string `json:"StartDate"`
	EndDate string `json:"EndDate"`
}

func main(){
	var router App
	router.Initialize()
	log.Fatal(http.ListenAndServe(":3000", router.Router))
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter().StrictSlash(true)
	a.InitializeRoutes()
}

func (a *App) InitializeRoutes() {
	a.Router = a.Router.PathPrefix("/AdStatisticsApp").Subrouter()
	a.Router.HandleFunc("/",a.Welcome).Methods("POST")
	a.Router.HandleFunc("/users/{*}",a.Users).Methods("POST")
	a.Router.HandleFunc("/adverts/{*}",a.Adverts).Methods("POST")
	a.Router.HandleFunc("/campaigns/{*}",a.Campaigns).Methods("POST")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Start AdStatisticsApp!")
}

func (a *App) Users(w http.ResponseWriter, r *http.Request) {
	add:= "http://localhost:3001"+r.URL.String()
	fmt.Println(add)
	err := r.ParseForm()
	if CheckError(err){
		fmt.Fprint(w, err.Error(),"\n")
		return
	}
	res := SetUser(ParseUser(r))

	response := SendAPIQuery(add,res)
	fmt.Fprint(w, string(response),"\n")
}

func ParseUser(r *http.Request) UserExternal {
	var result UserExternal
	result.Name = r.Form.Get("Name")
	result.Login = r.Form.Get("Login")
	result.Description = r.Form.Get("Description")
	return result
}

func SetUser(r UserExternal) UserInternal {
	var result UserInternal
	result.Name = r.Name
	result.Login = r.Login
	result.Description = r.Description
	return result
}

func (a *App) Adverts(w http.ResponseWriter, r *http.Request) {
	add:= "http://localhost:3002"+r.URL.String()
	err := r.ParseForm()
	CheckError(err)
	res := SetAdvert(ParseAdvert(r))
	val, err := json.Marshal(res)
	CheckError(err)
	response := SendAPIQuery(add,val)
	fmt.Fprint(w, string(response),"\n")
}

func ParseAdvert(r *http.Request) AdvertExternal {
	var result AdvertExternal
	result.Type = r.Form.Get("Type")
	result.Text = r.Form.Get("Text")
	result.Title = r.Form.Get("Title")
	result.CampaignName = r.Form.Get("CampaignName")
	return result
}

func SetAdvert(r AdvertExternal) AdvertInternal {
	var result AdvertInternal
	result.Type = r.Type
	result.Text = r.Text
	result.Title = r.Title
	return result
}

func (a *App) Campaigns(w http.ResponseWriter, r *http.Request) {
	add:= "http://localhost:3003"+r.URL.String()
	err := r.ParseForm()
	CheckError(err)
	res := SetCampaign(ParseCampaign(r))
	val, err := json.Marshal(res)
	CheckError(err)
	response := SendAPIQuery(add,val)
	fmt.Fprint(w, string(response),"\n")
}

func ParseCampaign(r *http.Request) CampaignExternal {
	var result CampaignExternal
	result.Login = r.Form.Get("Login")
	result.Email = r.Form.Get("Email")
	result.Name = r.Form.Get("Name")
	result.Type = r.Form.Get("Type")
	result.StartDate = r.Form.Get("StartDate")
	result.EndDate = r.Form.Get("EndDate")
	return result
}

func SetCampaign(r CampaignExternal) CampaignInternal {
	var result CampaignInternal
	result.Email = r.Email
	result.Name = r.Name
	result.Type = r.Type
	result.StartDate = r.StartDate
	result.EndDate = r.EndDate
	return result
}

func CheckError(err error) bool{
	if err!=nil{
		return true
	}else{
		return false
	}
}

func SendAPIQuery(add string, body interface{}) []byte{
	var response *http.Response
	v := url.Values{}
	for key,value := range structs.Map(body) {
		v.Add(key,value.(string))
	}
	response, err := http.PostForm(add,v)
	CheckError(err)
	defer response.Body.Close()
	result,err:= ioutil.ReadAll(response.Body)
	CheckError(err)
	return result
}