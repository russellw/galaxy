package main

import (
	"flag"
	"fmt"
)

func main() {
	serverMode := flag.Bool("server", false, "Run as server")
	flag.Parse()
	
	if *serverMode {
		runServer()
		return
	}
	
	runSimulation()
}

func runServer() {
	players := []Player{
		{ID: "player1", Name: "Terran Federation"},
		{ID: "player2", Name: "Zephyrian Empire"},
		{ID: "player3", Name: "Cosmic Alliance"},
		{ID: "player4", Name: "Nova Collective"},
	}
	
	server := NewGameServer(players, 20, 50, 30) // 30 second turns
	server.StartServer(8080)
}

func runSimulation() {
	fmt.Println("Galaxy Strategy Game - Turn-Based Test")
	fmt.Println("======================================")

	players := []Player{
		{ID: "p1", Name: "Terran Federation"},
		{ID: "p2", Name: "Zephyrian Empire"},
		{ID: "p3", Name: "Cosmic Alliance"},
	}

	fmt.Printf("Starting game with %d players:\n", len(players))
	for i, player := range players {
		fmt.Printf("%d. %s (%s)\n", i+1, player.Name, player.ID)
	}

	gameState := NewGameState(players, 15, 10)
	
	fmt.Printf("\nGame initialized: %d systems, %d turns max\n", len(gameState.Galaxy.StarSystems), gameState.MaxTurns)

	// Show initial state
	fmt.Println("\nInitial Player Status:")
	fmt.Println("=====================")
	for _, player := range players {
		fmt.Printf("%s: %s\n", player.Name, gameState.GetPlayerSummary(player.ID))
	}

	// Simulate a few turns
	for turn := 1; turn <= 3 && !gameState.GameOver; turn++ {
		fmt.Printf("\n============================================================")
		fmt.Printf("\nTURN %d\n", turn)
		fmt.Printf("============================================================\n")

		// Generate sample orders for each player
		addSampleOrders(&gameState)

		// Process the turn
		gameState.ProcessTurn()

		// Show player status after turn
		fmt.Println("\nPlayer Status After Turn:")
		for _, player := range players {
			fmt.Printf("%s: %s\n", player.Name, gameState.GetPlayerSummary(player.ID))
		}
	}

	fmt.Println("\n============================================================")
	fmt.Println("GAME SIMULATION COMPLETE")
	fmt.Printf("============================================================\n")
	
	if gameState.GameOver && gameState.Winner != "" {
		fmt.Printf("Winner: %s\n", gameState.Winner)
	} else {
		fmt.Println("Game still in progress")
	}
}

func addSampleOrders(gs *GameState) {
	fmt.Println("Generating sample orders for players...")
	
	for _, player := range gs.Players {
		systems := gs.Galaxy.GetSystemsByOwner(player.ID)
		if len(systems) > 0 {
			homeworld := systems[0].GetPlanetsByOwner(player.ID)[0]
			
			// Build facilities
			buildOrder := Order{
				PlayerID:  player.ID,
				OrderType: string(OrderBuildFacility),
				PlanetID:  homeworld.ID,
				Parameters: map[string]interface{}{
					"facility_type": "MetalMine",
				},
				Priority: 5,
			}
			gs.AddOrder(buildOrder)
			
			// Upgrade existing facility
			if len(homeworld.Facilities) > 0 {
				upgradeOrder := Order{
					PlayerID:  player.ID,
					OrderType: string(OrderUpgradeFacility),
					PlanetID:  homeworld.ID,
					Parameters: map[string]interface{}{
						"facility_type": homeworld.Facilities[0].Type,
					},
					Priority: 3,
				}
				gs.AddOrder(upgradeOrder)
			}
			
			// Build ships
			shipOrder := Order{
				PlayerID:  player.ID,
				OrderType: string(OrderBuildShip),
				PlanetID:  homeworld.ID,
				Parameters: map[string]interface{}{
					"ship_type": "Fighter",
				},
				Priority: 4,
			}
			gs.AddOrder(shipOrder)
		}
	}
}