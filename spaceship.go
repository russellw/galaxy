package main

type Spaceship struct {
	ID          string
	Name        string
	Owner       string
	Hull        int
	MaxHull     int
	Armor       int
	Shields     int
	MaxShields  int
	Attack      int
	Speed       int
	IsDestroyed bool
}

type Fleet struct {
	ID        string
	Owner     string
	Ships     []Spaceship
	Location  string
}

func NewSpaceship(id, name, owner string, hull, armor, shields, attack, speed int) Spaceship {
	return Spaceship{
		ID:          id,
		Name:        name,
		Owner:       owner,
		Hull:        hull,
		MaxHull:     hull,
		Armor:       armor,
		Shields:     shields,
		MaxShields:  shields,
		Attack:      attack,
		Speed:       speed,
		IsDestroyed: false,
	}
}

func NewFleet(id, owner, location string, ships []Spaceship) Fleet {
	return Fleet{
		ID:       id,
		Owner:    owner,
		Ships:    ships,
		Location: location,
	}
}

func (s *Spaceship) TakeDamage(damage int) {
	remainingDamage := damage
	
	if s.Shields > 0 {
		if remainingDamage >= s.Shields {
			remainingDamage -= s.Shields
			s.Shields = 0
		} else {
			s.Shields -= remainingDamage
			remainingDamage = 0
		}
	}
	
	if remainingDamage > 0 {
		effectiveDamage := remainingDamage - s.Armor
		if effectiveDamage > 0 {
			s.Hull -= effectiveDamage
			if s.Hull <= 0 {
				s.Hull = 0
				s.IsDestroyed = true
			}
		}
	}
}

func (s *Spaceship) IsAlive() bool {
	return !s.IsDestroyed && s.Hull > 0
}

func (f *Fleet) GetAliveShips() []Spaceship {
	var aliveShips []Spaceship
	for _, ship := range f.Ships {
		if ship.IsAlive() {
			aliveShips = append(aliveShips, ship)
		}
	}
	return aliveShips
}

func (f *Fleet) IsDefeated() bool {
	return len(f.GetAliveShips()) == 0
}