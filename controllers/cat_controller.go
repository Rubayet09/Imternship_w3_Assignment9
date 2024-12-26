package controllers

import (
    "CatVotingApp/models"
    "encoding/json"
    "github.com/beego/beego/v2/server/web"
)

type CatController struct {
    web.Controller
}

// Request and response channel types
type RequestChannel struct {
    ResponseChan chan interface{}
    ErrorChan   chan error
}

// Global channels for different operations
var (
    catsChan     = make(chan *RequestChannel)
    breedsChan   = make(chan *RequestChannel)
    breedDetailsChan = make(chan struct {
        ID string
        ReqChan *RequestChannel
    })
    voteChan     = make(chan struct {
        Vote models.VoteRequest
        ReqChan *RequestChannel
    })
    favoritesChan = make(chan *RequestChannel)
    removeFavChan = make(chan struct {
        ID string
        ReqChan *RequestChannel
    })
)

func init() {
    // Start the worker goroutines
    go catsWorker()
    go breedsWorker()
    go breedDetailsWorker()
    go voteWorker()
    go favoritesWorker()
    go removeFavoriteWorker()
}

func (c *CatController) Get() {
    c.TplName = "index.tpl"
}

func (c *CatController) GetCats() {
    reqChan := &RequestChannel{
        ResponseChan: make(chan interface{}),
        ErrorChan:    make(chan error),
    }
    catsChan <- reqChan
    c.handleWorkerResponse(reqChan.ResponseChan, reqChan.ErrorChan)
}


func (c *CatController) GetBreeds() {
    reqChan := &RequestChannel{
        ResponseChan: make(chan interface{}),
        ErrorChan:   make(chan error),
    }
    
    // Send request to worker
    breedsChan <- reqChan
    c.handleWorkerResponse(reqChan.ResponseChan, reqChan.ErrorChan)

    

}

func (c *CatController) GetBreedDetails() {
    breedID := c.GetString("id")
    if breedID == "" {
        c.Data["json"] = map[string]interface{}{
            "status":  "error",
            "message": "breed id is required",
        }
        c.ServeJSON()
        return
    }

    reqChan := &RequestChannel{
        ResponseChan: make(chan interface{}),
        ErrorChan:   make(chan error),
    }
    
    breedDetailsChan <- struct {
        ID string
        ReqChan *RequestChannel
    }{breedID, reqChan}
    
    select {
    case breed := <-reqChan.ResponseChan:
        c.Data["json"] = map[string]interface{}{
            "status": "success",
            "data":   breed,
        }
    case err := <-reqChan.ErrorChan:
        c.Data["json"] = map[string]interface{}{
            "status":  "error",
            "message": err.Error(),
        }
    }
    c.ServeJSON()
}

func (c *CatController) VoteCat() {
    var voteReq models.VoteRequest
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &voteReq); err != nil {
        c.Data["json"] = map[string]string{
            "status":  "error",
            "message": "Invalid request format",
        }
        c.ServeJSON()
        return
    }

    reqChan := &RequestChannel{
        ResponseChan: make(chan interface{}),
        ErrorChan:   make(chan error),
    }
    
    voteChan <- struct {
        Vote models.VoteRequest
        ReqChan *RequestChannel
    }{voteReq, reqChan}
    
    select {
    case <-reqChan.ResponseChan:
        c.Data["json"] = map[string]string{
            "status":  "success",
            "message": "Vote recorded",
        }
    case err := <-reqChan.ErrorChan:
        c.Data["json"] = map[string]string{
            "status":  "error",
            "message": err.Error(),
        }
    }
    c.ServeJSON()
}

func (c *CatController) GetFavorites() {
    reqChan := &RequestChannel{
        ResponseChan: make(chan interface{}),
        ErrorChan:   make(chan error),
    }
    
    favoritesChan <- reqChan
    c.handleWorkerResponse(reqChan.ResponseChan, reqChan.ErrorChan)

}

func (c *CatController) RemoveFavorite() {
    id := c.Ctx.Input.Param(":id")
    if id == "" {
        c.Data["json"] = map[string]string{
            "status":  "error",
            "message": "Image ID is required",
        }
        c.ServeJSON()
        return
    }

    reqChan := &RequestChannel{
        ResponseChan: make(chan interface{}),
        ErrorChan:   make(chan error),
    }
    
    removeFavChan <- struct {
        ID string
        ReqChan *RequestChannel
    }{id, reqChan}
    
    select {
    case <-reqChan.ResponseChan:
        c.Data["json"] = map[string]string{
            "status":  "success",
            "message": "Favorite removed successfully",
        }
    case err := <-reqChan.ErrorChan:
        c.Data["json"] = map[string]string{
            "status":  "error",
            "message": err.Error(),
        }
    }
    c.ServeJSON()
}

func (c *CatController) handleWorkerResponse(responseChan <-chan interface{}, errorChan <-chan error) {
    select {
    case data := <-responseChan:
        c.Data["json"] = map[string]interface{}{
            "status": "success",
            "data":   data,
        }
    case err := <-errorChan:
        c.Data["json"] = map[string]interface{}{
            "status":  "error",
            "message": err.Error(),
        }
    }
    c.ServeJSON()
}


// Worker functions
func catsWorker() {
    for reqChan := range catsChan {
        cats, err := models.FetchCats()
        if err != nil {
            reqChan.ErrorChan <- err
        } else {
            reqChan.ResponseChan <- cats
        }
    }
}

func breedsWorker() {
    for reqChan := range breedsChan {
        breeds, err := models.FetchBreeds()
        if err != nil {
            reqChan.ErrorChan <- err
        } else {
            reqChan.ResponseChan <- breeds
        }
    }
}

func breedDetailsWorker() {
    for req := range breedDetailsChan {
        breed, err := models.FetchBreedDetails(req.ID)
        if err != nil {
            req.ReqChan.ErrorChan <- err
        } else {
            req.ReqChan.ResponseChan <- breed
        }
    }
}

func voteWorker() {
    for req := range voteChan {
        var err error
        if req.Vote.Vote == "love" {
            err = models.SaveFavorite(req.Vote.ImageID, req.Vote.ImageURL)
        }
        
        if err != nil {
            req.ReqChan.ErrorChan <- err
        } else {
            req.ReqChan.ResponseChan <- true
        }
    }
}

func favoritesWorker() {
    for reqChan := range favoritesChan {
        favorites := models.GetFavorites()
        reqChan.ResponseChan <- favorites
    }
}

func removeFavoriteWorker() {
    for req := range removeFavChan {
        if err := models.RemoveFavorite(req.ID); err != nil {
            req.ReqChan.ErrorChan <- err
        } else {
            req.ReqChan.ResponseChan <- true
        }
    }
}