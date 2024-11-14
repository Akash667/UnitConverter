package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

var measurementUnits map[string][]string = make(map[string][]string)
var conversionFactors map[string]map[string]float32 = make(map[string]map[string]float32)

type valueError struct {
	errorString string
}

func (e valueError) Error() string {
	return e.errorString
}

var myError valueError

func init() {
	measurementUnits["length"] = []string{"m", "ft", "cm"}
	measurementUnits["weight"] = []string{"kg", "g", "oz"}
	measurementUnits["temperature"] = []string{"C", "K", "F"}
	myError = valueError{"invalid value"}

	conversionFactors["length"] = map[string]float32{"m": 1.000, "ft": 3.281, "cm": 100.000}
	conversionFactors["weight"] = map[string]float32{"kg": 1.000, "g": 1000.000, "oz": 35.274}
	// conversionFactors["temperature"] = map[string]float32{"C":0.000, "K":273.150, "F"}

}

type UnitConverter struct {
	IsSubmitted    bool
	SelectedUnit   string
	UnitValue      string
	GivenUnit      string
	ConvertToUnit  string
	ConvertedValue string
	UnitsOfValue   []string
}

func main() {

	tmpl := template.Must(template.ParseFiles("form1.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		params := r.URL.Query()
		unit := params.Get("unitToUse")
		var UnitValues []string
		if unit == "" {
			unit = "length"
		}
		UnitValues = measurementUnits[unit]
		if r.Method != http.MethodPost {
			// fmt.Println(UnitConverter{SelectedUnit: unit, UnitsOfValue: UnitValues})
			currentState := UnitConverter{SelectedUnit: unit, UnitsOfValue: UnitValues}
			err := tmpl.Execute(w,
				currentState)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		newUnitvalue := UnitConverter{
			IsSubmitted:   true,
			SelectedUnit:  unit,
			UnitValue:     r.FormValue("MeasurementValue"),
			GivenUnit:     r.FormValue("isUnit"),
			ConvertToUnit: r.FormValue("toUnit"),
		}
		computedValue, err := newUnitvalue.convertUnits()
		if err != nil {
			errValue := UnitConverter{
				IsSubmitted:  true,
				SelectedUnit: unit,
				UnitValue:    "Invalid",
				GivenUnit:    "Invalid",
			}
			tmpl.Execute(w, errValue)
			return
		}
		newUnitvalue.ConvertedValue = strconv.FormatFloat(float64(computedValue), 'f', 2, 64)
		tmpl.Execute(w, newUnitvalue)

	})

	http.ListenAndServe(":8080", nil)

}

func (conver *UnitConverter) convertUnits() (float32, error) {

	valueNum, err := strconv.ParseFloat(conver.UnitValue, 32)
	unitType := conver.SelectedUnit
	isUnit := conver.GivenUnit
	toUnit := conver.ConvertToUnit
	if err != nil {
		fmt.Println("Invalid value given for conversion")
		return 0, err
	}

	if unitType == "length" || unitType == "weight" {
		return float32(valueNum) * (conversionFactors[unitType][toUnit]) / (conversionFactors[unitType][isUnit]), nil
	} else if unitType == "temperature" {
		if isUnit == "C" {

			if toUnit == "K" {
				return float32(valueNum) + 273.15, nil
			} else if toUnit == "F" {
				return float32(valueNum)*9/5 + 32, nil
			} else if toUnit == "C" {
				return float32(valueNum), nil
			}

		} else if isUnit == "F" {

			if toUnit == "C" {
				return (float32(valueNum) - 32) * 5 / 9, nil
			} else if toUnit == "K" {
				return (float32(valueNum)-32)*5/9 + 273.15, nil
			} else if toUnit == "F" {
				return float32(valueNum), nil
			}
		} else if isUnit == "K" {
			if toUnit == "C" {
				return float32(valueNum) - 273.15, nil
			} else if toUnit == "F" {
				return (float32(valueNum)-273.15)*9/5 + 32, nil
			} else if toUnit == "K" {
				return float32(valueNum), nil
			}
		}

	}

	return 0, myError

}
