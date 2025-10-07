package post_logic

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

// PostCache holds the parsed posts and the state of the directory.
// It's safe for concurrent use.
type PostCache struct {
	mu             sync.RWMutex
	posts          []Post
	directoryState map[string]time.Time
	postsDir       string
}

// NewPostCache creates and initializes a new PostCache.
// It performs an initial scan and parse of the posts directory.
func NewPostCache(postsDir string) (*PostCache, error) {
	pc := &PostCache{
		postsDir: postsDir,
	}

	log.Printf("Performing initial scan of posts directory: %s", postsDir)
	if err := pc.refreshPosts(); err != nil {
		return nil, err
	}

	return pc, nil
}

// GetPosts returns the cached posts. If the directory has changed,
// it transparently re-parses the posts before returning.
func (pc *PostCache) GetPosts() []Post {
	// Check for changes with a read lock first for performance.
	pc.mu.RLock()
	currentState, err := getDirectoryState(pc.postsDir)
	if err != nil {
		log.Printf("Error checking directory state, returning stale data: %v", err)
		pc.mu.RUnlock()
		return pc.posts
	}

	// If states are the same, return the cached data.
	if reflect.DeepEqual(currentState, pc.directoryState) {
		pc.mu.RUnlock()
		return pc.posts
	}
	pc.mu.RUnlock() // Must unlock before acquiring a write lock

	// If states are different, acquire a write lock to update the cache.
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Re-check the state after acquiring the write lock, in case another
	// goroutine already updated the cache.
	currentState, _ = getDirectoryState(pc.postsDir)
	if reflect.DeepEqual(currentState, pc.directoryState) {
		return pc.posts
	}

	log.Println("Posts directory has changed. Refreshing posts...")
	if err := pc.refreshPosts(); err != nil {
		log.Printf("Error refreshing posts, returning stale data: %v", err)
	}

	return pc.posts
}

// GetPostByID retrieves a single post by its ID from the cache.
// It is thread-safe.
func (pc *PostCache) GetPostByID(id int) (Post, error) {
	pc.mu.RLock() // Lock for reading
	defer pc.mu.RUnlock()

	// Find the post with the matching ID
	for _, post := range pc.posts {
		if post.Id == id {
			return post, nil
		}
	}

	return Post{}, errors.New("post not found")
}

// refreshPosts scans the directory, parses all HTML files, and updates the cache.
// This function is NOT thread-safe and must be called within a write lock.
func (pc *PostCache) refreshPosts() error {
	var newPosts []Post

	entries, err := os.ReadDir(pc.postsDir)
	if err != nil {
		return err
	}

	postID := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		if ext == ".html" || ext == ".htm" {
			fullPath := filepath.Join(pc.postsDir, fileName)
			contentBytes, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("Error reading post file %s: %v", fullPath, err)
				continue
			}

			var post Post
			post.Id = postID
			// Assuming Parse function is updated to take an ID
			post.Parse(string(contentBytes))

			newPosts = append(newPosts, post)
			postID++
		}
	}

	pc.posts = newPosts
	// Update the directory state after successful parsing
	pc.directoryState, err = getDirectoryState(pc.postsDir)
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed posts. Found %d posts.", len(pc.posts))
	return nil
}
