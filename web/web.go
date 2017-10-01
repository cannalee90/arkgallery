package web

import (
	"reflect"
	"strings"
	"path/filepath"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/facebookgo/grace/gracehttp"
	"fmt"
	"arkgallery/web/util"
)


type HandleList struct {

}

func (h *HandleList) init() {

	elem := reflect.ValueOf(*h)

	for i := 0; i < elem.NumField(); i++ {
		f := elem.Field(i)
		ty := f.Type()
		cn := strings.ToLower(ty.Name())

		v := reflect.New(ty)

		for i := 0; i < v.NumMethod(); i++ {
			t := v.Type().Method(i)
			m := v.MethodByName(t.Name)
			name := t.Name

			if strings.ToUpper(name[:1]) != name[:1] {
				continue
			}

			mn := strings.ToLower(name)

			functions["/" + cn + "/" + mn] = m.Interface()
		}
	}
}

func Init() {
	initFunctions()
	initViews()
}

var functions map[string]interface{}

func initFunctions() {
	functions = make(map[string]interface{})

	list := HandleList{}
	list.init() // set functions
}

var views map[string]string

func initViews() {
	views = make(map[string]string)

	dsDir := filepath.Join(getCurDir(), "web", "view") + "/"

	getFiles(dsDir, "/")
}

func getFiles(path string, prefix string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		filename := f.Name()
		if f.IsDir() {
			getFiles(path + filename + "/", prefix + filename + "/")
		} else {
			targetFile := path + filename
			if idx := strings.LastIndex(f.Name(), "."); idx > -1 {
				filename = filename[:idx]
			}

			if filename == "index" {
				filename = ""
			}

			views[prefix + filename] = targetFile
		}
	}
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	//	logger.Debug("HandleFunc path : ", path)
	if f, ok := functions[path]; ok {
		session := NewSession(w, r)
		_, err := Call(f, session)
		if err != nil {
			session.ResponseEnd(nil, err)
		}
	} else {
		w.WriteHeader(404)
	}
}

func handleView(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if v, ok := views[path]; ok {
		text, err := ReadFile(v)
		if err == nil {
			w.Write(text)
		} else {
			w.WriteHeader(404)
		}
	} else {
		w.WriteHeader(404)
	}
}

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data := make([]byte, 1024 * 100)
	count, err := file.Read(data)
	if err != nil {
		return nil, err
	}

	return data[:count], nil
}

func Call(function interface{}, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(function)
	if len(params) != f.Type().NumIn() {
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func getCurDir() string {
	return filepath.Dir(os.Args[0])
}


func RunServer() {
	Init()

	mux := http.NewServeMux()
	staticDir := filepath.Join(getCurDir(), "web", "view") + "/"

	mux.Handle("/web/view/", http.StripPrefix("/web/view/", http.FileServer(http.Dir(staticDir))))
	for k, _ := range functions {
		mux.HandleFunc(k, handleFunc)
	}
	for k, _ := range views {
		if _, ok := functions[k]; !ok {
			mux.HandleFunc(k, handleView)
		}
	}

	gracehttp.Serve(&http.Server{
		Addr:    fmt.Sprintf(":%v", 3000),
		Handler: mux,
	})
}

func NewSession(w http.ResponseWriter, r *http.Request) util.Session {
	session := util.Session{}
	session.Init(w, r)
	return session
}