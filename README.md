# clickbait-classifier
Example of using a [naive Bayesian classifier](https://github.com/lytics/multibayes) to classify how clickbait-y article headlines. Sample clickbait headlines have been scraped from Buzzfeed. Sample non-clickbait headlines have been scraped from Reuters, Aljazeera, CNN and Bloomberg.

### Example

```{go}
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/lytics/multibayes"
)

func main() {
	// make a new classifier
	classifier := multibayes.NewClassifier()

	// read headlines from `buzzfeed_headlines.csv`
	buzzFile, err := os.Open("./train/buzzfeed_headlines.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer buzzFile.Close()
	reader := csv.NewReader(buzzFile)
	buzzData, err := reader.ReadAll()

	//  read headlines from `reuters_headlines.csv`
	reutersFile, err := os.Open("./train/reuters_headlines.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer reutersFile.Close()
	reader = csv.NewReader(reutersFile)
	reutersData, err := reader.ReadAll()

	// train the classifier
	for _, doc := range buzzData {
		classifier.Add(doc[0], []string{"clickbait"})
	}
	for _, doc := range reutersData {
		classifier.Add(doc[0], []string{"not_clickbait"})
	}

	// predict new classes
	probs := classifier.Posterior("50 ways to win big")
	fmt.Printf("Posterior Probabilities: %+v\n", probs)
	// Posterior Probabilities: map[clickbait:0.9649489234536172 not_clickbait:0.03505107654638282]

	// predict new classes
	probs = classifier.Posterior("Pope lands in US")
	fmt.Printf("Posterior Probabilities: %+v\n", probs)
	// Posterior Probabilities: map[clickbait:0.19909505931578056 not_clickbait:0.8009049406842195]
}
```
