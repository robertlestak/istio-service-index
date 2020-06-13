package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/rs/cors"
)

type Service struct {
	Name        string   `json:"name"`
	Hosts       []string `json:"hosts"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
}

type Annotations struct {
	Description string `json:"istio-service-index.v1.lestak.sh/description"`
	Category    string `json:"istio-service-index.v1.lestak.sh/category"`
}

type Resource struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations       Annotations `json:"annotations"`
		CreationTimestamp time.Time   `json:"creationTimestamp"`
		Generation        int         `json:"generation"`
		Name              string      `json:"name"`
		Namespace         string      `json:"namespace"`
		ResourceVersion   string      `json:"resourceVersion"`
		SelfLink          string      `json:"selfLink"`
		UID               string      `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		Gateways []string `json:"gateways"`
		Hosts    []string `json:"hosts"`
		HTTP     []struct {
			Match []struct {
				URI struct {
					Prefix string `json:"prefix"`
				} `json:"uri"`
			} `json:"match"`
			Route []struct {
				Destination struct {
					Host string `json:"host"`
				} `json:"destination"`
			} `json:"route"`
		} `json:"http"`
	} `json:"spec"`
}

type ResourceList struct {
	APIVersion string     `json:"apiVersion"`
	Items      []Resource `json:"items"`
	Kind       string     `json:"kind"`
	Metadata   struct {
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
	} `json:"metadata"`
}

func GetServices() ([]Resource, error) {
	rl := new(ResourceList)
	rls := new([]Resource)
	// TODO: change to use kube golang client
	// I know this is hacky and wrong. I'll fix it soon.
	cmd := exec.Command("kubectl", "get", "virtualservice", "-A", "-o", "json")
	out := new(bytes.Buffer)
	e := new(bytes.Buffer)
	cmd.Stdout = out
	cmd.Stderr = e
	serr := cmd.Start()
	if serr != nil {
		return *rls, serr
	}
	werr := cmd.Wait()
	if werr != nil {
		return *rls, werr
	}
	berr := json.Unmarshal(out.Bytes(), &rl)
	if berr != nil {
		return *rls, berr
	}
	return rl.Items, nil
}

func ServiceList() ([]*Service, error) {
	var sl []*Service
	svcs, e := GetServices()
	if e != nil {
		return sl, e
	}
	for _, v := range svcs {
		s := &Service{
			Name:        v.Metadata.Name,
			Hosts:       v.Spec.Hosts,
			Category:    v.Metadata.Annotations.Category,
			Description: v.Metadata.Annotations.Description,
		}
		sl = append(sl, s)
	}
	return sl, nil
}

func handleServiceList(w http.ResponseWriter, r *http.Request) {
	sl, err := ServiceList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jd, jerr := json.Marshal(sl)
	if jerr != nil {
		http.Error(w, jerr.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, string(jd))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/api", handleServiceList)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8070"
	}
	log.Println("Listening on port :", port)
	h := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":"+port, h))
}
