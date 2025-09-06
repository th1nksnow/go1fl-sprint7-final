package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"

	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()

		reqString := fmt.Sprintf("/cafe?search=%s&city=%s", v.search, city)

		req := httptest.NewRequest("GET", reqString, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)
		if response.Body.String() == "" {
			assert.Equal(t, v.wantCount, 0)
			continue
		}
		cafes := strings.Split(response.Body.String(), ",")
		for _, cafe := range cafes {
			require.Equal(t, true, strings.Contains(strings.ToLower(cafe), strings.ToLower(v.search)))
		}
		assert.Equal(t, v.wantCount, len(cafes))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, 100},
	}

	for town, cafes := range cafeList {
		cafesCount := len(cafes)
		for _, v := range requests {
			response := httptest.NewRecorder()

			reqString := fmt.Sprintf("/cafe?count=%d&city=%s", v.count, town)

			req := httptest.NewRequest("GET", reqString, nil)
			handler.ServeHTTP(response, req)

			require.Equal(t, http.StatusOK, response.Code)
			if v.count == 0 {
				assert.Equal(t, "", response.Body.String())
				continue
			}
			if cafesCount < v.want {
				assert.Equal(t, cafesCount, len(strings.Split(response.Body.String(), ",")))
				continue
			}
			assert.Equal(t, v.want, len(strings.Split(response.Body.String(), ",")))
		}
	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}
