package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	product struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Price float32 `json:"price"`
	}

	payment struct {
		ID int `json:"id"`
		Amount string `json:"amount"`
	}
)

var (
	products = map[int]*product{}
	payments = map[int]*payment{}
	seq   = 1
	lock  = sync.Mutex{}
	productsPath = "/products/:id"
)

func createProduct(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := &product{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	products[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func createPayment(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := &payment{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	payments[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func getProduct(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, products[id])
}

func updateProduct(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := new(product)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	products[id].Name = u.Name
	products[id].Price = u.Price
	return c.JSON(http.StatusOK, products[id])
}

func deleteProduct(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	delete(products, id)
	return c.NoContent(http.StatusNoContent)
}

func getAllProducts(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, products)
}

func main() {
	e := echo.New()
	
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/products", getAllProducts)
	e.POST("/products", createProduct)
	e.GET(productsPath, getProduct)
	e.PUT(productsPath, updateProduct)
	e.DELETE(productsPath, deleteProduct)
	e.POST("/payments", createPayment)

	e.Logger.Fatal(e.Start(":1323"))
}