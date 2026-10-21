package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

"strconv"
"fmt"
	"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

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


func TestCafeSearch(t *testing.T) {
	requests := []struct {
		search    string // Передаваемое значение search
		wantCount int    // Ожидаемое количество кафе в ответе
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}



	for _, tc := range requests {
		
		t.Run(fmt.Sprintf("search=%s", tc.search), func(t *testing.T) {
			url := fmt.Sprintf("/cafes?city=moscow&search=%s", tc.search)
			
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatalf("Не удалось создать запрос: %v", err)
			}

		rr := httptest.NewRecorder()
			
			mainHandle(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Ожидался статус 200 OK, получен %d", rr.Code)
			}

			bodyStr := strings.TrimSpace(rr.Body.String())

			var cafes []string
			if bodyStr != "" {
				cafes = strings.Split(bodyStr, ",")
			}

			if len(cafes) != tc.wantCount {
				t.Errorf("Ожидалось количество кафе %d, получено %d (ответ сервера: %q)", tc.wantCount, len(cafes), bodyStr)
			}


			for _, cafeName := range cafes {

				lowerCafeName := strings.ToLower(cafeName)
				lowerSearch := strings.ToLower(tc.search)

				if !strings.Contains(lowerCafeName, lowerSearch) {
					t.Errorf("Кафе %q не содержит поисковую строку %q", cafeName, tc.search)
				}
			}
		})
	}
}


func TestCafeCount(t *testing.T) {
city := "tula" 
	totalCafes := len(cafeList[city])

	requests := []struct {
        count int   // передаваемое значение count
        want  int   // ожидаемое количество кафе в ответе
    }{
       		{count: 0, want: 0},
		{count: 1, want: min(totalCafes, 1)},
		{count: 2, want: min(totalCafes, 2)},
		{count: 100, want: min(totalCafes, 100)},
    }
	for _, tt := range requests {
		t.Run("count_"+strconv.Itoa(tt.count), func(t *testing.T) {
			// Формируем URL-запрос с параметрами city и count
			target := fmt.Sprintf("/cafes?city=%s&count=%s", city, strconv.Itoa(tt.count))
			req, err := http.NewRequest(http.MethodGet, target, nil)
			require.NoError(t, err) // Проверяем, что запрос создался без ошибок

			responseRecorder := httptest.NewRecorder()
			
			mainHandle(responseRecorder, req)

			require.Equal(t, http.StatusOK, responseRecorder.Code)

			body := responseRecorder.Body.String()

			var gotCafesCount int

			if body != "" {
				cafesSlice := strings.Split(body, ",")
				gotCafesCount = len(cafesSlice)
			} else {
				gotCafesCount = 0
			}

			assert.Equal(t, tt.want, gotCafesCount, "Количество кафе в ответе не совпадает с ожидаемым")
		})}

}

