package main

import (
    "encoding/json"
    "net/http"
    "github.com/glebarez/sqlite"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
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

    _ = http.ListenAndServe(":8080", nil)
}
