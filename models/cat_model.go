package models

import (
    "encoding/json"
    "fmt"
    "github.com/beego/beego/v2/server/web"
    "io/ioutil"
    "net/http"
    "sync"
)

type Cat struct {
    ID     string `json:"id"`
    URL    string `json:"url"`
    Breeds []Breed `json:"breeds"`
}

type Breed struct {
    ID            string  `json:"id"`
    Name          string  `json:"name"`
    Description   string  `json:"description"`
    Origin        string  `json:"origin"`
    WikipediaURL  string  `json:"wikipedia_url"`
    ReferenceImageID string `json:"reference_image_id"`
    Images          []string `json:"images,omitempty"`
}

type FavoriteImage struct {
    ID  string `json:"id"`
    URL string `json:"url"`
}

type VoteRequest struct {
    ImageID  string `json:"image_id"`
    ImageURL string `json:"image_url"`
    Vote     string `json:"vote"` // "like", "dislike", or "love"
}

var (
    favorites []FavoriteImage
    favMutex  sync.Mutex
)

func FetchCats() ([]Cat, error) {
    catAPIURL, err := web.AppConfig.String("cat_api_url")
    if err != nil {
        return nil, err
    }

    url := fmt.Sprintf("%s/images/search?limit=10", catAPIURL)
    response, err := fetchFromAPI(url)
    if err != nil {
        return nil, err
    }
    
    var cats []Cat
    if err := json.Unmarshal(response, &cats); err != nil {
        return nil, err
    }
    return cats, nil
}

// Update the FetchBreeds function in your cat_model.go
func FetchBreeds() ([]Breed, error) {
    catAPIURL, err := web.AppConfig.String("cat_api_url")
    if err != nil {
        return nil, fmt.Errorf("error getting cat_api_url: %v", err)
    }
    apiKey, err := web.AppConfig.String("cat_api_key")
    if err != nil {
        return nil, fmt.Errorf("error getting cat_api_key: %v", err)
    }

    url := fmt.Sprintf("%s/breeds", catAPIURL)
    fmt.Printf("Fetching breeds from URL: %s\n", url) // Debug log

    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Set("x-api-key", apiKey)
    fmt.Printf("Using API key: %s\n", apiKey) // Debug log

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    fmt.Printf("API Response Status: %s\n", resp.Status) // Debug log

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    fmt.Printf("Response body: %s\n", string(body)) // Debug log

    var breeds []Breed
    if err := json.Unmarshal(body, &breeds); err != nil {
        return nil, fmt.Errorf("error unmarshaling breeds: %v", err)
    }

    fmt.Printf("Successfully parsed %d breeds\n", len(breeds)) // Debug log
    return breeds, nil
}


func FetchBreedDetails(breedId string) (*Breed, error) {
    catAPIURL, _ := web.AppConfig.String("cat_api_url")
    apiKey, _ := web.AppConfig.String("cat_api_key")
    
    // First get breed details
    url := fmt.Sprintf("%s/breeds/%s", catAPIURL, breedId)
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }
    
    req.Header.Set("x-api-key", apiKey)
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error fetching breed: %v", err)
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %v", err)
    }
    
    var breed Breed
    if err := json.Unmarshal(body, &breed); err != nil {
        return nil, fmt.Errorf("error unmarshaling breed: %v", err)
    }
    
    // Now fetch images for this breed
    imagesURL := fmt.Sprintf("%s/images/search?breed_ids=%s&limit=5", catAPIURL, breedId)
    req, err = http.NewRequest("GET", imagesURL, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating images request: %v", err)
    }
    
    req.Header.Set("x-api-key", apiKey)
    resp, err = client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error fetching images: %v", err)
    }
    defer resp.Body.Close()
    
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading images response: %v", err)
    }
    
    var imagesData []struct {
        URL string `json:"url"`
    }
    if err := json.Unmarshal(body, &imagesData); err != nil {
        return nil, fmt.Errorf("error unmarshaling images: %v", err)
    }
    
    // Ensure we have the images array
    breed.Images = make([]string, len(imagesData))
    for i, img := range imagesData {
        breed.Images[i] = img.URL
    }

    return &breed, nil
}

func SaveFavorite(id string, url string) error {
    favMutex.Lock()
    defer favMutex.Unlock()
    
    for _, fav := range favorites {
        if fav.ID == id {
            return fmt.Errorf("already in favorites")
        }
    }
    favorites = append(favorites, FavoriteImage{ID: id, URL: url})
    return nil
}

func GetFavorites() []FavoriteImage {
    favMutex.Lock()
    defer favMutex.Unlock()
    return append([]FavoriteImage{}, favorites...)
}

func RemoveFavorite(id string) error {
    favMutex.Lock()
    defer favMutex.Unlock()
    
    // Log the incoming ID and current favorites for debugging
    fmt.Printf("Attempting to remove favorite with ID: %s\n", id)
    fmt.Printf("Current favorites: %+v\n", favorites)
    
    for i, fav := range favorites {
        // Log each comparison
        fmt.Printf("Comparing with favorite ID: %s\n", fav.ID)
        
        if fav.ID == id {
            // Remove the item
            favorites = append(favorites[:i], favorites[i+1:]...)
            fmt.Printf("Successfully removed favorite with ID: %s\n", id)
            return nil
        }
    }
    
    return fmt.Errorf("favorite not found with ID: %s", id)
}


func fetchFromAPI(url string) ([]byte, error) {
    apiKey, err := web.AppConfig.String("cat_api_key")
    if err != nil {
        return nil, err
    }

    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("x-api-key", apiKey)
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
    }
    
    return ioutil.ReadAll(resp.Body)
}