package main

import(
	"fmt"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"time"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
}

type UserExternal struct {
	Name string
	Login string
	Description string
	Date string
}

type UserInternal struct {
	Name string `json:"Name,omitempty"`
	Login string `json:"Login,omitempty"`
	Description string `json:"Description,omitempty"`
	Date string `json:"Date,omitempty"`
}

type Response struct {
	Result Result `json:"Result,omitempty"`
	Error string `json:"Error,omitempty"`
}

type Result struct {
	AddResult UserInternal `json:"AddResult,omitempty"`
	UpdateResult UserInternal `json:"UpdateResult,omitempty"`
	User UserInternal `json:"User,omitempty"`
}

func main(){
	var router App
	router.Initialize("root","ksCnhtkjr_97","AdStatistics")
	api := router.Router.PathPrefix("/AdStatisticsApp/users").Subrouter()
	api.HandleFunc("/",router.StartPage).Methods("POST")
	api.HandleFunc("/add",router.Add).Methods("POST")
	api.HandleFunc("/get",router.Get).Methods("POST")
	api.HandleFunc("/update",router.Update).Methods("POST")
	log.Fatal(http.ListenAndServe(":3001", router.Router))
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
	a.Router = a.Router.PathPrefix("/AdStatisticsApp/users").Subrouter()
	a.Router.HandleFunc("/",a.StartPage).Methods("POST")
	a.Router.HandleFunc("/add",a.Add).Methods("POST")
	a.Router.HandleFunc("/get",a.Get).Methods("POST")
	a.Router.HandleFunc("/update",a.Update).Methods("POST")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) StartPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "By this service you can get info about users!")
	return
}

func (a *App) Add(w http.ResponseWriter, r *http.Request) {
	check,res := a.UserParseForm(w,r)
	if (!check){
		return
	}

	if res.Login != "" && res.Name != ""{
		if CkeckForInsertInDB(res.Login,a.DB){
			result := InsertNewUserInDB(res, a.DB)
			b, err := json.Marshal(result)
			CheckError(err)
			fmt.Fprint(w, string(b))
		}else{
			b := CreateResponseError("User with same login exist.")
			fmt.Fprint(w, string(b))
			return
		}
	}else{
		b := CreateResponseError("Error occured with parameters.")
		fmt.Fprint(w, string(b))
		return
	}
}

func (a *App) Get(w http.ResponseWriter, r *http.Request) {
	check,res := a.UserParseForm(w,r)
	if (!check){
		return
	}
	if res.Login !=""{
		if (!CkeckForInsertInDB(res.Login,a.DB)){
			result := GetUserFromDB(res.Login, a.DB)
			b, err := json.Marshal(result)
			CheckError(err)
			fmt.Fprint(w, string(b))
		}else{
			b := CreateResponseError("Incorrect input data: no such login.")
			fmt.Fprint(w, string(b))
			return
		}
	}else{
		b := CreateResponseError("Incorrect input data: no login for get data.")
		fmt.Fprint(w, string(b))
		return
	}
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	check,res := a.UserParseForm(w,r)
	if (!check){
		return
	}
	if res.Login != ""{
		if (!CkeckForInsertInDB(res.Login,a.DB)){
			result := UpdateUserInDB(res, a.DB)
			b, err := json.Marshal(result)
			CheckError(err)
			fmt.Fprint(w, string(b))
		}else{
			b := CreateResponseError("Incorrect input data: no such login.")
			fmt.Fprint(w, string(b))
			return
		}
	}else{
		b := CreateResponseError("Incorrect input data: no login for changes.")
		fmt.Fprint(w, string(b))
		return
	}
}

func (a *App) UserParseForm(w http.ResponseWriter,r *http.Request) (bool,UserInternal){
	var res UserInternal
	err := r.ParseForm()
	if CheckError(err){
		fmt.Fprint(w, err.Error())
		return false,res
	}
	res = SetUser(ParseUser(r))
	return true,res
}

func ParseUser(r *http.Request) UserExternal {
	var result UserExternal
	result.Name = r.Form.Get("Name")
	result.Login = r.Form.Get("Login")
	result.Description = r.Form.Get("Description")
	return result
}

func CreateResponseError(message string) []byte {
	var result Response
	result.Error = message
	b, err := json.Marshal(result)
	CheckError(err)
	return b
}

func SetUser(r UserExternal) UserInternal {
	var result UserInternal
	result.Name = r.Name
	result.Login = r.Login
	result.Description = r.Description
	result.Date = time.Now().Format("2006.01.02")
	return result
}

func CheckError(err error) bool{
	if err!=nil{
		return true
	}else{
		return false
	}
}

func CkeckForInsertInDB(Login string, db *sql.DB) bool{
	query, err := db.Prepare(`
    			SELECT 
					count(login)
				FROM Users
				WHERE Login like ?`)
    CheckError(err)
    defer query.Close()

    queryres,err := query.Query(Login)
    CheckError(err)

    defer queryres.Close()

    var count string

    for queryres.Next(){
        queryres.Scan(&count)
    }

    if count!="0"{
    	return false
    }
	return true
}

func InsertNewUserInDB(res UserInternal, db *sql.DB) Response{
	var result Response
	query, err := db.Prepare(`
    			INSERT INTO Users SET 
					Name = ?,
					Login = ?,
					Description = ?,
					Date = ?`)
    CheckError(err)
    defer query.Close()

    _,err = query.Exec(res.Name,
    						   res.Login,
    						   res.Description,
    						   res.Date)
    CheckError(err)

    //defer queryres.Close()
    result.Result.AddResult = res
    return result
}

func UpdateUserInDB(res UserInternal, db *sql.DB) Response {
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

func GetUserFromDB(Login string, db *sql.DB) Response {
	var result Response
	var User UserInternal
	query, err := db.Prepare(`
    			SELECT 
					Name,
					Description,
					Date
				FROM Users
				WHERE Login like ?`)
    CheckError(err)
    defer query.Close()

    queryres,err := query.Query(Login)
    CheckError(err)

    defer queryres.Close()
    
    User.Login = Login
    for queryres.Next(){
        queryres.Scan(&User.Name,&User.Description,&User.Date)
    }

    result.Result.User = User
    return result
}

func UpdateQuery(res UserInternal) (string,int){
	if res.Name != ""{
		if res.Description != ""{
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