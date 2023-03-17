package main

import (
	bytes2 "bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type FileName struct {
	Name string `json:"name"`
}

type Welcome struct {
	Sale string
	Time string
}

type Data struct {
	fileName string
	URL      string
}

type User struct {
	nome     string
	email    string
	password string
}

func main() {
	ipServerPub := "https://10.72.251.147:8443"
	ipServerSecure := "https://192.168.1.119:8443"

	fs := http.FileServer(http.Dir("./securityCopy"))
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	welcome := Welcome{"ola", time.Now().Format(time.Stamp)}
	template := template.Must(template.ParseFiles("template/template.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if sale := r.FormValue("sale"); sale != "" {
			welcome.Sale = sale
		}
		if err := template.ExecuteTemplate(w, "login.html", welcome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.Handle("/files/", http.StripPrefix("/files", fs))

	http.HandleFunc("/SD", func(w http.ResponseWriter, r *http.Request) {
		if err := template.ExecuteTemplate(w, "template.html", welcome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	db, err := sql.Open("postgres", "postgres://postgres:1234@localhost/?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE modc")
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", "postgres://postgres:1234@localhost/modc?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
CREATE TABLE Utilizadores (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255)
)
`)
	if err != nil {
		log.Fatal(err)
	}

	// Register a new user
	newUser := User{
		nome:     "John Doe",
		email:    "johndoe@example.com",
		password: "password123",
	}
	registerUser(db, newUser)

	// Log in an existing user
	existingUser := User{
		email:    "johndoe@example.com",
		password: "password123",
	}
	loggedIn, err := loginUser(db, existingUser)
	if err != nil {
		panic(err)
	}
	if loggedIn {
		fmt.Println("User logged in successfully")
	} else {
		fmt.Println("Incorrect email or password")
	}

	http.HandleFunc("/save", func(w http.ResponseWriter, response *http.Request) {

		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		response.Body.Close()
		log.Println(string(bytes))
		var fileResponse FileName
		errUnmarshal := json.Unmarshal(bytes, &fileResponse)

		if errUnmarshal != nil {
			log.Fatal(errUnmarshal)
		}
		log.Printf("%+v", fileResponse)

		data := &Data{
			fileName: fileResponse.Name,
			URL:      ipServerPub + "/files/" + fileResponse.Name,
		}

		log.Printf("%+v", data.download("securityCopy/"))

		w.WriteHeader(http.StatusOK)
		requestBody, err := json.Marshal(map[string]string{"name": data.fileName})
		serverResp, err := http.Post(ipServerSecure+"/save", "application/json", bytes2.NewBuffer(requestBody))
		if err != nil {
			log.Print(err)
		}
		log.Print(serverResp)
	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, response *http.Request) {
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		response.Body.Close()
		log.Println(string(bytes))
		var fileResponse FileName
		errUnmarshal := json.Unmarshal(bytes, &fileResponse)

		if errUnmarshal != nil {
			log.Fatal(errUnmarshal)
		}
		log.Printf("%+v", fileResponse)

		data := &Data{
			fileName: fileResponse.Name,
			URL:      ipServerPub + "/files/" + fileResponse.Name,
		}

		data.download("temp/")

	})

	http.HandleFunc("/retrieve", func(w http.ResponseWriter, response *http.Request) {

		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		response.Body.Close()
		log.Println(string(bytes))
		var fileResponse FileName
		errUnmarshal := json.Unmarshal(bytes, &fileResponse)

		if errUnmarshal != nil {
			log.Fatal(errUnmarshal)
		}
		log.Printf("%+v", fileResponse)

		data := &Data{
			fileName: fileResponse.Name,
			URL:      ipServerSecure + "/files/" + fileResponse.Name,
		}

		log.Printf("%+v", data.download("securityCopy/"))

	})

	fmt.Println(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil))

}

func (data *Data) download(Dir string) error {
	fmt.Println("abc")
	response, err := http.Get(data.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response status")
	}

	file, err := os.Create(Dir + data.fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)

	if err != nil {
		return err
	}

	return nil
}

func registerUser(db *sql.DB, user User) {
	_, err := db.Exec("INSERT INTO utilizadores (nome, email, password) VALUES ($1, $2, $3)", user.nome, user.email, user.password)
	if err != nil {
		panic(err)
	}
}

func loginUser(db *sql.DB, user User) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM utilizadores WHERE email = $1 AND password = $2", user.email, user.password).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}
