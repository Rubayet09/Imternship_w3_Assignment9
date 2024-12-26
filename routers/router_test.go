package routers

import (
    "CatVotingApp/controllers"
    "testing"
    "github.com/beego/beego/v2/server/web"
	 "github.com/beego/beego/v2/server/web/context"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
)

// Initialize routes for testing
func initTestRoutes() {
    web.Router("/", &controllers.CatController{})
    
    ns := web.NewNamespace("/api",
        web.NSRouter("/cats", &controllers.CatController{}, "get:GetCats"),
        web.NSRouter("/breeds", &controllers.CatController{}, "get:GetBreeds"),
        web.NSRouter("/breed", &controllers.CatController{}, "get:GetBreedDetails"),
        web.NSRouter("/vote", &controllers.CatController{}, "post:VoteCat"),
        web.NSRouter("/favorites", &controllers.CatController{}, "get:GetFavorites"),
        web.NSRouter("/favorites/:id", &controllers.CatController{}, "delete:RemoveFavorite"),
    )
    
    web.AddNamespace(ns)
    web.SetStaticPath("/static", "static")
}

func TestRouterInit(t *testing.T) {
    // Initialize router
    initTestRoutes()
    
    // Test cases for each route
    tests := []struct {
        name           string
        method         string
        path           string
        expectedStatus int
    }{
        {
            name:           "Home Route",
            method:         "GET",
            path:           "/",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Get Cats Route",
            method:         "GET",
            path:           "/api/cats",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Get Breeds Route",
            method:         "GET",
            path:           "/api/breeds",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Get Breed Details Route",
            method:         "GET",
            path:           "/api/breed",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Vote Cat Route",
            method:         "POST",
            path:           "/api/vote",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Get Favorites Route",
            method:         "GET",
            path:           "/api/favorites",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Remove Favorite Route",
            method:         "DELETE",
            path:           "/api/favorites/123",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Static Files Route",
            method:         "GET",
            path:           "/static/test.css",
            expectedStatus: http.StatusNotFound, // Will be 404 since file doesn't exist
        },
    }

    // Run tests for each route
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            r, err := http.NewRequest(tt.method, tt.path, nil)
            assert.NoError(t, err)

            w := httptest.NewRecorder()
            web.BeeApp.Handlers.ServeHTTP(w, r)
            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}


func TestNamespaceConfiguration(t *testing.T) {
    initTestRoutes()

    // Test API namespace routes
    apiRoutes := []string{
        "/api/cats",
        "/api/breeds",
        "/api/breed",
        "/api/vote",
        "/api/favorites",
        "/api/favorites/123",
    }

    for _, route := range apiRoutes {
        t.Run("Route "+route, func(t *testing.T) {
            // Create test request for the route
            r, err := http.NewRequest("GET", route, nil)
            assert.NoError(t, err)

            // Create response recorder to capture the response
            w := httptest.NewRecorder()

            // Create a new context using the correct package
            ctx := context.NewContext()
            ctx.Reset(w, r)

            // Verify the route exists by attempting to match it
            routeInfo, found := web.BeeApp.Handlers.FindRouter(ctx)
            assert.True(t, found, "Route should be registered: "+route)
            assert.NotNil(t, routeInfo, "Route info should not be nil: "+route)
        })
    }
}


func TestStaticPathConfiguration(t *testing.T) {
    initTestRoutes()
    
    // Create test request for static file
    r, err := http.NewRequest("GET", "/static/", nil)
    assert.NoError(t, err)
    
    w := httptest.NewRecorder()
    web.BeeApp.Handlers.ServeHTTP(w, r)
    
    // Should get directory listing or 403 depending on configuration
    assert.Contains(t, []int{http.StatusForbidden, http.StatusOK}, w.Code)
}

func TestMethodNotAllowed(t *testing.T) {
    initTestRoutes()
    
    tests := []struct {
        name           string
        method         string
        path           string
        expectedStatus []int // Allow multiple possible status codes
    }{
        {
            name:           "POST to GET route",
            method:         "POST",
            path:           "/api/cats",
            expectedStatus: []int{http.StatusMethodNotAllowed, http.StatusNotFound},
        },
        {
            name:           "GET to POST route",
            method:         "GET",
            path:           "/api/vote",
            expectedStatus: []int{http.StatusMethodNotAllowed, http.StatusNotFound},
        },
        {
            name:           "PUT to GET route",
            method:         "PUT",
            path:           "/api/breeds",
            expectedStatus: []int{http.StatusMethodNotAllowed, http.StatusNotFound},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            r, err := http.NewRequest(tt.method, tt.path, nil)
            assert.NoError(t, err)

            w := httptest.NewRecorder()
            web.BeeApp.Handlers.ServeHTTP(w, r)

            // Check if the response status is one of the expected statuses
            assert.Contains(t, tt.expectedStatus, w.Code, 
                "Expected status to be one of %v, but got %d", tt.expectedStatus, w.Code)
            
            // Also verify that the correct method works for this route
            correctMethod := "GET"
            if tt.path == "/api/vote" {
                correctMethod = "POST"
            }
            
            // Test the correct method
            r2, _ := http.NewRequest(correctMethod, tt.path, nil)
            w2 := httptest.NewRecorder()
            web.BeeApp.Handlers.ServeHTTP(w2, r2)
            
            assert.Equal(t, http.StatusOK, w2.Code,
                "Route should accept %s method", correctMethod)
        })
    }
}