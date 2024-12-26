package controllers

import (
    "CatVotingApp/models"
    "encoding/json"
    "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/server/web/context"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
    "errors"
)

func init() {
    web.AppConfig.Set("cat_api_url", "https://api.thecatapi.com/v1")
    web.AppConfig.Set("cat_api_key", "test_key")
}

func setupTestController(r *http.Request) (*CatController, *httptest.ResponseRecorder) {
    w := httptest.NewRecorder()
    c := &CatController{}
    
    ctx := context.NewContext()
    ctx.Reset(w, r)
    c.Init(ctx, "", "", nil)
    
    return c, w
}

// Helper function to mock request body
func setRequestBody(controller *CatController, body []byte) {
    controller.Ctx.Input.RequestBody = body
}

func TestCatController_GetCats(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/cats", nil)
    controller, w := setupTestController(r)
    
    // Mock worker response
    go func() {
        req := <-catsChan
        req.ResponseChan <- []models.Cat{{ID: "test123", URL: "http://example.com/cat.jpg"}}
    }()
    
    controller.GetCats()
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "success", response["status"])
}

// func TestCatController_GetCats_Error(t *testing.T) {
//     r, _ := http.NewRequest("GET", "/api/cats", nil)
//     controller, w := setupTestController(r)

//     go func() {
//         req := <-catsChan
//         req.ErrorChan <- ErrMock("Failed to fetch cats")
//     }()

//     controller.GetCats()

//     assert.Equal(t, http.StatusOK, w.Code)

//     var response map[string]interface{}
//     err := json.Unmarshal(w.Body.Bytes(), &response)
//     assert.NoError(t, err)
//     assert.Equal(t, "error", response["status"])
// }

func TestCatController_GetCats_Error(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/cats", nil)
    controller, w := setupTestController(r)

    // Simulate worker sending an error
    go func() {
        req := <-catsChan
        req.ErrorChan <- errors.New("Failed to fetch cats")
    }()

    controller.GetCats()

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "error", response["status"])
    assert.Equal(t, "Failed to fetch cats", response["message"])
}

func TestCatController_VoteCat(t *testing.T) {
    voteReq := models.VoteRequest{
        ImageID:  "test123",
        ImageURL: "http://example.com/cat.jpg",
        Vote:     "love",
    }
    
    body, _ := json.Marshal(voteReq)
    r, _ := http.NewRequest("POST", "/api/vote", nil)
    controller, w := setupTestController(r)
    
    setRequestBody(controller, body)
    
    // Mock worker response
    go func() {
        req := <-voteChan
        req.ReqChan.ResponseChan <- true
    }()
    
    controller.VoteCat()
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "success", response["status"])
}

func TestCatController_VoteCat_InvalidRequest(t *testing.T) {
    body := []byte(`{"invalid": "json"}`)
    r, _ := http.NewRequest("POST", "/api/vote", nil)
    controller, w := setupTestController(r)
    
    setRequestBody(controller, body)
    
    controller.VoteCat()
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "error", response["status"])
    assert.Equal(t, "Invalid request format", response["message"])
}

func TestCatController_VoteCat_Error(t *testing.T) {
    voteReq := models.VoteRequest{
        ImageID:  "test123",
        ImageURL: "http://example.com/cat.jpg",
        Vote:     "love",
    }
    
    body, _ := json.Marshal(voteReq)
    r, _ := http.NewRequest("POST", "/api/vote", nil)
    controller, w := setupTestController(r)
    
    setRequestBody(controller, body)
    
    // Mock worker error response
    go func() {
        req := <-voteChan
        req.ReqChan.ErrorChan <- ErrMock("Vote failed")
    }()
    
    controller.VoteCat()
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    
    assert.Equal(t, "error", response["status"])
    assert.Equal(t, "Vote failed", response["message"])
}

func TestCatController_RemoveFavorite(t *testing.T) {
    r, _ := http.NewRequest("DELETE", "/api/favorites/test123", nil)
    controller, w := setupTestController(r)
    
    // Set the ID parameter
    controller.Ctx.Input.SetParam(":id", "test123")
    
    // Mock worker response
    go func() {
        req := <-removeFavChan
        req.ReqChan.ResponseChan <- true
    }()
    
    controller.RemoveFavorite()
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    
    assert.Equal(t, "success", response["status"])
}

func TestCatController_RemoveFavorite_MissingID(t *testing.T) {
    r, _ := http.NewRequest("DELETE", "/api/favorites/", nil)
    controller, w := setupTestController(r)
    
    controller.RemoveFavorite()
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    
    assert.Equal(t, "error", response["status"])
    assert.Equal(t, "Image ID is required", response["message"])
}

func TestCatController_GetBreeds(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/breeds", nil)
    controller, w := setupTestController(r)
    
    // Mock worker response
    go func() {
        req := <-breedsChan
        req.ResponseChan <- []models.Breed{{ID: "test123", Name: "Test Breed"}}
    }()
    
    controller.GetBreeds()
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "success", response["status"])
}




func TestCatController_GetBreedDetails(t *testing.T) {
    // Create request with query parameter
    r, _ := http.NewRequest("GET", "/api/breed?id=abys", nil)
    controller, w := setupTestController(r)
    
    // Mock worker response
    go func() {
        req := <-breedDetailsChan
        req.ReqChan.ResponseChan <- models.Breed{ID: "abys", Name: "Abyssinian"}
    }()
    
    controller.GetBreedDetails()
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "success", response["status"])
    
    // Add more specific assertions to verify the response data
    data, ok := response["data"].(map[string]interface{})
    assert.True(t, ok)
    assert.Equal(t, "abys", data["ID"])
    assert.Equal(t, "Abyssinian", data["Name"])
}


func TestCatController_GetBreedDetails_MissingID(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/breed", nil)
    controller, w := setupTestController(r)

    controller.GetBreedDetails()

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "error", response["status"])
    assert.Equal(t, "breed id is required", response["message"])
}


func TestCatController_GetFavorites(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/favorites", nil)
    controller, w := setupTestController(r)
    
    // Mock worker response
    go func() {
        req := <-favoritesChan
        req.ResponseChan <- []models.FavoriteImage{{ID: "test123", URL: "http://example.com/cat.jpg"}}
    }()
    
    controller.GetFavorites()
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "success", response["status"])
}

func TestCatController_GetFavorites_EmptyResponse(t *testing.T) {
    r, _ := http.NewRequest("GET", "/api/favorites", nil)
    controller, w := setupTestController(r)

    go func() {
        req := <-favoritesChan
        req.ResponseChan <- []models.FavoriteImage{}
    }()

    controller.GetFavorites()

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "success", response["status"])
    assert.Equal(t, 0, len(response["data"].([]interface{})))
}


func ErrMock(message string) error {
    return errors.New(message)
}