package models

import (
    "github.com/beego/beego/v2/server/web"
    "github.com/stretchr/testify/assert"
    "testing"
)

func init() {
    web.AppConfig.Set("cat_api_url", "https://api.thecatapi.com/v1")
    web.AppConfig.Set("cat_api_key", "test_key")
}

func TestFetchCats(t *testing.T) {
    cats, err := FetchCats()
    assert.NoError(t, err)
    assert.NotNil(t, cats)
}

func TestFetchCatsErrors(t *testing.T) {
    // Save original config values
    originalURL, err := web.AppConfig.String("cat_api_url")
    if err != nil {
        t.Fatal("Failed to get original cat_api_url:", err)
    }

    // Restore config in a defer
    defer func() {
        err := web.AppConfig.Set("cat_api_url", originalURL)
        if err != nil {
            t.Error("Failed to restore cat_api_url:", err)
        }
    }()

    // Test missing API URL
    err = web.AppConfig.Set("cat_api_url", "")
    if err != nil {
        t.Fatal("Failed to set empty cat_api_url:", err)
    }
    
    cats, err := FetchCats()
    assert.Error(t, err)
    assert.Nil(t, cats)

    // Test invalid API URL
    err = web.AppConfig.Set("cat_api_url", "http://invalid-url")
    if err != nil {
        t.Fatal("Failed to set invalid cat_api_url:", err)
    }
    
    cats, err = FetchCats()
    assert.Error(t, err)
    assert.Nil(t, cats)
}


func TestFetchBreeds(t *testing.T) {
    breeds, err := FetchBreeds()
    assert.NoError(t, err)
    assert.NotNil(t, breeds)
}




func TestSaveAndGetFavorites(t *testing.T) {
    // Clear existing favorites
    favorites = []FavoriteImage{} // Reset the global favorites slice
    
    // Test saving a favorite
    err := SaveFavorite("test123", "http://example.com/cat.jpg")
    assert.NoError(t, err)
    
    // Test getting favorites
    favs := GetFavorites()
    assert.Equal(t, 1, len(favs))
    assert.Equal(t, "test123", favs[0].ID)
    assert.Equal(t, "http://example.com/cat.jpg", favs[0].URL)
    
    // Test duplicate save
    err = SaveFavorite("test123", "http://example.com/cat.jpg")
    assert.Error(t, err)
    assert.Equal(t, "already in favorites", err.Error())
}

func TestRemoveFavorite(t *testing.T) {
    // Setup test data
    favorites = []FavoriteImage{
        {ID: "test123", URL: "http://example.com/cat.jpg"},
    }
    
    // Test removing existing favorite
    err := RemoveFavorite("test123")
    assert.NoError(t, err)
    assert.Equal(t, 0, len(GetFavorites()))
    
    // Test removing non-existent favorite
    err = RemoveFavorite("nonexistent")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "favorite not found")
}

func TestFetchBreedDetails(t *testing.T) {
    breed, err := FetchBreedDetails("abys")
    assert.NoError(t, err)
    assert.NotNil(t, breed)
    assert.Equal(t, "abys", breed.ID)
}

func TestFetchBreedDetailsErrors(t *testing.T) {
    // Save original config values
    originalURL, err := web.AppConfig.String("cat_api_url")
    if err != nil {
        t.Fatal("Failed to get original cat_api_url:", err)
    }
    originalKey, err := web.AppConfig.String("cat_api_key")
    if err != nil {
        t.Fatal("Failed to get original cat_api_key:", err)
    }

    // Restore config in a defer
    defer func() {
        err := web.AppConfig.Set("cat_api_url", originalURL)
        if err != nil {
            t.Error("Failed to restore cat_api_url:", err)
        }
        err = web.AppConfig.Set("cat_api_key", originalKey)
        if err != nil {
            t.Error("Failed to restore cat_api_key:", err)
        }
    }()

    // Test with invalid breed ID
    breed, err := FetchBreedDetails("invalid-breed-id")
    assert.Error(t, err)
    assert.Nil(t, breed)

    // Test with missing API URL
    err = web.AppConfig.Set("cat_api_url", "")
    if err != nil {
        t.Fatal("Failed to set empty cat_api_url:", err)
    }
    
    breed, err = FetchBreedDetails("abys")
    assert.Error(t, err)
    assert.Nil(t, breed)

    // Test with missing API key
    err = web.AppConfig.Set("cat_api_key", "")
    if err != nil {
        t.Fatal("Failed to set empty cat_api_key:", err)
    }
    
    breed, err = FetchBreedDetails("abys")
    assert.Error(t, err)
    assert.Nil(t, breed)
}