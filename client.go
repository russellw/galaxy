package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GameClient struct {
	serverURL string
	playerID  string
	client    *http.Client
}

func NewGameClient(serverURL, playerID string) *GameClient {
	return &GameClient{
		serverURL: serverURL,
		playerID:  playerID,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (gc *GameClient) Connect() error {
	connectData := map[string]string{
		"player_id": gc.playerID,
	}
	
	jsonData, _ := json.Marshal(connectData)
	resp, err := gc.client.Post(gc.serverURL+"/connect", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var response APIResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if response.Success {
		fmt.Printf("Successfully connected as %s\n", gc.playerID)
	} else {
		fmt.Printf("Failed to connect: %s\n", response.Message)
	}
	
	return nil
}

func (gc *GameClient) SubmitOrder(orderType, planetID string, parameters map[string]interface{}, priority int) error {
	orderData := OrderRequest{
		PlayerID:   gc.playerID,
		OrderType:  orderType,
		PlanetID:   planetID,
		Parameters: parameters,
		Priority:   priority,
	}
	
	jsonData, _ := json.Marshal(orderData)
	resp, err := gc.client.Post(gc.serverURL+"/orders", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var response APIResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if response.Success {
		fmt.Printf("Order submitted: %s\n", response.Message)
	} else {
		fmt.Printf("Order failed: %s\n", response.Message)
	}
	
	return nil
}

func (gc *GameClient) GetGameStatus() error {
	resp, err := gc.client.Get(gc.serverURL + "/status")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var response APIResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if response.Success {
		statusData, _ := json.MarshalIndent(response.Data, "", "  ")
		fmt.Printf("Game Status:\n%s\n", statusData)
	}
	
	return nil
}

func (gc *GameClient) GetPlayerStatus() error {
	resp, err := gc.client.Get(gc.serverURL + "/player/" + gc.playerID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var response APIResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	if response.Success {
		playerData, _ := json.MarshalIndent(response.Data, "", "  ")
		fmt.Printf("Player Status:\n%s\n", playerData)
	}
	
	return nil
}

func (gc *GameClient) GetGameState() error {
	resp, err := gc.client.Get(gc.serverURL + "/game")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Full Game State:\n%s\n", string(body))
	
	return nil
}

// Example usage function
func DemoClient() {
	client := NewGameClient("http://localhost:8080", "player1")
	
	fmt.Println("=== Galaxy Game Client Demo ===")
	
	// Connect to the game
	client.Connect()
	
	// Get initial status
	fmt.Println("\n1. Getting game status...")
	client.GetGameStatus()
	
	// Get player status
	fmt.Println("\n2. Getting player status...")
	client.GetPlayerStatus()
	
	// Submit some orders
	fmt.Println("\n3. Submitting orders...")
	
	// Build a metal mine
	client.SubmitOrder("BUILD_FACILITY", "planet_player1_home", map[string]interface{}{
		"facility_type": "MetalMine",
	}, 5)
	
	// Build a ship
	client.SubmitOrder("BUILD_SHIP", "planet_player1_home", map[string]interface{}{
		"ship_type": "Fighter",
	}, 4)
	
	// Upgrade a facility
	client.SubmitOrder("UPGRADE_FACILITY", "planet_player1_home", map[string]interface{}{
		"facility_type": "Factory",
	}, 3)
	
	fmt.Println("\n4. Orders submitted. Check server logs for processing.")
	
	// Wait a bit and check status again
	time.Sleep(2 * time.Second)
	fmt.Println("\n5. Updated player status...")
	client.GetPlayerStatus()
}