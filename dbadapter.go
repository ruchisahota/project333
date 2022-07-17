package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func AddItem(w http.ResponseWriter, r *http.Request) {

	var newEntry Entry
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "wrong data")
	}
	json.Unmarshal(reqBody, &newEntry)

	//cassInit()
	if err := Session.Query("INSERT INTO Wardrobe(id, name, type, description) VALUES(?, ?, ?, ?)",
		newEntry.ID, newEntry.Name, newEntry.Type, newEntry.Description).Exec(); err != nil {
		fmt.Println("Error while inserting")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//defer closeCass()

	w.WriteHeader(http.StatusCreated)
}

func AddAllItems(w http.ResponseWriter, r *http.Request) {

	var Entries []Entry
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "wrong data")
	}
	json.Unmarshal(reqBody, &Entries)

	//cassInit()
	for _, newEntry := range Entries {
		if err := Session.Query("INSERT INTO Wardrobe(id, name, type, description) VALUES(?, ?, ?, ?)",
			newEntry.ID, newEntry.Name, newEntry.Type, newEntry.Description).Exec(); err != nil {
			fmt.Println("Error while inserting")
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//defer closeCass()

	w.WriteHeader(http.StatusCreated)
}

func GetAllItems(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Inside GetAllEntries - retrieving all Entries ")

	var Entries []Entry
	m := map[string]interface{}{}

	//cassInit()
	q := "SELECT * FROM Wardrobe"
	iter := Session.Query(q).Iter()
	defer iter.Close()

	for iter.MapScan(m) {
		tmp := Entry{
			ID:          m["id"].(int),
			Name:        m["name"].(string),
			Type:        m["type"].(string),
			Description: m["description"].(string),
		}
		Entries = append(Entries, tmp)
		m = map[string]interface{}{}
	}
	//defer closeCass()

	data, _ := json.Marshal(Entries)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetForPerson(w http.ResponseWriter, r *http.Request) {

	Name := r.URL.Query().Get("name")
	if Name == "" {
		http.Error(w, "The name query parameter is missing", http.StatusBadRequest)
		return
	}

	fmt.Println("Looking for Entries for Name %s ", Name)

	var Entries []Entry
	m := map[string]interface{}{}

	//cassInit()
	q := "SELECT * FROM Wardrobe WHERE name=? allow filtering"
	iter := Session.Query(q, Name).Iter()
	defer iter.Close()

	for iter.MapScan(m) {
		tmp := Entry{
			ID:          m["id"].(int),
			Name:        m["name"].(string),
			Type:        m["type"].(string),
			Description: m["description"].(string),
		}
		Entries = append(Entries, tmp)
		m = map[string]interface{}{}
	}
	//defer closeCass()

	data, _ := json.Marshal(Entries)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func PickOOTD(w http.ResponseWriter, r *http.Request) {

	Name := r.URL.Query().Get("name")
	if Name == "" {
		http.Error(w, "The name query parameter is missing", http.StatusBadRequest)
		return
	}

	fmt.Printf("Looking for OOTD for Name %s ", Name)

	var Entries []Entry
	m := map[string]interface{}{}
	var ootd OOTD
	rand.Seed(time.Now().UnixNano())

	//cassInit()
	//defer closeCass()

	// Pick Top
	iter1 := Session.Query("SELECT * FROM Wardrobe WHERE name=? and type =? allow filtering", Name, "top").Iter()
	defer iter1.Close()
	for iter1.MapScan(m) {
		tmp := Entry{
			ID:          m["id"].(int),
			Name:        m["name"].(string),
			Type:        m["type"].(string),
			Description: m["description"].(string),
		}
		Entries = append(Entries, tmp)
		m = map[string]interface{}{}
	}

	if len(Entries) == 0 {
		w.WriteHeader(http.StatusOK)
		data, _ := json.Marshal(ootd)
		w.Write(data)
		return
	}
	ootd.Top = Entries[rand.Intn(len(Entries))].Description

	// Pick Bottom
	Entries = []Entry{}
	iter2 := Session.Query("SELECT * FROM Wardrobe WHERE name=? and type =? allow filtering", Name, "bottom").Iter()
	defer iter2.Close()
	for iter2.MapScan(m) {
		tmp := Entry{
			ID:          m["id"].(int),
			Name:        m["name"].(string),
			Type:        m["type"].(string),
			Description: m["description"].(string),
		}
		Entries = append(Entries, tmp)
		m = map[string]interface{}{}
	}
	if len(Entries) == 0 {
		w.WriteHeader(http.StatusOK)
		data, _ := json.Marshal(ootd)
		w.Write(data)
		return
	}
	ootd.Bottom = Entries[rand.Intn(len(Entries))].Description

	// Pick Accessory
	Entries = []Entry{}
	iter3 := Session.Query("SELECT * FROM Wardrobe WHERE name=? and type =? allow filtering", Name, "accessory").Iter()
	defer iter3.Close()
	for iter3.MapScan(m) {
		tmp := Entry{
			ID:          m["id"].(int),
			Name:        m["name"].(string),
			Type:        m["type"].(string),
			Description: m["description"].(string),
		}
		Entries = append(Entries, tmp)
		m = map[string]interface{}{}
	}
	if len(Entries) == 0 {
		w.WriteHeader(http.StatusOK)
		data, _ := json.Marshal(ootd)
		w.Write(data)
		return
	}
	ootd.Accessory = Entries[rand.Intn(len(Entries))].Description

	data, _ := json.Marshal(ootd)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func CountAllItems(w http.ResponseWriter, r *http.Request) {
	var Count string

	//cassInit()
	err := Session.Query("SELECT count(*) FROM Wardrobe").Scan(&Count)
	if err != nil {
		panic(err)
	}
	//defer closeCass()

	data, _ := json.Marshal(Count)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// func DeleteOneEntry(w http.ResponseWriter, r *http.Request) {

// 	EntryID := r.URL.Query().Get("id")
// 	if EntryID == "" {
// 		http.Error(w, "The id query parameter is missing", http.StatusBadRequest)
// 		return
// 	}

// 	Name := r.URL.Query().Get("name")
// 	if Name == "" {
// 		http.Error(w, "The name query parameter is missing", http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Println("Deleting Entry with ID %d and Name %s ", EntryID, Name)

// 	//cassInit()
// 	if err := Session.Query("DELETE FROM Wardrobe WHERE id = ? and name = ? ", EntryID, Name).Exec(); err != nil {
// 		fmt.Println("Error while deleting")
// 		fmt.Println(err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	//defer closeCass()
// }

func DeleteAllItems(w http.ResponseWriter, r *http.Request) {

	//cassInit()
	if err := Session.Query("TRUNCATE Wardrobe").Exec(); err != nil {
		fmt.Println("Error while deleting all Entries")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//defer closeCass()

	w.WriteHeader(http.StatusOK)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {

	EntryID := r.URL.Query().Get("id")
	if EntryID == "" {
		http.Error(w, "The id query parameter is missing", http.StatusBadRequest)
		return
	}

	Name := r.URL.Query().Get("name")
	if Name == "" {
		http.Error(w, "The name query parameter is missing", http.StatusBadRequest)
		return
	}

	Type := r.URL.Query().Get("type")
	if Name == "" {
		http.Error(w, "The type query parameter is missing", http.StatusBadRequest)
		return
	}

	fmt.Println("Updating Entry with ID %d, Name %s, Type %s", EntryID, Name, Type)

	var UpdateEntry Entry
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data properly")
	}
	json.Unmarshal(reqBody, &UpdateEntry)

	//cassInit()
	if err := Session.Query("UPDATE Wardrobe SET type = ?, description = ? WHERE name = ? and id = ? and type = ?",
		UpdateEntry.Type, UpdateEntry.Description, Name, EntryID, Type).Exec(); err != nil {
		fmt.Println("Error while updating")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//defer closeCass()

	w.WriteHeader(http.StatusOK)
}

func GetEntireWardrobe(w http.ResponseWriter, r *http.Request) ([]Entry, error) {
	var Entries []Entry
	m := map[string]interface{}{}

	//cassInit()
	q := "SELECT * FROM Wardrobe"
	iter := Session.Query(q).Iter()
	defer iter.Close()

	for iter.MapScan(m) {
		tmp := Entry{
			ID:          m["id"].(int),
			Name:        m["name"].(string),
			Type:        m["type"].(string),
			Description: m["description"].(string),
		}
		Entries = append(Entries, tmp)
		m = map[string]interface{}{}
	}
	//defer closeCass()

	return Entries, nil
}
