package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"github.com/MaxHalford/gago"
	"encoding/csv"
	"os"
	"bufio"
	"io"
	"time"
	"strconv"
	"math"
)

var (
	corpus = strings.Split(RandStringRunes(6908*65), "")
	target = strings.Split(RandStringRunes(6908*65), "")
)
var curIndex = 0
var dataArray []parking

// Strings is a slice of strings.
type Strings []string

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("10")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Evaluate a Strings slice by counting the number of mismatches between itself
// and the target string.
func (X Strings) Evaluate() (mismatches float64, err error) {
	for j, val := range dataArray[curIndex].occupings {
		col := 0
		for _, s := range X[j*dataArray[curIndex].lots : (j+1)*dataArray[curIndex].lots] {
			if s == "1" {
				col++
			}
		}
		mismatches += math.Abs(float64(col - val.places))
	}
	for i := 0; i < len(dataArray[curIndex].occupings)-2; i++ {
		for j := 0; j < dataArray[curIndex].lots; j++ {
			if X[i*dataArray[curIndex].lots+j] != X[(i+1)*dataArray[curIndex].lots+j] && X[(i+2)*dataArray[curIndex].lots+j] != X[(i+1)*dataArray[curIndex].lots+j] {
				mismatches += 0.3
			}
		}
	}

	return
}

// Mutate a Strings slice by replacing it's elements by random characters
// contained in  a corpus.
func (X Strings) Mutate(rng *rand.Rand) {
	/*for i:=0;i<len(dataArray[0].occupings)-2;i++{
		for j:=0;j<dataArray[0].lots;j++ {
			if X[i*dataArray[0].lots+j] != X[(i+1)*dataArray[0].lots+j] && X[(i+2)*dataArray[0].lots+j] != X[(i+1)*dataArray[0].lots+j]{
				X[i*dataArray[0].lots+j]=X[(i+1)*dataArray[0].lots+j]
			}
		}
	}*/
	//gago.MutUniformString(X, target, 6908, rng)
	toReverse := true
	j := rng.Intn(len(dataArray[curIndex].occupings) - 1)
	val := dataArray[curIndex].occupings[j]
	displace := 0
	if dataArray[curIndex].lots > val.places {
		displace = rng.Intn(dataArray[curIndex].lots - val.places)
	}
	for i := 0; i < dataArray[curIndex].lots; i++ {
		toReverse = rng.Intn(20) != 3
		if (toReverse) && (i >= displace) && (i < val.places+displace) {
			X[j*dataArray[curIndex].lots+i] = "1"
		} else {
			X[j*dataArray[curIndex].lots+i] = "0"
		}
	}
}

// Crossover a Strings slice with another by applying 2-point crossover.
func (X Strings) Crossover(Y gago.Genome, rng *rand.Rand) {
	j := rng.Intn(len(dataArray[curIndex].occupings) - 1)
	for i := 0; i < dataArray[curIndex].lots; i++ {
		X[j*dataArray[curIndex].lots+i], Y.(Strings)[j*dataArray[curIndex].lots+i] = Y.(Strings)[j*dataArray[curIndex].lots+i], X[j*dataArray[curIndex].lots+i]
	}
	//gago.CrossGNXString(X, Y.(Strings), 5, rng)
}

// MakeStrings creates random Strings slices by picking random characters from a
// corpus.
func MakeStrings(rng *rand.Rand) gago.Genome {
	for j, val := range dataArray[curIndex].occupings {
		displace := 0
		if dataArray[curIndex].lots > val.places {
			displace = rng.Intn(dataArray[curIndex].lots - val.places)
		}
		for i := 0; i < dataArray[curIndex].lots; i++ {
			if i >= displace && i < val.places+displace {
				corpus[j*dataArray[curIndex].lots+i] = "1"
			} else {
				corpus[j*dataArray[curIndex].lots+i] = "0"
			}
		}
	}
	var XX = make(Strings, len(corpus))
	copy(XX, corpus)
	return XX
	//return Strings(gago.InitUnifString(len(target), corpus, rng))
}

// Clone a Strings slice..
func (X Strings) Clone() gago.Genome {
	var XX = make(Strings, len(X))
	copy(XX, X)
	return XX
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

type occuping struct {
	moment time.Time
	places int
}
type parking struct {
	name      string
	lots      int
	occupings []occuping
}

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	/*file, err := os.Open("data")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}*/
	dataArray = []parking{}
	f, _ := os.Open("data")

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		// Display record.
		// ... Display record length.
		// ... Display all individual elements of the slice.
		t1, e := time.Parse(time.RFC3339, record[1])
		check(e)
		places, e := strconv.Atoi(record[0])
		check(e)
		lots, e := strconv.Atoi(record[3])
		check(e)
		oc := occuping{
			moment: t1,
			places: places,
		}
		pos := -1
		for i, val := range dataArray {
			if val.name == record[4] {
				pos = i
				break
			}
		}
		if pos == -1 {
			pos = len(dataArray)
			dataArray = append(dataArray, parking{name: record[4], occupings: []occuping{}, lots: lots})
		}
		if oc.places > dataArray[pos].lots {
			oc.places = dataArray[pos].lots
		}
		dataArray[pos].occupings = append(dataArray[pos].occupings, oc)
		//fmt.Println(oc)
	}
	fmt.Println(len(dataArray[0].occupings))
	fmt.Println(dataArray[0].lots)
	//fmt.Println(dataArray)
	corpus = strings.Split(RandStringRunes(len(dataArray[0].occupings)*dataArray[0].lots), "")
	for curIndex=0;curIndex<len(dataArray);curIndex++ {
		corpus = strings.Split(RandStringRunes(len(dataArray[curIndex].occupings)*dataArray[curIndex].lots), "")
		var ga= gago.Generational(MakeStrings)
		ga.PopSize = 15
		ga.ParallelEval = true
		ga.NBest = 3
		ga.Initialize()

		for i := 1; i < 10000; i++ {
			ga.Evolve()
			// Concatenate the elements from the best individual and display the result
			var buffer bytes.Buffer
			/*for _, letter := range ga.HallOfFame[0].Genome.(Strings) {
			buffer.WriteString(letter)
		}*/
			fmt.Printf("Result %s %d -> %s (%.0f mismatches)\n", dataArray[curIndex].name, i, buffer.String(), ga.HallOfFame[0].Fitness)
		}
		var buffer bytes.Buffer
		curnum:=1
		currow:=0
		for _, letter := range ga.HallOfFame[0].Genome.(Strings) {
			if curnum==1{
				if currow<len(dataArray[curIndex].occupings) {
					buffer.WriteString(dataArray[curIndex].occupings[currow].moment.Format("2006-01-02 15:04:05") + "\t")
				}else{

				}
			}
			buffer.WriteString(letter+"\t")
			curnum++
			if curnum==dataArray[curIndex].lots+1 {
				curnum=1
				currow++
				buffer.WriteString("\n")
			}
		}
		WriteStringToFile("results/"+dataArray[curIndex].name+".xls",buffer.String())
		fmt.Printf("Result %s -> (%.0f mismatches)\n", dataArray[curIndex].name, ga.HallOfFame[0].Fitness)
	}
}
