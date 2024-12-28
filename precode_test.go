package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerCorrectRequest(t *testing.T) {
	//=== RUN   TestMainHandlerCorrectRequest
	//--- PASS: TestMainHandlerCorrectRequest (0.00s)
	//PASS
	require.NotEmpty(t, cafeList)
	require.NotEmpty(t, cafeList["moscow"])

	loop := len(cafeList["moscow"])
	for num := 1; num <= loop; num++ {

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cafe?count=%d&city=moscow", num), nil)

		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(mainHandle)
		handler.ServeHTTP(responseRecorder, req)

		require.NotEmpty(t, responseRecorder.Code)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.NotEmpty(t, responseRecorder.Body)
	}
}

func TestMainHandlerIncorrectCity(t *testing.T) {
	//=== RUN   TestMainHandlerIncorrectCity
	//--- PASS: TestMainHandlerIncorrectCity (0.00s)
	//PASS

	//for create random city
	letters := []string{"a", "b", "c", "d", "e", "f", "g"}
	loop := rand.Intn(100_000)

	for num := 1; num <= loop; num++ {

		rand.Shuffle(len(letters), func(i, j int) {
			letters[i], letters[j] = letters[j], letters[i]
		})
		city := strings.Join(letters, "")

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cafe?count=3&city=%s", city), nil)

		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(mainHandle)
		handler.ServeHTTP(responseRecorder, req)

		require.NotEmpty(t, responseRecorder.Code)
		require.NotEmpty(t, responseRecorder.Body)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.Equal(t, "wrong city value", responseRecorder.Body.String())
	}
}

func TestMainHandlerWhenCountMoreThanTotalCafeList(t *testing.T) {
	//=== RUN   TestMainHandlerWhenCountMoreThanTotalCafeList
	//--- PASS: TestMainHandlerWhenCountMoreThanTotalCafeList (7.80s)
	//PASS
	require.NotEmpty(t, cafeList)
	require.NotEmpty(t, cafeList["moscow"])

	start := len(cafeList["moscow"]) + 1
	loop := rand.Intn(1_000_000) + start

	for num := start; num < loop; num++ {

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cafe?count=%d&city=moscow", num), nil)

		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(mainHandle)
		handler.ServeHTTP(responseRecorder, req)

		require.NotEmpty(t, responseRecorder.Code)
		require.NotEmpty(t, responseRecorder.Body)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)

		bodyStr := responseRecorder.Body.String()
		listCafe := strings.Split(bodyStr, ",")

		require.LessOrEqual(t, len(listCafe), totalCount)
	}
}
