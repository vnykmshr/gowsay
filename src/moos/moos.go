package moos

import (
	"math/rand"
	"time"
)

var (
	moos []string
	rnd  = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func init() {
	moos = []string{
		"Holstein You Close",
		"I Love You Dairy Much",
		"I'll Never Love An-Udder",
		"Friends For-Heifer",
		"Be Mine For-Heifer",
		"Somebody Brand Moo",
		"Udderly Cool",
		"When cows bungee jump milk comes out their nose",
		"Simply MOO-velous!",
		"Cow's it going?",
		"I'm head over hooves for you!",
		"Party til the Cows Come Home!",
		"You are MOOsic to my ears!",
		"If you were a cow you'd be a belt by now.",
		"Hay! That's Moovelous",
		"Hay! Have a Moovin' Day!",
		"Moo, Baby.",
		"I'm MOO'ved by your generosity",
		"Thank You...You've MOO'd my day",
		"Deja Moo! ",
		"Deja Moo – the feeling you’ve heard this bull before",
		"Don’t mean to be a moo-sance",
		"I hate these “moo’d” swings",
		"Hope you’ll soon be feeling mooey bien again",
		"Home is where the herd is",
		"Holy Cow!",
		"Good moos travels fast",
		"Good MOOrning!",
		"Friends Never Steer You Wrong!",
		"Kitchen closed 'cause I'm not in the moooo-d to cook!",
		"Let’s COWmunicate",
		"Let’s go to the MOO-vies",
		"Manure happens!",
		"Me and my _udder_ -half",
		"Milkin’ It For All It’s Worth",
		"Moo Kids on the block",
		"Moosical chairs",
		"Overworked and _Udder_ paid!",
		"Precious Mooments",
		"The grass is always greener on the udder side",
		"Udderly Amazing!",
		"Udderly Adorable!",
		"Wake up in a happy mooooood!",
	}
}

func GetRandomMoo() string {
	return moos[rnd.Intn(len(moos))]
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
