package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/ellisbywater/gocial/internal/store"
)

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi", "ivan", "jack",
	"karen", "leo", "mike", "nancy", "oliver", "paul", "quincy", "rachel", "steve", "tom",
	"ursula", "victor", "wendy", "xander", "yvonne", "zach", "amanda", "brian", "claire", "daniel",
	"ella", "fay", "george", "hannah", "isaac", "julie", "kevin", "lily", "matthew", "nina",
	"oscar", "peter", "quinn", "rosa", "sam", "tina", "uma", "vicky", "will", "xena",
}

var titles = []string{
	"10 Tips for Better Code",
	"Why Go is Perfect for Backend Development",
	"The Future of Cloud Computing",
	"Mastering Go's Goroutines",
	"Understanding Go Interfaces",
	"How to Handle Errors in Go",
	"Building Scalable Web APIs with Go",
	"5 Go Libraries You Should Know",
	"Optimizing Your Go Application",
	"Go vs Python: Which Should You Learn?",
	"Exploring Go's Concurrency Model",
	"An Introduction to Go Generics",
	"Go Memory Management Explained",
	"Building RESTful APIs in Go",
	"Best Practices for Go Developers",
	"Working with Go Modules",
	"Creating Microservices with Go",
	"Deploying Go Apps on Kubernetes",
	"Testing in Go: A Beginner's Guide",
	"Go's Reflection Package: A Deep Dive",
}

var contents = []string{
	"Tips on maintaining focus, creating a home office setup, and tools for virtual collaboration.",
	"Steps to organize your inbox, files, social media accounts, and apps.",
	"A deep dive into the science of mindfulness, its benefits, and how to start practicing.",
	"A list of lesser-known travel spots that offer great experiences without breaking the bank.",
	"A guide to journaling, including different types of journals (gratitude, bullet, etc.) and their benefits.",
	"Book recommendations across genres that offer profound insights and lessons.",
	"Tips for choosing eco-friendly fabrics, ethical brands, and ways to shop mindfully.",
	"Strategies to prevent burnout and create boundaries between work and personal life.",
	"Explore the meaning behind different colors and how they can influence your day-to-day life.",
	"A step-by-step guide to creating a productive, balanced morning ritual.",
	"A look at AI’s potential in various industries and the ethical concerns that accompany its growth.",
	"Simplified advice for new investors on stocks, bonds, and how to get started with minimal risk.",
	"Simple changes to reduce your carbon footprint, from energy-efficient appliances to low-VOC paint.",
	"How online tools like LinkedIn and social media can enhance professional connections.",
	"Tips on writing compelling posts, using visuals effectively, and engaging with your audience.",
	"An easy-to-understand breakdown of blockchain technology, Bitcoin, and other cryptocurrencies.",
	"Insights from psychology, neuroscience, and positive psychology about what makes us happy.",
	"A guide to self-care routines that prioritize mental, physical, and emotional well-being.",
	"Practical tips for living with less and finding joy in simplicity.",
	"Techniques for preparing and delivering speeches, as well as tips to calm nerves",
}

var tags = []string{
	"GoLang", "Programming", "Backend", "Concurrency", "Web Development",
	"API", "REST", "Microservices", "Docker", "Kubernetes",
	"Cloud", "Database", "DevOps", "Software Engineering", "Design Patterns",
	"Performance", "Scalability", "Security", "Testing", "Deployment",
}

var comments = []string{
	"Great post! Really informative.",
	"I totally agree with your point.",
	"Thanks for sharing this.",
	"Can you explain this in more detail?",
	"I found this very helpful, thank you!",
	"Interesting perspective, I hadn't thought about it that way.",
	"This really resonates with me.",
	"Well written and insightful!",
	"Keep up the good work!",
	"Looking forward to your next post.",
}

func Seed(store store.Storage, db *sql.DB) {
	fmt.Println("Seeding has started...")
	ctx := context.Background()
	fmt.Println("...Generating users")
	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}
	tx.Commit()
	fmt.Println("...Generating Posts")
	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post: ", err)
			return
		}
	}
	fmt.Println("...Generating Comments")
	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment: ", err)
			return
		}
	}

	fmt.Println("...Seeding has completed")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		ran := rand.IntN(len(users))
		user := users[ran]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.IntN(len(titles))],
			Content: contents[rand.IntN(len(contents))],
			Tags: []string{
				tags[rand.IntN(len(tags))],
				tags[rand.IntN(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.IntN(len(posts))].ID,
			UserID:  users[rand.IntN(len(users))].ID,
			Content: comments[rand.IntN(len(comments))],
		}
	}
	return cms
}
