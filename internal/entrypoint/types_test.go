package entrypoint

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
			{},
			{"FeuchtHauptstraße 37(Standort)", "Ehemaliges Wirtschaftsgebäude", "Zweigeschossiger Satteldachbau mit Fachwerkobergeschoss, bezeichnet mit „1697“", "ehemals D-5-74-123-14 zugehörig", "weitere Bilder"},
		},
	}

	Issue85Matrix = [][][]string{
		{
			{"No.", "Title", "Original release date", "English release date"},
			{"1", "Dog \u0026 ChainsawInu to Chensō (犬とチェンソー)",
				"March 4, 2019[9]",
				"October 6, 2020[15]978-1-9747-0993-9"},
			{"\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\" 力（パワー）, Pawā)\n\n\"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
				"\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\" 力（パワー）, Pawā)\n\n\"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
				"\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\" 力（パワー）, Pawā)\n\n\"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
				"\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\" 力（パワー）, Pawā)\n\n\"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)"},
		},
	}

	Issue85KeyValue = [][]map[string]string{
		{
			{
				"English release date":  "October 6, 2020[15]978-1-9747-0993-9",
				"No.":                   "1",
				"Original release date": "March 4, 2019[9] 978-4-08-881780-4",
				"Title":                 "Dog \u0026 ChainsawInu to Chensō (犬とチェンソー)",
			},
			{
				"English release date":  "\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\" (力（パワー）, Pawā)\n                \n                \"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
				"No.":                   "\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\" (力（パワー）, Pawā)\n                \n                \"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
				"Original release date": "\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\" (力（パワー）, Pawā)\n                \n                \"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
				"Title":                 "\"Dog \u0026 Chainsaw\" (犬とチェンソー, Inu to Chensō)\"The Place Where Pochita Is\" (ポチタの行方, Pochita no Yukue)\"Arrival in Tokyo\" (東京到着, Tōkyō Tōchaku)\"Power\" (力（パワー）, Pawā)\n                \n                \"A Way to Touch Some Boobs\" (胸を揉む方法, Mune o Momu Hōhō)\"Service\" (使役, Shieki)\"Meowy's Whereabouts\" (ニャーコの行方, Nyāko no Yukue)",
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
)
