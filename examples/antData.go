/* LEM-IN : OUTLINE OF INPUT CHECKS

Valid room format: 'name' 'coord_x' 'coord_y' (e.g.:"nameoftheroom 1 2", "4 6 7")
Valid links format: 'name1'-'name2' (e.g.: "1-2", "2-5")

Format Rules:

- Each room can only contain one ant at a time, except ##start and ##end.
- Ants can not walk over fellow ants.
- The program displays content of each text file, and each ant as it moves once through a tunnel to an empty room.
- The program never quits in an unexpected manner.
- The program ignores unknown commands.


List of checks:
1) Is 0 < #ants < 5000? 															-->DONE
2) Path between ##start and ##end exists?											-->DONE
3) Both ##start and ##end rooms exist?         										-->DONE
4) Are there any duplicate rooms?													-->DONE
5) Are there links to unknown rooms?												-->DONE
6) Are there rooms with invalid (missing or negative) coordinates?				  	-->DONE
7) Room coordinates are always int?												  	-->DONE
8) Room names never start with the letter "L" or with "#" and have no spaces?	  	-->DONE
9) Are there infinite loops, i.e. rooms that link to themselves?                  	-->DONE
10) Are there tunnels that join more than 2 rooms?								  	-->DONE

*/
package examples

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//----------------------------------------------------------------------
//This section contains error formulas that are invoked by other forumlas

// errorString is a trivial implementation of error
type errorString struct {
	s string
}

// By attaching the Error() method to errorString, the latter becomes also type error, in addition to type struct.
func (e *errorString) Error() string { // we have a pointer to errorString because method Error() works with an interface, which only takes a pointer.
	return e.s
}

// New returns an error that formats as the given text
func New(text string) error {
	message := "ERROR: invalid data format - "
	return &errorString{message}
}

/* Three more examples of 'error' implementation:

a)	func standardErrorMsg(text string) error {
		return &errorString{text}
	}

b)	f, err : sqrt(-1)
		if err != nil {
		fmt.Println(err)
		}

c)	Simply calling the errors package Error method:
		Error()
*/
//findMessage returns the relevant error message when input file contains errors

var (
	AntsTotalNb  int
	Edges_Prelim []string
	Edges        [][]string
	Rooms        []string
	Rooms2D      [][]string
	RoomsNumber  int
	length       int
	tempEdges    [][]string
)

func FindMessage(s string) string {
	var message string
	// Map of error information
	errorInformation := map[string]string{
		"badCommand":          "ERROR: invalid data format - unknown or missing terminal command.",
		"openFile":            "ERROR: invalid data format - unable to open file.",
		"ants":                "ERROR: invalid data format - missing or too small/ large number of ants.",
		"createFile":          "ERROR: invalid data format - unable to create outputFile.txt.",
		"missingStartEndRoom": "ERROR: invalid data format - missing start or end room.",
		"missingStartEndPath": "ERROR: invalid data format - missing start to end path.",
		"duplicateRoom":       "ERROR: invalid data format - duplicate room.",
		"unknownRoom":         "ERROR: invalid data format - unknown room.",
		"invalidCoordinates":  "ERROR: invalid data format - bad room coordinates.",
		"intCoordinates":      "ERROR: invalid data format - room coordinates must be of type int.",
		"roomName":            "ERROR: invalid data format - bad room name.",
		"infiniteLoop":        "ERROR: invalid data format - room that links to itself generates infinite loop.",
		"badTunnel":           "ERROR: invalid data format - tunnel with links to more than two rooms not allowed.",
		"badRmNmCoord":        "ERROR: invalid data format - bad room name or invalid room coordinates.",
	}

	for k, v := range errorInformation {
		if k == s {
			message = v
		}
	}
	return message
}

//-------------------------------------------------------------------------

// Checks that terminal commands and example text files are in correct format
func IgnoreBadCommands(file string) (bool, error) { // correct terminal command e.g.: go run . example00.txt

	if len(os.Args[:]) != 2 {
		return false, errors.New(FindMessage("badCommand"))
	}

	// map of example files:
	validExampleFiles := map[string]string{
		"examples/badexample00.txt": "badexample00.txt",
		"examples/badexample01.txt": "badexample01.txt",
		"examples/example00.txt":    "example00.txt",
		"examples/example01.txt":    "example01.txt",
		"examples/example02.txt":    "example02.txt",
		"examples/example03.txt":    "example03.txt",
		"examples/example04.txt":    "example04.txt",
		"examples/example05.txt":    "example05.txt",
		"examples/example06.txt":    "example06.txt",
		"examples/example07.txt":    "example07.txt",
	}
	if _, ok := validExampleFiles[file]; !ok {
		return false, errors.New(FindMessage("openFile"))
	}

	return true, nil
}

// Make output file and copy content of example file into output file. Not used for lem-in
func MakeOutputFile(outputFile, exampleFile string) (written int64, err error) { // exampleFile is one of the 10 example.txt files provided as inpout
	file, err := os.Open(os.Args[1])
	if err != nil {
		return 0, errors.New(FindMessage("openFile"))
	}

	defer file.Close()

	output, err := os.Create("outputFile.txt")
	if err != nil {
		return 0, errors.New(FindMessage("createFile"))
	}

	defer output.Close()

	return io.Copy(output, file)
	// when adding ants' moves, use: options := os.O_RDWR |os.O_APPEND | os.O_CREATE in: file, err := os.OpenFile("example00.txt", options, os.FileMode(0777))
}

// Check number of ants in file's first line
func CheckNumAnts(inputFile string) (int, error) { // input file is the specific example.txt file provided
	file, err := os.Open(inputFile)
	if err != nil {
		return 0, errors.New(FindMessage("openFile"))
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	row1, _, err := reader.ReadLine()

	if err == io.EOF {
		os.Exit(0)
	}
	antsNumber, err := strconv.Atoi(string(row1))
	if err != nil || antsNumber <= 0 || antsNumber > 5000 {
		return 0, errors.New(FindMessage("ants"))
	}
	AntsTotalNb = antsNumber
	return AntsTotalNb, nil
}

// nil error if start/end rooms exist
func CheckStartEndRoomsExist(inputFile string) error {
	var counter int
	file, err := os.Open(inputFile)
	if err != nil {
		return errors.New(FindMessage("openFile"))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "##start" || line == "##end" {
			counter++
		}
	}
	if counter < 2 {
		err = errors.New(FindMessage("missingStartEndRoom"))
	} else {
		err = nil
	}
	return err
}

func Cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

// returns 'nil' if there are no infinite loops
func CheckInfiniteLoopsExist(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return errors.New(FindMessage("openFile"))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		before, after, found := Cut(line, "-")
		if found && before == after { // infinite loop if tunnel links to same room
			return errors.New(FindMessage("infiniteLoop"))
		}
	}
	return nil
}

// returns nil if tunnels join not more than two rooms
func CheckTunnels(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return errors.New(FindMessage("openFile"))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		_, after, found := Cut(line, "-")
		if found {
			_, _, found := Cut(after, "-") // now check if there is a link to a third room
			if found {
				return errors.New(FindMessage("badTunnel"))
			}
		}
	}
	return nil
}

// Youtube on GO reading files : https://www.youtube.com/watch?v=X1DV9ZXinaU
// items := strings.Split(line," ") -->to e.g. check for correct room names or correct coordinates
// fmt.Printf("Name: %s %s email: %v",items[1],items[2],items[3])

// returns nil error if room name has no spaces, does not start with # or L, and has valid coordinates
func ChkRoomNamesCoord(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return errors.New(FindMessage("openFile"))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, " ")
		switch len(items) {
		case 1: // first line contains the number of ants
			continue
		case 2: // if only two items, then the y coordinate is missing
			return errors.New(FindMessage("invalidCoordinates"))
		case 3: // If 3 items, check that room names dont start with 'L' or '#'
			if string(items[0][0]) == "L" {
				return errors.New(FindMessage("roomName"))
			}
			if string(items[0]) != "##start" && string(items[0]) != "##end" {
				if string(items[0][0]) == "#" {
					return errors.New(FindMessage("roomName"))
				}
			}
		default: // if more than 3 items, then either there is a space in room name or there are three coordinates instead of two
			return errors.New(FindMessage("badRmNmCoord"))
		}

		x, _ := strconv.Atoi(items[1][0:])
		y, _ := strconv.Atoi(items[2][0:])

		if x < 0 || y < 0 { // negative coordinates are not admissible
			return errors.New(FindMessage("invalidCoordinates"))
		}
	}
	return nil // No error if room name and room coordinates are correct
}

// returns nil error if no duplicate room names or duplicate room coordinates are found
func ChkDuplicateRooms(inputFile string) error {
	var msg1 error
	var msg2 error
	duplicate_rNames := make(map[string]int)
	duplicate_rCoordinates := make(map[string]int)

	file, err := os.Open(inputFile)
	if err != nil {
		return errors.New(FindMessage("openFile"))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		before, after, found := Cut(line, " ")
		if found {
			_, ok := duplicate_rNames[before] // check if the room name exist in the 'duplicate_rNames' map
			if ok {
				duplicate_rNames[before] += 1 // increase counter by 1 if already in the map
			} else {
				duplicate_rNames[before] = 1 // else start counting from 1
			}
			_, exist := duplicate_rCoordinates[after] // check if room coordinates already exist in the 'duplicate_rCoordinates' map
			if exist {
				duplicate_rCoordinates[after] += 1 // increase counter by 1 if already in the map
			} else {
				duplicate_rCoordinates[after] = 1 // else start counting from 1
			}

		}
	}
	for _, v := range duplicate_rNames { // map keys with v > 1 are duplicate rooms
		if v > 1 {
			msg1 = errors.New(FindMessage("duplicateRoom"))
			return msg1
		}
	}
	for _, c := range duplicate_rCoordinates { // map keys with c > 1 are duplicate rooms
		if c > 1 {
			msg2 = errors.New(FindMessage("duplicateRoom"))
			return msg2
		}
	}
	return nil
}

// Returns true if a string is included in a slice of strings
func IsValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// returns nil error if tunnels don't link unknown rooms, i.e. rooms without coordinates.
// Also returns slice of room names, # of rooms, slice of edges, used for building antFarm
func ChkUnknownRooms(inputFile string) error {
	var tunnelNames []string
	file, err := os.Open(inputFile)
	if err != nil {
		return errors.New(FindMessage("openFile"))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		before, _, found := Cut(line, " ")
		if found {
			Rooms = append(Rooms, before) // rNames is a slice of room names whose coordinates are known
		}
		before, after, found := Cut(line, "-")
		if found {

			if !IsValueInList(before, tunnelNames) {
				tunnelNames = append(tunnelNames, before) // tunnelName is a slice of rooms linked by tunnels
			}
			if !IsValueInList(after, tunnelNames) {
				tunnelNames = append(tunnelNames, after) // tunnelName is a slice of rooms linked by tunnels
			}
		}
	}

	for _, v := range tunnelNames {
		if !IsValueInList(v, Rooms) {
			return errors.New(FindMessage("unknownRoom")) // a room is considered 'unknown' if it is not included in the 'Rooms' slice
		}
	}
	return nil
}

func CmpareStringSlices(a, b []string) bool {
	/*if len(a) != len(b) {
		return false
	}*/
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

/*
RoomsEdges is a function that returns a slice
of room names, the total number of rooms, a two-dimensional array of edges,
and an error if anything is amiss*/
func RoomsEdges(inputFile string) ([]string, int, [][]string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, 0, nil, errors.New(FindMessage("openFile"))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		before, _, found := Cut(line, " ")
		if found {
			if !IsValueInList(before, Rooms) {
				Rooms = append(Rooms, before) // rNames is a slice of rooms whose coordinates are known
			}
		}

		before, after, found := Cut(line, "-")
		if found {
			var temp []string
			temp = append(temp, before, after)
			tempEdges = append(tempEdges, temp)
			length = len(tempEdges) / 3 // Take one third of the tempEdges slice because it triplicates the edges.
		}
	}

	for i, v := range Rooms {
		var temp []string
		temp = append(temp, v)
		// Populate the Edges composite literal with room names
		if i == len(Rooms) {
			break // Doesn't work: it still appends each room three times!
		} else {
			Edges = append(Edges, temp)
		}

	}

	for i, v := range Rooms { // Solution: take one third of Edges composite literal
		for _, edg := range tempEdges[0:length] { // same here
			// fmt.Printf("v is: %v, edg is: %v\n", v, edg)
			// if isValueInList(edg[0], v) {
			if edg[0] == v { // if edge belongs to room name
				if !IsValueInList(edg[1], Edges[i]) { // check if edge has been added already
					Edges[i] = append(Edges[i], edg[1]) // if not, append edge
				}
			}
		}
	}
	var ResultEdges [][]string
	// remove Edges slices that contain room name only
	for i, v := range Edges {
		if len(v) > 1 {
			if i < len(Edges)-1 {
				ResultEdges = append(ResultEdges, v) // append slices that contain edges
			}
		}
	}
	RoomsNumber = len(Rooms)
	return Rooms, RoomsNumber, ResultEdges, nil
}

// RoomTypes function assigns either 'start', 'end', or 'intermediate' type to each room in the graph
func RoomTypes() map[string]string {
	RTypes := map[string]string{}
	inputFile := "examples/" + os.Args[1]
	TestFile := inputFile
	temp := []string{}
	s, _ := os.Open(TestFile) // open the file
	defer s.Close()
	f, _ := ioutil.ReadAll(s)
	data := strings.Split(string(f), "\n")

	for i := 1; i < len(data); i++ {
		if data[i] == "##start" {
			temp = strings.Fields(data[i+1])
			RTypes[temp[0]] = "start"
			i++
		} else if data[i] == "##end" {
			temp = strings.Fields(data[i+1])
			RTypes[temp[0]] = "end"
			i++
		} else {
			temp = strings.Fields(data[i])
			if len(temp) == 3 {
				RTypes[temp[0]] = "intermediate"
			}

		}
	}

	return RTypes
}

// FindStartEnd function returns the actual room name for the start or the end rooms in the graph
func FindStartEnd(startend string) string {
	theMap := RoomTypes()
	value := ""
	for rm := range theMap {
		// fmt.Print("se is:\t", se)
		// fmt.Println("\trm is:", rm)
		if theMap[rm] == startend {
			value = rm
		}
	}
	return value
}

/*
SortPaths is a function that sort the slice inside the slice.
This function will return the ordered paths. The algorithm that is used in
this function is bubble sort.
*/
func SortPaths(path []string) []string {
	for i := 0; i < len(path)-1; i++ {
		// if path[i] is less than path[i+1] swap them
		if len(path[i]) < len(path[i+1]) {
			path[i], path[i+1] = path[i+1], path[i]
		}
	}
	return path
}

/*
OrganisePaths in function that sorts the slice inside the slice based on their length.
This function will return the ordered paths. Once again, the algoriothms that is used in this
function is bubble sort.
*/
func OrganisePaths(path [][]string) [][]string {
	for i := range path {
		for j := range path {
			// if path[i] is less than path[j] swap them
			if len(path[i]) < len(path[j]) {
				path[i], path[j] = path[j], path[i]
			}
		}
	}
	return path
}

// Turns a string into an array of strings by using white space delimiter
func TurnintoArray(s string) []string {
	arr := strings.Fields(s[1 : len(s)-1])
	return arr
}
