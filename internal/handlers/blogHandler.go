package handlers

import (
	"context"
	"net/http"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	request "github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/auth"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetAllBlogs(c *gin.Context) {
	options := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).Find(context.TODO(), bson.M{}, options)
	if err != nil {
		logrus.Errorf("Error getting all blogs: GetAllBlogs API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	var blogs []models.Blog
	if err := cursor.All(context.TODO(), &blogs); err != nil {
		logrus.Errorf("Error decoding the blogs: GetAllBlogs API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Blogs fetched successfully", blogs)
}

func GetBlogById(c *gin.Context) {
	blogId := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		logrus.Errorf("Invalid blog id: GetBlogById API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid blog id", nil)
		return
	}

	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).FindOne(context.TODO(), bson.M{"_id": objectId})
	if result.Err() != nil {
		logrus.Errorf("Blog not found: GetBlogById API: %v", result.Err())
		response.HandleResponse(c, http.StatusNotFound, "Blog not found", nil)
		return
	}

	var blog models.Blog
	if err := result.Decode(&blog); err != nil {
		logrus.Errorf("Error decoding the blog: GetBlogById API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Blog fetched successfully", blog)
}

func UpdateBlog(c *gin.Context) {}

func DeleteBlog(c *gin.Context) {}

func SearchBlogs(c *gin.Context) {
	query := c.Query("query")

	searchStage := bson.D{
		{Key: "$search", Value: bson.D{
			{Key: "index", Value: "text"},
			{Key: "text", Value: bson.D{
				{Key: "query", Value: query},
				{Key: "path", Value: "title"},
			}},
		}},
	}

	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).Aggregate(context.TODO(), mongo.Pipeline{searchStage})
	if err != nil {
		logrus.Errorf("Error searching blogs: SearchBlogs API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	var blogs []models.Blog
	if err := cursor.All(context.TODO(), &blogs); err != nil {
		logrus.Errorf("Error decoding the blogs: SearchBlogs API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Blogs fetched successfully", blogs)
}

func GetBlogsByUser(c *gin.Context) {
	// get the data from context
	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	options := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).Find(context.TODO(), bson.M{"authorId": decodeUser.ID}, options)
	if err != nil {
		logrus.Errorf("Error getting all blogs: GetAllBlogs API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	var blogs []models.Blog
	if err := cursor.All(context.TODO(), &blogs); err != nil {
		logrus.Errorf("Error decoding the blogs: GetAllBlogs API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Blogs fetched successfully", blogs)
}
