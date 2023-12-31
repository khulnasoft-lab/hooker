package actions

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClient_Init(t *testing.T) {
	ec := HTTPClient{}
	require.NoError(t, ec.Init())
}

func TestHTTPClient_GetName(t *testing.T) {
	ec := HTTPClient{Name: "my-http-action"}
	require.NoError(t, ec.Init())
	require.Equal(t, "my-http-action", ec.GetName())
}

func TestHTTPClient_Send(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		inputEvent     string
		bodyFile       string
		bodyContent    string
		testServerFunc http.HandlerFunc
		expectedError  string
	}{
		{
			name:   "happy path method get",
			method: http.MethodGet,
			testServerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, []string{"bar", "baz"}, r.Header.Values("fookey"))
				assert.Empty(t, r.Header.Get("Hooker-Event")) // no event sent
			},
		},
		{
			name:       "happy path method post with body file, string input event",
			method:     http.MethodPost,
			bodyFile:   "goldens/validbody.txt",
			inputEvent: "foo bar baz header",
			testServerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, []string{"bar", "baz"}, r.Header.Values("fookey"))
				assert.Equal(t, "Zm9vIGJhciBiYXogaGVhZGVy", r.Header.Get("Hooker-Event"))

				b, _ := ioutil.ReadAll(r.Body)
				assert.Equal(t, "foo bar baz body", string(b))
			},
		},
		{
			name:        "happy path method post with body content, string input event",
			method:      http.MethodPost,
			bodyContent: "foo bar baz body",
			inputEvent:  "foo bar baz header",
			testServerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, []string{"bar", "baz"}, r.Header.Values("fookey"))
				assert.Equal(t, "Zm9vIGJhciBiYXogaGVhZGVy", r.Header.Get("Hooker-Event"))

				b, _ := ioutil.ReadAll(r.Body)
				assert.Equal(t, "foo bar baz body", string(b))
			},
		},
		{
			name:     "happy path method post, json input event",
			method:   http.MethodPost,
			bodyFile: "goldens/validbody.txt",
			inputEvent: `{
	"argsNum": 2
}`,
			testServerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, []string{"bar", "baz"}, r.Header.Values("fookey"))
				assert.Equal(t, "ewoJImFyZ3NOdW0iOiAyCn0=", r.Header.Get("Hooker-Event"))

				b, _ := ioutil.ReadAll(r.Body)
				assert.Equal(t, "foo bar baz body", string(b))
			},
		},
		{
			name:   "happy path method post, with relative body content, json input event",
			method: http.MethodPost,
			bodyContent: `argsNum: event.input.argsNum
eventID: event.input.eventID`,
			inputEvent: `{
	"argsNum": 2,
	"eventID": "TRC-2"
}`,
			testServerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, []string{"bar", "baz"}, r.Header.Values("fookey"))
				assert.Equal(t, "ewoJImFyZ3NOdW0iOiAyLAoJImV2ZW50SUQiOiAiVFJDLTIiCn0=", r.Header.Get("Hooker-Event"))

				b, _ := ioutil.ReadAll(r.Body)
				assert.Equal(t, `argsNum: 2
eventID: TRC-2`, string(b))
			},
		},
		{
			name:   "happy path method post, with relative body content, string input event",
			method: http.MethodPost,
			bodyContent: `event1: event.input
event1: event.input`,
			inputEvent: `"argsNum": 2, "eventID": "TRC-2"`,
			testServerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, []string{"bar", "baz"}, r.Header.Values("fookey"))
				assert.Equal(t, "ImFyZ3NOdW0iOiAyLCAiZXZlbnRJRCI6ICJUUkMtMiI=", r.Header.Get("Hooker-Event"))

				b, _ := ioutil.ReadAll(r.Body)
				assert.Equal(t, `event1: "argsNum": 2, "eventID": "TRC-2"
event1: "argsNum": 2, "eventID": "TRC-2"`, string(b))
			},
		},
		{
			name:   "sad path method get - server unavailable",
			method: http.MethodGet,
			testServerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("internal server error"))
			},
			expectedError: "http status NOT OK: HTTP 500 Internal Server Error, response: internal server error",
		},
		{
			name:          "sad path method get - bad url",
			method:        http.MethodGet,
			expectedError: `Get "path-to-nowhere": unsupported protocol scheme ""`,
		},
		{
			name:          "sad path, body file not found",
			method:        http.MethodPost,
			bodyFile:      "invalid.txt",
			expectedError: "unable to read body file: invalid.txt, err: open invalid.txt: no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var testUrl *url.URL
			if tc.testServerFunc != nil {
				ts := httptest.NewServer(tc.testServerFunc)
				testUrl, _ = url.Parse(ts.URL)
			} else {
				testUrl, _ = url.Parse("path-to-nowhere")
			}

			ec := HTTPClient{
				URL:         testUrl,
				Method:      tc.method,
				Headers:     map[string][]string{"fookey": {"bar", "baz"}},
				BodyFile:    tc.bodyFile,
				BodyContent: tc.bodyContent,
			}

			switch {
			case tc.expectedError != "":
				require.EqualError(t, ec.Send(map[string]string{"description": "foo bar baz header"}), tc.expectedError, tc.name)
			default:
				require.NoError(t, ec.Send(map[string]string{"description": tc.inputEvent}), tc.name)
			}
		})
	}
}
