package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Car struct {
	ID             int    `json:"id"`
	NameOfMark     string `json:"name_of_mark"`
	NameOfModel    string `json:"name_of_model"`
	Mileage        int    `json:"mileage"`
	NumberOfOwners int    `json:"number_of_owners"`
}

type Response struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
}

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password = 123 dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	db.AutoMigrate(&Car{})
}

func GetHandler(c echo.Context) error {
	var cars []Car
	if err := db.Find(&cars).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Could not get the car message",
		})
	}

	return c.JSON(http.StatusOK, &cars)
}

func GetCarByIDHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Wrong ID",
		})
	}

	var car Car
	if err := db.First(&car, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, Response{
				Status: 404,
				Text:   "Not found: Car not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, Response{
			Status: 500,
			Text:   "Internal server error: Could not retrieve the car",
		})
	}

	return c.JSON(http.StatusOK, car)
}

func PostHandler(c echo.Context) error {
	var car Car
	if err := c.Bind(&car); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Wrong car message",
		})
	}

	if err := db.Create(&car).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Could not create the car message",
		})
	}

	return c.JSON(http.StatusCreated, Response{
		Status: 201,
		Text:   "Created",
	})
}

func PatchHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Wrong ID",
		})
	}

	var updatedCarMessage Car
	if err := c.Bind(&updatedCarMessage); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Invalid input",
		})
	}

	if err := db.Model(&Car{}).Where("id = ?", id).Updates(updatedCarMessage).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Could not update the message",
		})
	}

	return c.JSON(http.StatusNoContent, Response{
		Status: 204,
		Text:   "No content",
	})
}

func PutHandler(c echo.Context) error {
	var cars []Car
	if err := c.Bind(&cars); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Invalid input",
		})
	}

	for _, car := range cars {
		if err := db.Model(&Car{}).Where("id = ?", car.ID).Updates(car).Error; err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status: 400,
				Text:   "Bad request: Could not update the car message",
			})
		}
	}

	return c.JSON(http.StatusOK, Response{
		Status: 200,
		Text:   "OK",
	})
}

func DeleteHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Wrong ID",
		})
	}

	var carMessage Car
	if err := db.Delete(&carMessage, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status: 400,
			Text:   "Bad request: Could not delete the message",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status: 204,
		Text:   "No content",
	})
}

func main() {
	initDB()

	e := echo.New()

	e.GET("/cars", GetHandler)
	e.GET("/cars/:id", GetCarByIDHandler)
	e.POST("/cars", PostHandler)
	e.PATCH("/cars/:id", PatchHandler)
	e.PUT("/cars", PutHandler)
	e.DELETE("/cars/:id", DeleteHandler)

	e.Start(":8080")
}
