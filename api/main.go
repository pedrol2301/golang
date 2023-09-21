package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const apiurl = "https://pokeapi.co/api/v2/"

// Criar TimeOut para a requisição
var testClient = http.Client{
	Timeout: time.Second * 5,
}

type PokemonList struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []PokemonFirstInfo `json:"results"`
}

type PokemonFirstInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	ID   int
}

type Pokemon struct {
	Abilities              []AbilityDetail `json:"abilities"`
	BaseExperience         int             `json:"base_experience"`
	Forms                  []PokemonForm   `json:"forms"`
	GameIndices            []GameIndex     `json:"game_indices"`
	Height                 int             `json:"height"`
	HeldItems              []interface{}   `json:"held_items"`
	ID                     int             `json:"id"`
	IsDefault              bool            `json:"is_default"`
	LocationAreaEncounters string          `json:"location_area_encounters"`
	Moves                  []MoveDetail    `json:"moves"`
	Name                   string          `json:"name"`
	Order                  int             `json:"order"`
	PastTypes              []interface{}   `json:"past_types"`
	Species                PokemonSpecies  `json:"species"`
	Sprites                PokemonSprites  `json:"sprites"`
	Stats                  []StatDetail    `json:"stats"`
	Types                  []TypeDetail    `json:"types"`
	Weight                 int             `json:"weight"`
	MainSprite             map[string]string
}

type AbilityDetail struct {
	Ability  Ability `json:"ability"`
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
}

type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonForm struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type GameIndex struct {
	GameIndex int         `json:"game_index"`
	Version   VersionName `json:"version"`
}

type VersionName struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type MoveDetail struct {
	Move                Move                 `json:"move"`
	VersionGroupDetails []VersionGroupDetail `json:"version_group_details"`
}

type Move struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionGroupDetail struct {
	LevelLearnedAt  int              `json:"level_learned_at"`
	MoveLearnMethod MoveLearnMethod  `json:"move_learn_method"`
	VersionGroup    VersionGroupName `json:"version_group"`
}

type MoveLearnMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionGroupName struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonSpecies struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonSprites struct {
	BackDefault      string                  `json:"back_default"`
	BackFemale       interface{}             `json:"back_female"`
	BackShiny        interface{}             `json:"back_shiny"`
	BackShinyFemale  interface{}             `json:"back_shiny_female"`
	FrontDefault     string                  `json:"front_default"`
	FrontFemale      interface{}             `json:"front_female"`
	FrontShiny       interface{}             `json:"front_shiny"`
	FrontShinyFemale interface{}             `json:"front_shiny_female"`
	Other            map[string]SpriteDetail `json:"other"`
	Versions         map[string]SpriteDetail `json:"versions"`
}

type SpriteDetail struct {
	FrontDefault     string `json:"front_default"`
	FrontFemale      string `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale string `json:"front_shiny_female"`
}

type StatDetail struct {
	BaseStat int      `json:"base_stat"`
	Effort   int      `json:"effort"`
	Stat     StatName `json:"stat"`
}

type StatName struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type TypeDetail struct {
	Slot int      `json:"slot"`
	Type TypeName `json:"type"`
}

type TypeName struct {
	Name  string `json:"name"`
	Color string
	URL   string `json:"url"`
}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func treatForPresentation(j Pokemon) Pokemon {

	for stat := range j.Stats {
		j.Stats[stat].Stat.Name = capitalizeFirstLetter(j.Stats[stat].Stat.Name)
	}
	for tp := range j.Types {
		j.Types[tp].Type.Name = capitalizeFirstLetter(j.Types[tp].Type.Name)
		switch j.Types[tp].Type.Name {
		case "Bug":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-success"
		case "Dark":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-black"
		case "Dragon":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-info"
		case "Electric":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-warning"
		case "Fairy":
			j.Types[tp].Type.Color = "has-text-black-bis has-background-danger-light"
		case "Fighting":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-danger"
		case "Fire":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-danger"
		case "Flying":
			j.Types[tp].Type.Color = "has-text-black-bis has-background-link-light"
		case "Ghost":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-link-dark"
		case "Grass":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-success"
		case "Ground":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-warning-dark"
		case "Ice":
			j.Types[tp].Type.Color = "has-text-black-bis has-background-link-light"
		case "Normal":
			j.Types[tp].Type.Color = "has-text-black-bis"
		case "Poison":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-primary-dark"
		case "Psychic":
			j.Types[tp].Type.Color = "has-text-black-bis has-background-danger-light"
		case "Rock":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-warning-dark"
		case "Steel":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-link-dark"
		case "Water":
			j.Types[tp].Type.Color = "has-text-white-bis has-background-info"

		default:
			j.Types[tp].Type.Color = "has-background-success"
		}
	}

	j.Name = capitalizeFirstLetter(j.Name)
	j.MainSprite = make(map[string]string)
	j.MainSprite["Front"] = j.Sprites.FrontDefault
	j.MainSprite["Back"] = j.Sprites.BackDefault

	return j
}

func main() {

	r := gin.Default()

	r.Static("/img", "../frontend/assets/")

	r.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		req, errreq := http.NewRequest(http.MethodGet, apiurl+"pokemon?limit=152&offset=0", nil)
		if errreq != nil {
			http.Error(c.Writer, "Error when starting communication: "+errreq.Error(), http.StatusBadRequest)
			return
		}

		res, errres := testClient.Do(req)
		if errres != nil {
			http.Error(c.Writer, "Error when starting communication: "+errreq.Error(), http.StatusBadRequest)
			return
		}

		var responseData PokemonList

		err := json.NewDecoder(res.Body).Decode(&responseData)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		for p := range responseData.Results {
			responseData.Results[p].Name = capitalizeFirstLetter(responseData.Results[p].Name)
			urlParts := strings.Split(responseData.Results[p].URL, "/")
			idStr := urlParts[len(urlParts)-2] // O ID é a penúltima parte da URL
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Erro ao converter ID:", err)
			}
			responseData.Results[p].ID = id
		}

		// c.Writer.Header().Set("Content-Type", "application/json")
		// err = json.NewEncoder(w).Encode(responseData)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return nil, err
		// }

		tmpl := template.Must(template.ParseFiles("../frontend/index.html"))
		tmpl.Execute(c.Writer, responseData)
	})

	r.GET("/pokemon/:id", func(c *gin.Context) {
		pokeId := c.Params.ByName("id")
		req, errreq := http.NewRequest(http.MethodGet, apiurl+"pokemon/"+pokeId+"/", nil)
		if errreq != nil {
			http.Error(c.Writer, "Error when starting communication: "+errreq.Error(), http.StatusBadRequest)
			return
		}

		res, errres := testClient.Do(req)
		if errres != nil {
			http.Error(c.Writer, "Error when starting communication: "+errreq.Error(), http.StatusBadRequest)
			return
		}
		var responseData Pokemon
		err := json.NewDecoder(res.Body).Decode(&responseData)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		responseData = treatForPresentation(responseData)
		tmpl, err := template.ParseFiles("../frontend/components/pokemon.html")
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		// fmt.Print(responseData)

		// Execute o modelo com os dados e escreva a resposta
		err = tmpl.Execute(c.Writer, responseData)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	r.GET("/hi", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.Run("0.0.0.0:3000")
}
