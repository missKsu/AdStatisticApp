package main

import(
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/fatih/structs"
	"net/http"
	"log"
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"time"
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
	AdvertId int
	Type string
	Text string
	Title string
	Href string
	//AdGroupName string
	CampaignName string
}

type AdvertInternal struct {
	AdvertId int `json:"AdvertId"`
	Type string `json:"Type"`
	Text string `json:"Text"`
	Title string `json:"Title"`
	Href string `json:"Href"`
	//AdGroupId int `json:"AdGroupId"`
	CampaignId int `json:"CampaignId"`
}

type CampaignExternal struct {
	CampaignId int
	Login string
	Email string
	Name string
	Type string
	StartDate string
	EndDate string
}

type CampaignInternal struct {
	Login string `json:"Login"`
	CampaignId int `json:"CampaignId,string,omitempty"`
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
	a.Router.HandleFunc("/updateInfoAboutUserForCampaigns",a.UpdateInfoAboutUserForCampaigns).Methods("POST")
	a.Router.HandleFunc("/updateType",a.UpdateType).Methods("POST")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Start AdStatisticsApp!")
	LogMessage("Request for info.")
}

func (a *App) Users(w http.ResponseWriter, r *http.Request) {
	LogMessage("Request for users' service.")
	add:= "http://localhost:3001"+r.URL.String()
	fmt.Println(add)
	err := r.ParseForm()
	if CheckError(err){
		fmt.Fprint(w, err.Error(),"\n")
		return
	}
	res := SetUser(ParseUser(r))
	val, err := json.Marshal(res)
	CheckError(err)
	response := SendAPIQuery(add,val)
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
	LogMessage("Request for adverts' service.")
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
	LogMessage("Try to parse advert request.")
	var result AdvertExternal
	var id int
	var err error
	if r.Form.Get("AdvertId") != ""{
		id,err = strconv.Atoi(r.Form.Get("AdvertId"))
		CheckError(err)
	}else{
		id = 0
	}
	result.AdvertId = id
	result.Type = r.Form.Get("Type")
	result.Text = r.Form.Get("Text")
	result.Title = r.Form.Get("Title")
	result.Href = r.Form.Get("Href")
	result.CampaignName = r.Form.Get("CampaignName")
	return result
}

func SetAdvert(r AdvertExternal) AdvertInternal {
	LogMessage("Try to set internal advert request.")
	var result AdvertInternal
	result.AdvertId = r.AdvertId
	result.Type = r.Type
	result.Text = r.Text
	result.Title = r.Title
	result.Href = r.Href
	result.CampaignId = GetCampaignIdByName(r.CampaignName)
	return result
}

func GetCampaignIdByName(CampaignName string) int{
	LogMessage("Try to get campaign id by name.")
	add:= "http://localhost:3002/campaigns/getCampaignIdByName"
	var que CampaignInternal
	que.Name = CampaignName
	val, err := json.Marshal(que)
	LogMessage(string(val))
	CheckError(err)
	response := SendAPIQuery(add,val)
	err = json.Unmarshal(response, &que)
	return que.CampaignId
}

func (a *App) Campaigns(w http.ResponseWriter, r *http.Request) {
	LogMessage("Request for campaigns' service.")
	add:= "http://localhost:3002"+r.URL.String()
	err := r.ParseForm()
	CheckError(err)
	res := SetCampaign(ParseCampaign(r))
	val, err := json.Marshal(res)
	CheckError(err)
	fmt.Println(string(val))
	response := SendAPIQuery(add,val)
	fmt.Fprint(w, string(response),"\n")
}

func ParseCampaign(r *http.Request) CampaignExternal {
	var result CampaignExternal
	var id int
	var err error
	if r.Form.Get("CampaignId") != ""{
		id,err = strconv.Atoi(r.Form.Get("CampaignId"))
		CheckError(err)
	}else{
		id = 0
	}
	result.CampaignId = id
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
	result.CampaignId = r.CampaignId
	result.Login = r.Login
	result.Email = r.Email
	result.Name = r.Name
	result.Type = r.Type
	result.StartDate = r.StartDate
	result.EndDate = r.EndDate
	return result
}

func CheckError(err error) bool{
	if err!=nil{
		LogMessage("Error: "+err.Error()+".")
		return true
	}else{
		return false
	}
}

func SendAPIQuery(add string, body []byte) []byte{
	LogMessage("Send api query to "+add+".")
	var response *http.Response
	v := url.Values{}
	bodyArray := make(map[string]string)
	err := json.Unmarshal(body,&bodyArray)
	CheckError(err)
	for key,value := range bodyArray {
		if value != ""{
			v.Add(key,value)
		}
		
	}
	response, err = http.PostForm(add,v)
	CheckError(err)
	defer response.Body.Close()
	result,err:= ioutil.ReadAll(response.Body)
	CheckError(err)
	return result
}

func (a *App) UpdateInfoAboutUserForCampaigns(w http.ResponseWriter, r *http.Request) {

}

func (a *App) UpdateType(w http.ResponseWriter, r *http.Request) {

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