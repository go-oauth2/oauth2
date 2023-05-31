package server

import (
    "fmt"
    "net/http"
    "net/url"
    "testing"
)

var clientBasicHandlerTests = []struct {
    Username         string
    Password         string
    ExpectedUsername string
    ExpectedPassword string
    ExpectError      bool
}{
    {
        Username:         "username",
        Password:         "password",
        ExpectedUsername: "username",
        ExpectedPassword: "password",
        ExpectError:      false,
    },
    {
        Username:         url.QueryEscape("+%25%26%2B%C2%A3%E2%82%AC"),
        Password:         url.QueryEscape("+%25%26%2B%C2%A3%E2%82%AC"),
        ExpectedUsername: "+%25%26%2B%C2%A3%E2%82%AC",
        ExpectedPassword: "+%25%26%2B%C2%A3%E2%82%AC",
        ExpectError:      false,
    },
    {
        Username:         "% +'/€$",
        Password:         "% +'/€$",
        ExpectedUsername: "",
        ExpectedPassword: "",
        ExpectError:      true,
    },
    {
        Username:         "+%25%26%2B%C2%A3%E2%82%AC",
        Password:         "+%25%26%2B%C2%A3%E2%82%AC",
        ExpectedUsername: " %&+£€",
        ExpectedPassword: " %&+£€",
        ExpectError:      false,
    },
}

func TestClientBasicHandler(t *testing.T) {
    for i, test := range clientBasicHandlerTests {
        t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
            req, _ := http.NewRequest("", "", nil)
            req.SetBasicAuth(test.Username, test.Password)

            username, password, err := ClientBasicHandler(req)
            if test.ExpectError && err == nil {
                t.Error("expected error, got nil")
            } else if !test.ExpectError && err != nil {
                t.Errorf("unexpected error: %s", err.Error())
            }

            if test.ExpectedUsername != username {
                t.Errorf("unexpected username, got %s, want %s", username, test.ExpectedUsername)
            }

            if test.ExpectedPassword != password {
                t.Errorf("unexpected password, got %s, want %s", password, test.ExpectedPassword)
            }
        })
    }
}
