package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"search-api/dtos"
)

const solrURL = "http://localhost:8983/solr/hotels_core"

// AddHotel se encarga de agregar un hotel a Solr
func AddHotel(hotelDTO dtos.HotelDTO) error {
	// Preparar el hotel en formato JSON
	doc := map[string]interface{}{
		"name":      hotelDTO.Name,
		"location":  hotelDTO.Location,
		"amenities": hotelDTO.Amenities,
		"rating":    hotelDTO.Rating,
	}

	// Convertir a JSON
	data, err := json.Marshal([]map[string]interface{}{doc})
	if err != nil {
		return fmt.Errorf("error marshaling hotel: %v", err)
	}

	// Hacer la solicitud a Solr para indexar el documento
	req, err := http.NewRequest("POST", solrURL+"/update?commit=true", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request to Solr: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Solr: %v", err)
	}
	defer resp.Body.Close()

	// Verificar si la respuesta es exitosa
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error indexing hotel, status code: %d", resp.StatusCode)
	}

	return nil
}

// SearchHotels busca hoteles en Solr por ciudad o nombre
func SearchHotels(query string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/select?q=name:%s OR location:%s&wt=json", solrURL, query, query)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al hacer la solicitud a Solr: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error en la respuesta de Solr, código: %d", resp.StatusCode)
	}

	// Leer respuesta de Solr
	var result struct {
		Response struct {
			Docs []map[string]interface{} `json:"docs"`
		} `json:"response"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error al parsear la respuesta de Solr: %v", err)
	}

	return result.Response.Docs, nil
}
