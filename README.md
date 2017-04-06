#this is a simple framwork with golang.

you can build a simple restful project with it.

##how to use:

package main

import (

    "net/http"

    "github.com/yang-f/beauty/utils/log"

    "github.com/yang-f/beauty/router"

    "github.com/yang-f/beauty/settings"

    "github.com/yang-f/beauty/controllers"

    "github.com/yang-f/beauty/decorates"

)

func main() {

    log.Printf("start server on port %s", settings.Listen)

    router.BRoutes = router.Routes{

	router.Route{

	    "getConfig",

	    "GET",

	    "/",

	    false,

	    controllers.Config,

	    "application/json;charset=utf-8",

	},

    }

    settings.Listen = ":8080"

    router := router.NewRouter()

    log.Fatal(http.ListenAndServe(settings.Listen, router))

}

##support:

token 

db

cors

log