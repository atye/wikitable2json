package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var (
	port = "36661"
)

func TestEntrypoint(t *testing.T) {
	c := Config{
		Client: &mockTableGetter{getMatrix: [][][]string{{{"test"}}}},
		Cache:  NewCache(1, 500*time.Millisecond, 500*time.Millisecond),
		Port:   port,
	}

	go Run(c)
	waitforServer()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/page", port))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected *, got %s", resp.Header.Get("Access-Control-Allow-Origin"))
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected *, got %s", resp.Header.Get("Content-Type"))
	}

	want := [][][]string{{{"test"}}}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var got [][][]string
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func waitforServer() {
	for {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s", port))
		if err != nil || resp.StatusCode != http.StatusOK {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		return
	}
}
