package models

type Document struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    Content string `json:"content"`
    // Add more fields as needed
}
