package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	//Create / POST
	router.POST("/users", postUser)
	//Retrieve / GET
	router.GET("/users/:id", getUser)
	//Update
	router.PUT("/users/:id", updateUser)
	//Delete
	router.DELETE("users/:id", deleteUser)

	router.Start(":8080")
}

//FUNCTION
func getUser(c echo.Context) error {
	db := Connect()
	defer db.Close()

	idGet := c.Param("id")

	var id int
	var name string
	var age int
	var address string
	var password string

	fmt.Println("Get All Users")

	err := db.QueryRow("SELECT * from users WHERE id= ?", idGet).Scan(&id, &name, &age, &address, &password)
	if err != nil {
		log.Println(err)
		errorResponse(c.Response().Writer)
		return err
	}

	response := User{ID: id, Name: name, Password: password, Age: age, Address: address}
	return c.JSON(http.StatusOK, response)
}

func deleteUser(c echo.Context) error {
	db := Connect()
	defer db.Close()

	idDelete := c.Param("id")
	query := "DELETE FROM users WHERE id = ?"
	queryStatement, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err2 := queryStatement.Exec(idDelete)
	if err2 != nil {
		fmt.Println(err2)
		return err2
	}

	return c.JSON(http.StatusOK, "user dengan id = "+idDelete+"sudah dihapus")
}

func postUser(c echo.Context) error {
	db := Connect()
	defer db.Close()

	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}

	query := "INSERT INTO users(Name,Age,Address,Password) VALUES (?,?,?,?)"
	queryStatement, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer queryStatement.Close()

	queryResult, err2 := queryStatement.Exec(user.Name, user.Age, user.Address, user.Password)

	if err2 != nil {
		fmt.Println(err2)
		return err2
	}
	fmt.Println(queryResult.LastInsertId())

	return c.JSON(http.StatusCreated, user.Name)
}

func updateUser(c echo.Context) error {
	db := Connect()
	defer db.Close()

	idUpdate := c.Param("id")
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}

	query := "UPDATE users SET Name = ?, Age = ?, Address = ?, Password = ? WHERE id = ?"
	queryStatement, err := db.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer queryStatement.Close()

	queryResult, err2 := queryStatement.Exec(user.Name, user.Age, user.Address, user.Password, idUpdate)

	if err2 != nil {
		fmt.Println(err2)
		return err2
	}
	fmt.Println(queryResult.LastInsertId())

	return c.JSON(http.StatusOK, "user dengan id = "+idUpdate+"berhasil di update")
}

func successResponse(rw http.ResponseWriter) {
	var successResponse InsertResponse
	successResponse.Status = 200
	successResponse.Message = "Success"
	json.NewEncoder(rw).Encode(successResponse)
}

func errorResponse(rw http.ResponseWriter) {
	var errorResponse ErrorResponse
	errorResponse.Status = 404
	errorResponse.Message = "ERROR"
	json.NewEncoder(rw).Encode(errorResponse)
}

//CLASS
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

type InsertResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//DATABASE
func Connect() *sql.DB {
	fmt.Println("error CONNECT")
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/db_latihan_pbp?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	return db
}
