package database

import(
	"fmt"
	"database/sql"
	"time"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	Result Result `json:"Result,omitempty"`
	Error string `json:"Error,omitempty"`
}

type UserStruct struct {
	Name string `json:"Name,omitempty"`
	Login string `json:"Login,omitempty"`
	Description string `json:"Description,omitempty"`
	Date string `json:"Date,omitempty"`
}

type Result struct {
	AddResult UserStruct `json:"AddResult,omitempty"`
	UpdateResult UserStruct `json:"UpdateResult,omitempty"`
	User UserStruct `json:"User,omitempty"`
}

type Database struct{
	DB *sql.DB
}

func (db *Database) CkeckForInsertInDB(Login string) bool{
	LogMessage("Try to find user by login.")
	query, err := db.DB.Prepare(`
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

func (db *Database) InsertNewUserInDB(res UserStruct) Response{
	LogMessage("Try to insert user into DB.")
	var result Response
	query, err := db.DB.Prepare(`
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

    result.Result.AddResult = res
    return result
}

func (db *Database) UpdateUserInDB(res UserStruct) Response {
	LogMessage("Try to update user in DB.")
	var result Response
	stmt, ans := UpdateQuery(res)
	query, err := db.DB.Prepare(stmt)
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

    result.Result.UpdateResult = res
    return result
}

func (db *Database) GetUserFromDB(Login string) Response {
	LogMessage("Try to get user by login from DB.")
	var result Response
	var User UserStruct
	query, err := db.DB.Prepare(`
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

func CheckError(err error) bool{
	if err!=nil{
		LogMessage("Error: "+err.Error())
		return true
	}else{
		return false
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

func UpdateQuery(res UserStruct) (string,int){
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