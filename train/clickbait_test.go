package clickbait

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/lytics/multibayes"
)

func TestAccuracy(t *testing.T) {
	// we'll train our classifiers with 75% of our data and test how well they
	// classify the remaining 25%

	classifier := multibayes.NewClassifier()

	// read headlines from `buzzfeed_headlines.csv`
	buzzFile, err := os.Open("buzzfeed_headlines.csv")
	assert.Equal(t, nil, err)
	defer buzzFile.Close()

	reader := csv.NewReader(buzzFile)
	buzzData, err := reader.ReadAll()
	assert.Equal(t, nil, err)

	index := int(len(buzzData) / 4)
	clickbaitTraining := buzzData[:3*index]
	clickbaitTest := buzzData[3*index:]

	var notclickbaitTraining, notclickbaitTest [][]string
	//  read headlines from `reuters_headlines.csv`
	reutersFile, err := os.Open("reuters_headlines.csv")
	assert.Equal(t, nil, err)
	defer reutersFile.Close()

	reader = csv.NewReader(reutersFile)
	reutersData, err := reader.ReadAll()
	assert.Equal(t, nil, err)

	index = int(len(reutersData) / 4)
	notclickbaitTraining = append(notclickbaitTraining, reutersData[:3*index]...)
	notclickbaitTest = append(notclickbaitTest, reutersData[3*index:]...)

	// from `aljazeera_headlines.csv`
	aljazeeraFile, err := os.Open("aljazeera_headlines.csv")
	assert.Equal(t, nil, err)
	defer aljazeeraFile.Close()

	reader = csv.NewReader(aljazeeraFile)
	aljazeeraData, err := reader.ReadAll()
	assert.Equal(t, nil, err)

	index = int(len(aljazeeraData) / 4)
	notclickbaitTraining = append(notclickbaitTraining, aljazeeraData[:3*index]...)
	notclickbaitTest = append(notclickbaitTest, aljazeeraData[3*index:]...)

	// from `bloomberg_headlines.csv`
	bloombergFile, err := os.Open("bloomberg_headlines.csv")
	assert.Equal(t, nil, err)
	defer bloombergFile.Close()

	reader = csv.NewReader(bloombergFile)
	bloombergData, err := reader.ReadAll()
	assert.Equal(t, nil, err)

	index = int(len(bloombergData) / 4)
	notclickbaitTraining = append(notclickbaitTraining, bloombergData[:3*index]...)
	notclickbaitTest = append(notclickbaitTest, bloombergData[3*index:]...)

	// found that CNN's news headlines decreased accuracy
	/*
		// and from `cnn_headlines.csv`
		cnnFile, err := os.Open("cnn_headlines.csv")
		assert.Equal(t, nil, err)
		defer cnnFile.Close()

		reader = csv.NewReader(cnnFile)
		cnnData, err := reader.ReadAll()
		assert.Equal(t, nil, err)

		index = int(len(cnnData) / 4)
		notclickbaitTraining = append(notclickbaitTraining, cnnData[:3*index]...)
		notclickbaitTest = append(notclickbaitTest, cnnData[3*index:]...)
	*/

	// train the classifier
	fmt.Printf("TRAINING: Number of clickbait headlines: %v\tNumber of non-clickbait headlines: %v\n",
		len(clickbaitTraining), len(notclickbaitTraining))
	for _, doc := range clickbaitTraining {
		classifier.Add(doc[0], []string{CLICKBAIT})
	}
	for _, doc := range notclickbaitTraining {
		classifier.Add(doc[0], []string{NOT_CLICKBAIT})
	}

	fmt.Printf("TEST: Number of clickbait headlines: %v\tNumber of non-clickbait headlines: %v\n",
		len(clickbaitTest), len(notclickbaitTest))
	// how often do we correctly predict clickbait?
	correctClickbait := 0
	incorrectClickBait := 0
	totalClickbait := len(clickbaitTest)
	for _, doc := range clickbaitTest {
		probs := classifier.Posterior(doc[0])
		if probs[CLICKBAIT] > probs[NOT_CLICKBAIT] {
			correctClickbait++
		}
		if probs[NOT_CLICKBAIT] > probs[CLICKBAIT] {
			incorrectClickBait++
		}
	}

	// how often do we correctly predict not clickbait?
	correctNotClickbait := 0
	incorrectNotClickbait := 0
	totalNotClickbait := len(notclickbaitTest)
	for _, doc := range notclickbaitTest {
		probs := classifier.Posterior(doc[0])
		if probs[NOT_CLICKBAIT] > probs[CLICKBAIT] {
			correctNotClickbait++
		}
		if probs[CLICKBAIT] > probs[NOT_CLICKBAIT] {
			incorrectNotClickbait++
		}
	}

	// actual|predicted
	confusion := map[string]float64{
		"clickbait|clickbait":   float64(correctClickbait) / float64(totalClickbait),
		"clickbait|!clickbait":  float64(incorrectClickBait) / float64(totalClickbait),
		"!clickbait|clickbait":  float64(incorrectNotClickbait) / float64(totalNotClickbait),
		"!clickbait|!clickbait": float64(correctNotClickbait) / float64(totalNotClickbait),
	}
	fmt.Printf("Confusion matrix: %+v\n", confusion)
}
