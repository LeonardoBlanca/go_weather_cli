package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"time"
)

// Mapeando o JSON para o nosso struct
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		ForecastDay []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?q=-20.538477139654926%2C%20-47.40334825706073&days=1&key=8041e41ace5141f6be831419250606&lang=pt")
	if err != nil {
		// Se houver algum erro, ele encerra a aplicação a printa o erro.
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	//	Vamos ler o corpo da resposta
	body, err := io.ReadAll(res.Body)
	// Vamos verificar por erro igual antes
	if err != nil {
		// Se houver algum erro, ele encerra a aplicação a printa o erro.
		panic(err)
	}

	// Colocando o struct na variável
	var weather Weather
	// Usando o pacote JSON (Passamos o body e o endereço da variável)
	// Vai converter o nosso body no tipo da variável abaixo
	err = json.Unmarshal(body, &weather)
	if err != nil {
		// Se houver algum erro, ele encerra a aplicação a printa o erro.
		panic(err)
	}

	// Agora vamos imprimir a nossa variável (já aplicada o JSON nela)
	// fmt.Println(weather)

	// Vamos extrair variáveis para facilitar trabalhar
	location, current, hours := weather.Location, weather.Current, weather.Forecast.ForecastDay[0].Hour

	// Vamos substituir a string abaixo pelos valores das variáveis
	// fmt.Printf("Franca, Brazil, tempC, condition")
	fmt.Printf("%s, %s: %.0fC, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	// Vamos fazer um loop para imprimir a hora
	for _, hour := range hours {
		// Transformando o timeEpoch em dia
		date := time.Unix(hour.TimeEpoch, 0)

		// Adicionando uma condição para mostrar apenas as horas futuras
		if date.Before(time.Now()) {
			// Pular o loop adicionando continue
			continue
		}

		// Assim como no anterior, vamos alterar o valor da string
		// fmt.Printf("hour - temp, rain, condition")

		// adicionar na variável para adicionar cor
		message := fmt.Sprintf(
			"%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"), // O 15:04 é o formato da data. Converteu automático hh para 15 e mm para 04
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			// Se não for chover, printa de amarelo
			color.Green(message)
		} else {
			// Aqui adicionamos o pacote das cores https://github.com/fatih/color
			// No terminal digite: go get github.com/fatih/color
			// Vamos aplicar o red quando a probabilidade for maior que 40%
			color.Red(message)
		}
	}
}
