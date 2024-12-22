// -*- eval: (hs-minor-mode 1); -*-
package utils

type Creature struct {
	Name CreatureName
	Days int
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
	{Days: 1, Name: "Drak ohně"},
	{Days: 1, Name: "Drak země"},
	{Days: 1, Name: "Drak života"},
	{Days: 1, Name: "Drak vody"},
	{Days: 1, Name: "Drak dřeva"},
	{Days: 1, Name: "Drak smrti"},
	{Days: 1, Name: "Drak vzduchu"},
	{Days: 1, Name: "Drak chaosu"},
}

var DragonsAfterCreatureIndex = []int{0, 2, 5, 8, 14, 18, 23, 29}

var (
	Chimera       = Creatures[0]
	NUM_CREATURES = len(Creatures)
	NUM_DRAGONS   = len(Dragons)
)
