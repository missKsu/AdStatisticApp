package main

import(
	"fmt"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"strconv"
	"encoding/json"
	"os"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
}

type CampaignExternal struct {
	Login string
	CampaignId int
	Email string
	Name string
	Type string
	StartDate string
	EndDate string
}

type CampaignInternal struct {
	Login string `json:"Login,omitempty"`
	CampaignId int `json:"CampaignId,omitempty"`
	Email string `json:"Email,omitempty"`
	Name string `json:"Name,omitempty"`
	Type int `json:"Type,omitempty"`
	StartDate string `json:"StartDate,omitempty"`
	EndDate string `json:"EndDate,omitempty"`
}

type Response struct {
	Result Result `json:"Result,omitempty"`
	Error string `json:"Error,omitempty"`
}

type Result struct {
	AddResult CampaignInternal `json:"AddResult,omitempty"`
	UpdateResult CampaignInternal `json:"UpdateResult,omitempty"`
	Campaign CampaignInternal `json:"User,omitempty"`
}

func main(){
	var router App
	router.Initialize("root","ksCnhtkjr_97","AdStatistics")
	log.Fatal(http.ListenAndServe(":3002", router.Router))
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
	a.Router = a.Router.PathPrefix("/AdStatisticsApp/campaigns").Subrouter()
	a.Router.HandleFunc("/",a.StartPage).Methods("POST")
	a.Router.HandleFunc("/add",a.Add).Methods("POST")
	a.Router.HandleFunc("/get",a.Get).Methods("POST")
	a.Router.HandleFunc("/update",a.Update).Methods("POST")
	a.Router.HandleFunc("/getCampaignIdByName",a.GetCampaignIdByName).Methods("POST")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) StartPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "By this service you can add, update or get info about campaigns!")
	return
}

func (a *App) Add(w http.ResponseWriter, r *http.Request) {
	LogMessage("Request for add data.")
	check,res := a.CampaignParseForm(w,r)
	if (!check){
		return
	}

	if (CheckForInsertInDBLoginAndName(res.Login, res.Name, a.DB)==0){
		fmt.Println(res.Email,res.Login,res.Name,res.Type,res.StartDate,res.EndDate)
		if (res.Login != "" && res.Name != "" && res.Type != 0 && res.StartDate != "" && res.EndDate != ""){
			result := AddCampaignToDB(res, a.DB)
			b, err := json.Marshal(result)
			CheckError(err)
			fmt.Fprint(w, string(b))
			LogMessage("Success add data.")
			return
		}else{
			b := CreateResponseError("Incorrect input data.")
			fmt.Fprint(w, string(b))
			LogMessage("Error: incorrect input data.")
			return
		}
	}else{
		b := CreateResponseError("Incorrect login or name for add campaign.")
		fmt.Fprint(w, string(b))
		LogMessage("Error: incorrect login or name for add campaign.")
		return
	}
}

func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	LogMessage("Request for get data.")
	check,res := a.CampaignParseForm(w,r)
	if (!check){
		return
	}

	if res.CampaignId != 0 && res.Login != ""{
		result := GetCampaignFromDB(res, a.DB)
		b, err := json.Marshal(result)
		CheckError(err)
		fmt.Fprint(w, string(b))
		LogMessage("Success get data.")
		return
	}else{
		b := CreateResponseError("No id or login for get data.")
		fmt.Fprint(w, string(b))
		LogMessage("Error: no id or login for get data.")
		return
	}
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	LogMessage("Try to update campaign.")
	check,res := a.UserParseForm(w,r)
	if (!check){
		return
	}
	if res.CampaignId != 0{
		if (!CheckForInsertInDBCampaignId(res.CampaignId,a.DB)){
			result := UpdateCampaignInDB(res, a.DB)
			b, err := json.Marshal(result)
			CheckError(err)
			fmt.Fprint(w, string(b))
		}else{
			b := CreateResponseError("Incorrect input data: no such id.")
			fmt.Fprint(w, string(b))
			return
		}
	}else{
		b := CreateResponseError("Incorrect input data: no id for changes.")
		fmt.Fprint(w, string(b))
		return
	}
}

func (a *App) GetCampaignIdByName(w http.ResponseWriter, r *http.Request) {
	LogMessage("Request for get id by name.")
	check,res := a.CampaignParseForm(w,r)
	if (!check){
		return
	}

	if res.Name != ""{
		result := GetCampaignIdFromDB(res.Name, a.DB)
		b, err := json.Marshal(result)
		CheckError(err)
		fmt.Fprint(w, string(b))
		LogMessage("Success get id by name.")
		return
	}else{
		b := CreateResponseError("No name for get id.")
		fmt.Fprint(w, string(b))
		LogMessage("Error: no name for get id.")
		return
	}
}


func (a *App) CampaignParseForm(w http.ResponseWriter,r *http.Request) (bool,CampaignInternal){
	LogMessage("Try to parse request params.")
	var res CampaignInternal
	err := r.ParseForm()
	if CheckError(err){
		fmt.Fprint(w, err.Error())
		return false,res
	}
	res = SetCampaign(ParseCampaign(r))
	return true,res
}

func ParseCampaign(r *http.Request) CampaignExternal {
	var result CampaignExternal
	if r.Form.Get("CampaignId") != ""{
		id,err := strconv.Atoi(r.Form.Get("CampaignId"))
		CheckError(err)
		result.CampaignId = id
	}
	result.Name = r.Form.Get("Name")
	result.Login = r.Form.Get("Login")
	result.Email = r.Form.Get("Email")
	result.Type = r.Form.Get("Type")
	result.StartDate = r.Form.Get("StartDate")
	result.EndDate = r.Form.Get("EndDate")
	return result
}

func SetCampaign(r CampaignExternal) CampaignInternal {
	var result CampaignInternal
	result.CampaignId = r.CampaignId
	result.Name = r.Name
	result.Login = r.Login
	result.Email = r.Email
	result.Type = SetType(r.Type)
	result.StartDate = r.StartDate
	result.EndDate = r.EndDate
	return result
}

func SetType(TypeName string) int{
	var result int
	switch TypeName{
		case "TEXT_CAMPAIGN":
			result = 1
		case "MOBILE_APP_CAMPAIGN":
			result = 2
		case "DYNAMIC_TEXT_CAMPAIGN":
			result = 3
		case "CPM_BANNER_CAMPAIGN":
			result = 4
	}
	return result
}

func CheckForInsertInDBLoginAndName(Login string, Name string, db *sql.DB) int{
	query, err := db.Prepare(`
    			SELECT 
					CampaignId
				FROM Campaigns
				WHERE Login like ? AND Name like ?`)
    CheckError(err)
    defer query.Close()

    queryres,err := query.Query(Login,
    							Name)
    CheckError(err)

    defer queryres.Close()

    var id int

    for queryres.Next(){
        queryres.Scan(&id)
    }

    return id
}

func CheckForInsertInDBCampaignId(CampaignId int, db *sql.DB) bool{
	query, err := db.Prepare(`
    			SELECT 
					count(CampaignId)
				FROM Campaigns
				WHERE CampaignId like ?`)
    CheckError(err)
    defer query.Close()

    queryres,err := query.Query(CampaignId)
    CheckError(err)

    defer queryres.Close()

    var count int

    for queryres.Next(){
        queryres.Scan(&count)
    }

    if count == 0{
    	return false
    }else
    	return true
}

func CheckError(err error) bool{
	if err!=nil{
		LogMessage("Error: "+err.Error()+".")
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

func AddCampaignToDB (res CampaignInternal, db *sql.DB) Response{
	var result Response
	query, err := db.Prepare(`
    			INSERT INTO Campaigns SET 
					Name = ?,
					Login = ?,
					Email = ?,
					Type = ?,
					StartDate = ?,
					EndDate = ?`)
    CheckError(err)
    defer query.Close()

    _,err = query.Exec( res.Name,
    					res.Login,
    					res.Email,
    					res.Type,
    					res.StartDate,
    					res.EndDate)
    CheckError(err)

    //defer queryres.Close()
    result.Result.AddResult = res
    return result
}

func GetCampaignFromDB(res CampaignInternal, db *sql.DB) Response{
	var result Response
	var Campaign CampaignInternal
	query, err := db.Prepare(`
    			SELECT 
					Name,
					Email,
					Type,
					StartDate,
					EndDate
				FROM Campaigns
				WHERE Login like ? AND CampaignId like ?`)
    CheckError(err)
    defer query.Close()

    queryres,err := query.Query(res.Login, res.CampaignId)
    CheckError(err)

    defer queryres.Close()
    
    for queryres.Next(){
        queryres.Scan(&Campaign.Name, &Campaign.Email, &Campaign.Type, &Campaign.StartDate, &Campaign.EndDate)
    }

    Campaign.CampaignId = res.CampaignId

    result.Result.Campaign = Campaign
    return result
}

func GetCampaignIdFromDB(Name string, db *sql.DB) Response{
	var result Response
	var Campaign CampaignInternal
	query, err := db.Prepare(`
    			SELECT 
					CampaignId
				FROM Campaigns
				WHERE Name like ?`)
    CheckError(err)
    defer query.Close()

    queryres,err := query.Query(Name)
    CheckError(err)

    defer queryres.Close()
    
    for queryres.Next(){
        queryres.Scan(&Campaign.CampaignId)
    }

    Campaign.Name = Name

    result.Result.Campaign = Campaign
    return result
}

func UpdateCampaignInDB(res UserInternal, db *sql.DB) Response {
	LogMessage("Try to update campaign in DB.")
	var result Response
	stmt, ans := UpdateQuery(res)
	query, err := db.Prepare(stmt)
    CheckError(err)
    defer query.Close()

    switch ans{
    case 1:
    	_,err = query.Exec(res.Name,
    					   res.Login)
    case 2:
    	_,err = query.Exec(res.Name,
    					   res.Description,
    					   res.Login)
    case 3:
    	_,err = query.Exec(res.Description,
    					   res.Login)
    }
    CheckError(err)

    //defer queryres.Close()
    result.Result.UpdateResult = res
    return result
}

func UpdateQuery(res UserInternal) (string,int){
	if res.Name != ""{
		if res.Description == ""{
			return `UPDATE Users SET
						Name = ?
					WHERE Login = ?`,1
		}else{
			return `UPDATE Users SET
						Name = ?,
						Description = ?
					WHERE Login = ?`,2
		}
	}else{
		return `UPDATE Users SET
					Description = ?
				WHERE Login = ?`,3
	}
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