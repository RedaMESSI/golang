package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenURL(t *testing.T) {
	// Crée une requête POST de test avec une URL longue
	reqBody := strings.NewReader("long_url=https://example.com")
	req := httptest.NewRequest("POST", "/shorten", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Appelle la fonction de raccourcissement d'URL
	ShortenURL(w, req)

	// Vérifie si la réponse est un code de statut HTTP 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Vérifie si la réponse contient l'URL courte générée
	expectedURL := "URL courte générée: https://google.com/short_hash"
	if !strings.Contains(w.Body.String(), expectedURL) {
		t.Errorf("Expected response body to contain %s", expectedURL)
	}
}

func TestRedirect(t *testing.T) {
	// Crée une requête GET de test avec une URL courte
	req := httptest.NewRequest("GET", "/short_hash", nil)
	w := httptest.NewRecorder()

	// Appelle la fonction de redirection d'URL
	Redirect(w, req)

	// Vérifie si la réponse est un code de statut HTTP 302 Found (redirection)
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got %d", http.StatusFound, w.Code)
	}
}

// Vous pouvez ajouter d'autres tests pour GetStats et d'autres fonctions si nécessaire
