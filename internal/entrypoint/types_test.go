package entrypoint

import "github.com/atye/wikitable2json/pkg/client"

var (
	GoldenMatrix = [][][]string{
		{
			{"Column 1", "Column 2", "Column 3"},
			{"A", "B", "B"},
			{"A", "C", "D"},
			{"E", "F", "F"},
			{"G", "F", "F"},
			{"H", "H", "H"},
		},
	}

	GoldenMatrixDouble = [][][]string{
		{
			{"Column 1", "Column 2", "Column 3"},
			{"A", "B", "B"},
			{"A", "C", "D"},
			{"E", "F", "F"},
			{"G", "F", "F"},
			{"H", "H", "H"},
		},
	}

	GoldenKeyValue = [][]map[string]string{
		{
			{
				"Column 1": "A",
				"Column 2": "B",
				"Column 3": "B",
			},
			{
				"Column 1": "A",
				"Column 2": "C",
				"Column 3": "D",
			},
			{
				"Column 1": "E",
				"Column 2": "F",
				"Column 3": "F",
			},
			{
				"Column 1": "G",
				"Column 2": "F",
				"Column 3": "F",
			},
			{
				"Column 1": "H",
				"Column 2": "H",
				"Column 3": "H",
			},
		},
	}

	IssueOneMatrix = [][][]string{
		{
			{"Jeju", "South Korea", "official, in Jeju Island"},
			{"Jeju"},
		},
	}

	DataSortValueMatrix [][][]string = [][][]string{
		{
			{"Abu Dhabi, United Arab Emirates", "N/A"},
		},
	}

	Issue34Matrix = [][][]string{
		{
			{"18 May 2019 election", "", "", "", "test"},
		},
	}

	Issue56Matrix = [][][]string{
		{
			{"test", "2035", "2035"},
		},
	}

	Issue105Matrix = [][][]string{
		{
			{"test0\ntest1", "test2"},
		},
	}

	Issue105MatrixVerbose = [][][]client.Verbose{
		{
			[]client.Verbose{
				{
					Text: "test0\ntest1",
					Links: []client.Link{
						{
							Href: "./test1",
							Text: "test1",
						},
					},
				},
				{
					Text: "test2",
				},
			},
		},
	}

	ReferenceMatrix = [][][]string{
		{
			{"Roy Morgan"},
		},
	}

	ComplexMatrix = [][][]string{
		{
			{"Date", "Brand", "Interview mode", "Sample size", "Primary vote", "Primary vote", "Primary vote", "Primary vote", "Primary vote", "Primary vote", "UND", "2pp vote", "2pp vote"},
			{"Date", "Brand", "Interview mode", "Sample size", "L/NP", "ALP", "GRN", "ONP", "UAP", "OTH", "UND", "L/NP", "ALP"},
			{"18–24 April 2022", "Roy Morgan", "Telephone/online", "1393", "35.5%", "35%", "12%", "4.5%", "1.5%", "11.5%", "–", "45.5%", "54.5%"},
			{"20–23 April 2022", "Newspoll-YouGov", "Online", "1538", "36%", "37%", "11%", "3%", "4%", "9%", "–", "47%", "53%"},
			{"20–23 April 2022", "Ipsos", "Telephone/online", "2302", "32%", "34%", "12%", "4%", "3%", "8%", "8%", "45%", "55%"},
		},
	}

	Issue77Matrix = [][][]string{
		{
			{"Lage", "Objekt", "Beschreibung", "Akten-Nr.", "Bild"},
			{"FeuchtHauptstraße 37(Standort)", "Ehemaliges Wirtschaftsgebäude", "Zweigeschossiger Satteldachbau mit Fachwerkobergeschoss, bezeichnet mit „1697“", "ehemals D-5-74-123-14 zugehörig", "weitere Bilder"},
		},
	}

	Issue77MatrixVerbose = [][][]client.Verbose{
		{
			{
				{
					Text: "Lage",
				},
				{
					Text: "Objekt",
				},
				{
					Text: "Beschreibung",
				},
				{
					Text: "Akten-Nr.",
				},
				{
					Text: "Bild",
				},
			},
			{
				{
					Text: "FeuchtHauptstraße 37(Standort)",
					Links: []client.Link{
						{
							Text: "Standort",
							Href: "https://geohack.toolforge.org/geohack.php?pagename=Liste_der_Baudenkm%C3%A4ler_in_Feucht&language=de&params=49.37546_N_11.21422_E_region:DE-BY_type:building&title=Feucht%2C+Hauptstra%C3%9Fe+37%2C+Ehemaliges+Wirtschaftsgeb%C3%A4ude",
						},
					},
				},
				{
					Text: "Ehemaliges Wirtschaftsgebäude",
				},
				{
					Text: "Zweigeschossiger Satteldachbau mit Fachwerkobergeschoss, bezeichnet mit „1697“",
				},
				{
					Text: "ehemals D-5-74-123-14 zugehörig",
				},
				{
					Text: "weitere Bilder",
					Links: []client.Link{
						{
							Href: "./Datei:2018_Feucht_Hauptstraße_37_02.jpg",
						},
						{
							Text: "weitere Bilder",
							Href: "https://commons.wikimedia.org/wiki/Category:Hauptstraße%2037%20(Ehemaliges%20Wirtschaftsgebäude,%20D-5-74-123-14)",
						},
					},
				},
			},
		},
	}

	Issue77KeyValueVerbose = [][]map[string]client.Verbose{
		{
			{
				"Lage": {
					Text: "FeuchtHauptstraße 37(Standort)",
					Links: []client.Link{
						{
							Text: "Standort",
							Href: "https://geohack.toolforge.org/geohack.php?pagename=Liste_der_Baudenkm%C3%A4ler_in_Feucht&language=de&params=49.37546_N_11.21422_E_region:DE-BY_type:building&title=Feucht%2C+Hauptstra%C3%9Fe+37%2C+Ehemaliges+Wirtschaftsgeb%C3%A4ude",
						},
					},
				},
				"Objekt": {
					Text: "Ehemaliges Wirtschaftsgebäude",
				},
				"Beschreibung": {
					Text: "Zweigeschossiger Satteldachbau mit Fachwerkobergeschoss, bezeichnet mit „1697“",
				},
				"Akten-Nr.": {
					Text: "ehemals D-5-74-123-14 zugehörig",
				},
				"Bild": {
					Text: "weitere Bilder",
					Links: []client.Link{
						{
							Href: "./Datei:2018_Feucht_Hauptstraße_37_02.jpg",
						},
						{
							Text: "weitere Bilder",
							Href: "https://commons.wikimedia.org/wiki/Category:Hauptstraße%2037%20(Ehemaliges%20Wirtschaftsgebäude,%20D-5-74-123-14)",
						},
					},
				},
			},
		},
	}

	Issue85KeyValue = [][]map[string]string{
		{
			{
				"English release date":  "October 6, 2020[15] 978-1-9747-0993-9",
				"No.":                   "1",
				"Original release date": "March 4, 2019[9] 978-4-08-881780-4",
				"Title":                 "Dog \u0026 ChainsawInu to Chensō (犬とチェンソー)",
			},
			{
				"English release date":  "\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\"（力（パワー) , Pawā)\"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
				"No.":                   "\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\"（力（パワー) , Pawā)\"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
				"Original release date": "\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\"（力（パワー) , Pawā)\"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
				"Title":                 "\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\"（力（パワー) , Pawā)\"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
			},
		},
	}

	MismatchRowsKeyValue = [][]map[string]string{
		{
			{
				"Rank":    "1",
				"Account": "Alpha",
			},
			{
				"Rank":    "1",
				"Account": "Alpha",
				"null2":   "Extra",
			},
		},
	}

	AllTableClasses = [][][]string{
		GoldenMatrix[0],
		GoldenMatrix[0],
		GoldenMatrix[0],
	}

	SimpleKeyValue = [][]map[string]string{
		{
			{
				"Rank":    "1",
				"Account": "Alpha",
			},
		},
	}

	ComplexKeyValue = [][]map[string]string{
		{
			{
				"Date":              "18–24 April 2022",
				"Brand":             "Roy Morgan",
				"Interview mode":    "Telephone/online",
				"Sample size":       "1393",
				"Primary vote L/NP": "35.5%",
				"Primary vote ALP":  "35%",
				"Primary vote GRN":  "12%",
				"Primary vote ONP":  "4.5%",
				"Primary vote UAP":  "1.5%",
				"Primary vote OTH":  "11.5%",
				"UND":               "–",
				"2pp vote L/NP":     "45.5%",
				"2pp vote ALP":      "54.5%",
			},
			{
				"Date":              "20–23 April 2022",
				"Brand":             "Newspoll-YouGov",
				"Interview mode":    "Online",
				"Sample size":       "1538",
				"Primary vote L/NP": "36%",
				"Primary vote ALP":  "37%",
				"Primary vote GRN":  "11%",
				"Primary vote ONP":  "3%",
				"Primary vote UAP":  "4%",
				"Primary vote OTH":  "9%",
				"UND":               "–",
				"2pp vote L/NP":     "47%",
				"2pp vote ALP":      "53%",
			},
			{
				"Date":              "20–23 April 2022",
				"Brand":             "Ipsos",
				"Interview mode":    "Telephone/online",
				"Sample size":       "2302",
				"Primary vote L/NP": "32%",
				"Primary vote ALP":  "34%",
				"Primary vote GRN":  "12%",
				"Primary vote ONP":  "4%",
				"Primary vote UAP":  "3%",
				"Primary vote OTH":  "8%",
				"UND":               "8%",
				"2pp vote L/NP":     "45%",
				"2pp vote ALP":      "55%",
			},
		},
	}

	Issue93MatrixVerbose = [][][]client.Verbose{
		{
			{
				{
					Text: "header1",
				},
				{
					Text: "header2",
				},
			},
			{
				{
					Text: "test",
				},
				{
					Text: "Bolivia, Plurinational State of",
					Links: []client.Link{
						{
							Text: "Bolivia, Plurinational State of",
							Href: "./Bolivia",
						},
					},
				},
			},
		},
	}

	Issue93KeyValueVerbose = [][]map[string]client.Verbose{
		{
			{
				"header1": {
					Text: "test",
				},
				"header2": {
					Text: "Bolivia, Plurinational State of",
					Links: []client.Link{
						{
							Text: "Bolivia, Plurinational State of",
							Href: "./Bolivia",
						},
					},
				},
			},
		},
	}
)
