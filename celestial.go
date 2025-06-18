package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Star struct {
	ID          string
	Name        string
	StarType    string
	Size        float64
	Temperature int
	Luminosity  float64
	Age         int64
	Coordinates Coordinates
}

type Planet struct {
	ID           string
	Name         string
	StarSystemID string
	Owner        string
	PlanetType   string
	Size         float64
	Population   int64
	Resources    Resources
	Facilities   []Facility
	OrbitalPos   int
	Habitable    bool
	Atmosphere   string
	Temperature  int
}

type StarSystem struct {
	ID          string
	Name        string
	Star        Star
	Planets     []Planet
	Coordinates Coordinates
	Explored    bool
	ControlledBy string
}

type Coordinates struct {
	X float64
	Y float64
	Z float64
}

type Resources struct {
	Metals     int
	Energy     int
	Minerals   int
	Food       int
	Technology int
}

type Facility struct {
	ID       string
	Type     string
	Level    int
	Output   int
	PlanetID string
}

type Galaxy struct {
	ID          string
	Name        string
	StarSystems []StarSystem
	Size        int
}

func NewStar(id, name, starType string, size, luminosity float64, temperature int, age int64, coords Coordinates) Star {
	return Star{
		ID:          id,
		Name:        name,
		StarType:    starType,
		Size:        size,
		Temperature: temperature,
		Luminosity:  luminosity,
		Age:         age,
		Coordinates: coords,
	}
}

func NewPlanet(id, name, systemID, owner, planetType string, size float64, orbitalPos int, habitable bool) Planet {
	return Planet{
		ID:           id,
		Name:         name,
		StarSystemID: systemID,
		Owner:        owner,
		PlanetType:   planetType,
		Size:         size,
		Population:   0,
		Resources:    Resources{},
		Facilities:   []Facility{},
		OrbitalPos:   orbitalPos,
		Habitable:    habitable,
		Atmosphere:   "None",
		Temperature:  -200,
	}
}

func NewStarSystem(id, name string, star Star, coords Coordinates) StarSystem {
	return StarSystem{
		ID:           id,
		Name:         name,
		Star:         star,
		Planets:      []Planet{},
		Coordinates:  coords,
		Explored:     false,
		ControlledBy: "",
	}
}

func (p *Planet) AddFacility(facilityType string, level int) {
	facility := Facility{
		ID:       p.ID + "_" + facilityType,
		Type:     facilityType,
		Level:    level,
		Output:   level * 10,
		PlanetID: p.ID,
	}
	p.Facilities = append(p.Facilities, facility)
}

func (p *Planet) GetTotalProduction(resourceType string) int {
	total := 0
	for _, facility := range p.Facilities {
		if facility.Type == resourceType {
			total += facility.Output
		}
	}
	return total
}

func (s *StarSystem) AddPlanet(planet Planet) {
	s.Planets = append(s.Planets, planet)
}

func (s *StarSystem) GetHabitablePlanets() []Planet {
	var habitable []Planet
	for _, planet := range s.Planets {
		if planet.Habitable {
			habitable = append(habitable, planet)
		}
	}
	return habitable
}

func (s *StarSystem) GetPlanetsByOwner(owner string) []Planet {
	var owned []Planet
	for _, planet := range s.Planets {
		if planet.Owner == owner {
			owned = append(owned, planet)
		}
	}
	return owned
}

func (g *Galaxy) AddStarSystem(system StarSystem) {
	g.StarSystems = append(g.StarSystems, system)
}

func (g *Galaxy) GetSystemByID(id string) *StarSystem {
	for i := range g.StarSystems {
		if g.StarSystems[i].ID == id {
			return &g.StarSystems[i]
		}
	}
	return nil
}

func (g *Galaxy) GetSystemsByOwner(owner string) []StarSystem {
	var controlled []StarSystem
	for _, system := range g.StarSystems {
		if system.ControlledBy == owner {
			controlled = append(controlled, system)
		}
	}
	return controlled
}

func CalculateDistance(coord1, coord2 Coordinates) float64 {
	dx := coord1.X - coord2.X
	dy := coord1.Y - coord2.Y
	dz := coord1.Z - coord2.Z
	return (dx*dx + dy*dy + dz*dz) * 0.5
}

type Player struct {
	ID   string
	Name string
}

func InitializeGalaxy(players []Player, galaxySize int) Galaxy {
	rand.Seed(time.Now().UnixNano())
	
	galaxy := Galaxy{
		ID:          "galaxy_1",
		Name:        "New Galaxy",
		StarSystems: []StarSystem{},
		Size:        galaxySize,
	}
	
	playerCount := len(players)
	if playerCount == 0 {
		return galaxy
	}
	
	homeworlds := make([]StarSystem, playerCount)
	
	for i, player := range players {
		coords := generateHomeworldCoordinates(i, playerCount, galaxySize)
		
		star := NewStar(
			fmt.Sprintf("star_%s", player.ID),
			fmt.Sprintf("%s Prime", player.Name),
			"G-Class",
			1.0,
			1.0,
			5778,
			4600000000,
			coords,
		)
		
		system := NewStarSystem(
			fmt.Sprintf("system_%s", player.ID),
			fmt.Sprintf("%s System", player.Name),
			star,
			coords,
		)
		
		homeworld := NewPlanet(
			fmt.Sprintf("planet_%s_home", player.ID),
			fmt.Sprintf("%s Prime", player.Name),
			system.ID,
			player.ID,
			"Terrestrial",
			1.0,
			3,
			true,
		)
		
		homeworld.Population = 1000000
		homeworld.Atmosphere = "Oxygen-Nitrogen"
		homeworld.Temperature = 15
		homeworld.Resources = Resources{
			Metals:     100,
			Energy:     50,
			Minerals:   75,
			Food:       200,
			Technology: 25,
		}
		
		homeworld.AddFacility("MetalMine", 2)
		homeworld.AddFacility("PowerPlant", 2)
		homeworld.AddFacility("Farm", 3)
		homeworld.AddFacility("Factory", 1)
		
		system.AddPlanet(homeworld)
		system.Explored = true
		system.ControlledBy = player.ID
		
		for j := 1; j <= 5; j++ {
			if j == 3 {
				continue
			}
			
			planet := generateRandomPlanet(fmt.Sprintf("planet_%s_%d", player.ID, j), system.ID, j)
			system.AddPlanet(planet)
		}
		
		homeworlds[i] = system
	}
	
	for _, system := range homeworlds {
		galaxy.AddStarSystem(system)
	}
	
	neutralSystemCount := galaxySize - playerCount
	for i := 0; i < neutralSystemCount; i++ {
		coords := generateRandomCoordinates(galaxySize)
		
		star := generateRandomStar(fmt.Sprintf("star_neutral_%d", i), coords)
		system := NewStarSystem(
			fmt.Sprintf("system_neutral_%d", i),
			fmt.Sprintf("System-%d", i+1),
			star,
			coords,
		)
		
		planetCount := rand.Intn(6) + 2
		for j := 1; j <= planetCount; j++ {
			planet := generateRandomPlanet(fmt.Sprintf("planet_neutral_%d_%d", i, j), system.ID, j)
			system.AddPlanet(planet)
		}
		
		galaxy.AddStarSystem(system)
	}
	
	return galaxy
}

func generateHomeworldCoordinates(playerIndex, totalPlayers, galaxySize int) Coordinates {
	angle := float64(playerIndex) * 2.0 * 3.14159 / float64(totalPlayers)
	radius := float64(galaxySize) * 0.3
	
	return Coordinates{
		X: radius * float64(rand.Float64()*0.2+0.9) * (angle * 0.159155),
		Y: radius * float64(rand.Float64()*0.2+0.9) * (angle * 0.318310),
		Z: float64(rand.Intn(20) - 10),
	}
}

func generateRandomCoordinates(galaxySize int) Coordinates {
	maxCoord := float64(galaxySize)
	return Coordinates{
		X: (rand.Float64() - 0.5) * maxCoord,
		Y: (rand.Float64() - 0.5) * maxCoord,
		Z: (rand.Float64() - 0.5) * maxCoord * 0.2,
	}
}

func generateRandomStar(id string, coords Coordinates) Star {
	starTypes := []string{"G-Class", "K-Class", "M-Class", "F-Class", "A-Class"}
	starNames := []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon", "Zeta", "Eta", "Theta"}
	
	starType := starTypes[rand.Intn(len(starTypes))]
	name := starNames[rand.Intn(len(starNames))]
	
	var temp int
	var size, luminosity float64
	
	switch starType {
	case "M-Class":
		temp = 3000 + rand.Intn(1000)
		size = 0.3 + rand.Float64()*0.4
		luminosity = 0.01 + rand.Float64()*0.09
	case "K-Class":
		temp = 4000 + rand.Intn(1200)
		size = 0.7 + rand.Float64()*0.3
		luminosity = 0.1 + rand.Float64()*0.4
	case "G-Class":
		temp = 5200 + rand.Intn(800)
		size = 0.9 + rand.Float64()*0.2
		luminosity = 0.8 + rand.Float64()*0.4
	case "F-Class":
		temp = 6000 + rand.Intn(1000)
		size = 1.1 + rand.Float64()*0.3
		luminosity = 1.5 + rand.Float64()*1.0
	case "A-Class":
		temp = 7500 + rand.Intn(2500)
		size = 1.4 + rand.Float64()*0.6
		luminosity = 5.0 + rand.Float64()*20.0
	}
	
	return NewStar(id, name, starType, size, luminosity, temp, int64(rand.Intn(10000000000)), coords)
}

func generateRandomPlanet(id, systemID string, orbitalPos int) Planet {
	planetTypes := []string{"Rocky", "Gas Giant", "Ice World", "Desert", "Ocean World"}
	planetType := planetTypes[rand.Intn(len(planetTypes))]
	
	size := 0.5 + rand.Float64()*2.0
	habitable := false
	
	if orbitalPos >= 2 && orbitalPos <= 4 && (planetType == "Rocky" || planetType == "Ocean World") {
		habitable = rand.Float64() < 0.3
	}
	
	planet := NewPlanet(id, fmt.Sprintf("Planet-%d", orbitalPos), systemID, "", planetType, size, orbitalPos, habitable)
	
	if habitable {
		planet.Atmosphere = "Oxygen-Nitrogen"
		planet.Temperature = -10 + rand.Intn(40)
		planet.Resources = Resources{
			Metals:     rand.Intn(150),
			Energy:     rand.Intn(100),
			Minerals:   rand.Intn(200),
			Food:       rand.Intn(50),
			Technology: 0,
		}
	} else {
		planet.Atmosphere = "None"
		planet.Temperature = -200 + rand.Intn(600)
		planet.Resources = Resources{
			Metals:     rand.Intn(300),
			Energy:     rand.Intn(200),
			Minerals:   rand.Intn(400),
			Food:       0,
			Technology: 0,
		}
	}
	
	return planet
}