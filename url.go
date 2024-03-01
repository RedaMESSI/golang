package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var ctx context.Context

func main() {
	// Initialiser le client MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	ctx = context.Background()

	r := mux.NewRouter()
	r.HandleFunc("/shorten", ShortenURL).Methods("POST")
	r.HandleFunc("/{shortURL}", Redirect).Methods("GET")
	r.HandleFunc("/stats/{shortURL}", GetStats).Methods("GET")

	http.Handle("/", r)

	fmt.Println("Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type URL struct {
	LongURL    string    `bson:"long_url"`
	ShortURL   string    `bson:"short_url"`
	CreatedAt  time.Time `bson:"created_at"`
	Expiration time.Time `bson:"expiration"`
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	longURL := r.FormValue("long_url")
	expiration := time.Now().AddDate(0, 0, 7) // Expiration dans 7 jours

	// Générer une URL courte, par exemple, en utilisant une fonction de hachage
	shortURL := "https://google.com/" + "short_hash"

	url := URL{
		LongURL:    longURL,
		ShortURL:   shortURL,
		CreatedAt:  time.Now(),
		Expiration: expiration,
	}

	// Insérer l'URL dans la base de données MongoDB
	urlCollection := client.Database("url_shortener").Collection("urls")
	_, err := urlCollection.InsertOne(ctx, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "URL courte générée: %s", shortURL)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	var url URL
	urlCollection := client.Database("url_shortener").Collection("urls")
	err := urlCollection.FindOne(ctx, bson.M{"short_url": shortURL}).Decode(&url)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Vérifier si le lien a expiré
	if url.Expiration.Before(time.Now()) {
		http.Error(w, "Ce lien a expiré", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusFound)
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	// Obtenir les statistiques basiques à partir de la base de données
	// par exemple, nombre de clics sur ce lien

	fmt.Fprintf(w, "Statistiques pour %s: nombre de clics = X", shortURL)
}
