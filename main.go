package main

import "fmt"

func main() {
	fmt.Println("Galaxy Space Battle Test")
	fmt.Println("========================")

	ship1 := NewSpaceship("f1s1", "Destroyer Alpha", "Player1", 100, 10, 50, 25, 15)
	ship2 := NewSpaceship("f1s2", "Cruiser Beta", "Player1", 150, 15, 75, 20, 12)
	
	ship3 := NewSpaceship("f2s1", "Fighter Gamma", "Player2", 80, 5, 30, 30, 20)
	ship4 := NewSpaceship("f2s2", "Battleship Delta", "Player2", 200, 20, 100, 35, 8)

	fleet1 := NewFleet("fleet1", "Player1", "Sector A", []Spaceship{ship1, ship2})
	fleet2 := NewFleet("fleet2", "Player2", "Sector A", []Spaceship{ship3, ship4})

	fmt.Printf("Fleet 1 (%s): %d ships\n", fleet1.Owner, len(fleet1.Ships))
	for _, ship := range fleet1.Ships {
		fmt.Printf("  %s - Hull:%d Armor:%d Shields:%d Attack:%d Speed:%d\n", 
			ship.Name, ship.Hull, ship.Armor, ship.Shields, ship.Attack, ship.Speed)
	}

	fmt.Printf("\nFleet 2 (%s): %d ships\n", fleet2.Owner, len(fleet2.Ships))
	for _, ship := range fleet2.Ships {
		fmt.Printf("  %s - Hull:%d Armor:%d Shields:%d Attack:%d Speed:%d\n", 
			ship.Name, ship.Hull, ship.Armor, ship.Shields, ship.Attack, ship.Speed)
	}

	fmt.Println("\nStarting battle...")
	result := RunSpaceBattle(fleet1, fleet2)

	fmt.Println("\n==================================================")
	PrintBattleResult(result)

	fmt.Println("\nDetailed battle log:")
	for _, round := range result.Rounds {
		fmt.Printf("Round %d:\n", round.RoundNumber)
		for _, attack := range round.Attacks {
			if attack.Hit {
				fmt.Printf("  %s attacks %s for %d damage\n", attack.Attacker, attack.Target, attack.Damage)
			} else {
				fmt.Printf("  %s attacks %s but misses\n", attack.Attacker, attack.Target)
			}
		}
	}
}