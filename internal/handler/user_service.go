package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"
)

type UserCredentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type UserService interface {
    OnCreateUser(credentials *UserCredentials)
    OnUpdateUser(credentials *UserCredentials)
}

type httpUserService struct {
    baseURL    string
    httpClient *http.Client
}

func NewHTTPUserService(baseURL string) UserService {
    return &httpUserService{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 5 * time.Second,
        },
    }
}

func (s *httpUserService) post(path string, payload any) {
    body, _ := json.Marshal(payload)
    req, err := http.NewRequest(http.MethodPost, s.baseURL+path, bytes.NewReader(body))
    if err != nil {
        return
    }
    req.Header.Set("Content-Type", "application/json")
    resp, err := s.httpClient.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()
}

func (s *httpUserService) OnCreateUser(credentials *UserCredentials) {
    s.post("/users/create", credentials)
}

func (s *httpUserService) OnUpdateUser(credentials *UserCredentials) {
    s.post("/users/update", credentials)
}

