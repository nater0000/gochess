package gostuff

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/freeeve/pgn.v1"
)

type allPgnGames []GameAnalysis

func GameAnalysisByPgn(w http.ResponseWriter, r *http.Request) {
	valid := ValidateCredentials(w, r)
	if valid == false {
		return
	}

	// Get the gameID specified in the front end
	pgnData := r.FormValue("pgnData")
	depth, err := strconv.Atoi(template.HTMLEscapeString(r.FormValue("depth")))

	if err != nil {
		fmt.Println("Could not convert string to int in GetEngineAnalysisById", err)
		return
	}

	if depth < 1 || depth > 7 {
		fmt.Println("Depth is not in a valid range: ", depth)
		return
	}

	var allGames allPgnGames
	allGames.readPgn(pgnData)

	jsonGamesAnalysis, err := json.Marshal(allGames)
	if err != nil {
		fmt.Println("Could not marshal gameAnalysis", err)
		return
	}
	w.Write([]byte((jsonGamesAnalysis)))
}

// Reads pgn and converts it to a game
func (allGames *allPgnGames) readPgn(pgnData string) {

	reader := strings.NewReader(pgnData)
	ps := pgn.NewPGNScanner(reader)

	engine := startEngine(nil)

	for ps.Next() {
		game, err := ps.Scan()
		if err != nil {
			fmt.Println("Can't read pgn string for", pgnData, err)
			return
		}
		fmt.Println(game.Tags)

		// TODO:
		// cannot use game.Moves (type []pgn.Move) as type []chess.Move in argument to allGames.analyzePgnGames
		// convert pgn.Move into chess.Move
		allGames.analyzePgnGames(game.Moves, engine)
	}

	engine.Quit()
	fmt.Println("Done scanning pgn games")
}
