package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Game struct {
	Ants             int
	StartIndex       int
	HasStart         bool
	EndIndex         int
	HasEnd           bool
	RoomNames        []string
	Coordinates      [][]int
	Connections      [][]int
	MaxFlow          int
	PathQueue        [][]int
	PathsFound       bool
	ValidPaths       [][]int
	InputConnections string
}

type Paths struct {
	Comb     [][]int
	CombFlow int
	CombLen  int
}

type PathStorage struct {
	Paths []Paths
}

// Struct for best path combinationa found
type FinalPath struct {
	Len  int
	Ants int
	Path []int
}

// Structs for result stage/printinga nd calculations
type Result struct {
	ActiveRow int
	Result    []Row
}
type Row struct {
	Output     string
	PathsTaken []int
}

// reads given data
func ReadFile(fn string) Game {
	var newGame Game
	connectionFound := false

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	dataCollection := strings.Split(string(data), "\n")

	roomPattern := "^[a-zA-Z0-9]+\\s[0-9]+\\s[0-9]+$"
	tunnelPattern := "^[a-zA-Z0-9]+\\-[a-zA-Z0-9]+$"

	for i := 0; i < len(dataCollection); i++ {
		// ant count
		if i == 0 {
			newGame.Ants, err = strconv.Atoi(dataCollection[i])
			if err != nil {
				fmt.Println("ERROR: invalid data format")
				os.Exit(0)
			}

			if newGame.Ants < 0 {
				fmt.Println("ERROR: invalid data format, invalid number of Ants")
				os.Exit(0)
			}
			continue
		}
		// find start and end
		switch dataCollection[i] {
		case "##start":
			newGame.HasStart = true
			newGame.StartIndex = len(newGame.RoomNames)
			newGame.RoomNames = append(newGame.RoomNames, dataCollection[i+1][:strings.IndexByte(dataCollection[i+1], ' ')])
			i++
			// Save coordinates
			coordinates := strings.Split(dataCollection[i], " ")
			x, err := strconv.Atoi(coordinates[1])
			if err != nil || x < 0 {
				fmt.Print("ERROR: Invalid room coordinates")
				os.Exit(0)
			}
			y, err := strconv.Atoi(coordinates[2])
			if err != nil || y < 0 {
				fmt.Print("ERROR: Invalid room coordinates")
				os.Exit(0)
			}
			newGame.Coordinates = append(newGame.Coordinates, []int{x, y})
			continue
		case "##end":
			newGame.HasEnd = true
			newGame.EndIndex = len(newGame.RoomNames)
			newGame.RoomNames = append(newGame.RoomNames, dataCollection[i+1][:strings.IndexByte(dataCollection[i+1], ' ')])
			i++
			// Save coordinates
			coordinates := strings.Split(dataCollection[i], " ")
			x, err := strconv.Atoi(coordinates[1])
			if err != nil || x < 0 {
				fmt.Print("ERROR: Invalid room coordinates")
				os.Exit(0)
			}
			y, err := strconv.Atoi(coordinates[2])
			if err != nil || y < 0 {
				fmt.Print("ERROR: Invalid room coordinates")
				os.Exit(0)
			}
			newGame.Coordinates = append(newGame.Coordinates, []int{x, y})
			continue
		// Save rooms
		default:
			roomMatched, _ := regexp.MatchString(roomPattern, dataCollection[i])
			tunnelMatched, _ := regexp.MatchString(tunnelPattern, dataCollection[i])
			if roomMatched {
				if connectionFound == true {
					fmt.Print("ERROR: Invalid input order")
					os.Exit(0)
				}
				newGame.RoomNames = append(newGame.RoomNames, dataCollection[i][:strings.IndexByte(dataCollection[i], ' ')])
				// Save coordinates
				coordinates := strings.Split(dataCollection[i], " ")
				x, err := strconv.Atoi(coordinates[1])
				if err != nil || x < 0 {
					fmt.Print("ERROR: Invalid room coordinates")
					os.Exit(0)
				}
				y, err := strconv.Atoi(coordinates[2])
				if err != nil || y < 0 {
					fmt.Print("ERROR: Invalid room coordinates")
					os.Exit(0)
				}
				newGame.Coordinates = append(newGame.Coordinates, []int{x, y})
			} else if tunnelMatched {
				newGame.InputConnections += dataCollection[i] + "\n"
				if connectionFound == false {
					newGame.Connections = make([][]int, len(newGame.RoomNames))
					connectionFound = true
				}
				temp1 := dataCollection[i][:strings.IndexByte(dataCollection[i], '-')]
				temp2 := dataCollection[i][strings.IndexByte(dataCollection[i], '-')+1:]
				temp1Index, temp2Index := -1, -1
				for i := 0; i < len(newGame.RoomNames); i++ {
					if temp1 == newGame.RoomNames[i] {
						temp1Index = i
						if temp2Index >= 0 {
							break
						}
					}
					if temp2 == newGame.RoomNames[i] {
						temp2Index = i
						if temp1Index >= 0 {
							break
						}
					}
				}
				if temp1Index == -1 || temp2Index == -1 {
					fmt.Println("ERROR: connection to room that does not exist")
					os.Exit(0)
				}
				newGame.Connections[temp1Index] = append(newGame.Connections[temp1Index], temp2Index)
				newGame.Connections[temp2Index] = append(newGame.Connections[temp2Index], temp1Index)
			}
		}

	}
	return newGame
}

func addFirstLevel(game *Game, storage *PathStorage) {
	for _, connection := range game.Connections[game.StartIndex] {
		var newPath []int
		newPath = append(newPath, game.StartIndex)
		newPath = append(newPath, connection)
		if hitEnd(newPath, game.EndIndex) {
			result := Paths{}
			result.Comb = append(result.Comb, newPath)
			result.CombFlow = 1
			result.CombLen = 1
			storage.Paths = append(storage.Paths, result)
			// game.ValidPaths = append(game.ValidPaths, newPath)
			game.PathsFound = true
		} else {
			game.PathQueue = append(game.PathQueue, newPath)
		}
	}
}

func findPaths(game *Game, storage *PathStorage) {
	if game.PathsFound {
		return
	}
	if len(game.PathQueue) == 0 {
		addFirstLevel(game, storage)
	}
	currCount := len(game.PathQueue)
	for _, path := range game.PathQueue {
		for _, room := range game.Connections[path[len(path)-1]] {
			if contains(path, room) {
				continue
			}
			newPath := append(path, room)
			test := make([]int, len(newPath))
			for i := 0; i < len(newPath); i++ {
				test[i] = newPath[i]
			}
			// check if path contains the room
			// if doesnt, make a new path with tunnel in the end
			if hitEnd(test, game.EndIndex) {
				game.ValidPaths = append(game.ValidPaths, test)
				findNonOverlapping(game, storage)
				// Stops the game - best result found
				if game.PathsFound == true {
					return
				}
				continue
			}
			game.PathQueue = append(game.PathQueue, test)
		}
	}
	// clean up Queue
	game.PathQueue = game.PathQueue[currCount:]
	// If empty, all of the paths are found
	if len(game.PathQueue) == 0 {
		game.PathsFound = true
	}
	findPaths(game, storage)
}

// Calculate path flow
func calcPathFlow(game *Game, path Paths) (res int) {
	var pathLength int
	for i := 0; i < len(path.Comb); i++ {
		pathLength = pathLength + len(path.Comb[i])
	}
	temp := float64(pathLength - len(path.Comb) + game.Ants)
	length := float64(len(path.Comb))
	res = int(math.Ceil(temp/length) - 1)
	return
}

// calculate the possible combinations with new path
func findNonOverlapping(game *Game, storage *PathStorage) {
	if len(game.ValidPaths) == 1 {
		var bestResult Paths
		bestResult.Comb = game.ValidPaths
		bestResult.CombFlow = calcPathFlow(game, bestResult)
		bestResult.CombLen = 1
		storage.Paths = append(storage.Paths, bestResult)
	} else {
		bestPath := storage.Paths[0]
		// first check if new path is the answer
		newPath := game.ValidPaths[len(game.ValidPaths)-1]
		// Best comb always first in array
		// Check if this will be better based on flow
		test := Paths{}
		test.Comb = append(bestPath.Comb, newPath)
		test.CombFlow = calcPathFlow(game, test)
		newBest := false
		// Estimated flow
		if test.CombFlow < bestPath.CombFlow {
			// fmt.Printf("Possibility: %v\n", test)
			newBest = true
			for i := 0; i < len(bestPath.Comb); i++ {
				if !noOverlap(bestPath.Comb[i], newPath) {
					newBest = false
					break
				}
			}
			// If true no overlap and shorter than best so far
			if newBest {
				// fmt.Printf("Old: %v\n", bestPath)
				// fmt.Printf("New best: %v\n", test)
				test.CombLen = len(test.Comb)
				storage.Paths = append([]Paths{test}, storage.Paths[:]...)
				if test.CombLen == game.MaxFlow {
					game.PathsFound = true
					return
				}
			}
		}
		// Go and find no overlap, add new ones to storage
		// loop over all other paths in storage
		loopLength := len(storage.Paths)
		for i := 1; i < loopLength; i++ {
			// Check if could be better based on flow
			test := Paths{}
			test.Comb = append(storage.Paths[i].Comb, newPath)
			test.CombFlow = calcPathFlow(game, test)
			// If could check overlap
			if test.CombFlow < storage.Paths[i].CombFlow {
				better := true
				for l := 0; l < len(storage.Paths[i].Comb); l++ {
					if !noOverlap(storage.Paths[i].Comb[l], newPath) {
						better = false
						break
					}
				}
				if better {
					// If no overlap add new comb in this place, shift old one down,skip one in the loop
					test.CombLen = len(test.Comb)
					// additiona check if this comb is better than current best
					if test.CombFlow < bestPath.CombFlow {
						storage.Paths = append([]Paths{test}, storage.Paths[:]...)
					} else {
						temp := storage.Paths[i]
						storage.Paths[i] = test
						storage.Paths = append(storage.Paths, temp)
						continue
					}
				}
			}
		}
		// Always add new path to the end of storage
		singlePath := Paths{}
		singlePath.Comb = append(singlePath.Comb, newPath)
		singlePath.CombFlow = calcPathFlow(game, singlePath)
		singlePath.CombLen = len(singlePath.Comb)
		storage.Paths = append(storage.Paths, singlePath)
	}
	if storage.Paths[0].CombLen == game.MaxFlow {
		game.PathsFound = true
	}
}

// Check if two given paths don't overlap
func noOverlap(pOne []int, pTwo []int) bool {
	for i := 1; i < len(pOne)-1; i++ {
		for l := 1; l < len(pTwo)-1; l++ {
			if pOne[i] == pTwo[l] {
				return false
			}
		}
	}
	return true
}

// Catch if path goes into loop
func contains(slice []int, i int) bool {
	for _, v := range slice {
		if v == i {
			return true
		}
	}
	return false
}

// Check if next room is the final one
func hitEnd(path []int, end int) bool {
	return path[len(path)-1] == end
}

// calculate num of paths in best case scenario
func calcMaxFlow(game *Game) {
	if len(game.Connections[game.StartIndex]) >= len(game.Connections[game.EndIndex]) {
		game.MaxFlow = len(game.Connections[game.EndIndex])
	} else {
		game.MaxFlow = len(game.Connections[game.StartIndex])
	}
}

// Validate input
func validInput(game *Game) bool {
	// Ants are more than 0
	if game.Ants == 0 {
		return false
	}
	// Start and end have at least 1 tunnel
	if len(game.Connections[game.StartIndex]) == 0 || len(game.Connections[game.EndIndex]) == 0 {
		return false
	}
	return true
}

// Temp function for printing out the right paths
func printResult(game *Game, storage *PathStorage) {
	result := storage.Paths[0].Comb
	for i := 0; i < len(result); i++ {
		fmt.Println(i)
		for l := 0; l < len(result[i]); l++ {
			fmt.Printf(game.RoomNames[result[i][l]] + " ")
		}
		fmt.Println()
	}
}

func calcAnts(game *Game, pathComb *Paths) {
	// Take the fastes combination put it into seperate structures
	var paths []FinalPath
	for l := 0; l < len(pathComb.Comb); l++ {
		path := FinalPath{Len: len(pathComb.Comb[l]) - 1, Path: pathComb.Comb[l]}
		paths = append(paths, path)
	}
	// Get ready result template
	var output Result
	output.ActiveRow = 0
	for l := 0; l < pathComb.CombFlow; l++ {
		row := Row{}
		output.Result = append(output.Result, row)
	}
	// Calculations for ants distribution
	currPath := 0
	var nextPath int
	for i := 1; i <= game.Ants; i++ {
		if i == 1 {
			paths[currPath].Ants++
			fillInOutput(i, &output, game, paths, paths[currPath])
			continue
		}
		currFlow := paths[currPath].Len + paths[currPath].Ants
		nextPath = currPath + 1
		if currPath == len(paths)-1 {
			nextPath = 0
		}
		nextFlow := paths[nextPath].Len + paths[nextPath].Ants
		if currFlow > nextFlow {
			paths[nextPath].Ants++
			fillInOutput(i, &output, game, paths, paths[nextPath])
			currPath = nextPath
		} else {
			paths[currPath].Ants++
			fillInOutput(i, &output, game, paths, paths[currPath])
		}

	}
	printInput(game)
	fmt.Println()
	for i := 0; i < len(output.Result); i++ {
		fmt.Println(output.Result[i].Output)
	}
}

func printInput(game *Game) {
	fmt.Println(game.Ants)
	for i := 0; i < len(game.RoomNames); i++ {
		switch i {
		case game.StartIndex:
			fmt.Println("##start")
		case game.EndIndex:
			fmt.Println("##end")
		}
		fmt.Printf("%v", game.RoomNames[i])
		for l := 0; l < len(game.Coordinates[i]); l++ {
			fmt.Printf(" %v", game.Coordinates[i][l])
		}
		fmt.Println()
	}
	// Print connections
	fmt.Printf("%v", game.InputConnections)
}

// Saving data for printing
func fillInOutput(ant int, output *Result, game *Game, paths []FinalPath, path FinalPath) {
	// Determine in which row the ant will start
	// Take the active row, check the paths taken, if none,
	//  than add current paths 1 index
	// Populate the output
	// Check if Row is saturated, than change active row to next
	if len(paths) == 1 && paths[0].Len == 1 {
		for i := 1; i < path.Len+1; i++ {
			s := fmt.Sprintf("L%v-%v ", ant, game.RoomNames[path.Path[i]])
			output.Result[output.ActiveRow+i-1].Output += s
		}
		return
	}
	if compare(output.Result[output.ActiveRow].PathsTaken, path.Path[1]) {
		// Compare true - we can put ant in this row
		output.Result[output.ActiveRow].PathsTaken = append(output.Result[output.ActiveRow].PathsTaken, path.Path[1])
		for i := 1; i < path.Len+1; i++ {
			s := fmt.Sprintf("L%v-%v ", ant, game.RoomNames[path.Path[i]])
			output.Result[output.ActiveRow+i-1].Output += s
		}
		// If current row full, change to next
		if len(output.Result[output.ActiveRow].PathsTaken) == len(paths) {
			output.ActiveRow++
		}
	} else {
		output.ActiveRow++
		output.Result[output.ActiveRow].PathsTaken = append(output.Result[output.ActiveRow].PathsTaken, path.Path[1])
		for i := 1; i < path.Len+1; i++ {
			s := fmt.Sprintf("L%v-%v ", ant, game.RoomNames[path.Path[i]])
			output.Result[output.ActiveRow+i-1].Output += s
		}
		// If current row full, change to next
		output.ActiveRow--
		if len(output.Result[output.ActiveRow].PathsTaken) == len(paths) {
			output.ActiveRow++
		}
	}
}

// Helper func
func compare(s []int, i int) bool {
	for l := 0; l < len(s); l++ {
		if s[l] == i {
			return false
		}
	}
	return true
}

func main() {
	game := ReadFile("examples/" + os.Args[1])
	var pathStorage PathStorage
	if !game.HasStart || !game.HasEnd {
		fmt.Println("ERROR: no start or end room found")
		os.Exit(0)
	}
	if !validInput(&game) {
		fmt.Println("ERROR: invalid data format")
		os.Exit(0)
	}
	calcMaxFlow(&game)
	findPaths(&game, &pathStorage)
	// printResult(&game, &pathStorage)
	if len(pathStorage.Paths) == 0 {
		fmt.Println("ERROR: no paths connecting start and end rooms")
		os.Exit(0)
	}
	calcAnts(&game, &pathStorage.Paths[0])
}

