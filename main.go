package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

/* 50:36 min*/

const Version = "1.0.1"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, option *Options) (*Driver, error) {

	dir = filepath.Clean(dir)

	opts := Options{}

	if option != nil {
		opts = *option
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n")
		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s'...", dir)

	return &driver, os.MkdirAll(dir, 0755)

}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	     
	if collection == ""{
		fmt.Errorf("Missing resource - no place to save record (no name)")

	}
	
	if resource ==  "" {
			fmt.Errorf("Missing resource - unable to save record (no name)")
		
	}

    mutex := d.geOrCreateMutex()

	mutex.Lock()

	defer mutex.Unlock()

	dir := filepath.Join(d.dir,collection)

	fnlPath := filepath.Join(dir,resource+".json")

	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir,0755); err != nil {
		return err
	}

	return nil 
}

func (d *Driver) Read() {}

func (d *Driver) ReadAll() {}

func (d *Driver) geOrCreateMutex() *sync.Mutex {}

func stat(path string)(fi os.FileInfo,error){
	if fi, err = os.Stat(path); os.IsNotExist(err){
		fi,err = os.Stat(path+".json")
	}
	return
}

var lum = lumber.APPEND

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"

	db, err := New(dir, nil)

	if err != nil {
		fmt.Println("Error ", err)
	}

	employes := []User{{"Jhon", "23", "23344333 ", "Myrl Tech", Address{"bangalore", "karnatala", "india", "410013"}},
		{"Paul", "24", "23344333", "Google", Address{"bangalore", "karnatala", "india", "410013"}},
		{"Robert", "33", "23344333", "Microsoft", Address{"bangalore", "karnatala", "india", "410013"}},
		{"Vince", "45", "23344333", "Facebook", Address{"bangalore", "karnatala", "india", "410013"}},
		{"Neo", "26", "23344333", "Remote-Teams", Address{"bangalore", "karnatala", "india", "410013"}},
		{"Albert", "24", "23344333", "Dominate", Address{"bangalore", "karnatala", "india", "410013"}}}

	for _, v := range employes {

		db.Write(v.Name, User{
			Name:    v.Name,
			Age:     v.Age,
			Contact: v.Contact,
			Company: v.Company,
			Address: v.Address,
		})

	}

	record, err := db.ReadAll("Users")

	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println(record)

	allusers := []User{}

	for _, f := range record {

		employeedFound := User{}
		if err := json.Unmarshal([]byte(f), employeedFound); err != nil {
			fmt.Println("Error", err)
		}

		allusers = append(allusers, employeedFound)
	}
	fmt.Println(allusers)

	if err := db.Delete("user", "jhon"); err != nil {
		fmt.Println("Error", err)
	}

	if err := db.Delete("user", ""); err != nil {
		fmt.Println("Error", err)
	}
}

/*
0102
04245253695
17468622
#8250
*/
