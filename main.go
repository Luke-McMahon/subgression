package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	loadEnv()
	connectToDatabase()

	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUserByID)
	router.PATCH("/users/:id", updateUser)
	router.POST("/users", addUser)
	router.DELETE("/users/:id", deleteUser)

	router.Run("localhost:8080")
}

func loadEnv() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalf("Error loading .env file")
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	}
}

func connectToDatabase() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "subgression",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to DB!")
}

func getUsers(c *gin.Context) {
	var dbUsers []User
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.FullName, &u.BeltRank, &u.Degree); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		dbUsers = append(dbUsers, u)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, dbUsers)
}

func updateUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Endpoint not implemented"})
}

func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	dbUser := dbGetUserByID(id)
	if dbUser == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, dbUser)
}

func dbGetUserByID(ID int) *User {
	var dbUser User

	row := db.QueryRow("SELECT * FROM users WHERE id = ?", ID)
	if err := row.Scan(&dbUser.ID, &dbUser.FullName, &dbUser.BeltRank, &dbUser.Degree); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	return &dbUser
}

func addUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	_, err := db.Exec("INSERT INTO users (fullname, beltrank, degree) VALUES (?, ?, ?)", newUser.FullName, newUser.BeltRank, newUser.Degree)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

func deleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	res, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err == nil {
		count, err := res.RowsAffected()
		if count == 1 {
			getUsers(c)
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
}

type User struct {
	ID       string `json:"id,omitempty"`
	FullName string `json:"full_name"`
	BeltRank string `json:"belt_rank"`
	Degree   int    `json:"degree"`
}
