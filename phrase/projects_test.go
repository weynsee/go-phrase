package phrase

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProjectsService_Current(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/projects/current", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": 1,"name":"project name","slug":"project-name"}`)
	})

	project, err := client.Projects.Current()
	if err != nil {
		t.Errorf("Projects.Current returned error: %v", err)
	}

	want := Project{ID: 1, Name: "project name", Slug: "project-name"}
	if !reflect.DeepEqual(*project, want) {
		t.Errorf("Projects.Current returned %+v, want %+v", project, want)
	}
}

func TestProjectsService_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Projects.Current()
		return err
	})
}
