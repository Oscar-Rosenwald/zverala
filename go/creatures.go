// -*- eval: (hs-minor-mode 1); -*-
package main

type Creature struct {
	name CreatureName
	days int
}

type CreatureName string

var Creatures []CreatureName = []CreatureName{
	"Chiméra",
	"Fénix",
	"Toch Amogaši",
	"Sfinga",
	"Vlkodlak",
	"Zlovlk",
	"Jednorožec",
	"Griffin",
	"Lví želva",
	"Kraken",
	"Kyklop",
	"Syréna",
	"Yeti",
	"Nessie",
	"Vyjící chluporyba",
	"Olifant",
	"Ždiboň",
	"Ent",
	"Labuť",
	"Kerberos",
	"Bazilišek",
	"Akromantule",
	"Goa'uld",
	"Vetřelec",
	"Létající bizon",
	"Pegas",
	"Mothra",
	"Sleipnir",
	"Velká A'tuin",
	"Horus",
	"Cthulhu",
	"Hydra",
	"Balrog",
	"Odgru Jahad",
}

type Dragon Creature

var Dragons = []Dragon{
	{days: 1, name: "Drak ohně"},
	{days: 1, name: "Drak země"},
	{days: 1, name: "Drak života"},
	{days: 1, name: "Drak vody"},
	{days: 1, name: "Drak dřeva"},
	{days: 1, name: "Drak smrti"},
	{days: 1, name: "Drak vzduchu"},
	{days: 1, name: "Drak chaosu"},
}

var DragonsAfterCreatureIndex = []int{0, 2, 5, 8, 14, 18, 23, 29}

var (
	Chimera       = Creatures[0]
	NUM_CREATURES = len(Creatures)
	NUM_DRAGONS   = len(Dragons)
)
