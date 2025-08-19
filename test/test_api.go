package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"todo_mcp/internal/models"
)

func main() {
	baseURL := "http://localhost:8080/api/v1"

	// Wait a moment for server to start
	time.Sleep(1 * time.Second)

	fmt.Println("Testing Todo API...")

	// Test 1: Create a todo
	fmt.Println("\n1. Creating a todo...")
	createReq := models.CreateTodoRequest{
		Title:       "Learn Go",
		Description: "Complete the Go tutorial",
	}

	jsonData, _ := json.Marshal(createReq)
	resp, err := http.Post(baseURL+"/todos", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", string(body))

	// Test 2: Get all todos
	fmt.Println("\n2. Getting all todos...")
	resp, err = http.Get(baseURL + "/todos")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", string(body))

	// Test 3: Get todo by ID
	fmt.Println("\n3. Getting todo by ID (1)...")
	resp, err = http.Get(baseURL + "/todos/1")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", string(body))
}
