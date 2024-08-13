package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/svetlana-mel/url-shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/svetlana-mel/url-shortener/internal/lib/logger/discardlogger"
	"github.com/svetlana-mel/url-shortener/internal/repository"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name           string
		alias          string
		url            string
		urlHasAlias    bool
		expectedError  string
		saveMockError  error
		aliasMockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Success: empty alias",
			alias: "",
			url:   "https://pinterest.com",
		},
		{
			name:          "Empty URL",
			url:           "",
			expectedError: "field URL is a required field",
		},
		// {
		// 	name:          "Invalid URL",
		// 	url:           "some invalid URL",
		// 	alias:         "some_alias",
		// 	expectedError: "field URL is not a valid URL",
		// },
		// // errors in mocks
		{
			name:          "SaveURL Error: error alias already exists",
			alias:         "test_alias",
			url:           "https://google.com",
			expectedError: "generator error: get dublicate alias",
			saveMockError: repository.ErrAliasAlreadyExists,
		},
		{
			name:          "SaveURL Error: error unexpected",
			alias:         "test_alias",
			url:           "https://google.com",
			expectedError: "failed to save url-alias pair",
			saveMockError: errors.New("unexpected error"),
		},
		{
			name:        "GetAlias success, url already has alias",
			alias:       "new_alias",
			urlHasAlias: true,
			url:         "https://google.com",
		},
		{
			name:           "GetAlias: error unexpected",
			alias:          "test_alias",
			url:            "https://google.com",
			expectedError:  "failed to check alias existence",
			aliasMockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			urlSaverMock := mocks.NewURLSaver(t)

			// задаем ожидаемые параметры для возвращаемых значений моков
			if tc.url != "" {
				if tc.urlHasAlias {
					urlSaverMock.On("GetAlias", tc.url, mock.AnythingOfType("string")).
						Return("existed_alias", nil)
				} else if tc.aliasMockError != nil {
					urlSaverMock.On("GetAlias", tc.url, mock.AnythingOfType("string")).
						Return("", tc.aliasMockError)
				} else {
					urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
						Return(func() error {
							if tc.saveMockError != nil {
								return tc.saveMockError
							}
							return nil
						}())
					urlSaverMock.On("GetAlias", tc.url, mock.AnythingOfType("string")).
						Return("", nil)
				}
			}

			handler := New(urlSaverMock, discardlogger.NewDiscardLogger())

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			request, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			responseRecorder := httptest.NewRecorder()

			handler.ServeHTTP(responseRecorder, request)

			// unmarshal body
			body := responseRecorder.Body.String()
			var saveResponse SaveResponse
			err = json.Unmarshal([]byte(body), &saveResponse)
			require.NoError(t, err)

			// checks for the handler result
			require.Equal(t, http.StatusOK, responseRecorder.Code)
			require.Equal(t, tc.expectedError, saveResponse.Error)
		})
	}
}
