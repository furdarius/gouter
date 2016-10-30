package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNode(t *testing.T) {
	n := newNode("", nil)

	n.insert("/list", func(w http.ResponseWriter, r *http.Request, p Params) {
		fmt.Fprintf(w, "Call /list")
	})

	n.insert("/list/:id/show", func(w http.ResponseWriter, r *http.Request, p Params) {
		fmt.Fprintf(w, "Call /list/:id/show")
	})

	n.insert("/list/:id/edit", func(w http.ResponseWriter, r *http.Request, p Params) {
		fmt.Fprintf(w, "Call /list/:id/edit")
	})

	n.insert("/", func(w http.ResponseWriter, r *http.Request, p Params) {
		fmt.Fprintf(w, "Call /")
	})

	n.insert("/user/:name", func(w http.ResponseWriter, r *http.Request, p Params) {
		fmt.Fprintf(w, "Call /user/:name")
	})

	n.insert("/user/:name/edit", func(w http.ResponseWriter, r *http.Request, p Params) {
		fmt.Fprintf(w, "Call /user/:name/edit")
	})

	n.insert("/cars/model_:type", func(w http.ResponseWriter, r *http.Request, p Params) {
		fmt.Fprintf(w, "Call /cars/model_:type")
	})

	tests := []struct {
		route string
		want  string
		exist bool
	}{
		{
			route: "/user/flaker",
			want:  "/user/:name",
			exist: true,
		},
		{
			route: "/",
			want:  "/",
			exist: true,
		},
		{
			route: "/future",
			want:  "",
			exist: false,
		},
		{
			route: "/user/flaker/edit",
			want:  "/user/:name/edit",
			exist: true,
		},
		{
			route: "/list/show",
			want:  "",
			exist: false,
		},
		{
			route: "/list/23/show",
			want:  "/list/:id/show",
			exist: true,
		},
		{
			route: "/list",
			want:  "/list",
			exist: true,
		},
		{
			route: "/cars/model_s",
			want:  "/cars/model_:type",
			exist: true,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)
		writer := httptest.NewRecorder()

		handler, params := n.lookup(test.route)

		if test.exist && handler == nil {
			t.Fatalf("Route \"%s\" must exist, but lookup return %v!", test.route, handler)
		}

		if !test.exist && handler != nil {
			t.Fatalf("Route \"%s\" must not exist, but lookup return %v!", test.route, handler)
		}

		if test.exist {
			handler(writer, req, params)

			wantBody := fmt.Sprintf("Call %s", test.want)
			gotBody, err := ioutil.ReadAll(writer.Body)

			if err != nil {
				t.Fatal(err)
			}

			if string(gotBody) != wantBody {
				t.Fatalf("Route \"%s\" got \"%s\" want \"%s\"", test.route, string(gotBody), wantBody)
			}
		}
	}

	//outputTrie(n, 0)
}

func outputTrie(n *node, padding int) {
	for key, value := range n.edges {
		paramStatus := "Param Exist"
		if value.param == nil {
			paramStatus = "Param Not Exist"
		}

		fmt.Println(strings.Repeat(" ", padding), fmt.Sprintf("[%s]", string(key)), string(value.path), value.value, paramStatus)

		if value.param != nil {
			paramStatus = "Param Exist"
			if value.param.param == nil {
				paramStatus = "Param Not Exist"
			}

			fmt.Println(strings.Repeat(" ", padding+8), ":"+string(value.param.path), value.param.value, paramStatus)

			outputTrie(value.param, padding+12)
		}

		outputTrie(value, padding+4)
	}
}
