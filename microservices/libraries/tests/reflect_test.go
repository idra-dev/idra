package tests

import (
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"microservices/libraries/reflection"
	"os"
	"reflect"
	"testing"
)

func TestReflect1(t *testing.T) {
	var v cdc_shared.Connector
	fieldTags := reflection.GetFieldTags(v)

	for fieldName, tags := range fieldTags {
		fmt.Printf("%s:\n", fieldName)
		for tagName, tagValue := range tags {
			fmt.Printf("%s: %s\n", tagName, tagValue)
		}
	}

	jsonResult, err := json.Marshal(fieldTags)
	if err != nil {
		fmt.Println("Errore nella conversione in JSON:", err)
		return
	}

	fmt.Println("Risultato in formato JSON:")
	fmt.Println(string(jsonResult))
	fmt.Println(reflection.GetColumnName(&v, &v.ConnectorName))
}

func TestValidate(*testing.T) {
	var v cdc_shared.Connector
	res := reflection.ValidateStruct(v)
	println(res)
}

func TestValidateJson(*testing.T) {
	var v cdc_shared.Connector
	jsonString, err := ReadFile("/Users/antonioradesca/Code/cdc-microservices/microservices/libraries/tests/connector_sample_correct.json")
	if err != nil {
		fmt.Println("Errore nella lettura del file:", err)
	} else {
		fmt.Println("Contenuto del file:", jsonString)
	}

	// Decodifica la stringa JSON nella struttura Persona
	if err := json.Unmarshal([]byte(jsonString), &v); err != nil {
		fmt.Println("Errore nella decodifica JSON:", err)
		return
	}
	res := reflection.ValidateStruct(v)
	println(res)
}

func TestValidateInvalidJson(*testing.T) {
	var v cdc_shared.Connector
	jsonString, _ := ReadFile("/Users/antonioradesca/Code/cdc-microservices/microservices/libraries/tests/connector_sample_not_correct.json")

	// Decodifica la stringa JSON nella struttura Persona
	if err := json.Unmarshal([]byte(jsonString), &v); err != nil {
		fmt.Println("Errore nella decodifica JSON:", err)
		return
	}
	res := reflection.ValidateStruct(v)
	println(res)
}

func ReadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Leggi il contenuto del file
	content := ""
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		content += string(buffer[:n])
	}

	return content, nil
}

func TestCreateInstance(*testing.T) {
	exampleStruct := data.KafkaConnector{}
	structType := reflect.TypeOf(exampleStruct)
	shapeType := reflect.TypeOf((*cdc_shared.ConnectorProvider)(nil)).Elem()

	// Verificare se la struttura Circle implementa l'interfaccia Shape
	if structType.Implements(shapeType) {
		fmt.Println("La struttura implementa l'interfaccia.")
	} else {
		fmt.Println("La struttura non implementa l'interfaccia.")
	}
	sample := data.RetrieveProvider("KafkaConnector")
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		// Otteniamo il tag JSON del campo
		jsonTag := field.Tag.Get("json")
		fmt.Printf("Campo: %s, Tag JSON: %s\n", field.Name, jsonTag)
	}
	fmt.Println(sample.Name())

	structName := "KafkaConnector" // Puoi impostare il nome della struct come stringa

	// Mappa dei nomi delle struct ai tipi delle struct
	structTypes := map[string]reflect.Type{
		"KafkaConnector": reflect.TypeOf(data.KafkaConnector{}),
		"Immudb":         reflect.TypeOf(data.ImmudbIdraDriver{}),
	}

	// Ottenere il tipo della struct dal nome
	structType, ok := structTypes[structName]
	if !ok {
		fmt.Println("Struct non trovata")
		return
	}

	// Creare un'istanza vuota della struct
	instance := reflect.New(structType).Elem()
	// Ottenere il valore dell'istanza
	value := instance.Interface()
	fmt.Println(value)
	fmt.Println(value.(cdc_shared.ConnectorProvider).Name())
}
