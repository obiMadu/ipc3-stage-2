package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type args struct {
	sqlStatus            string
	sqlStatement         string
	sqlErrorText         string
	sqlReturnRows        [][]any
	httpRequestURL       string
	httpRequestMethod    string
	httpRequestBody      map[string]any
	httpResponseBodyCode int
	httpResponseBody     jsonResponse
}

type test struct {
	name string
	args args
}

func InitDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unable to init mock db %s", err.Error())
	}

	db, err := gorm.Open(postgres.New(
		postgres.Config{
			Conn: dbMock,
		}), &gorm.Config{})

	if err != nil {
		t.Fatalf("Failed to connect to db; gorm; %s", err.Error())
	}

	return db, mock
}

func TestCreateUser(t *testing.T) {
	tests := []test{
		{
			name: "Create user accurately",
			args: args{
				sqlStatement:      ".*",
				httpRequestURL:    "/",
				httpRequestMethod: "POST",
				httpRequestBody: gin.H{
					"username": "user1",
					"email":    "user1@example.com",
				},
				httpResponseBodyCode: http.StatusOK,
				httpResponseBody: jsonResponse{
					Status:  "success",
					Message: "User created successfully.",
				},
			},
		},
		{
			name: "Create user without username",
			args: args{
				sqlStatement:      ".*",
				httpRequestURL:    "/",
				httpRequestMethod: "POST",
				httpRequestBody: gin.H{
					"email": "user1@example.com",
				},
				httpResponseBodyCode: http.StatusBadRequest,
				httpResponseBody: jsonResponse{
					Status:  "error",
					Message: "You must specify both a username and an email.",
				},
			},
		},
		{
			name: "Create user without email",
			args: args{
				sqlStatement:      ".*",
				httpRequestURL:    "/",
				httpRequestMethod: "POST",
				httpRequestBody: gin.H{
					"username": "user1",
				},
				httpResponseBodyCode: http.StatusBadRequest,
				httpResponseBody: jsonResponse{
					Status:  "error",
					Message: "You must specify both a username and an email.",
				},
			},
		},
		{
			name: "Create user with duplicate username_EXPOSE SQLMOCK BUG",
			args: args{
				sqlStatement:      ".*",
				sqlStatus:         "error",
				sqlErrorText:      "duplicate key value violates unique constraint username",
				httpRequestURL:    "/",
				httpRequestMethod: "POST",
				httpRequestBody: gin.H{
					"username": "user1",
					"email":    "user2@example.com",
				},
				httpResponseBodyCode: http.StatusBadRequest,
				httpResponseBody: jsonResponse{
					Status:  "error",
					Message: "Username has been taken!",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// !important: get raw db and defer close, to avoid intra-test db leaks
			// Trust me those suck.
			db, mock := InitDB(t)
			rawDB, err := db.DB()
			if err != nil {
				t.Fatalf("Unable to get sql.DB from gorm.DB, %v", err)
			}
			defer rawDB.Close()

			switch test.args.sqlStatus {
			case "error":
				mock.ExpectBegin()
				mock.ExpectQuery(".*").WillReturnError(fmt.Errorf(test.args.sqlErrorText))
				mock.ExpectRollback()
			default:
				mock.ExpectBegin()
				mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows(nil))
				mock.ExpectCommit()
			}

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.Use(gin.Logger())

			// Setup listen route & handler
			r.POST("/", func(c *gin.Context) {
				CreateUser(c, db)
			})

			// Parse the URL to ensure query parameters are correctly included
			parsedURL, err := url.Parse(test.args.httpRequestURL)
			if err != nil {
				t.Fatalf("Failed to parse request URL %v", err)
			}

			requestBody, err := json.Marshal(test.args.httpRequestBody)
			if err != nil {
				t.Fatalf("Unable to marshal request body %v", err)
			}

			req, err := http.NewRequest(test.args.httpRequestMethod, parsedURL.String(), bytes.NewReader(requestBody))
			if err != nil {
				t.Fatalf("Failed to create request %v", err)
			}
			defer req.Body.Close()

			if requestBody != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var responseBody jsonResponse
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Fatalf("Failed to Unmarshall response body: %v", err)
			}

			assert.Equal(t, test.args.httpResponseBodyCode, w.Code)
			assert.Equal(t, test.args.httpResponseBody, responseBody)
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []test{
		{
			name: "Get all 2 rows",
			args: args{
				sqlStatement: ".*",
				sqlReturnRows: [][]any{
					{1, "Obi", "obi@example.com"},
					{2, "Marry", "marry@example.com"},
				},
				httpRequestURL:       "/",
				httpRequestMethod:    "GET",
				httpResponseBodyCode: http.StatusOK,
				httpResponseBody: jsonResponse{
					Status:  "success",
					Message: "Retrieved all users.",
					Data: gin.H{
						"users": []any{
							map[string]any{"id": 1.0, "username": "Obi", "email": "obi@example.com"},
							map[string]any{"id": 2.0, "username": "Marry", "email": "marry@example.com"},
						},
					},
				},
			},
		},
		{
			name: "Get user by username",
			args: args{
				sqlStatement: ".*",
				sqlReturnRows: [][]any{
					{1, "Obi", "obi@example.com"},
				},
				httpRequestURL:       "/?username=Obi",
				httpRequestMethod:    "GET",
				httpResponseBodyCode: http.StatusOK,
				httpResponseBody: jsonResponse{
					Status:  "success",
					Message: "User retrieved successfully.",
					Data: gin.H{
						"user": map[string]any{"id": 1.0, "username": "Obi", "email": "obi@example.com"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := InitDB(t)
			rawDB, err := db.DB()
			if err != nil {
				t.Fatalf("Unable to get sql.DB from gorm.DB, %v", err)
			}
			defer rawDB.Close()

			rows := sqlmock.NewRows([]string{"id", "username", "email"})
			for _, row := range test.args.sqlReturnRows {
				rows.AddRow(row[0], row[1], row[2])
			}

			mock.ExpectQuery(test.args.sqlStatement).WillReturnRows(rows)

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.Use(gin.Logger())
			r.GET("/", func(c *gin.Context) {
				GetAll(c, db)
			})

			// Parse the URL to ensure query parameters are correctly included
			parsedURL, err := url.Parse(test.args.httpRequestURL)
			if err != nil {
				t.Fatalf("Failed to parse request URL %v", err)
			}

			req, err := http.NewRequest(test.args.httpRequestMethod, parsedURL.String(), nil)
			if err != nil {
				t.Fatalf("Failed to create request %v", err)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var responseBody jsonResponse
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Fatalf("Failed to Unmarshall response body: %v", err)
			}

			assert.Equal(t, test.args.httpResponseBodyCode, w.Code)
			assert.Equal(t, test.args.httpResponseBody, responseBody)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []test{
		{
			name: "Get user by ID",
			args: args{
				sqlStatement: ".*",
				sqlReturnRows: [][]any{
					{1, "Obi", "obi@example.com"},
				},
				httpRequestURL:       "/1",
				httpRequestMethod:    "GET",
				httpResponseBodyCode: http.StatusOK,
				httpResponseBody: jsonResponse{
					Status:  "success",
					Message: "User retrieved successfully.",
					Data: gin.H{
						"user": map[string]any{"id": 1.0, "username": "Obi", "email": "obi@example.com"},
					},
				},
			},
		},
		{
			name: "Get user by string ID",
			args: args{
				sqlStatement: ".*",
				sqlReturnRows: [][]any{
					{1, "Obi", "obi@example.com"},
				},
				httpRequestURL:       "/one",
				httpRequestMethod:    "GET",
				httpResponseBodyCode: http.StatusBadRequest,
				httpResponseBody: jsonResponse{
					Status:  "error",
					Message: "UserID must be a positive interger.",
					Data:    nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := InitDB(t)
			rawDB, err := db.DB()
			if err != nil {
				t.Fatalf("Unable to get sql.DB from gorm.DB, %v", err)
			}
			defer rawDB.Close()

			rows := sqlmock.NewRows([]string{"id", "username", "email"})
			for _, row := range test.args.sqlReturnRows {
				rows.AddRow(row[0], row[1], row[2])
			}

			mock.ExpectQuery(test.args.sqlStatement).WillReturnRows(rows)

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.Use(gin.Logger())
			r.GET("/:userID", func(c *gin.Context) {
				GetUserByID(c, db)
			})

			// Parse the URL to ensure query parameters are correctly included
			parsedURL, err := url.Parse(test.args.httpRequestURL)
			if err != nil {
				t.Fatalf("Failed to parse request URL %v", err)
			}

			req, err := http.NewRequest(test.args.httpRequestMethod, parsedURL.String(), nil)
			if err != nil {
				t.Fatalf("Failed to create request %v", err)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var responseBody jsonResponse
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Fatalf("Failed to Unmarshall response body: %v", err)
			}

			assert.Equal(t, test.args.httpResponseBodyCode, w.Code)
			assert.Equal(t, test.args.httpResponseBody, responseBody)
		})
	}
}

func TestDeleteUserByID(t *testing.T) {
	tests := []test{
		{
			name: "Delete user by ID",
			args: args{
				sqlStatement: ".*",
				sqlReturnRows: [][]any{
					{1, "Obi", "obi@example.com"},
				},
				httpRequestURL:       "/1",
				httpRequestMethod:    "DELETE",
				httpResponseBodyCode: http.StatusOK,
				httpResponseBody: jsonResponse{
					Status:  "success",
					Message: "User deleted successfully.",
					Data:    nil,
				},
			},
		},
		{
			name: "Delete user by non-existent ID",
			args: args{
				sqlStatement:         ".*",
				sqlReturnRows:        [][]any{},
				httpRequestURL:       "/5",
				httpRequestMethod:    "DELETE",
				httpResponseBodyCode: http.StatusBadRequest,
				httpResponseBody: jsonResponse{
					Status:  "error",
					Message: "User does not exist.",
					Data:    nil,
				},
			},
		},
		{
			name: "Delete user by string ID",
			args: args{
				sqlStatement:         ".*",
				httpRequestURL:       "/one",
				httpRequestMethod:    "DELETE",
				httpResponseBodyCode: http.StatusBadRequest,
				httpResponseBody: jsonResponse{
					Status:  "error",
					Message: "UserID must be a positive interger.",
					Data:    nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := InitDB(t)
			rawDB, err := db.DB()
			if err != nil {
				t.Fatalf("Unable to get sql.DB from gorm.DB, %v", err)
			}
			defer rawDB.Close()

			rows := sqlmock.NewRows([]string{"id", "username", "email"})
			for _, row := range test.args.sqlReturnRows {
				rows.AddRow(row[0], row[1], row[2])
			}

			mock.ExpectQuery(test.args.sqlStatement).WillReturnRows(rows)
			mock.ExpectBegin()
			mock.ExpectExec(test.args.sqlStatement).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.Use(gin.Logger())
			r.DELETE("/:userID", func(c *gin.Context) {
				DeleteUserByID(c, db)
			})

			// Parse the URL to ensure query parameters are correctly included
			parsedURL, err := url.Parse(test.args.httpRequestURL)
			if err != nil {
				t.Fatalf("Failed to parse request URL %v", err)
			}

			req, err := http.NewRequest(test.args.httpRequestMethod, parsedURL.String(), nil)
			if err != nil {
				t.Fatalf("Failed to create request %v", err)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var responseBody jsonResponse
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Fatalf("Failed to Unmarshall response body: %v", err)
			}

			assert.Equal(t, test.args.httpResponseBodyCode, w.Code)
			assert.Equal(t, test.args.httpResponseBody, responseBody)
		})
	}
}

func TestDeleteUserByUsername(t *testing.T) {
	tests := []test{
		{
			name: "Delete user by Username",
			args: args{
				sqlStatement: ".*",
				sqlReturnRows: [][]any{
					{1, "Obi", "obi@example.com"},
				},
				httpRequestURL:       "/?username=Obi",
				httpRequestMethod:    "DELETE",
				httpResponseBodyCode: http.StatusOK,
				httpResponseBody: jsonResponse{
					Status:  "success",
					Message: "User deleted successfully.",
				},
			},
		},
		{
			name: "Delete non-existent user by username",
			args: args{
				sqlStatement:         ".*",
				sqlReturnRows:        [][]any{},
				httpRequestURL:       "/?username=two",
				httpRequestMethod:    "DELETE",
				httpResponseBodyCode: http.StatusBadRequest,
				httpResponseBody: jsonResponse{
					Status:  "error",
					Message: "User does not exist.",
				},
			},
		},
		{
			name: "Call delete without username query",
			args: args{
				sqlStatement:         ".*",
				sqlReturnRows:        [][]any{},
				httpRequestURL:       "/",
				httpRequestMethod:    "DELETE",
				httpResponseBodyCode: http.StatusBadRequest,
				httpResponseBody: jsonResponse{
					Status:  "error",
					Message: "You must/can only specify a user to delete.",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := InitDB(t)
			rawDB, err := db.DB()
			if err != nil {
				t.Fatalf("Unable to get sql.DB from gorm.DB, %v", err)
			}
			defer rawDB.Close()

			rows := sqlmock.NewRows([]string{"id", "username", "email"})
			for _, row := range test.args.sqlReturnRows {
				rows.AddRow(row[0], row[1], row[2])
			}

			mock.ExpectQuery(test.args.sqlStatement).WillReturnRows(rows)
			mock.ExpectBegin()
			mock.ExpectExec(test.args.sqlStatement).WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectCommit()

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.Use(gin.Logger())
			r.DELETE("/", func(c *gin.Context) {
				DeleteUserByUsername(c, db)
			})

			// Parse the URL to ensure query parameters are correctly included
			parsedURL, err := url.Parse(test.args.httpRequestURL)
			if err != nil {
				t.Fatalf("Failed to parse request URL %v", err)
			}

			req, err := http.NewRequest(test.args.httpRequestMethod, parsedURL.String(), nil)
			if err != nil {
				t.Fatalf("Failed to create request %v", err)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var responseBody jsonResponse
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Fatalf("Failed to Unmarshall response body: %v", err)
			}

			assert.Equal(t, test.args.httpResponseBodyCode, w.Code)
			assert.Equal(t, test.args.httpResponseBody, responseBody)
		})
	}
}
