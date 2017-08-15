package main

import "github.com/matishsiao/goInfo"

import "net"
import "encoding/json"
import "net/http"
import "log"
import "os"
import "time"
import "errors"
import "io/ioutil"
import "fmt"

const (
	volumePath     = "/vmworld"
	volumeFlagFile = "vmworld_2017.flag"
)

type metadata struct {
	Timestamp     string               `json:"timestamp"`
	OS            *goInfo.GoInfoObject `json:"system"`
	IPs           []string             `json:"IPs"`
	VolumeName    string               `json:"volume_name,omitempty"`
	FilesInVolume []string             `json:"files_in_volume,omitempty"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		info, err := getSystemInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Only write once
		writeFlagFile(volumePath)

		w.WriteHeader(http.StatusOK)
		w.Write(info)
	})
	log.Println("Server is listening on 62017")
	log.Fatal(http.ListenAndServe(":62017", nil))
}

func getSystemInfo() ([]byte, error) {
	meta := &metadata{}
	meta.OS = goInfo.GetInfo()
	meta.IPs = getIPAddresses()
	meta.Timestamp = time.Now().UTC().String()

	files, err := listFilesInVolume(volumePath)
	if err != nil {
		meta.VolumeName = err.Error()
	} else {
		meta.VolumeName = volumePath
		meta.FilesInVolume = files
	}

	data, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getIPAddresses() []string {
	ips := []string{}
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				if ip != nil && !ip.IsLoopback() {
					ip4 := ip.To4()
					if ip4 != nil {
						ips = append(ips, ip4.String())
					}
				}
			}
		}
	}

	return ips
}

func fileExisting(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return !os.IsNotExist(err)
}

func listFilesInVolume(path string) ([]string, error) {
	files := []string{}
	if !fileExisting(path) {
		return files, errors.New("Path " + path + " does not exist")
	}

	osFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return files, err
	}

	for _, f := range osFiles {
		files = append(files, f.Name())
	}

	return files, nil
}

func writeFlagFile(path string) {
	//Volume not existing
	if !fileExisting(path) {
		return
	}

	fname := fmt.Sprintf("%s/%s", volumePath, volumeFlagFile)
	if fileExisting(fname) {
		//No need to write again
		return
	}

	log.Printf("Write flag file failed: %s\n", ioutil.WriteFile(fname, []byte("VMworld 2017"), 0777))
}
