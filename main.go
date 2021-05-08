package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func goDotEnvVar(key string) string {
	err := godotenv.Load("db.env")
	if err != nil {
		log.Fatalf("error loading env")
	}
	return os.Getenv(key)
}

// var pgConnStr = goDotEnvVar("POSTGRES_URL")

var (
	db     *gorm.DB
	pgUser = goDotEnvVar("POSTGRES_USER")
	pgPW   = goDotEnvVar("POSTGRES_PASSWORD")
	pgHost = goDotEnvVar("POSTGRES_HOST")
	pgPort = goDotEnvVar("POSTGRES_PORT")
	pgDB   = goDotEnvVar("POSTGRES_DB")
)

func init() {
	var err error
	var numPort, err1 = strconv.Atoi(pgPort)
	if err1 != nil {
		log.Fatal(err1)
	}
	dsn := url.URL{
		User:     url.UserPassword(pgUser, pgPW),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", pgHost, numPort),
		Path:     pgDB,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	println(dsn.String())
	db, err = gorm.Open("postgres", dsn.String())
	if err != nil {
		panic("failed to connect database")
	}

	//Migrate the schema
	db.AutoMigrate(&todoModel{})
}

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}
	router.Run()
}

type (
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}

	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}

func fetchAllTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

func updateTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "NNo todo found!"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}
