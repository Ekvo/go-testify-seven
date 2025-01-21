package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

const (
	count = "0"
	city  = "moscow"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cafe?count=%s&city=%s", count, city), nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.NotEmpty(t, responseRecorder.Code)
	require.NotEmpty(t, responseRecorder.Body)

	bodySTR := responseRecorder.Body.String()
	countSTR := req.URL.Query().Get("count")
	citySTR := req.URL.Query().Get("city")

	if !assert.Equal(t, http.StatusOK, responseRecorder.Code) {

		switch {

		case assert.Empty(t, countSTR):
			require.Equal(t, "count missing", bodySTR)

		case assert.Error(t, itNumber(countSTR)):
			require.Equal(t, "wrong count value", bodySTR)

		default:
			assert.NotEmpty(t, citySTR)
			assert.Equal(t, "moscow", citySTR)
			require.Equal(t, "wrong city value", bodySTR)

		}

		require.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	} else {

		require.NotEmpty(t, cafeList["moscow"])

		totalCount := len(cafeList["moscow"])
		cafes := strings.Split(bodySTR, ",")

		require.LessOrEqual(t, len(cafes), totalCount)
	}
}

func itNumber(number string) (err error) {
	_, err = strconv.Atoi(number)
	return err
}
