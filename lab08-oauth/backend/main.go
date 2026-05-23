package main

import (
    "context"
    "encoding/hex"
    "encoding/json"
    "net/http"

    "github.com/glebarez/sqlite"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "gorm.io/gorm"
)

const (
    googleClientID     = "I-WONT-LEAK-THAT"
    googleClientSecret = "CAUSE-I-AM-NOT-AI"
    googleRedirectURL  = "http://localhost:8080/auth/google/callback"
    frontendURL        = "http://localhost:8000"
)

type User struct {
    gorm.Model
    Email        string `gorm:"uniqueIndex" json:"email"`
    PasswordHash string `json:"-"`
    Provider     string `json:"provider"`
    ProviderID   string `json:"provider_id"`
}

var db *gorm.DB

func withCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next(w, r)
    }
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid payload", http.StatusBadRequest)
        return
    }
    if req.Email == "" || req.Password == "" {
        http.Error(w, "email and password required", http.StatusBadRequest)
        return
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "server error", http.StatusInternalServerError)
        return
    }
    user := User{Email: req.Email, PasswordHash: string(hash)}
    if err := db.Create(&user).Error; err != nil {
        http.Error(w, "could not create user", http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid payload", http.StatusBadRequest)
        return
    }
    var user User
    if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
}

func randomToken() (string) {
    b := make([]byte, 24)
    return hex.EncodeToString(b)
}

func authGoogleHandler(w http.ResponseWriter, r *http.Request) {
    cfg := &oauth2.Config{
        ClientID:     googleClientID,
        ClientSecret: googleClientSecret,
        RedirectURL:  googleRedirectURL,
        Scopes:       []string{"email", "profile"},
        Endpoint:     google.Endpoint,
    }
    w.Header().Set("Cache-Control", "no-store")
    w.Header().Set("Pragma", "no-cache")
    state := randomToken()

    http.SetCookie(w, &http.Cookie{
        Name:     "oauth_state",
        Value:    state,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
    })
    authURL := cfg.AuthCodeURL(
        state,
        oauth2.AccessTypeOffline,
        oauth2.SetAuthURLParam("prompt", "select_account"),
    )
    http.Redirect(w, r, authURL, http.StatusFound)
}

func authGoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
    cfg := &oauth2.Config{
        ClientID:     googleClientID,
        ClientSecret: googleClientSecret,
        RedirectURL:  googleRedirectURL,
        Scopes:       []string{"email", "profile"},
        Endpoint:     google.Endpoint,
    }
    code := r.URL.Query().Get("code")
    state := r.URL.Query().Get("state")
    cookie, err := r.Cookie("oauth_state")
    if err != nil || cookie.Value != state {
        http.Error(w, "invalid state", http.StatusBadRequest)
        return
    }
    token, err := cfg.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "oauth exchange failed", http.StatusBadRequest)
        return
    }
    client := cfg.Client(context.Background(), token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        http.Error(w, "user info failed", http.StatusBadRequest)
        return
    }
    defer resp.Body.Close()
    var info struct {
        ID    string `json:"id"`
        Email string `json:"email"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
        http.Error(w, "invalid user info", http.StatusBadRequest)
        return
    }
    var user User
    err = db.Where("provider = ? AND provider_id = ?", "google", info.ID).First(&user).Error
    if err != nil {
        user = User{Email: info.Email, Provider: "google", ProviderID: info.ID}
        if err := db.Create(&user).Error; err != nil {
            http.Error(w, "could not create user", http.StatusBadRequest)
            return
        }
    }
    http.Redirect(w, r, frontendURL+"/?oauth=ok", http.StatusFound)
}

func main() {
    var err error
    db, err = gorm.Open(sqlite.Open("/app/data/app.db"), &gorm.Config{})
    if err != nil {
        return
    }
    if err := db.AutoMigrate(&User{}); err != nil {
        return
    }

    http.HandleFunc("/register", withCORS(registerHandler))
    http.HandleFunc("/login", withCORS(loginHandler))
    http.HandleFunc("/auth/google", authGoogleHandler)
    http.HandleFunc("/auth/google/callback", authGoogleCallbackHandler)

    _ = http.ListenAndServe(":8080", nil)
}
