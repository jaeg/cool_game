package lore

import "math/rand"

func RandomGoblinName() (name string) {
	firstNames := []string{"Garl", "Gup", "Herr", "Mur", "Fat", "Belch", "Barf", "Brig", "Alan", "One eye", "Old", "Broken", "Gulp", "Snot", "Snee"}
	lastNames := []string{"Dur", "Twig", "Smellytoes", "Crook'ed", "The Unwashed", "Bile"}

	name = firstNames[getRandom(0, len(firstNames))] + " " + lastNames[getRandom(0, len(lastNames))]
	return
}

func getRandom(low int, high int) int {
	return (rand.Intn((high - low))) + low
}
