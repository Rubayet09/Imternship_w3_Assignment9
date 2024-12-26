package routers

import (
    "CatVotingApp/controllers"
    "github.com/beego/beego/v2/server/web"
)

func init() {
    web.Router("/", &controllers.CatController{})
    
    ns := web.NewNamespace("/api",
        web.NSRouter("/cats", &controllers.CatController{}, "get:GetCats"),
        web.NSRouter("/breeds", &controllers.CatController{}, "get:GetBreeds"),
        web.NSRouter("/breed", &controllers.CatController{}, "get:GetBreedDetails"),
        web.NSRouter("/vote", &controllers.CatController{}, "post:VoteCat"),
        web.NSRouter("/favorites", &controllers.CatController{}, "get:GetFavorites"),
        web.NSRouter("/favorites/:id", &controllers.CatController{}, "delete:RemoveFavorite")    )
    
    web.AddNamespace(ns)
    web.SetStaticPath("/static", "static")
}