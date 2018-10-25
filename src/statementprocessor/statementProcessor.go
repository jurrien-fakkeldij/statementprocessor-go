package main

import (
	"bufio"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/shopspring/decimal"
)

//StatementState enumeration
type StatementState int

const (
	//CORRECT State no issue with the Statement
	CORRECT StatementState = 0
	//BADENDBALANCE State there is something wrong with the end balance
	BADENDBALANCE StatementState = 1
	//DUPLICATEREFERENCE State there is another Statement with the same reference
	DUPLICATEREFERENCE StatementState = 2
)

func (state StatementState) String() string {
	// declare an array of strings
	names := [...]string{
		"Correct statement",
		"Bad endbalance in statement",
		"Duplicate reference found"}

	//If it does not exist, return Unknown.
	if state < CORRECT || state > DUPLICATEREFERENCE {
		return "Unknown"
	}
	// return the name of a Weekday
	// constant from the names array
	// above.
	return names[state]
}

// XMLStatements structure to read from the xml file.
type XMLStatements struct {
	XMLName    xml.Name       `xml:"records"`
	Statements []XMLStatement `xml:"record"`
}

// XMLStatement structure to read from the xml file.
type XMLStatement struct {
	XMLName       xml.Name `xml:"record"`
	Reference     string   `xml:"reference,attr"`
	AccountNumber string   `xml:"accountNumber"`
	Description   string   `xml:"description"`
	StartBalance  string   `xml:"startBalance"`
	Mutation      string   `xml:"mutation"`
	EndBalance    string   `xml:"endBalance"`
}

// Statement structure to read from the csv file and to output to a csv file.
type Statement struct {
	Reference     int
	AccountNumber string
	Description   string
	StartBalance  decimal.Decimal
	Mutation      decimal.Decimal
	EndBalance    decimal.Decimal
	State         StatementState
}

/*
 * Statement Proceesser application for the Luminis Assignment.
**/
func main() {
	csvFile, error := os.Open("../records.csv")
	checkError("Unable to open file", error, true)

	xmlFile, error := os.Open("../records.xml")
	checkError("Unable to open file", error, true)

	// defer the closing of our xmlFile so that we can parse it later on, so it gets closed at the end of the function.
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Statements array
	var xmlStatements XMLStatements

	// unmarshal byteArray which contains our xml string content into xmlStatements for later parsing.
	xml.Unmarshal(byteValue, &xmlStatements)

	var statements []Statement

	//Parse all the xml statements and add them to the Statements array.
	for i := 0; i < len(xmlStatements.Statements); i++ {
		statements = parseAndAppendStatement(xmlStatements.Statements[i].Reference,
			xmlStatements.Statements[i].AccountNumber,
			xmlStatements.Statements[i].Description,
			xmlStatements.Statements[i].StartBalance,
			xmlStatements.Statements[i].Mutation,
			xmlStatements.Statements[i].EndBalance,
			statements)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	//Parse all the csv statements and add them to the Statements array.
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			checkError("Error reading csv file", error, true)
		} else {
			statements = parseAndAppendStatement(line[0], line[1], line[2], line[3], line[4], line[5], statements)
		}
	}

	//Create output file.
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err, true)
	//wait with closing until at the end of the function.
	defer file.Close()

	writer := csv.NewWriter(file)
	//wait with closing until at the end of the function.
	defer writer.Flush()

	//Write header
	err = writer.Write([]string{"Reference", "Description", "State"})
	checkError("Cannot write to file", err, true)

	//Go through all the statements and write them to the file if the state is not CORRECT.
	for i := 0; i < len(statements); i++ {
		if statements[i].State != CORRECT {
			fmt.Println(statements[i])
			err := writer.Write([]string{strconv.Itoa(statements[i].Reference), statements[i].Description, statements[i].State.String()})
			checkError("Cannot write to file", err, true)
		}
	}
}

//Function to check the error messages and log a fatal when it is fatal.
func checkError(message string, err error, fatal bool) {
	if err != nil {
		if fatal {
			log.Fatal(message, err)
		} else {
			log.Println(message, err)
		}
	}
}

//Function to parse all the incoming strings (from both the csv file and xml file) in the same manner and add them to the same array.
func parseAndAppendStatement(reference string, accountNumber string, description string, startBalance string, mutation string, endBalance string, statements []Statement) []Statement {
	if reference == "Reference" {
		//First line in  csv file, don't parse.
		return statements
	}
	//Convert from string to int.
	referenceID, error := strconv.Atoi(reference)
	checkError("Unable to parse reference id", error, false)

	//convert from string to decimal (float gave floating point errors in the balance match later)
	startBalanceVal, error := decimal.NewFromString(startBalance)
	checkError("Unable to parse start_balance", error, false)

	//convert from string to decimal (float gave floating point errors in the balance match later)
	mutationVal, error := decimal.NewFromString(mutation)
	checkError("Unable to parse mutation", error, false)

	//convert from string to decimal (float gave floating point errors in the balance match later)
	endBalanceVal, error := decimal.NewFromString(endBalance)
	checkError("Unable to parse end_balance", error, false)

	if error == nil {
		state := CORRECT
		//Check if the reference is already parsed before. If so, set both to duplicate reference state.
		for i := 0; i < len(statements); i++ {
			if statements[i].Reference == referenceID {
				statements[i].State = DUPLICATEREFERENCE
				state = DUPLICATEREFERENCE
			}
		}

		//Check if the end balance matches the addition of the start balance and mutation.
		if !startBalanceVal.Add(mutationVal).Equal(endBalanceVal) {
			state = BADENDBALANCE
		}

		//Everything is correctly parsed and the correct state is set. Add to statements array.
		statements = append(statements, Statement{
			Reference:     referenceID,
			AccountNumber: accountNumber,
			Description:   description,
			StartBalance:  startBalanceVal,
			Mutation:      mutationVal,
			EndBalance:    endBalanceVal,
			State:         state,
		})
	}
	return statements
}
