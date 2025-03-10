package main

import (
	"fmt"
	"log"
	"notex/api/repository"
	"notex/config"
	"notex/model"
	"notex/pkg/database"
	"time"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig("../config/config.yaml"); err != nil {
		log.Printf("Warning: Failed to load config file: %v", err)
		log.Println("Using default configuration")
	}

	cfg := config.GetConfig()

	// 初始化数据库连接
	if err := database.Initialize(cfg.Database); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 创建仓库实例
	userRepo := repository.NewUserRepository()
	categoryRepo := repository.NewCategoryRepository()
	tagRepo := repository.NewTagRepository()
	postRepo := repository.NewPostRepository()

	// 创建测试用户
	users := getOrCreateUsers(userRepo)
	fmt.Println("Created users:", len(users))

	// 创建分类
	categories := getOrCreateCategories(categoryRepo)
	fmt.Println("Created categories:", len(categories))

	// 创建标签
	tags := getOrCreateTags(tagRepo)
	fmt.Println("Created tags:", len(tags))

	// 创建文章
	posts := createPosts(postRepo, users, categories, tags)
	fmt.Println("Created posts:", len(posts))
}

func getOrCreateUsers(repo *repository.UserRepository) []*model.User {
	users := make([]*model.User, 0, 3)
	usernames := []string{"admin", "editor", "user"}

	// 尝试获取已存在的用户
	for _, username := range usernames {
		if user, err := repo.FindByUsername(username); err == nil {
			users = append(users, user)
		}
	}

	// 如果没有找到足够的用户，创建新用户
	if len(users) < 3 {
		newUsers := []*model.User{
			{
				Username: "admin",
				Email:    "admin@example.com",
				Role:     model.RoleAdmin,
				Status:   "active",
				Bio:      "系统管理员",
			},
			{
				Username: "editor",
				Email:    "editor@example.com",
				Role:     model.RoleEditor,
				Status:   "active",
				Bio:      "内容编辑",
			},
			{
				Username: "user",
				Email:    "user@example.com",
				Role:     model.RoleUser,
				Status:   "active",
				Bio:      "普通用户",
			},
		}

		for _, u := range newUsers {
			exists := false
			for _, existingUser := range users {
				if existingUser.Username == u.Username {
					exists = true
					break
				}
			}
			if !exists {
				u.SetPassword("123456") // 设置默认密码
				if err := repo.Create(u); err != nil {
					log.Printf("Failed to create user %s: %v", u.Username, err)
				} else {
					// 重新获取用户以获得完整的信息（包括ID）
					if createdUser, err := repo.FindByUsername(u.Username); err == nil {
						users = append(users, createdUser)
					}
				}
			}
		}
	}

	return users
}

func getOrCreateCategories(repo *repository.CategoryRepository) []*model.Category {
	categories := make([]*model.Category, 0, 5)
	categoryNames := []string{"技术", "生活", "学习", "工具", "其他"}

	// 获取所有分类
	existingCategories, _, err := repo.List(1, 100, "")
	if err != nil {
		log.Printf("Failed to get categories: %v", err)
	}

	// 将已存在的分类添加到结果中
	for _, name := range categoryNames {
		var found bool
		for _, c := range existingCategories {
			if c.Name == name {
				cat := c // 创建一个新的变量来存储分类
				categories = append(categories, &cat)
				found = true
				break
			}
		}
		if !found {
			category := &model.Category{
				Name:        name,
				Description: name + "相关文章",
			}
			if err := repo.Create(category); err != nil {
				log.Printf("Failed to create category %s: %v", name, err)
			} else {
				// 重新获取分类以获得完整的信息（包括ID）
				if createdCategory, err := repo.FindByID(category.ID); err == nil {
					categories = append(categories, createdCategory)
				}
			}
		}
	}

	return categories
}

func getOrCreateTags(repo *repository.TagRepository) []*model.Tag {
	tags := make([]*model.Tag, 0, 12)
	tagNames := []string{
		"Go", "Python", "JavaScript", "Vue", "React",
		"Docker", "Linux", "Git", "数据库", "前端", "后端", "全栈",
	}

	// 获取所有标签
	existingTags, _, err := repo.List(1, 100, "")
	if err != nil {
		log.Printf("Failed to get tags: %v", err)
	}

	// 将已存在的标签添加到结果中
	for _, name := range tagNames {
		var found bool
		for _, t := range existingTags {
			if t.Name == name {
				tag := t // 创建一个新的变量来存储标签
				tags = append(tags, &tag)
				found = true
				break
			}
		}
		if !found {
			tag := &model.Tag{Name: name}
			if err := repo.Create(tag); err != nil {
				log.Printf("Failed to create tag %s: %v", name, err)
			} else {
				// 重新获取标签以获得完整的信息（包括ID）
				if createdTag, err := repo.FindByID(tag.ID); err == nil {
					tags = append(tags, createdTag)
				}
			}
		}
	}

	return tags
}

func createPosts(repo *repository.PostRepository, users []*model.User, categories []*model.Category, tags []*model.Tag) []*model.Post {
	posts := []*model.Post{
		{
			Title:    "Go语言入门指南",
			Content:  "Go（又称 Golang）是 Google 开发的一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言。",
			Summary:  "Go语言入门基础知识介绍",
			Status:   "published",
			UserID:   users[1].ID,                      // editor
			Category: *categories[0],                   // 技术
			Tags:     []model.Tag{*tags[0], *tags[10]}, // Go, 后端
			Slug:     "go-language-guide",
		},
		{
			Title:    "Vue.js 3.0 新特性解析",
			Content:  "Vue.js 3.0 带来了许多令人兴奋的新特性，包括 Composition API、更好的 TypeScript 支持等。",
			Summary:  "详解 Vue.js 3.0 的主要新特性",
			Status:   "published",
			UserID:   users[1].ID,                     // editor
			Category: *categories[0],                  // 技术
			Tags:     []model.Tag{*tags[3], *tags[9]}, // Vue, 前端
			Slug:     "vue-3-new-features",
		},
		{
			Title:    "Docker 容器化实践",
			Content:  "Docker 是一个开源的应用容器引擎，让开发者可以打包他们的应用以及依赖包到一个可移植的容器中。",
			Summary:  "Docker 基础知识和实践经验分享",
			Status:   "published",
			UserID:   users[1].ID,                     // editor
			Category: *categories[3],                  // 工具
			Tags:     []model.Tag{*tags[5], *tags[6]}, // Docker, Linux
			Slug:     "docker-containerization",
		},
		{
			Title:    "Git 版本控制入门",
			Content:  "Git 是一个开源的分布式版本控制系统，可以有效、高速地处理从很小到非常大的项目版本管理。",
			Summary:  "学习 Git 的基础知识和常用命令",
			Status:   "published",
			UserID:   users[2].ID,           // user
			Category: *categories[2],        // 学习
			Tags:     []model.Tag{*tags[7]}, // Git
			Slug:     "git-version-control",
		},
		{
			Title:    "程序员的生活感悟",
			Content:  "作为一名程序员，除了技术提升，也要注意生活质量和身心健康。",
			Summary:  "分享程序员的日常生活和感悟",
			Status:   "published",
			UserID:   users[2].ID,    // user
			Category: *categories[1], // 生活
			Tags:     []model.Tag{},  // 无标签
			Slug:     "programmer-life",
		},
	}

	for _, post := range posts {
		post.PublishedAt = time.Now()
		if err := repo.Create(post); err != nil {
			log.Printf("Failed to create post %s: %v", post.Title, err)
		}
	}

	return posts
}
