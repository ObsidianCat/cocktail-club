package common

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"runtime"

	//"cocktail-club/collection"
	"fmt"
	"github.com/oleiade/reflections"
	"io/ioutil"
	"net/http"
	"strconv"
)

type apiSearchResultsByName struct {
	Cocktails []APICocktail `json:"drinks"`
}
type apiSearchResultsByIngredient struct {
	Cocktails []CocktailPreview `json:"drinks"`
}

// ProxyRequest proxies request to external cocktail API
func ProxyRequest(cocktailsAPIUrl string) []byte {
	res, _ := http.Get(cocktailsAPIUrl)
	var bytes []byte
	bytes, _ = ioutil.ReadAll(res.Body)
	return bytes
}

func stringToInt(st string) int {
	result, _ := strconv.ParseInt(st, 10, 0)
	return int(result)
}

func intToString(num int) string {
	result := strconv.Itoa(num)
	return result
}

func apiCtailToCtail(input APICocktail) Cocktail {
	var ingedientsList []string
	for i := 1; i <= 5; i++ {
		value, _ := reflections.GetField(input, "StrIngredient"+intToString(i))
		if value != "" {
			ingedientsList = append(ingedientsList, value.(string))
		}
	}
	cTail := Cocktail{
		Name:        input.Name,
		Ingredients: ingedientsList,
		ID:          stringToInt(input.ID),
		Preparation: input.Preparations,
		Image:       input.Image,
	}
	return cTail
}

// TransformAPIBytesToCtails convert external API data to cocktails
func TransformAPIBytesToCtails(cocktailBytes []byte) []Cocktail {
	var cocktailsListFromAPI apiSearchResultsByName
	var cocktailsResults []Cocktail

	err := json.Unmarshal(cocktailBytes, &cocktailsListFromAPI)
	if err != nil {
		fmt.Println("error:", err)
	}

	for _, apiCTail := range cocktailsListFromAPI.Cocktails {
		cocktailsResults = append(cocktailsResults, apiCtailToCtail(apiCTail))
	}
	return cocktailsResults
}

// TransformAPIBytesToCtailPreview convert external API data to cocktail previews
func TransformAPIBytesToCtailPreview(cocktailBytes []byte) []CocktailPreview {
	var cocktailsListFromAPI apiSearchResultsByIngredient
	//var cocktailsResults []CocktailPreview

	err := json.Unmarshal(cocktailBytes, &cocktailsListFromAPI)
	if err != nil {
		fmt.Println("error:", err)
	}

	return cocktailsListFromAPI.Cocktails
}

// ReadDataFileWithPathFromCallerFile allow to create path to file start from file which call the function instead of project root
func ReadDataFileWithPathFromCallerFile(fileName string) ([]byte, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil, errors.New("cannot read file")
	}

	prefixPath := filepath.Dir(file)
	path := prefixPath + "/" + fileName

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
		return nil, errors.New("cannot read file")
	}

	return data, nil
}
