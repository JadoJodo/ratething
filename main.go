package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	thing "github.com/jadojodo/ratething/ratething"
	_ "github.com/mattn/go-sqlite3"
)

const fileName = "sqlite.db"

func main() {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	thingRepository := thing.NewSQLiteRepository(db)

	if err := thingRepository.MigrateThings(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/things", func(c *gin.Context) {
		things, err := thingRepository.AllThings()

		if nil != err {
			log.Fatal(err)
		}

		c.IndentedJSON(http.StatusOK, things)
	})

	r.GET("/things/:thing", func(c *gin.Context) {
		thingId := c.Param("thing")

		thing, err := thingRepository.GetThingByUUID(thingId)

		if nil != err {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusOK, thing)
	})

	r.POST("/things", func(c *gin.Context) {
		var json thing.Thing
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		thing, err := thingRepository.CreateThing(json)

		if nil != err {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusCreated, thing)
	})

	r.Run()
}
