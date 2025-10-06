package post_logic

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Goes through the ./posts directory and returns a parsed list of posts 
func Get_Posts() []Post {
	var posts []Post
	entries, err := os.ReadDir("./posts/")
	if err != nil {
		log.Fatalf("Failed to read the directory: %v", err)
	}
	for i, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		if ext == ".html" || ext == ".htm" {
			fullPath := filepath.Join("./posts/", fileName)
			contentBytes, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("Error reading file %s: %v", fullPath, err)
				continue
			}
			post := Post{Id: i}
			contentString := string(contentBytes)
			post.Parse(contentString)
			posts = append(posts, post)
		}
	}
	return posts
}
