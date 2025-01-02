package handlers

import (
	"net/http"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	request "github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/auth"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateBlog(c *gin.Context) {
	var body request.CreateBlogRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Invalid request body: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err := utils.ValidateRequest(body)
	if err != nil {
		logrus.Errorf("Error validating the request body: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// get the data from context
	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	result, err := models.CreateBlog(&models.Blog{
		Title:           body.Title,
		Body:            body.Body,
		IsBlogPublished: body.IsBlogPublished,
		AuthorID:        decodeUser.ID,
	})
	if err != nil {
		logrus.Errorf("Error creating blog: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusCreated, "Blog created successfully", result)
}

func GetAllBlogs(c *gin.Context) {}

func GetBlogById(c *gin.Context) {}

func UpdateBlog(c *gin.Context) {}

func DeleteBlog(c *gin.Context) {}

func SearchBlogs(c *gin.Context) {}

func GetBlogsByUserId(c *gin.Context) {}
