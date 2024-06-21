package handlers

import (
	"encoding/json"
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
	sqlStatement      string
	returnRows        [][]any
	httpRequestURL    string
	httpRequestMethod string
	httpCode          int
	httpResponse      jsonResponse
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

func TestGetAll(t *testing.T) {
	db, mock := InitDB(t)

	tests := []test{
		{
			name: "Get all 2 rows",
			args: args{
				sqlStatement: ".*",
				returnRows: [][]any{
					{1, "Obi", "obi@example.com"},
					{2, "Marry", "marry@example.com"},
				},
				httpRequestURL:    "/",
				httpRequestMethod: "GET",
				httpCode:          http.StatusOK,
				httpResponse: jsonResponse{
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
				returnRows: [][]any{
					{1, "Obi", "obi@example.com"},
				},
				httpRequestURL:    "/?username=Obi",
				httpRequestMethod: "GET",
				httpCode:          http.StatusOK,
				httpResponse: jsonResponse{
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
			rows := sqlmock.NewRows([]string{"id", "username", "email"})
			for _, row := range test.args.returnRows {
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

			assert.Equal(t, test.args.httpCode, w.Code)
			assert.Equal(t, test.args.httpResponse, responseBody)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	db, mock := InitDB(t)

	tests := []test{
		{
			name: "Get user by ID",
			args: args{
				sqlStatement: ".*",
				returnRows: [][]any{
					{1, "Obi", "obi@example.com"},
				},
				httpRequestURL:    "/1",
				httpRequestMethod: "GET",
				httpCode:          http.StatusOK,
				httpResponse: jsonResponse{
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
				returnRows: [][]any{
					{1, "Obi", "obi@example.com"},
				},
				httpRequestURL:    "/one",
				httpRequestMethod: "GET",
				httpCode:          http.StatusBadRequest,
				httpResponse: jsonResponse{
					Status:  "error",
					Message: "UserID must be a positive interger.",
					Data:    nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"id", "username", "email"})
			for _, row := range test.args.returnRows {
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

			assert.Equal(t, test.args.httpCode, w.Code)
			assert.Equal(t, test.args.httpResponse, responseBody)
		})
	}
}
