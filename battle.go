package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type BattleResult struct {
	Winner    string
	Survivors []Spaceship
	Rounds    []BattleRound
}

type BattleRound struct {
	RoundNumber int
	Attacks     []Attack
}

type Attack struct {
	Attacker string
	Target   string
	Damage   int
	Hit      bool
}

func RunSpaceBattle(fleet1, fleet2 Fleet) BattleResult {
	rand.Seed(time.Now().UnixNano())
	
	result := BattleResult{
		Winner:    "",
		Survivors: []Spaceship{},
		Rounds:    []BattleRound{},
	}
	
	roundNumber := 1
	maxRounds := 100
	
	for roundNumber <= maxRounds {
		if fleet1.IsDefeated() || fleet2.IsDefeated() {
			break
		}
		
		round := BattleRound{
			RoundNumber: roundNumber,
			Attacks:     []Attack{},
		}
		
		allShips := append(fleet1.GetAliveShips(), fleet2.GetAliveShips()...)
		sort.Slice(allShips, func(i, j int) bool {
			return allShips[i].Speed > allShips[j].Speed
		})
		
		for _, attacker := range allShips {
			if !attacker.IsAlive() {
				continue
			}
			
			var enemies []Spaceship
			if attacker.Owner == fleet1.Owner {
				enemies = fleet2.GetAliveShips()
			} else {
				enemies = fleet1.GetAliveShips()
			}
			
			if len(enemies) == 0 {
				break
			}
			
			target := enemies[rand.Intn(len(enemies))]
			
			hitChance := 0.7
			hit := rand.Float64() < hitChance
			
			attack := Attack{
				Attacker: attacker.ID,
				Target:   target.ID,
				Damage:   0,
				Hit:      hit,
			}
			
			if hit {
				attack.Damage = attacker.Attack
				
				for i := range fleet1.Ships {
					if fleet1.Ships[i].ID == target.ID {
						fleet1.Ships[i].TakeDamage(attack.Damage)
						break
					}
				}
				for i := range fleet2.Ships {
					if fleet2.Ships[i].ID == target.ID {
						fleet2.Ships[i].TakeDamage(attack.Damage)
						break
					}
				}
			}
			
			round.Attacks = append(round.Attacks, attack)
		}
		
		result.Rounds = append(result.Rounds, round)
		roundNumber++
	}
	
	if fleet1.IsDefeated() && fleet2.IsDefeated() {
		result.Winner = "Draw"
	} else if fleet1.IsDefeated() {
		result.Winner = fleet2.Owner
		result.Survivors = fleet2.GetAliveShips()
	} else if fleet2.IsDefeated() {
		result.Winner = fleet1.Owner
		result.Survivors = fleet1.GetAliveShips()
	} else {
		result.Winner = "Draw"
	}
	
	return result
}

func PrintBattleResult(result BattleResult) {
	fmt.Printf("Battle Result: Winner is %s\n", result.Winner)
	fmt.Printf("Battle lasted %d rounds\n", len(result.Rounds))
	fmt.Printf("Survivors: %d ships\n", len(result.Survivors))
	
	for _, survivor := range result.Survivors {
		fmt.Printf("  %s (%s) - Hull: %d/%d, Shields: %d/%d\n", 
			survivor.Name, survivor.ID, survivor.Hull, survivor.MaxHull, 
			survivor.Shields, survivor.MaxShields)
	}
}