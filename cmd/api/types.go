package main

type jsonRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname,omitempty"`
}

type jsonResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
	Error   map[string]any `json:"error,omitempty"`
}
