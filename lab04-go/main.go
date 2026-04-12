package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"category"`
}

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Cart struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("products.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&Category{}, &Product{}, &Cart{}); err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/products", getProducts)
	e.GET("/products/:id", getProduct)
	e.POST("/products", createProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)
	e.GET("/categories", getCategories)
	e.POST("/categories", createCategory)
	e.GET("/carts", getCarts)
	e.POST("/carts", createCart)

	e.Logger.Fatal(e.Start(":8080"))
}

func getProducts(c echo.Context) error {
	var products []Product
	if err := db.Preload("Category").Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}
	return c.JSON(http.StatusOK, products)
}

func getProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var product Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	return c.JSON(http.StatusOK, product)
}

func createProduct(c echo.Context) error {
	var input struct {
		Name       string  `json:"name"`
		Price      float64 `json:"price"`
		CategoryID uint    `json:"category_id"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	var category Category
	if err := db.First(&category, input.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "category not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	product := Product{Name: input.Name, Price: input.Price, CategoryID: input.CategoryID}
	if err := db.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	if err := db.Preload("Category").First(&product, product.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	return c.JSON(http.StatusCreated, product)
}

func updateProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var input struct {
		Name       string  `json:"name"`
		Price      float64 `json:"price"`
		CategoryID uint    `json:"category_id"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	var category Category
	if err := db.First(&category, input.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "category not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	product.Name = input.Name
	product.Price = input.Price
	product.CategoryID = input.CategoryID
	if err := db.Save(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	if err := db.Preload("Category").First(&product, product.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	return c.JSON(http.StatusOK, product)
}

func deleteProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	result := db.Delete(&Product{}, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
	}

	return c.NoContent(http.StatusNoContent)
}

func getCategories(c echo.Context) error {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}
	return c.JSON(http.StatusOK, categories)
}

func createCategory(c echo.Context) error {
	var input struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	category := Category{Name: input.Name}
	if err := db.Create(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	return c.JSON(http.StatusCreated, category)
}

func getCarts(c echo.Context) error {
	var carts []Cart
	if err := db.Preload("Product").Preload("Product.Category").Find(&carts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}
	return c.JSON(http.StatusOK, carts)
}

func createCart(c echo.Context) error {
	var input struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	var product Product
	if err := db.First(&product, input.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	cart := Cart{ProductID: input.ProductID, Quantity: input.Quantity}
	if err := db.Create(&cart).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	if err := db.Preload("Product").Preload("Product.Category").First(&cart, cart.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	return c.JSON(http.StatusCreated, cart)
}