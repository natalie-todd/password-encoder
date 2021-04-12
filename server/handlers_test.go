package server

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"password-encoder/mocks"
	"password-encoder/service"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHandler_CreateHash(t *testing.T) {

	start := time.Now()

	type args struct {
		URL          string
		passWord     string
		startTime    time.Time
		fiveSecTimer *time.Timer
	}
	type mox struct {
		hashedPasswords []string
		calls           int
	}
	type expected struct {
		response   int
		httpStatus int
	}
	tests := []struct {
		name     string
		args     args
		mox      mox
		expected expected
	}{
		{
			name:     "Single Password 200 Success",
			args:     args{URL: "/hash", passWord: "angryMonkey", startTime: start, fiveSecTimer: time.NewTimer(5 * time.Second)},
			mox:      mox{hashedPasswords: []string{"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="}, calls: 1},
			expected: expected{response: 1, httpStatus: http.StatusOK},
		},
		{
			name:     "Multiple Passwords 200 Success",
			args:     args{URL: "/hash", passWord: "angryMonkey", startTime: start, fiveSecTimer: time.NewTimer(5 * time.Second)},
			mox:      mox{hashedPasswords: []string{"hash1", "hash2", "hash3"}, calls: 1},
			expected: expected{response: 4, httpStatus: http.StatusOK},
		},
		{
			name:     "400 Bad Request",
			args:     args{URL: "/hash", passWord: ""},
			mox:      mox{calls: 0},
			expected: expected{response: 0, httpStatus: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		mockController := gomock.NewController(t)
		mockService := mocks.NewMockServicer(mockController)
		testHandler := InitializeHandler(mockService)

		testRouter := mux.NewRouter()
		testRouter.HandleFunc(tt.args.URL, testHandler.CreateHash).Methods(http.MethodPost)

		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set("password", tt.args.passWord)
			responseWriter := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, tt.args.URL, strings.NewReader(data.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
			//TODO: remove gomock.Any()
			mockService.EXPECT().CalculateHashAndDuration(gomock.Any(), gomock.Any(), tt.args.passWord).Return().Times(tt.mox.calls)
			mockService.EXPECT().GetHashedPasswords().Return(tt.mox.hashedPasswords).Times(tt.mox.calls)

			testRouter.ServeHTTP(responseWriter, r)
			assert.Equal(t, tt.expected.httpStatus, responseWriter.Code)
			decoder := json.NewDecoder(responseWriter.Body)
			decoder.Decode(&tt.expected.response)
			assert.Equal(t, tt.expected.response, tt.expected.response)
		})
		mockController.Finish()
	}
}

func TestHandler_GetHash(t *testing.T) {
	type args struct {
		URL        string
		queryParam string
	}
	type mox struct {
		hashedPasswords []string
		calls           int
	}
	type expected struct {
		response   string
		httpStatus int
	}
	tests := []struct {
		name     string
		args     args
		mox      mox
		expected expected
	}{
		{
			name:     "First Hash 200 Success",
			args:     args{URL: "/hash", queryParam: "1"},
			mox:      mox{hashedPasswords: []string{"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="}, calls:1},
			expected: expected{response: "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==", httpStatus: http.StatusOK},
		},
		{
			name:     "Last Hash 200 Success",
			args:     args{URL: "/hash", queryParam: "3"},
			mox:      mox{hashedPasswords: []string{"hash1", "hash2", "hash3"}, calls:1},
			expected: expected{response: "hash3", httpStatus: http.StatusOK},
		},
		{
			name:     "400 Bad Request",
			args:     args{URL: "/hash", queryParam: "three"},
			mox:      mox{calls:0},
			expected: expected{response: "", httpStatus: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		mockController := gomock.NewController(t)
		mockService := mocks.NewMockServicer(mockController)
		testHandler := InitializeHandler(mockService)

		testRouter := mux.NewRouter()
		testRouter.HandleFunc(fmt.Sprintf("%s/{id}", tt.args.URL), testHandler.GetHash).Methods(http.MethodGet)

		t.Run(tt.name, func(t *testing.T) {
			responseWriter := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", tt.args.URL, tt.args.queryParam), nil)

			mockService.EXPECT().GetHashedPasswords().Return(tt.mox.hashedPasswords).Times(tt.mox.calls)

			testRouter.ServeHTTP(responseWriter, r)
			assert.Equal(t, tt.expected.httpStatus, responseWriter.Code)
			decoder := json.NewDecoder(responseWriter.Body)
			decoder.Decode(&tt.expected.response)
			assert.Equal(t, tt.expected.response, tt.expected.response)
		})
		mockController.Finish()
	}
}

func TestHandler_CalculateStats(t *testing.T) {
	type args struct {
		URL string
	}
	type mox struct {
		stats *service.Stats
	}
	type expected struct {
		response   *service.Stats
		httpStatus int
	}
	tests := []struct {
		name     string
		args     args
		mox      mox
		expected expected
	}{
		{
			name: "200 Success Case",
			args: args{URL: "/stats"},
			mox: mox{stats: &service.Stats{
				Total:   3,
				Average: 0,
			}},
			expected: expected{response: &service.Stats{
				Total:   3,
				Average: 0,
			}, httpStatus: http.StatusOK},
		},
	}
	for _, tt := range tests {
		mockController := gomock.NewController(t)
		mockService := mocks.NewMockServicer(mockController)
		testHandler := InitializeHandler(mockService)

		testRouter := mux.NewRouter()
		testRouter.HandleFunc(tt.args.URL, testHandler.GetStats).Methods(http.MethodGet)

		t.Run(tt.name, func(t *testing.T) {
			responseWriter := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.args.URL, nil)
			mockService.EXPECT().CalculateStats().Return(tt.mox.stats).Times(1)

			testRouter.ServeHTTP(responseWriter, r)
			assert.Equal(t, tt.expected.httpStatus, responseWriter.Code)
		})
		mockController.Finish()
	}
}
