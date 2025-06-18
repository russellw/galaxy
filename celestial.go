package main

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