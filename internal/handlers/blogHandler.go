package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/services"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	request "github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/auth"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateBlog(c *gin.Context) {
	var body request.CreateBlogRequest
	if err := c.ShouldBind(&body); err != nil {
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

	// get the image from the request
	imageFile, err := c.FormFile("image")
	if err != nil {
		logrus.Errorf("Error getting image: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Error getting image", nil)
		return
	}

	// upload the image locally
	err = c.SaveUploadedFile(imageFile, "./assets/uploads/"+imageFile.Filename)
	if err != nil {
		logrus.Errorf("Error uploading image: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	// validate the image size and type
	err = utils.ValidateImageFile(imageFile)
	if err != nil {
		logrus.Errorf("Error validating image: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	relativeFilePath := "./assets/uploads/" + imageFile.Filename
	imageUrl, err := services.UploadImageToCloudinary(relativeFilePath)
	if err != nil {
		logrus.Errorf("Error uploading image to cloudinary: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	result, err := models.CreateBlog(&models.Blog{
		Title:           body.Title,
		Body:            body.Body,
		IsBlogPublished: body.IsBlogPublished,
		AuthorID:        decodeUser.ID,
		ImageURL:        imageUrl,
	})
	if err != nil {
		logrus.Errorf("Error creating blog: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	// update user collection
	userObjectId, err := primitive.ObjectIDFromHex(decodeUser.ID)
	if err != nil {
		logrus.Errorf("Could not convert user id into object id: CreateBlog API: %v", nil)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	res := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.USER_COLLECTION).FindOneAndUpdate(
		context.TODO(), 
		bson.M{
			"_id": userObjectId,
		}, 
		bson.M{
			"$inc": bson.M{
				"stats.blogs_created": 1,
			},
		},
	)
	if res.Err() != nil {
		logrus.Errorf("Error updating user collection: CreateBlog API: %v", res.Err())
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusCreated, "Blog created successfully", result)
}

func GetAllBlogs(c *gin.Context) {
	query := c.Query("q")

	var filter interface{}

	if query != "" {
		filter = bson.M{
			"title": bson.M{
				"$regex":   query, // The substring you're looking for
				"$options": "i",   // Makes the search case-insensitive (optional)
			},
		}
	} else {
		filter = bson.M{}
	}

	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).Find(context.TODO(), filter)
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

func UpdateBlog(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("Invalid blog id: UpdateBlog API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid blog id", nil)
		return
	}

	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).FindOne(context.TODO(), bson.M{"_id": objectId})
	if result.Err() != nil {
		logrus.Errorf("Blog not found: UpdateBlog API: %v", result.Err())
		response.HandleResponse(c, http.StatusNotFound, "Blog not found", nil)
		return
	}

	var blogToUpdate models.Blog
	if err := result.Decode(&blogToUpdate); err != nil {
		logrus.Errorf("Error decoding the blog: UpdateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	var blog request.UpdateBlogRequest
	if err := c.ShouldBindJSON(&blog); err != nil {
		logrus.Errorf("Invalid request body: UpdateBlog API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if blog.Title == "" && blog.Body == "" && blog.IsBlogPublished == blogToUpdate.IsBlogPublished {
		logrus.Error("No fields to update: UpdateBlog API")
		response.HandleResponse(c, http.StatusBadRequest, "No fields to update", nil)
		return
	}

	if blog.Title != "" && blogToUpdate.Title != blog.Title {
		blogToUpdate.Title = blog.Title

		// change the slug
		words := strings.Split(blog.Title, " ")
		slug := strings.Join(words, "-")
		blogToUpdate.Slug = slug
	}

	if blog.Body != "" && blogToUpdate.Body != blog.Body {
		blogToUpdate.Body = blog.Body
	}

	if blog.IsBlogPublished != blogToUpdate.IsBlogPublished {
		blogToUpdate.IsBlogPublished = blog.IsBlogPublished
	}

	updateStage := bson.M{
		"$set": bson.M{
			"title":           blogToUpdate.Title,
			"body":            blogToUpdate.Body,
			"isBlogPublished": blogToUpdate.IsBlogPublished,
			"slug":            blogToUpdate.Slug,
		},
	}

	res, updateErr := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).UpdateOne(context.TODO(), bson.M{"_id": objectId}, updateStage)
	if updateErr != nil {
		logrus.Errorf("Error updating blog: UpdateBlog API: %v", updateErr)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	if res.ModifiedCount == 0 {
		logrus.Error("Did not update any fields: UpdateBlog API")
		response.HandleResponse(c, http.StatusBadRequest, "Did not update any fields", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Blog updated successfully", nil)
}

func DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("Invalid blog id: DeleteBlog API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid blog id", nil)
		return
	}

	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).FindOneAndDelete(context.TODO(), bson.M{"_id": objectId})
	if result.Err() != nil {
		logrus.Error("Blog not found: DeleteBlog API")
		response.HandleResponse(c, http.StatusNotFound, "Blog not found", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Blog deleted successfully", nil)
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
