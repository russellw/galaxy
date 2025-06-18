package main

import "fmt"

func main() {
	fmt.Println("Galaxy Initialization Test")
	fmt.Println("==========================")

	players := []Player{
		{ID: "p1", Name: "Terran Federation"},
		{ID: "p2", Name: "Zephyrian Empire"},
		{ID: "p3", Name: "Cosmic Alliance"},
		{ID: "p4", Name: "Nova Collective"},
		{ID: "p5", Name: "Stellar Republic"},
		{ID: "p6", Name: "Void Consortium"},
		{ID: "p7", Name: "Galactic Union"},
	}

	fmt.Printf("Initializing galaxy with %d players:\n", len(players))
	for i, player := range players {
		fmt.Printf("%d. %s (%s)\n", i+1, player.Name, player.ID)
	}

	galaxySize := 20
	galaxy := InitializeGalaxy(players, galaxySize)

	fmt.Printf("\nGalaxy '%s' created with %d star systems\n", galaxy.Name, len(galaxy.StarSystems))

	fmt.Println("\nPlayer Homeworlds:")
	fmt.Println("==================")
	for _, player := range players {
		systems := galaxy.GetSystemsByOwner(player.ID)
		if len(systems) > 0 {
			system := systems[0]
			homeworld := system.GetPlanetsByOwner(player.ID)[0]
			
			fmt.Printf("\n%s (%s):\n", player.Name, player.ID)
			fmt.Printf("  System: %s at (%.1f, %.1f, %.1f)\n", 
				system.Name, system.Coordinates.X, system.Coordinates.Y, system.Coordinates.Z)
			fmt.Printf("  Star: %s (%s, %.1f solar masses, %dK)\n", 
				system.Star.Name, system.Star.StarType, system.Star.Size, system.Star.Temperature)
			fmt.Printf("  Homeworld: %s\n", homeworld.Name)
			fmt.Printf("    Population: %d\n", homeworld.Population)
			fmt.Printf("    Resources: Metals=%d Energy=%d Minerals=%d Food=%d Tech=%d\n",
				homeworld.Resources.Metals, homeworld.Resources.Energy, 
				homeworld.Resources.Minerals, homeworld.Resources.Food, homeworld.Resources.Technology)
			fmt.Printf("    Facilities: %d total\n", len(homeworld.Facilities))
			for _, facility := range homeworld.Facilities {
				fmt.Printf("      %s (Level %d) - Output: %d\n", facility.Type, facility.Level, facility.Output)
			}
			fmt.Printf("    Other planets in system: %d\n", len(system.Planets)-1)
		}
	}

	fmt.Println("\nNeutral Systems:")
	fmt.Println("================")
	neutralCount := 0
	for _, system := range galaxy.StarSystems {
		if system.ControlledBy == "" {
			neutralCount++
		}
	}
	fmt.Printf("Total neutral systems: %d\n", neutralCount)

	if neutralCount > 0 {
		fmt.Println("\nSample neutral systems:")
		count := 0
		for _, system := range galaxy.StarSystems {
			if system.ControlledBy == "" && count < 3 {
				fmt.Printf("  %s: %s star, %d planets at (%.1f, %.1f, %.1f)\n", 
					system.Name, system.Star.StarType, len(system.Planets),
					system.Coordinates.X, system.Coordinates.Y, system.Coordinates.Z)
				count++
			}
		}
	}

	fmt.Println("\nGalaxy generation complete!")
}