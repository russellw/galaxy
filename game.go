package main

import (
	"fmt"
	"sort"
)

type GameState struct {
	Galaxy      Galaxy
	Players     []Player
	CurrentTurn int
	MaxTurns    int
	Orders      map[string][]Order
	GameOver    bool
	Winner      string
}

type Order struct {
	PlayerID    string
	OrderType   string
	PlanetID    string
	SystemID    string
	Parameters  map[string]interface{}
	Priority    int
}

type OrderType string

const (
	OrderBuildFacility    OrderType = "BUILD_FACILITY"
	OrderUpgradeFacility  OrderType = "UPGRADE_FACILITY"
	OrderBuildShip        OrderType = "BUILD_SHIP"
	OrderMoveFleet        OrderType = "MOVE_FLEET"
	OrderColonizePlanet   OrderType = "COLONIZE_PLANET"
	OrderResearch         OrderType = "RESEARCH"
)

func NewGameState(players []Player, galaxySize int, maxTurns int) GameState {
	galaxy := InitializeGalaxy(players, galaxySize)
	
	return GameState{
		Galaxy:      galaxy,
		Players:     players,
		CurrentTurn: 1,
		MaxTurns:    maxTurns,
		Orders:      make(map[string][]Order),
		GameOver:    false,
		Winner:      "",
	}
}

func (gs *GameState) AddOrder(order Order) {
	if gs.Orders[order.PlayerID] == nil {
		gs.Orders[order.PlayerID] = []Order{}
	}
	gs.Orders[order.PlayerID] = append(gs.Orders[order.PlayerID], order)
}

func (gs *GameState) ProcessTurn() {
	fmt.Printf("\n=== Processing Turn %d ===\n", gs.CurrentTurn)
	
	// Sort orders by priority
	allOrders := []Order{}
	for _, playerOrders := range gs.Orders {
		allOrders = append(allOrders, playerOrders...)
	}
	sort.Slice(allOrders, func(i, j int) bool {
		return allOrders[i].Priority > allOrders[j].Priority
	})
	
	// Process production orders first
	gs.processProductionOrders()
	
	// Process movement orders
	gs.processMovementOrders()
	
	// Process construction orders
	gs.processConstructionOrders()
	
	// Update resources
	gs.updateResources()
	
	// Clear orders for next turn
	gs.Orders = make(map[string][]Order)
	
	// Check win conditions
	gs.checkWinConditions()
	
	gs.CurrentTurn++
	if gs.CurrentTurn > gs.MaxTurns {
		gs.GameOver = true
		gs.determineWinner()
	}
	
	fmt.Printf("Turn %d completed.\n", gs.CurrentTurn-1)
}

func (gs *GameState) processProductionOrders() {
	fmt.Println("Processing production orders...")
	
	for playerID, orders := range gs.Orders {
		for _, order := range orders {
			switch OrderType(order.OrderType) {
			case OrderBuildShip:
				gs.processBuildShipOrder(order)
			case OrderBuildFacility:
				gs.processBuildFacilityOrder(order)
			case OrderUpgradeFacility:
				gs.processUpgradeFacilityOrder(order)
			}
		}
		_ = playerID
	}
}

func (gs *GameState) processMovementOrders() {
	fmt.Println("Processing movement orders...")
	
	for _, orders := range gs.Orders {
		for _, order := range orders {
			if OrderType(order.OrderType) == OrderMoveFleet {
				gs.processMoveFleetOrder(order)
			}
		}
	}
}

func (gs *GameState) processConstructionOrders() {
	fmt.Println("Processing construction orders...")
	
	for _, orders := range gs.Orders {
		for _, order := range orders {
			if OrderType(order.OrderType) == OrderColonizePlanet {
				gs.processColonizeOrder(order)
			}
		}
	}
}

func (gs *GameState) processBuildShipOrder(order Order) {
	planet := gs.findPlanet(order.PlanetID)
	if planet == nil || planet.Owner != order.PlayerID {
		return
	}
	
	shipType, ok := order.Parameters["ship_type"].(string)
	if !ok {
		return
	}
	
	cost := gs.getShipCost(shipType)
	if planet.Resources.Metals >= cost.Metals && planet.Resources.Energy >= cost.Energy {
		planet.Resources.Metals -= cost.Metals
		planet.Resources.Energy -= cost.Energy
		
		// Create ship (simplified - would normally add to fleet)
		fmt.Printf("Player %s built %s on %s\n", order.PlayerID, shipType, planet.Name)
	}
}

func (gs *GameState) processBuildFacilityOrder(order Order) {
	planet := gs.findPlanet(order.PlanetID)
	if planet == nil || planet.Owner != order.PlayerID {
		return
	}
	
	facilityType, ok := order.Parameters["facility_type"].(string)
	if !ok {
		return
	}
	
	cost := gs.getFacilityCost(facilityType)
	if planet.Resources.Metals >= cost.Metals && planet.Resources.Energy >= cost.Energy {
		planet.Resources.Metals -= cost.Metals
		planet.Resources.Energy -= cost.Energy
		
		planet.AddFacility(facilityType, 1)
		fmt.Printf("Player %s built %s on %s\n", order.PlayerID, facilityType, planet.Name)
	}
}

func (gs *GameState) processUpgradeFacilityOrder(order Order) {
	planet := gs.findPlanet(order.PlanetID)
	if planet == nil || planet.Owner != order.PlayerID {
		return
	}
	
	facilityType, ok := order.Parameters["facility_type"].(string)
	if !ok {
		return
	}
	
	for i := range planet.Facilities {
		if planet.Facilities[i].Type == facilityType {
			cost := gs.getFacilityUpgradeCost(facilityType, planet.Facilities[i].Level)
			if planet.Resources.Metals >= cost.Metals && planet.Resources.Energy >= cost.Energy {
				planet.Resources.Metals -= cost.Metals
				planet.Resources.Energy -= cost.Energy
				
				planet.Facilities[i].Level++
				planet.Facilities[i].Output = planet.Facilities[i].Level * 10
				fmt.Printf("Player %s upgraded %s to level %d on %s\n", 
					order.PlayerID, facilityType, planet.Facilities[i].Level, planet.Name)
			}
			break
		}
	}
}

func (gs *GameState) processMoveFleetOrder(order Order) {
	// Simplified fleet movement
	fmt.Printf("Player %s moving fleet from %s to %s\n", 
		order.PlayerID, order.Parameters["from"], order.Parameters["to"])
}

func (gs *GameState) processColonizeOrder(order Order) {
	planet := gs.findPlanet(order.PlanetID)
	if planet == nil || planet.Owner != "" {
		return
	}
	
	if planet.Habitable {
		planet.Owner = order.PlayerID
		planet.Population = 10000
		planet.AddFacility("Colony", 1)
		fmt.Printf("Player %s colonized %s\n", order.PlayerID, planet.Name)
	}
}

func (gs *GameState) updateResources() {
	fmt.Println("Updating resources...")
	
	for i := range gs.Galaxy.StarSystems {
		system := &gs.Galaxy.StarSystems[i]
		for j := range system.Planets {
			planet := &system.Planets[j]
			if planet.Owner != "" {
				// Add resource production from facilities
				planet.Resources.Metals += planet.GetTotalProduction("MetalMine")
				planet.Resources.Energy += planet.GetTotalProduction("PowerPlant")
				planet.Resources.Food += planet.GetTotalProduction("Farm")
				planet.Resources.Technology += planet.GetTotalProduction("Laboratory")
			}
		}
	}
}

func (gs *GameState) checkWinConditions() {
	// Check if any player controls majority of systems
	playerSystemCount := make(map[string]int)
	totalSystems := len(gs.Galaxy.StarSystems)
	
	for _, system := range gs.Galaxy.StarSystems {
		if system.ControlledBy != "" {
			playerSystemCount[system.ControlledBy]++
		}
	}
	
	for playerID, count := range playerSystemCount {
		if count > totalSystems/2 {
			gs.GameOver = true
			gs.Winner = playerID
			return
		}
	}
}

func (gs *GameState) determineWinner() {
	playerSystemCount := make(map[string]int)
	
	for _, system := range gs.Galaxy.StarSystems {
		if system.ControlledBy != "" {
			playerSystemCount[system.ControlledBy]++
		}
	}
	
	maxSystems := 0
	for playerID, count := range playerSystemCount {
		if count > maxSystems {
			maxSystems = count
			gs.Winner = playerID
		}
	}
}

func (gs *GameState) findPlanet(planetID string) *Planet {
	for i := range gs.Galaxy.StarSystems {
		for j := range gs.Galaxy.StarSystems[i].Planets {
			if gs.Galaxy.StarSystems[i].Planets[j].ID == planetID {
				return &gs.Galaxy.StarSystems[i].Planets[j]
			}
		}
	}
	return nil
}

func (gs *GameState) getShipCost(shipType string) Resources {
	costs := map[string]Resources{
		"Fighter":    {Metals: 50, Energy: 25, Minerals: 0, Food: 0, Technology: 0},
		"Destroyer":  {Metals: 100, Energy: 50, Minerals: 25, Food: 0, Technology: 0},
		"Cruiser":    {Metals: 200, Energy: 100, Minerals: 50, Food: 0, Technology: 0},
		"Battleship": {Metals: 400, Energy: 200, Minerals: 100, Food: 0, Technology: 0},
	}
	
	if cost, exists := costs[shipType]; exists {
		return cost
	}
	return Resources{Metals: 100, Energy: 50, Minerals: 0, Food: 0, Technology: 0}
}

func (gs *GameState) getFacilityCost(facilityType string) Resources {
	costs := map[string]Resources{
		"MetalMine":   {Metals: 50, Energy: 25, Minerals: 0, Food: 0, Technology: 0},
		"PowerPlant":  {Metals: 75, Energy: 0, Minerals: 25, Food: 0, Technology: 0},
		"Farm":        {Metals: 25, Energy: 10, Minerals: 0, Food: 0, Technology: 0},
		"Factory":     {Metals: 100, Energy: 50, Minerals: 50, Food: 0, Technology: 0},
		"Laboratory":  {Metals: 150, Energy: 75, Minerals: 25, Food: 0, Technology: 0},
	}
	
	if cost, exists := costs[facilityType]; exists {
		return cost
	}
	return Resources{Metals: 50, Energy: 25, Minerals: 0, Food: 0, Technology: 0}
}

func (gs *GameState) getFacilityUpgradeCost(facilityType string, currentLevel int) Resources {
	baseCost := gs.getFacilityCost(facilityType)
	multiplier := currentLevel + 1
	
	return Resources{
		Metals:     baseCost.Metals * multiplier,
		Energy:     baseCost.Energy * multiplier,
		Minerals:   baseCost.Minerals * multiplier,
		Food:       baseCost.Food * multiplier,
		Technology: baseCost.Technology * multiplier,
	}
}

func (gs *GameState) GetPlayerSummary(playerID string) string {
	systems := gs.Galaxy.GetSystemsByOwner(playerID)
	totalPlanets := 0
	totalPopulation := int64(0)
	totalResources := Resources{}
	
	for _, system := range systems {
		planets := system.GetPlanetsByOwner(playerID)
		totalPlanets += len(planets)
		
		for _, planet := range planets {
			totalPopulation += planet.Population
			totalResources.Metals += planet.Resources.Metals
			totalResources.Energy += planet.Resources.Energy
			totalResources.Minerals += planet.Resources.Minerals
			totalResources.Food += planet.Resources.Food
			totalResources.Technology += planet.Resources.Technology
		}
	}
	
	return fmt.Sprintf("Systems: %d, Planets: %d, Population: %d, Resources: M=%d E=%d Min=%d F=%d T=%d",
		len(systems), totalPlanets, totalPopulation,
		totalResources.Metals, totalResources.Energy, totalResources.Minerals,
		totalResources.Food, totalResources.Technology)
}