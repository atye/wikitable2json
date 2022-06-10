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
)
