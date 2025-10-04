package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"seesharpsi/htmx_quickstart/config"
)

// Session holds the state for a single user's story.
type Session struct {
	ID           string
	LastAccessed time.Time
}

// Manager handles the creation, storage, and retrieval of sessions.
type Manager struct {
	sessions map[string]*Session
	mutex    sync.Mutex
	config   *config.Config
}

// NewManager creates a new session manager with configuration.
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		config:   cfg,
	}
}

// CreateSession creates a new session and returns its ID.
func (m *Manager) CreateSession() string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Generate a random, secure session ID.
	b := make([]byte, 16)
	rand.Read(b)
	id := hex.EncodeToString(b)

	m.sessions[id] = &Session{
		ID: id,
		//GameState:    &story.GameState{},
		//StoryHistory: []story.StoryPage{},
		LastAccessed: time.Now(),
	}
	return id
}

// GetSession retrieves a session by its ID.
func (m *Manager) GetSession(id string) *Session {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	session, ok := m.sessions[id]
	if !ok {
		return nil
	}
	session.LastAccessed = time.Now()
	return session
}

// GetOrCreateSession retrieves an existing session or creates a new one.
func (m *Manager) GetOrCreateSession(r *http.Request) (*Session, http.Cookie) {
	cookie, err := r.Cookie(m.config.Session.CookieName)
	if err == nil {
		session := m.GetSession(cookie.Value)
		if session != nil {
			return session, *cookie
		}
	}

	// If no valid session is found, create a new one.
	id := m.CreateSession()
	newCookie := http.Cookie{
		Name:     m.config.Session.CookieName,
		Value:    id,
		Expires:  time.Now().Add(m.config.Session.MaxAge),
		HttpOnly: m.config.Session.HttpOnly,
		Secure:   m.config.Session.Secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // Default to Lax for now
	}
	return m.GetSession(id), newCookie
}
