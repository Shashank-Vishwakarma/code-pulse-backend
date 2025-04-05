package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/services"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	request "github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/auth"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/blog"
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

	userObjectId, err := primitive.ObjectIDFromHex(decodeUser.ID)
	if err != nil {
		logrus.Errorf("Could not convert user id into object id: CreateBlog API: %v", nil)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
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
		AuthorID:        userObjectId,
		ImageURL:        imageUrl,
	})
	if err != nil {
		logrus.Errorf("Error creating blog: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	// update user collection
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
			"isBlogPublished": true,
		}
	} else {
		filter = bson.M{
			"isBlogPublished": true,
		}
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

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"_id": objectId}}},

		// lookup for author
		{{"$lookup", bson.M{
			"from":         "users",       // Name of the users collection
			"localField":   "authorId",     // Field in the blogs collection
			"foreignField": "_id",         // Field in the users collection
			"as":           "author",  // Output array field for user data
		}}},
		{{"$unwind", bson.M{"path": "$author", "preserveNullAndEmptyArrays": true}}}, // Flatten author array

		bson.D{{"$lookup", bson.M{
			"from": "comments",
			"let": bson.M{"comment_ids": "$comment_ids"},
			"pipeline": mongo.Pipeline{
				{{"$match", bson.M{
					"$expr": bson.M{"$in": bson.A{"$_id", "$$comment_ids"}},
				}}},
				// Lookup for comment's user
				{{"$lookup", bson.M{
					"from":         "users",
					"localField":   "userId",
					"foreignField": "_id",
					"as":           "user",
				}}},
				{{"$unwind", bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}}},
				// Optional: remove sensitive fields
				{{"$project", bson.M{
					"user.password": 0,
					"user._id":      0,
				}}},
				{{"$sort", bson.M{"createdAt": -1}}},
			},
			"as": "comments",
		}}},

		{{"$project", bson.M{
			"author.password": 0,   // Exclude password from author data
			"author._id":      0,   // Optional: Exclude MongoDB's _id field for author		
		}}},
	}

	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).Aggregate(context.TODO(), pipeline)
	if err != nil {
		logrus.Errorf("Blog not found: GetBlogById API: %v", err.Error())
		response.HandleResponse(c, http.StatusNotFound, "Blog not found", nil)
		return
	}

	var blogs []struct {
		ID              string             `json:"id" bson:"_id"`
		Title           string             `json:"title" bson:"title"`
		Body            string             `json:"body" bson:"body"`
		ImageURL        string             `json:"imageUrl" bson:"imageUrl"`
		Slug            string             `json:"slug" bson:"slug"`
		Comments        []struct{
			ID        primitive.ObjectID `json:"id" bson:"_id"`
			Body      string             `json:"body" bson:"body"`
			UserID    primitive.ObjectID `json:"userId" bson:"userId"`
			User struct {
				Name     string `json:"name" bson:"name"`
				Username string `json:"username" bson:"username"`
			} `json:"user" bson:"user"`
			BlogID    primitive.ObjectID `json:"blogId" bson:"blogId"`
			CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
		}          `json:"comments,omitempty" bson:"comments"`
		AuthorID        primitive.ObjectID `json:"authorId" bson:"authorId"`
		Author struct{
			Name     string             `json:"name" bson:"name"`
			Username string             `json:"username" bson:"username"`
		} `json:"author" bson:"author"`
		CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	}
	if err := cursor.All(context.TODO(), &blogs); err != nil {
		logrus.Errorf("Error decoding the blog: GetBlogById API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Blog fetched successfully", blogs[0])
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
	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(decodeUser.ID)
	if err != nil {
		logrus.Errorf("Could not convert user id into object id: CreateBlog API: %v", nil)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	options := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).Find(context.TODO(), bson.M{"authorId": userObjectId}, options)
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

func CreateComment(c *gin.Context ) {
	blogId := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		logrus.Errorf("Invalid blog id: CreateComment API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid blog id", nil)
		return
	}

	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: CreateComment API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(decodeUser.ID)
	if err != nil {
		logrus.Errorf("Error getting user id: CreateComment API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var comment blog.CommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		logrus.Errorf("Invalid request body: CreateComment API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	result, err := models.InsertDocumentInComments(&models.Comment{
		Body: comment.Body,
		UserID: userObjectId,
		BlogID: objectId,
	})
	if err != nil {
		logrus.Errorf("Error creating comment: CreateComment API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	// update blog
	updateStage := bson.M{
		"$push": bson.M{
			"comment_ids": result.InsertedID,
		},
	}

	_, updateErr := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).UpdateOne(context.TODO(), bson.M{"_id": objectId}, updateStage)
	if updateErr != nil {
		logrus.Errorf("Error updating blog: CreateComment API: %v", updateErr)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Comment created successfully", result.InsertedID)
}

func GetAllCommentsOnABlog(c *gin.Context) {
	
}
