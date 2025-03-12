package router

import (
	"log"
	"notex/api/handler"
	"notex/api/repository"
	"notex/api/service"
	"notex/config"
	"notex/middleware"
	"notex/pkg/storage"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 全局中间件
	r.Use(middleware.IPRateLimit())  // IP限流
	r.Use(middleware.APIRateLimit()) // API限流

	// 静态文件服务（当配置了本地存储路径和URL前缀时）
	if cfg.Storage.Local.URLPrefix != "" && cfg.Storage.Local.UploadDir != "" {
		r.Static(cfg.Storage.Local.URLPrefix, cfg.Storage.Local.UploadDir)
	}

	adminService := service.NewAdminService()
	authService := service.NewAuthService()
	categoryService := service.NewCategoryService()
	commentService := service.NewCommentService()
	postService := service.NewPostService()
	tagService := service.NewTagService()
	verificationService := service.NewVerificationService()

	// 创建存储实例
	storageInstance, err := storage.DefaultFactory.CreateStorage(&cfg.Storage)
	if err != nil {
		log.Fatal("Failed to create storage instance:", err)
	}

	// 创建上传处理器
	uploadHandler := handler.NewUploadHandler(storageInstance, &cfg.Storage)

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关路由（无需认证）
		auth := api.Group("/auth")
		{
			authHandler := handler.NewAuthHandler(authService, postService)
			verificationHandler := handler.NewVerificationHandler(verificationService)

			auth.POST("/register", middleware.LoginRateLimit(), authHandler.Register)
			auth.POST("/login", middleware.LoginRateLimit(), authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/password-reset/send", verificationHandler.SendPasswordReset)
			auth.POST("/password-reset/verify", verificationHandler.ResetPassword)
		}

		// 公开接口，无需认证
		postHandler := handler.NewPostHandler(postService)
		commentHandler := handler.NewCommentHandler(commentService)
		categoryHandler := handler.NewCategoryHandler(categoryService)
		tagHandler := handler.NewTagHandler(tagService)
		authHandler := handler.NewAuthHandler(authService, postService)

		// 公开接口组
		public := api.Group("/public")
		{
			// 文章相关的公开接口
			public.GET("/posts", postHandler.ListPublicPosts)
			public.GET("/posts/:id", postHandler.GetPost)
			public.GET("/posts/:id/comments", commentHandler.ListComments)
			public.GET("/posts/archives", postHandler.GetArchives)
			public.GET("/posts/archives/:yearMonth", postHandler.GetPostsByArchive)

			// 分类和标签的公开接口
			public.GET("/categories", categoryHandler.ListCategories)
			public.GET("/categories/top", categoryHandler.GetTopCategories)
			public.GET("/tags", tagHandler.ListTags)
			public.GET("/tags/top", tagHandler.GetTopTags)

			// 用户公开信息接口
			public.GET("/users/:id/home", authHandler.GetUserHome)
			public.GET("/users/:id/comments", commentHandler.ListUserComments)
		}

		// 需要认证的路由组
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// 认证相关路由
			authHandler := handler.NewAuthHandler(authService, postService)
			verificationHandler := handler.NewVerificationHandler(verificationService)
			authenticated.GET("/auth/profile", authHandler.GetProfile)
			authenticated.PUT("/auth/profile", authHandler.UpdateProfile)
			authenticated.POST("/auth/change-password", authHandler.ChangePassword)
			authenticated.POST("/auth/logout", authHandler.Logout)
			authenticated.POST("/auth/email/update", verificationHandler.UpdateEmail)
			authenticated.POST("/auth/email/send-verification", verificationHandler.SendEmailVerification)
			authenticated.POST("/auth/email/verify", verificationHandler.VerifyEmail)

			// 用户相关路由
			users := authenticated.Group("/users")
			{
				// 用户评论相关路由
				users.GET("/comments", commentHandler.ListUserComments)
			}

			// 文件上传相关路由
			upload := authenticated.Group("/upload")
			{
				upload.POST("/file", middleware.RequireEditor(), uploadHandler.Upload)
				upload.GET("/config", middleware.RequireEditor(), uploadHandler.GetUploadConfig)
				upload.GET("/credentials", middleware.RequireEditor(), uploadHandler.GetCredentials)
			}

			// 文章相关路由（需要认证）
			posts := authenticated.Group("/posts")
			{
				posts.GET("/recent", postHandler.GetRecentPosts)
				posts.GET("", postHandler.ListPosts)
				posts.POST("", middleware.RequireEditor(), postHandler.CreatePost)
				posts.PUT("/:id", middleware.RequireEditor(), postHandler.UpdatePost)
				posts.DELETE("/:id", middleware.RequireEditor(), postHandler.DeletePost)

				// 评论相关路由（需要认证）
				posts.POST("/:id/comments", commentHandler.CreateComment)
				posts.DELETE("/:id/comments/:commentId", commentHandler.DeleteComment)
			}

			// 分类相关路由（需要认证）
			categories := authenticated.Group("/categories")
			{
				categories.GET("/:id", categoryHandler.GetCategory)
				categories.POST("", middleware.RequireEditor(), categoryHandler.CreateCategory)
				categories.PUT("/:id", middleware.RequireEditor(), categoryHandler.UpdateCategory)
				categories.DELETE("/:id", middleware.RequireEditor(), categoryHandler.DeleteCategory)
			}

			// 标签相关路由（需要认证）
			tags := authenticated.Group("/tags")
			{
				tags.GET("/:id", tagHandler.GetTag)
				tags.POST("", middleware.RequireEditor(), tagHandler.CreateTag)
				tags.PUT("/:id", middleware.RequireEditor(), tagHandler.UpdateTag)
				tags.DELETE("/:id", middleware.RequireEditor(), tagHandler.DeleteTag)
			}

			// 草稿相关路由（需要认证）
			draftService := service.NewDraftService(repository.NewDraftRepository())
			draftHandler := handler.NewDraftHandler(draftService)
			drafts := authenticated.Group("/drafts")
			{
				drafts.GET("", draftHandler.ListDrafts)
				drafts.GET("/:id", draftHandler.GetDraft)
				drafts.POST("", middleware.RequireEditor(), draftHandler.CreateDraft)
				drafts.PUT("/:id", middleware.RequireEditor(), draftHandler.UpdateDraft)
				drafts.DELETE("/:id", middleware.RequireEditor(), draftHandler.DeleteDraft)
				drafts.POST("/:id/publish", middleware.RequireEditor(), draftHandler.PublishDraft)
			}

			// 管理员路由
			adminHandler := handler.NewAdminHandler(adminService)
			adminHandler.RegisterRoutes(authenticated)
		}
	}

	return r
}
