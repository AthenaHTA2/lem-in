package main

import (
	"fmt"
	"io/ioutil"
	"lem-in/examples"
	"os"
	"strconv"
)

///////////

// This is the graph data structure that lem-in uses
type GraphData struct {
	Rooms       map[string]*RoomData
	StartRoom   string
	EndRoom     string
	AntsTotalNb int               // Total # of ants
	GraphRooms  []string          // str room names from the 'example##.txt' files
	RoomTypes   map[string]string // map of room names and corresponding room types: start/end/intermediate
	RoomsNumber int               // total number of rooms in graph, used to assign int room names
	Edges       [][]string        // all tunnels in the graph expressed as room names
	Paths       [][]*RoomData     // Paths is a two dimensional slice of room pointers representing paths from start room to end room
}

// This is the structure containing room information
type RoomData struct {
	NameStr string
	// x, y        int     // room coordinates
	// Parent     *Room   // 'before' room
	Edges    []string    // slice of 'after' rooms
	PtrEdges []*RoomData // slice of pointers to after ooms
	Occupied bool        // used to marshal ants and to derive graph paths
	RoomType string      //"start", "end", "itermediate"
	// Distance   int     // distance from "start", for selecting shortest path
	// AntName    string  // 1 to #of ants, used for terminal output
	// AntsNumber int     // Total # of ants in room, used for 'end' room
}

var pathArrPrelim = []string{}

/*
AddRoom function adds a Room to the Graph struct. Function takes in
name of the room, x coordinates and y coordinates
and append it to the &Room
*/
func (g *GraphData) AddRoom(name string) {
	folder := "examples/"
	userInput := os.Args[1]
	file := folder + userInput
	var rEdges []string
	rMap := examples.RoomTypes()
	roomType := rMap[name]
	_, _, gEdges, _ := examples.RoomsEdges(file)

	for _, v := range gEdges {
		if v[0] == name {
			for j := 1; j < len(v); j++ {
				rEdges = append(rEdges, string(v[j]))
			}
		}
	}
	// fmt.Println(rEdges)
	// Now we add room and its corresponding data to the graph
	g.Rooms[name] = &RoomData{
		RoomType: roomType,
		NameStr:  name,
		// X:         x,
		// Y:         y,
		Occupied: false,
		Edges:    rEdges, // a string slice of adjacent room names
	}
	// fmt.Println(g.Rooms[name])
}

// This function populates the graph structure with remaining graph data
func (g *GraphData) PopulateGraphStruct(file string) *GraphData {
	// use functions from 'antData' file to get Graph struct's data
	var (
		startRoom   string
		endRoom     string
		roomTypes   map[string]string
		graphRooms  []string
		roomsNumber int
		edges       [][]string
		antsTotalNb int
		// rooms       map[string][]*RoomData
	)

	graphRooms, _, _, _ = examples.RoomsEdges(file)
	_, roomsNumber, _, _ = examples.RoomsEdges(file)
	_, _, edges, _ = examples.RoomsEdges(file)
	antsTotalNb, _ = examples.CheckNumAnts(file)
	roomTypes = examples.RoomTypes()
	startRoom = examples.FindStartEnd("start")
	endRoom = examples.FindStartEnd("end")

	graph := &GraphData{
		AntsTotalNb: antsTotalNb,
		GraphRooms:  graphRooms,
		RoomsNumber: roomsNumber,
		Edges:       edges,
		RoomTypes:   roomTypes,
		StartRoom:   startRoom,
		EndRoom:     endRoom,
	}
	return graph
}

// below function adds edges to the 'GraphData' struct. Note to self: But surely it adds links to the RoomData struct that contains 'PtrEdges = map[string]*RoomData'
func (Graph *GraphData) AddLinks(from, to string) {
	fromRoom := Graph.Rooms[from]
	toRoom := Graph.Rooms[to]
	fromRoom.PtrEdges = append(fromRoom.PtrEdges, toRoom)
}

/*
PrintGraphStruct is a graph function that
prints the Room and their links for visualization
*/
func (Graph *GraphData) PrintGraphStruct() {
	for _, v := range Graph.Rooms {
		fmt.Printf("\nRoom name %v : ", v.NameStr)
		for _, v := range v.PtrEdges {
			fmt.Printf(" %v ", v.NameStr)
		}
	}
	fmt.Println()
}

/*
NewGraph is a function that starts a new graph by creating an empty room.
*/
func NewGraph() *GraphData {
	return &GraphData{
		Rooms:       map[string]*RoomData{},
		StartRoom:   "",
		EndRoom:     "",
		AntsTotalNb: 0,
		GraphRooms:  []string{},          // str room names from the 'example##.txt' files
		RoomTypes:   map[string]string{}, // map of room names and corresponding room types: start/end/intermediate
		RoomsNumber: 0,                   // total number of rooms in graph, used to assign int room names
		Edges:       [][]string{},        // all tunnels in the graph expressed as room names
		Paths:       [][]*RoomData{},
	}
}

/*
Generate is a function that populates the new graph with data using the AddRoom, PopulateGraphStruct, and AddLinks functions
and AddLinks functions.
*/
func (g *GraphData) Generate(file string) *GraphData {
	// The code below loops through the GraphRooms string slice
	// Gets the room name, and adds it to the graph using the AddRoom function.

	var graphRooms []string
	graphRooms, _, _, _ = examples.RoomsEdges(file)
	// fmt.Print("here : ")
	// fmt.Println(graphRooms)
	for _, rm := range graphRooms {
		room := rm // room name
		// fmt.Println("room : " + room)
		g.AddRoom(room)
	}

	var (
		edges [][]string // The code below loops through the two-dimensional string slice 'Edges'
		from  string     // Gets the 0 index and subsequent indexes of each slice
		to    string     // Builds links for each room
	) // And adds the links to the graph using the AddLinks method
	_, _, edges, _ = examples.RoomsEdges(file) // After that, in 'func main()' will use the g.PopulateGraphStruct function
	// to add remaining graph data.
	for i, ed := range edges {
		howManyEdges := len(edges[i])
		switch {
		case howManyEdges == 2:
			from = ed[0]
			to = ed[1]
			// The code below makes the graph undirected
			g.AddLinks(from, to)
			g.AddLinks(to, from)

		case howManyEdges > 2:
			for j := 1; j < howManyEdges; j++ {
				from = ed[0]
				to = ed[j]
				// The code below makes the graph undirected
				g.AddLinks(from, to)
				g.AddLinks(to, from)
			}
		}
	}
	return g
}

/*IsColonized is a function of *GraphData that checks if the Room is Occupied.
It will return true if space is Occupied and false if not.
*/
func (g *GraphData) IsOccupied(name string) bool {
	return g.Rooms[name].Occupied
}

/*
Below is a variable Array, which is an slice of string.
Initialize a global slice to create a method into slice.
*/
type Array []string

// MakeOccupied is a *Graph struct function that make the path Occupied in the graph
func (g *GraphData) MakeOccupied(start, end string, path Array, make bool) {
	for _, name := range path {
		if start != name && end != name {
			g.Rooms[name].Occupied = make
		}
	}
}

// FindPath is a method of the *Graph struct that find paths from the start to end.
func (Graph *GraphData) FindPath(start, end string, path Array, swtch bool) []string {
	var newPath Array
	shortest := make(Array, 0)

	if _, exist := Graph.Rooms[start]; !exist {
		return path
	}

	path = append(path, start)
	if start == end {
		return path
	}

	for _, node := range Graph.Rooms[start].PtrEdges {
		/*
			the if statement is to check wheather the current node is occupied or not,
			and if the current path have the same room or not
		*/
		if !(Graph.IsOccupied(node.NameStr)) && !examples.IsValueInList(node.NameStr, path) {
			newPath = Graph.FindPath(node.NameStr, end, path, swtch) // recursion, calling 'FindPath' function again
			if len(newPath) > 0 {
				// if the swtch is true it will find the shortest path in the graph
				if swtch {
					if examples.IsValueInList(start, newPath) && examples.IsValueInList(end, newPath) {
						pathArrPrelim = append(pathArrPrelim, fmt.Sprint(newPath))
						if len(shortest) == 0 { // if newPath is the first path to be found
							shortest = newPath
						}
						if len(newPath) < len(shortest) {
							shortest = newPath
						}

					}
				}

				// if the switch is false it will return the first path it finds
				if !(swtch) {
					if examples.IsValueInList(start, newPath) && examples.IsValueInList(end, newPath) {
						return newPath
					}
				}

			}
		}
	}

	return shortest // this will be returned if the swtch is true
}

// GetPathsList is a method of the *Graph struct that return a 2-dimensional slice of paths
func (Graph *GraphData) GetPathsList(start, end string, swtch bool) [][]string {
	antsPaths := [][]string{} // container for the paths
	var p Array               // init p for the parameter of the ShortestPath method
	var path Array            // container for the shortest path
	cnt := 0
	c := 0

	// the for loop below will loop until cnt is not equal to the length of the list of edges for the start room
	for cnt != len(Graph.Rooms[start].PtrEdges) {
		path = Graph.FindPath(start, end, p, swtch) // look for the path
		pathArrPrelim = examples.SortPaths(pathArrPrelim)
		if len(pathArrPrelim) > 1 {
			path = examples.TurnintoArray(pathArrPrelim[len(pathArrPrelim)-1])
		}
		pathArrPrelim = []string{}
		if len(path) != 0 {
			if len(antsPaths) == 0 {
				antsPaths = append(antsPaths, path)
			} else if len(antsPaths[cnt-1]) == len(path) {
				// if the current path and the previous path have the same length
				// checks if the path is not similar to the previous paths
				// if it is not similar append the path into the AntsPaths
				for i := 0; i < len(path); i++ {
					if antsPaths[cnt-1][i] == path[i] {
						c++
					}
				}
				if c != len(path) {
					antsPaths = append(antsPaths, path)
				}
			} else {
				antsPaths = append(antsPaths, path)
			}
		}
		Graph.MakeOccupied(start, end, path, true) // make the paths Colonized
		cnt++
	}
	return antsPaths
}

/*
BestPath is a function that receives two slices of paths and compare which one is the
best path to use. This function will return the best paths to use in asceding order based on path length
*/
func BestPath(AntsPaths1, AntsPaths2 [][]string) [][]string {
	// this function will check which one had more paths and return it
	if len(AntsPaths1) > len(AntsPaths2) {
		return AntsPaths1
	} else if len(AntsPaths1) < len(AntsPaths2) {
		return examples.OrganisePaths(AntsPaths2)
	} else {
		/*
			If length of both slices of slice of path is equal to each other it will check
			which one had the less room inside the slices of slice of path and return it
		*/
		antp1 := 0
		antp2 := 0
		for _, paths := range AntsPaths1 {
			antp1 = antp1 + len(paths)
			antp2 = antp2 + len(paths)
		}
		if antp1 < antp2 {
			return examples.OrganisePaths(AntsPaths1)
		} else {
			return examples.OrganisePaths(AntsPaths2)
		}
	}
}

/*
AntMoves is function that will receive the name of the ant and the path it is
taking and return a slice that Consists of the ants movement. For example,
the ant name is 1 and the path is [A0 A1 A2 end] the result will be:
[L1-A0 L1-A1 L1-A2 L1-end]
*/
func AntMoves(nameOfAnt int, paths []string) []string {
	result := []string{}
	str := ""
	antName := strconv.Itoa(nameOfAnt)
	for _, room := range paths {
		str = "L" + antName + "-" + room
		result = append(result, str)
		str = ""
	}
	return result
}

/*
AntMovesLen is the function that cacluates the exact
number of line it will take to print the movements of the ants downward.
For example, the container [][][]string ConsistOf the slices below.
[[[L1-2 L1-3 L1-1] [L2-2 L2-3 L2-1] [L3-2 L3-3 L3-1] [L4-2 L4-3 L4-1]]]
The container ConsistOf one path so the maxlen will be len of the container[0] to get
the number of ants inside that path and the lenOfMove will be the len of the container[0][0]
to get the number of moves of each ants in that path.
*/
func AntMovesLen(container [][][]string) int {
	result := 0
	maxLen := 0
	pos := 0
	if len(container) > 1 {
		for i := range container {
			for j := range container {
				if len(container[i]) < len(container[j]) {
					maxLen = len(container[j])
					pos = j
				}
			}
		}
	} else {
		maxLen = len(container[0]) // assign to the first path if the number of path in the container is one
	}
	lenOfMove := len(container[pos][0])
	// the number of ants in the path - 1 + the len of ants movemnts in that path
	// will give you the exact amount of line to print it downwards
	result = (maxLen - 1) + (lenOfMove)
	return result
}

/*
PrintAntsMoves is the function that prints the ants movements
This will receive [][][]string that ConsistOf the movement of that ant and print it downwards
for example, the container [][][]string ConsistOf
[[L1-2 L1-3 L1-1] [L2-2 L2-3 L2-1] [L3-2 L3-3 L3-1] [L4-2 L4-3 L4-1]]]
it will be converted to the result below.
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3 L4-2
L3-1 L4-3
L4-1
*/
func PrintAntsMoves(container [][][]string, AntMovesLen int) string {
	result := ""
	antsMoves := make([][]string, AntMovesLen)
	for _, c := range container {
		for j, paths := range c {
			for k, p := range paths {
				antsMoves[j+k] = append(antsMoves[j+k], p)
			}
		}
	}
	for _, a := range antsMoves {
		for _, v := range a {
			result += v + " "
		}
		result += "\n"
	}
	return result
}

// PathsSeletion pick the correct amount of paths and the correct amount of ants for each paths
func PathsSeletion(nAnts int, pahts [][]string) [][][]string {
	container := make([][][]string, len(pahts))

	if len(pahts) > 1 {
		cnt := 0
		i := 1
		for i != nAnts+1 {
			if cnt == len(pahts)-1 {
				cnt = 0
			}
			x := len(pahts[cnt]) + len(container[cnt])
			y := len(pahts[cnt+1]) + len(container[cnt+1])
			if !(x > y) {
				container[cnt] = append(container[cnt], AntMoves(i, pahts[cnt][1:]))
			} else {
				if cnt == len(pahts)-1 {
					cnt = 0
					container[0] = append(container[0], AntMoves(i, pahts[0][1:]))
				} else {
					cnt++
					container[cnt] = append(container[cnt], AntMoves(i, pahts[cnt][1:]))
				}
			}
			i++
		}
	} else {
		i := 1
		for i != nAnts+1 {
			container[0] = append(container[0], AntMoves(i, pahts[0][1:]))
			i++
		}
	}
	return container
}

////////////////

func main() {
	filePath := "examples/"
	terminalInput := os.Args[1]
	file := filePath + terminalInput
	// fmt.Println(file)
	canProceed := true
	boo, err := examples.IgnoreBadCommands(file)
	if !boo && err != nil {
		fmt.Println(err)
		canProceed = false
		os.Exit(0)
	} else {
		// start := time.Now()

		_, err := examples.CheckNumAnts(file)
		if err != nil {
			fmt.Println(err)
			canProceed = false
		}
		err = examples.CheckStartEndRoomsExist(file)
		if err != nil {
			fmt.Println(err)
			canProceed = false
		}
		err = examples.CheckInfiniteLoopsExist(file)
		if err != nil {
			fmt.Println(err)
			canProceed = false
		}
		err = examples.CheckTunnels(file)
		if err != nil {
			fmt.Println(err)
			canProceed = false
		}
		err = examples.ChkRoomNamesCoord(file)
		if err != nil {
			fmt.Println(err)
			canProceed = false
		}
		err = examples.ChkDuplicateRooms(file)
		if err != nil {
			fmt.Println(err)
			canProceed = false
		}
		err = examples.ChkUnknownRooms(file)
		if err != nil {
			fmt.Println(err)
			canProceed = false
		}

		// elapsed := time.Since(start)
		// fmt.Printf("it took %v milliseconds to run\n", elapsed)
		if canProceed {
			lemin := NewGraph()                   // initialise lemin as an empty graph
			leminGraphMap := lemin.Generate(file) // Generate function populates the following lemin fields: 1) a slice of pointers to rooms; 2) a 2-d slice of pointers to graph edges
			// fmt.Println("Lemin content1:", lemin)
			lemin = lemin.PopulateGraphStruct(file) // PopulateGraphStruct function adds remaining data to the lemin graph
			lemin.Rooms = leminGraphMap.Rooms
			AntsPaths1 := lemin.GetPathsList(lemin.StartRoom, lemin.EndRoom, true)
			lemin.MakeOccupied(lemin.StartRoom, lemin.EndRoom, lemin.GraphRooms, false)
			AntsPaths2 := lemin.GetPathsList(lemin.StartRoom, lemin.EndRoom, false)
			a := BestPath(AntsPaths1, AntsPaths2)
			if len(a) == 0 {
				fmt.Println("ERROR: invalid data format - missing start to end path.")
				return
			}
			s, _ := os.Open(file)     // open the file
			f, _ := ioutil.ReadAll(s) // read the file
			fmt.Println(string(f))
			fmt.Println()
			container := PathsSeletion(lemin.AntsTotalNb, a)
			fmt.Print(PrintAntsMoves(container, AntMovesLen(container)))
			fmt.Println("$")
		}
	}
}
