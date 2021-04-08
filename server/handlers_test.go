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
			name:     "Single Password 200 Success Case",
			args:     args{URL: "/hash", passWord: "angryMonkey", startTime: start, fiveSecTimer: time.NewTimer(5 * time.Second)},
			mox:      mox{hashedPasswords: []string{"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="}},
			expected: expected{response: 1, httpStatus: 200},
		},
		{
			name:     "Multiple Passwords 200 Success Case",
			args:     args{URL: "/hash", passWord: "angryMonkey", startTime: start, fiveSecTimer: time.NewTimer(5 * time.Second)},
			mox:      mox{hashedPasswords: []string{"hash1", "hash2", "hash3"}},
			expected: expected{response: 4, httpStatus: 200},
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
			mockService.EXPECT().CalculateHashAndDuration(gomock.Any(), gomock.Any(), tt.args.passWord).Return().Times(1)
			mockService.EXPECT().GetHashedPasswords().Return(tt.mox.hashedPasswords).Times(1)

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
		queryParam int
	}
	type mox struct {
		hashedPasswords []string
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
			name:     "First Hash 200 Success Case",
			args:     args{URL: "/hash", queryParam: 1},
			mox:      mox{hashedPasswords: []string{"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="}},
			expected: expected{response: "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==", httpStatus: 200},
		},
		{
			name:     "Last Hash 200 Success Case",
			args:     args{URL: "/hash", queryParam: 3},
			mox:      mox{hashedPasswords: []string{"hash1", "hash2", "hash3"}},
			expected: expected{response: "hash3", httpStatus: 200},
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
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", tt.args.URL, tt.args.queryParam), nil)
			mockService.EXPECT().GetHashedPasswords().Return(tt.mox.hashedPasswords).Times(1)

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
			name:     "200 Success Case",
			args:     args{URL: "/stats"},
			mox:      mox{stats:&service.Stats{
				Total:   3,
				Average: 0,
			}},
			expected: expected{response:&service.Stats{
				Total:   3,
				Average: 0,
			}, httpStatus: 200},
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
