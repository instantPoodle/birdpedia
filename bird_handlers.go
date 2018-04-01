package main

import (
    "net/http"
    "fmt"
    "encoding/json"
)

// Bird is the struct that defines bird attributes.
type Bird struct {
	Species     string `json:"species"`
	Description string `json:"description"`
}

var birds []Bird

func getBirdHandler(w http.ResponseWriter, r http.Request){
    //Convert the "birds" variable to json."
    birdListBytes, err := json.Marshal(birds)

    // If there is an error, print it to the console, and return a server
    // error response to the user
    if err != nil {
        fmt.Println(fmt.Errorf("Error: %v", err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    // If all goes well, write the JSON list of birds to the response.
    w.Write(birdListBytes)
}

func createBirdHandler(w http.ResponseWriter, r *http.Request) {
    // Create a new instance of Bird
    bird := Bird{}

    // We send all our data as HTML form data
    // the `ParseForm` method of the request, parsaes the
    // form values.
    err := r.ParseForm()

    // In case of any error, we reespond with an error to the user
    if err != nil {
        fmt.Println(fmt.Errorf("Error: %v", err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Get the information about the vird from the form info
    bird.Species = r.Form.Get("species")
    bird.Description = r.Form.Get("description")

    // Append our existing list of birds with a new entry
    birds = append(birds, bird)

    // Finally, we redirect the user to the original HTML page
    // (located at `/assets/`), using the http libraries `Redirect` method
    http.Redirect(w, r, "/assets/", http.StatusFound)
}






