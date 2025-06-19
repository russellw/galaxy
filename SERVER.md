# Galaxy Game Server

A multiplayer turn-based strategy game server where players compete to control star systems.

## Usage

### Build the game
```bash
go build -o galaxy
```

### Run simulation mode (offline)
```bash
./galaxy
```

### Run server mode
```bash
./galaxy -server
```

The server will start on port 8080 with 4 players and 30-second turns.

## API Endpoints

### GET /
Server information and available endpoints

### GET /status
Game status including current turn, connected players, and game state
```json
{
  "success": true,
  "data": {
    "current_turn": 1,
    "max_turns": 50,
    "game_over": false,
    "connected_players": 2,
    "total_players": 4,
    "turn_duration": 30
  }
}
```

### POST /connect
Connect as a player
```json
{
  "player_id": "player1"
}
```

### POST /orders
Submit orders for the current turn
```json
{
  "player_id": "player1",
  "order_type": "BUILD_FACILITY",
  "planet_id": "planet_player1_home",
  "parameters": {
    "facility_type": "MetalMine"
  },
  "priority": 5
}
```

### GET /player/{id}
Get player-specific information

### GET /game
Get full game state

### POST /turn
Manual turn control (admin)
```json
{
  "action": "process"
}
```

## Order Types

- `BUILD_FACILITY` - Build a new facility on a planet
- `UPGRADE_FACILITY` - Upgrade an existing facility
- `BUILD_SHIP` - Build a spaceship
- `MOVE_FLEET` - Move ships between systems
- `COLONIZE_PLANET` - Colonize an uninhabited planet

## Players

Default players:
- player1: Terran Federation
- player2: Zephyrian Empire  
- player3: Cosmic Alliance
- player4: Nova Collective

## Testing

Run the test script to verify server functionality:
```bash
./test_server.sh
```

Or test manually with curl:
```bash
# Get status
curl http://localhost:8080/status

# Connect as player
curl -X POST http://localhost:8080/connect \
  -H "Content-Type: application/json" \
  -d '{"player_id": "player1"}'

# Submit order
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"player_id": "player1", "order_type": "BUILD_FACILITY", "planet_id": "planet_player1_home", "parameters": {"facility_type": "MetalMine"}, "priority": 5}'
```