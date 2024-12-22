package main

import (
	"calc/enums"
	"calc/pkg"
	"calc/service"
	"encoding/json"
	"log"
	"net/http"
)

func startServer() {
	http.HandleFunc("/api/v1/calculate", calcHandler)

	log.Println("Server started")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Error starting server:", err)
	}
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req pkg.CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding JSON body:", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(pkg.CalcResponse{Error: enums.ErrInternalServerError})
		return
	}
	defer r.Body.Close()

	expression := req.Expression
	result, err := service.Calc(expression)

	if err != nil {
		log.Println("Error while calculating:", err.Error())
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(pkg.CalcResponse{Error: enums.ErrUnprocessableEntity})
		return
	} else {
		log.Println("Calculated successfully:", result)
	}

	resp := pkg.CalcResponse{Result: result}
	jsonResponse, err := json.Marshal(resp)

	if err != nil {
		log.Println("Error marshaling JSON:", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(pkg.CalcResponse{Error: enums.ErrInternalServerError})
		return
	}

	w.WriteHeader(200)
	w.Write(jsonResponse)
}

func main() {
	startServer()
}
