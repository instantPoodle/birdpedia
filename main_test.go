//main_test.go

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("Get", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(handler)

	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v",
			status, http.StatusOK)
	}

	//Check that the response is what we expect.
	expected := "Hello world!"
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v, want %v", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	r := newRouter()

	mockServer := httptest.NewServer(r)

	resp, err := http.Get(mockServer.URL + "/hello")

	// Handle any unexpected error.
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	// read the body into a bunch of bytes(b)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// convert the bytes to a string
	respString := string(b)
	expected := "Hello world!"

	// We want our response to match the one we defined in our handler.
	// If it does happen to be "Hello world!", then it confirms, that the
	// route is correct.
	if respString != expected {
		t.Errorf("Response should bve %s, got %s", expected, respString)
	}
}

func TestRouterForNonExistentRoute(t *testing.T){
    r := newRouter()
    mockServer := httptest.NewServer(r)
    // most of the code is similar . The only difference is that we now make a 
    // request to a route that we didn't refine, like the 'POST/hello' route.

    resp, err :=  http.Post(mockServer.URL+"/hello", "", nil)

    if err != nil {
        t.Fatal(err)
    }

    // We want our status to be 405 (method not allowed)
    if resp.StatusCode !=  http.StatusMethodNotAllowed {
        t.Errorf("Status should be 405, got %d", resp.StatusCode)
    }

    // The code to test the body is also mostly the same, except this time, we expect 
    // an empty body.
    defer resp.Body.Close()
    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        t.Fatal(err)
    }
    respString := string(b)
    expected := ""

    if respString != expected {
        t.Errorf("Response should be %s, got %s", expected, respString)
    }
}

func TestStaticFileServer(t *testing.T){
    r := newRouter()
    mockServer := httptest.NewServer(r)

    // We want to his the 'GET /assets/' route to get the index.html file response
    resp, err := http.Get(mockServer.URL + "/assets/")
    if err != nil {
        t.Fatal(err)
    }

    // We want our status to be 200 (ok)
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Status should be 200, got %d", resp.StatusCode)
    }

    // It isn't wise to test the entire content of the HTML file.
    // Instead, we test that the content-type header is "text/html"; charset=utf-8"
    // so that we know that an html file has been served.
    contentType := resp.Header.Get("Content-Type")
    expectedContentType := "text/html; charset=utf-8"

    if expectedContentType != contentType {
        t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
    }
}
