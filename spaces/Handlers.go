package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Services
func GetFile(w http.ResponseWriter, r *http.Request) {
	// Check unauthorized. Replace this Authorization token by a valid one
	// by automatic generation and / or a new and dedicated web service
	/*
		if r.Header.Get("Token") != "Jkd855c6x9Aqcf" {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusForbidden)
			panic("Non authorized access detected")
		}
	*/

	//vars := mux.Vars(r)

	vars := mux.Vars(r)
	space := vars["space"]

	if strings.Contains(space, "..") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	concretePath := getPrivateFolders(space)

	qChannel := make(chan FileModel)
	var wg sync.WaitGroup
	wg.Add(1)
	go grDiscoverFiles(concretePath, "", qChannel, &wg)

	if err := json.NewEncoder(w).Encode(<-qChannel); err != nil {
		panic(err)
	}
}

// Services
func GetFiles(w http.ResponseWriter, r *http.Request) {

	// Check unauthorized. Replace this Authorization token by a valid one
	// by automatic generation and / or a new and dedicated web service
	/*
		if r.Header.Get("Token") != "Jkd855c6x9Aqcf" {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusForbidden)
			panic("Non authorized access detected")
		}
	*/

	//vars := mux.Vars(r)

	vars := mux.Vars(r)
	space := vars["space"]

	if strings.Contains(space, "..") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	concretePath := getPrivateFolders(space)

	qChannel := make(chan FileModel)
	var wg sync.WaitGroup
	wg.Add(1)
	go grDiscoverFiles(concretePath, "", qChannel, &wg)

	var files []FileModel
	for file := range qChannel {
		files = append(files, file)
	}

	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.Printf(err.Error())
		panic(err)
	}

}

// Make that function asynchronous
func CreateSpace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var post map[string]interface{}
	err := decoder.Decode(&post)

	if err != nil {
		panic(err)
	}

	constructHeaders(w)
	id := post["id"].(string)
	go createUserSpace(id)
}

func Download(w http.ResponseWriter, r *http.Request) {
	/*
		vars := mux.Vars(r)
		Module := vars["file"]

		log.Printf("Downloading module %s ", Module)
		h := md5.New()
		io.WriteString(h, Module)

		out, err := os.Create("E:\\Projects\\ProjectFIles\\private\\" + fmt.Sprintf("%x", h.Sum(nil)) + "\\" + Module + ".dll")
		if err != nil {
			panic(err)
		}

		defer out.Close()

		w.WriteHeader(http.StatusOK)

		FileHeader := make([]byte, 512)
		out.Read(FileHeader)
		FileContentType := http.DetectContentType(FileHeader)

		//Get the file size
		FileStat, _ := out.Stat()                          //Get info from file
		FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

		//Send the headers
		w.Header().Set("Content-Disposition", "attachment; filename="+Module+".dll")
		w.Header().Set("Content-Type", FileContentType)
		w.Header().Set("Content-Length", FileSize)

		//Send the file
		//We read 512 bytes from the file already so we reset the offset back to 0
		out.Seek(0, 0)
		io.Copy(w, out) //'Copy' the file to the client
		//io.Copy(out, resp.Body)
	*/
}

func createUserSpace(id string) {

	notificationId := "1"
	notificationMessage := "1. Checking for requirements : "
	sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_NEW, PRIORITY_STD, TYPE_INFO, notificationMessage+"Initializing"))
	time.Sleep(500 * time.Millisecond)
	sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_ONGOING, PRIORITY_STD, TYPE_INFO, notificationMessage+"In progress"))
	time.Sleep(10 * time.Second)
	sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_DONE, PRIORITY_STD, TYPE_INFO, notificationMessage+"Done"))

	notificationId = "2"

	notificationMessage = "2. Preparing the space : "
	sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_NEW, PRIORITY_STD, TYPE_INFO, notificationMessage+"Initializing"))
	time.Sleep(500 * time.Millisecond)
	sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_ONGOING, PRIORITY_STD, TYPE_INFO, notificationMessage+"In progress"))
	time.Sleep(5 * time.Second)
	sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_DONE, PRIORITY_STD, TYPE_INFO, notificationMessage+"Done"))

	userSpace := getPrivateFolders(id)
	sharedFolder := getSharedFolders()

	notificationId = "3"
	notificationMessage = "3. Copying necessary data : "
	sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_NEW, PRIORITY_STD, TYPE_INFO, notificationMessage+"Initializing"))
	time.Sleep(500 * time.Millisecond)

	os.MkdirAll(userSpace, 0755)
	err := CopyDir(sharedFolder, userSpace)

	sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_ONGOING, PRIORITY_STD, TYPE_INFO, notificationMessage+"In progress"))
	time.Sleep(5 * time.Second)

	if err != nil {
		failOnError(err, "Failed to copy a directory")
		sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_ERROR, PRIORITY_CRITICAL, TYPE_ERROR, notificationMessage+"Failed"))
	} else {
		sendMessage("user-notification", false, constructNotification(notificationId, id, "CreateSpace", STATUS_DONE, PRIORITY_STD, TYPE_INFO, notificationMessage+"Done"))
	}
	time.Sleep(20 * time.Second)
	sendMessage("user-notification", false, constructNotification("OK", id, "CreateSpace", -1, -1, -1, ""))

}
