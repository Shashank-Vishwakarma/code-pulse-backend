package routes

import (
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/handlers"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/middlewares"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"github.com/gin-gonic/gin"
)

func BlogRoutes(r *gin.Engine) {
	blogRouteGroup := r.Group(constants.BLOG_API_BASE_ENDPOINT)

	blogRouteGroup.Use(middlewares.Authorization())

	blogRouteGroup.POST(constants.BLOG_API_CREATE_ENDPOINT, handlers.CreateBlog)
	blogRouteGroup.GET(constants.BLOG_API_GET_ALL_ENDPOINT, handlers.GetAllBlogs)
	blogRouteGroup.GET(constants.BLOG_API_GET_BY_ID_ENDPOINT, handlers.GetBlogById)
	blogRouteGroup.PUT(constants.BLOG_API_UPDATE_ENDPOINT, handlers.UpdateBlog)
	blogRouteGroup.DELETE(constants.BLOG_API_DELETE_ENDPOINT, handlers.DeleteBlog)
	blogRouteGroup.GET(constants.BLOG_API_SEARCH_ENDPOINT, handlers.SearchBlogs)
	blogRouteGroup.GET(constants.BLOG_API_GET_BY_USER_ID_ENDPOINT, handlers.GetBlogsByUser)
}
