package easteregg

import (
	"math/rand"
	"time"
)

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

	"Embrace the day with enthusiasm and let's make some progress!",
	"New day, new opportunities! Let's dive into work!",
	"Let's unlock our potential and make today extraordinary!",
	"A new day to make a difference! Let's get started!",
	"Let's turn our plans into action and make today count!",
	"Time to embrace challenges and create some remarkable work!",
	"Fuel up with positivity and let's conquer today's tasks!",
	"Let's bring our A-game and make today amazingly productive!",
	"Today holds endless possibilities. Let's explore them!",
	"Gather your energy and enthusiasm—it's time to shine at work!",
	"Time to weave creativity and effort into the fabric of our day!",
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

	"Great job today! Now go enjoy some well-deserved rest!",
	"You've earned a relaxing evening after a day of hard work!",
	"Clocking out for now, but ready to conquer more tomorrow!",
	"Mission accomplished for today! Enjoy your free time!",
	"Time to switch off work mode and embrace relaxation!",
	"You've conquered today's challenges—now enjoy the evening!",
	"Celebrate the day's achievements and get ready for a new tomorrow!",
	"Job well done! Now go rejuvenate and unwind!",
	"Leaving the battlefield of work—time for some peace and quiet!",
	"Close the work chapter for today and open the book of relaxation!",
	"Unplug from work and dive into an evening of joy and rest!",
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

	"Brace yourself, work is coming, whether we like it or not!",
	"Survived the morning? Now survive the rest of the day!",
	"Let's stumble through another day of chaotic work!",
	"Back to the grind, where chaos and deadlines reign!",
	"Another day to fake enthusiasm and suffer through tasks!",
	"Time to slave away in the relentless cycle of work!",
	"Brace yourself for a storm of stress and impossible demands!",
	"Drown in emails, meetings, and unending tasks—welcome back!",
	"Prepare to wrestle with the beast of workload and frustration!",
	"Time to get lost in the labyrinth of never-ending chores!",
	"Embrace the torture of monotony and unachievable expectations!",
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

	"Escaped work for today! Ready to dive into some mischief?",
	"Freedom at last! Time to unleash the inner party animal!",
	"Work's over! Let's go cause some delightful chaos!",
	"Survived the ordeal! Time to drown the sorrows!",
	"Freed from the chains of work—time to wreak havoc!",
	"Work's tyranny is over! Time for some rebellious fun!",
	"Escape the prison of tasks and dive into lawless leisure!",
	"Unleash the wild side and let the night's anarchy begin!",
	"Celebrate surviving another brutal day of toil and trouble!",
	"Burst out of work's confinement and plunge into chaotic delight!",
	"Endure no more! The realms of relaxation and chaos await!",
}

func pickRandomItemFromList(list []string) string {
	return list[rand.Intn(len(list))]
}

// Seed the random number generator with current time.
func Seed() {
	rand.Seed(time.Now().UnixNano())
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
