package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
}

type Recipe struct {
	Id           string    `json:"id,omitempty"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt,omitempty"`
}

func NewRecipeHandler(c *gin.Context) {

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	recipe.PublishedAt = time.Now()
	recipe.Id = xid.New().String()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)

}

func ListRecipeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].Id == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}
	recipe.Id = id
	recipe.PublishedAt = time.Now()

	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].Id == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted"})
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.Run()
}
