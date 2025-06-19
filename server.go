package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type GameServer struct {
	gameState    *GameState
	mutex        sync.RWMutex
	turnDuration time.Duration
	turnTimer    *time.Timer
	clients      map[string]*PlayerConnection
}

type PlayerConnection struct {
	PlayerID   string
	LastSeen   time.Time
	Connected  bool
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type OrderRequest struct {
	PlayerID   string                 `json:"player_id"`
	OrderType  string                 `json:"order_type"`
	PlanetID   string                 `json:"planet_id,omitempty"`
	SystemID   string                 `json:"system_id,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Priority   int                    `json:"priority"`
}

func NewGameServer(players []Player, galaxySize int, maxTurns int, turnDurationSeconds int) *GameServer {
	gameState := NewGameState(players, galaxySize, maxTurns)
	
	server := &GameServer{
		gameState:    &gameState,
		turnDuration: time.Duration(turnDurationSeconds) * time.Second,
		clients:      make(map[string]*PlayerConnection),
	}
	
	// Initialize player connections
	for _, player := range players {
		server.clients[player.ID] = &PlayerConnection{
			PlayerID:  player.ID,
			LastSeen:  time.Now(),
			Connected: false,
		}
	}
	
	return server
}

func (gs *GameServer) StartServer(port int) {
	http.HandleFunc("/", gs.handleRoot)
	http.HandleFunc("/status", gs.handleStatus)
	http.HandleFunc("/game", gs.handleGameState)
	http.HandleFunc("/orders", gs.handleOrders)
	http.HandleFunc("/player/", gs.handlePlayerStatus)
	http.HandleFunc("/connect", gs.handleConnect)
	http.HandleFunc("/turn", gs.handleTurnControl)
	
	fmt.Printf("Galaxy Game Server starting on port %d\n", port)
	fmt.Printf("Turn duration: %v\n", gs.turnDuration)
	fmt.Printf("Players: %v\n", gs.getPlayerNames())
	
	// Start the turn timer
	gs.startTurnTimer()
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func (gs *GameServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	response := `
Galaxy Strategy Game Server

Endpoints:
- GET  /status           - Server and game status
- GET  /game             - Full game state
- GET  /player/{id}      - Player-specific information
- POST /connect          - Connect as a player
- POST /orders           - Submit orders
- POST /turn             - Manual turn control (admin)

Game Status: ` + gs.getGameStatus()
	
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, response)
}

func (gs *GameServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	
	connectedPlayers := 0
	for _, client := range gs.clients {
		if client.Connected {
			connectedPlayers++
		}
	}
	
	status := map[string]interface{}{
		"current_turn":       gs.gameState.CurrentTurn,
		"max_turns":         gs.gameState.MaxTurns,
		"game_over":         gs.gameState.GameOver,
		"winner":            gs.gameState.Winner,
		"connected_players": connectedPlayers,
		"total_players":     len(gs.gameState.Players),
		"turn_duration":     gs.turnDuration.Seconds(),
		"systems_count":     len(gs.gameState.Galaxy.StarSystems),
	}
	
	gs.sendJSON(w, APIResponse{Success: true, Data: status})
}

func (gs *GameServer) handleGameState(w http.ResponseWriter, r *http.Request) {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	
	// Return simplified game state (not full internal state for security)
	gameData := map[string]interface{}{
		"current_turn": gs.gameState.CurrentTurn,
		"max_turns":   gs.gameState.MaxTurns,
		"game_over":   gs.gameState.GameOver,
		"winner":      gs.gameState.Winner,
		"players":     gs.getPlayerSummaries(),
		"systems":     gs.getSystemSummaries(),
	}
	
	gs.sendJSON(w, APIResponse{Success: true, Data: gameData})
}

func (gs *GameServer) handleOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Method not allowed"})
		return
	}
	
	var orderReq OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Invalid JSON"})
		return
	}
	
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	// Validate player
	if _, exists := gs.clients[orderReq.PlayerID]; !exists {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Invalid player ID"})
		return
	}
	
	// Check if game is over
	if gs.gameState.GameOver {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Game is over"})
		return
	}
	
	// Create and add order
	order := Order{
		PlayerID:   orderReq.PlayerID,
		OrderType:  orderReq.OrderType,
		PlanetID:   orderReq.PlanetID,
		SystemID:   orderReq.SystemID,
		Parameters: orderReq.Parameters,
		Priority:   orderReq.Priority,
	}
	
	gs.gameState.AddOrder(order)
	
	gs.sendJSON(w, APIResponse{
		Success: true, 
		Message: fmt.Sprintf("Order added for player %s", orderReq.PlayerID),
	})
}

func (gs *GameServer) handlePlayerStatus(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Path[len("/player/"):]
	
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	
	if _, exists := gs.clients[playerID]; !exists {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Player not found"})
		return
	}
	
	playerData := map[string]interface{}{
		"player_id":     playerID,
		"summary":       gs.gameState.GetPlayerSummary(playerID),
		"systems":       gs.getPlayerSystems(playerID),
		"current_turn":  gs.gameState.CurrentTurn,
		"orders_count":  len(gs.gameState.Orders[playerID]),
	}
	
	gs.sendJSON(w, APIResponse{Success: true, Data: playerData})
}

func (gs *GameServer) handleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Method not allowed"})
		return
	}
	
	var connectReq struct {
		PlayerID string `json:"player_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&connectReq); err != nil {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Invalid JSON"})
		return
	}
	
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	if client, exists := gs.clients[connectReq.PlayerID]; exists {
		client.Connected = true
		client.LastSeen = time.Now()
		gs.sendJSON(w, APIResponse{
			Success: true,
			Message: fmt.Sprintf("Connected as %s", connectReq.PlayerID),
		})
	} else {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Invalid player ID"})
	}
}

func (gs *GameServer) handleTurnControl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Method not allowed"})
		return
	}
	
	var turnReq struct {
		Action string `json:"action"` // "process" or "reset_timer"
	}
	
	if err := json.NewDecoder(r.Body).Decode(&turnReq); err != nil {
		gs.sendJSON(w, APIResponse{Success: false, Message: "Invalid JSON"})
		return
	}
	
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	
	switch turnReq.Action {
	case "process":
		gs.processTurn()
		gs.sendJSON(w, APIResponse{Success: true, Message: "Turn processed manually"})
	case "reset_timer":
		gs.resetTurnTimer()
		gs.sendJSON(w, APIResponse{Success: true, Message: "Turn timer reset"})
	default:
		gs.sendJSON(w, APIResponse{Success: false, Message: "Invalid action"})
	}
}

func (gs *GameServer) startTurnTimer() {
	gs.turnTimer = time.AfterFunc(gs.turnDuration, func() {
		gs.mutex.Lock()
		defer gs.mutex.Unlock()
		
		if !gs.gameState.GameOver {
			gs.processTurn()
			gs.startTurnTimer() // Start next turn timer
		}
	})
}

func (gs *GameServer) resetTurnTimer() {
	if gs.turnTimer != nil {
		gs.turnTimer.Stop()
	}
	gs.startTurnTimer()
}

func (gs *GameServer) processTurn() {
	fmt.Printf("Processing turn %d automatically...\n", gs.gameState.CurrentTurn)
	gs.gameState.ProcessTurn()
	
	if gs.gameState.GameOver {
		fmt.Printf("Game over! Winner: %s\n", gs.gameState.Winner)
	}
}

func (gs *GameServer) sendJSON(w http.ResponseWriter, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (gs *GameServer) getGameStatus() string {
	if gs.gameState.GameOver {
		return fmt.Sprintf("Game Over - Winner: %s", gs.gameState.Winner)
	}
	return fmt.Sprintf("Turn %d/%d - Active", gs.gameState.CurrentTurn, gs.gameState.MaxTurns)
}

func (gs *GameServer) getPlayerNames() []string {
	names := make([]string, len(gs.gameState.Players))
	for i, player := range gs.gameState.Players {
		names[i] = player.Name
	}
	return names
}

func (gs *GameServer) getPlayerSummaries() map[string]string {
	summaries := make(map[string]string)
	for _, player := range gs.gameState.Players {
		summaries[player.ID] = gs.gameState.GetPlayerSummary(player.ID)
	}
	return summaries
}

func (gs *GameServer) getSystemSummaries() []map[string]interface{} {
	systems := make([]map[string]interface{}, len(gs.gameState.Galaxy.StarSystems))
	for i, system := range gs.gameState.Galaxy.StarSystems {
		systems[i] = map[string]interface{}{
			"id":           system.ID,
			"name":         system.Name,
			"controlled_by": system.ControlledBy,
			"planet_count": len(system.Planets),
			"coordinates":  system.Coordinates,
		}
	}
	return systems
}

func (gs *GameServer) getPlayerSystems(playerID string) []map[string]interface{} {
	systems := gs.gameState.Galaxy.GetSystemsByOwner(playerID)
	result := make([]map[string]interface{}, len(systems))
	
	for i, system := range systems {
		planets := system.GetPlanetsByOwner(playerID)
		result[i] = map[string]interface{}{
			"id":           system.ID,
			"name":         system.Name,
			"planet_count": len(planets),
			"coordinates":  system.Coordinates,
		}
	}
	return result
}

