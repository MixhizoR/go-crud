package controllers

import (
	"github.com/MixhizoR/go-crud/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var users = []models.User{
	{
		ID:    "1",
		Name:  "John Doe",
		Email: "john.doe@example.com",
	},
}

func GetUsers(c *gin.Context) {
	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	for _, u := range users {
		if u.ID == id {
			c.JSON(200, u)
			return
		}
	}
	c.JSON(404, gin.H{"error": "User not found"})
}

func CreateUser(c *gin.Context) {
	var user models.User
	user.ID = uuid.New().String()
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	users = append(users, user)
	c.JSON(201, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for i, u := range users {
		if u.ID == id {
			user.ID = id
			users[i] = user
			c.JSON(200, user)
			return
		}
	}
	c.JSON(404, gin.H{"error": "User not found"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(204, nil)
			return
		}
	}
	c.JSON(404, gin.H{"error": "User not found"})
}
