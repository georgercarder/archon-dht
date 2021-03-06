package archon_dht

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

func (a *ArchonDHTs) RunTestTriggerServer() {
	go func() {
		a.runTestTriggerServer()
	}()
}

// runTestTriggerServer is useful for dev to test file versioning system
// as network grows
func (a *ArchonDHTs) runTestTriggerServer() {
	http.HandleFunc("/Stored", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["store"]
		if !ok || len(keys[0]) < 1 {
			fmt.Println("debug store key is missing")
			return
		}
		key := keys[0]

		// testing only on ETH for now
		dhtLayer := a.Layers["ETH"]
		versionData, err := dhtLayer.Config.PermissionLayer.NewVersionData()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		a.Stored(key, versionData)

		fmt.Fprintf(w, "Stored, %q\n%q\n", html.EscapeString(r.URL.Path), key) // response
	})

	http.HandleFunc("/GetValue", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["get"]
		if !ok || len(keys[0]) < 1 {
			fmt.Println("debug store key is missing")
			return
		}
		key := keys[0]

		dhtLayer := a.Layers["ETH"]
		keyAsCid, err := StringToCid(key)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		downloadKey := "/archondl/" + keyAsCid.String()
		value, err := dhtLayer.dHT.GetValue(context.Background(), downloadKey)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		uv := new(UrlsVersionedStruct)
		err = json.Unmarshal(value, &uv)

		fmt.Fprintf(w, "GetValue, %q\n%q\n%q\n", html.EscapeString(r.URL.Path), key, uv) // response
	})

	http.ListenAndServe(":9999", nil)
}
