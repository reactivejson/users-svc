// internal/app/handler.go
package app

import (
	"github.com/reactivejson/usr-svc/internal/domain"
	"github.com/reactivejson/usr-svc/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddUserHandler handles the request to add a new user.
func (s *UserService) AddUserHandler(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.AddUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "user added"})
}

// UpdateUserHandler handles the request to update an existing user.
func (s *UserService) UpdateUserHandler(c *gin.Context) {
	userID := c.Param("id")

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = userID

	err := s.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "user updated"})
}

// DeleteUserHandler handles the request to delete a user.
func (s *UserService) DeleteUserHandler(c *gin.Context) {
	userID := c.Param("id")

	err := s.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "user deleted"})
}

// GetUsersHandler handles the request to get a paginated list of users.
func (s *UserService) GetUsersHandler(c *gin.Context) {
	country := c.Query("country")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, pageSize, err := utils.ParsePageAndPageSize(pageStr, pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := s.GetUsers(country, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
