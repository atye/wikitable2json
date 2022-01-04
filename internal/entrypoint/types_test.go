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

	ReferenceMatrix = [][][]string{
		{
			{"Roy Morgan"},
		},
	}
)
