package auth

import (
    "errors" // Import the errors package
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    jwt "github.com/dgrijalva/jwt-go"
)

// Define the LoginRequest struct for handling login credentials
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Define the JWTClaims struct for JWT payload
type JWTClaims struct {
    Username    string   `json:"username"`
    Role        string   `json:"role"`
    Permissions []string `json:"permissions"`
    // Add more claims as needed
}

// Implement the Valid method of jwt.Claims interface
func (c *JWTClaims) Valid() error {
    // Check if the username is empty
    if c.Username == "" {
        return errors.New("username is empty")
    }

    // Check if the role is empty
    if c.Role == "" {
        return errors.New("role is empty")
    }

    // Perform additional validation checks as needed

    // Return nil if the token claims are valid
    return nil
}

// Endpoint for admin login
func LoginHandler(c *gin.Context) {
    var loginRequest LoginRequest
    if err := c.BindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    // Simulated user authentication
    if authenticateUser(loginRequest.Username, loginRequest.Password) {
        // Generate JWT with permissions (example)
        permissions := []string{"read", "write"} // You can customize permissions based on roles
        token := generateJWT(loginRequest.Username, "admin", permissions)

        // Send JWT in response
        c.JSON(http.StatusOK, gin.H{"token": token})
        return
    }

    // Authentication failed
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

// Simulated user authentication function
func authenticateUser(username, password string) bool {
    // Simulated user credentials
    correctUsername := "admin"
    correctPassword := "password"

    // Check if provided username and password match the correct credentials
    return username == correctUsername && password == correctPassword
}

// Generate JWT
func generateJWT(username, role string, permissions []string) string {
    claims := JWTClaims{
        Username:    username,
        Role:        role,
        Permissions: permissions,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims) // Use a pointer to the claims struct
    signedToken, err := token.SignedString([]byte("your-secret-key"))
    if err != nil {
        log.Println("Error generating JWT:", err)
        return ""
    }

    return signedToken
}
