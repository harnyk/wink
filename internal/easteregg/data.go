package easteregg

import "math/rand"

// 1. Checkin phrases
// 2. Checkout phrases
// 3. Checkin phrases (rude)
// 4. Checkout phrases (rude)

// CheckinPhrases is a list of phrases to use when checking in.
var CheckinPhrases = []string{
	"Good morning, friend! Time to get the day started!",
	"Hey there, ready to tackle the day ahead?",
	"Rise and shine, it's time to get to work and crush it!",
	"Another day, another opportunity to be our best selves at work!",
	"Time to hit the ground running and make the most of this day!",
	"Here we go again, let's make today a great one!",
	"Good morning, sunshine! It's time to shine at work!",
	"Ready to make the most of this day and be our best selves?",
	"Another day, another chance to be awesome at work!",
	"Let's make today a productive and successful one!",
}

// CheckoutPhrases is a list of phrases to use when checking out.
var CheckoutPhrases = []string{
	"Time to clock out and enjoy the rest of the day, buddy!",
	"See you tomorrow, have a great evening!",
	"Time to call it a day, see you bright and early tomorrow!",
	"Well done today, now it's time to relax and unwind!",
	"Another day down, time to put work aside and enjoy the rest of the evening!",
	"It's been a long day, time to rest and recharge for tomorrow!",
	"Time to say goodbye to work for today, have a great evening!",
	"See you tomorrow, have a good night!",
	"It's been a productive day, time to kick back and relax a bit!",
	"Time to say goodbye to work for now, have a great evening!",
}

// CheckinPhrasesRude is a list of rude phrases to use when checking in.
var CheckinPhrasesRude = []string{
	"Good morning, time to start another day of pretending to be sober!",
	"Ready to face another day with a pounding headache?",
	"Rise and shine, it's time to drag ourselves into work after a night of partying!",
	"Another day, another chance to try to hide our hangover at the office!",
	"Time to hit the ground running, or at least try to walk straight!",
	"Here we go again, another day of trying to look like we didn't drink too much last night!",
	"Let's make today a great one, or at least try to not look like we're still drunk from the night before!",
	"Good morning, sunshine! Time to start the day with a smile (or at least try to not look like we're about to puke)!",
	"Time to shine and make the most of this day, or at least try to stay awake!",
	"Another day, another opportunity to learn and grow (or at least try to not throw up at work)!",
}

// CheckoutPhrasesRude is a list of rude phrases to use when checking out.
var CheckoutPhrasesRude = []string{
	"Time to clock out and hit the bar!",
	"See you tomorrow, after I sleep off this hangover!",
	"Time to call it a day, and head to happy hour!",
	"Well done today, time to crack open a cold one!",
	"Another day down, time to get wasted!",
	"It's been a long day, time to get trashed!",
	"Time to say goodbye to work for today, and start drinking!",
	"See you tomorrow, after I sleep off this alcohol-induced coma!",
	"It's been a productive day, time to get hammered!",
	"Time to say goodbye to work for now, and start partying!",
}

func pickRandomItemFromList(list []string) string {
	return list[rand.Intn(len(list))]
}

// GetRandomCheckinPhrase returns a random checkin phrase.
// The rudeProbability is a number between 0 and 1 that determines
// the probability of returning a rude checkin phrase.
func GetRandomCheckinPhrase(rudeProbability float64) string {
	if rand.Float64() < rudeProbability {
		return pickRandomItemFromList(CheckinPhrasesRude)
	}
	return pickRandomItemFromList(CheckinPhrases)
}

// GetRandomCheckoutPhrase returns a random checkout phrase.
// The rudeProbability is a number between 0 and 1 that determines
// the probability of returning a rude checkout phrase.
func GetRandomCheckoutPhrase(rudeProbability float64) string {
	if rand.Float64() < rudeProbability {
		return pickRandomItemFromList(CheckoutPhrasesRude)
	}
	return pickRandomItemFromList(CheckoutPhrases)
}
