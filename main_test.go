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
