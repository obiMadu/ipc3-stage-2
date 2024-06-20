package main

type jsonRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname,omitempty"`
}

type jsonResponse struct {
	Status  string         `json:"status" binding:"required"`
	Message string         `json:"message" binding:"required"`
	Data    map[string]any `json:"data" binding:"required"`
	Error   map[string]any `json:"error,omitempty"`
}
