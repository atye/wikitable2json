package entrypoint

var (
	/*ArhaanKhanVerbose = []server.Verbose{
		{
			0: {
				0: "Arhaan Khan",
				1: "Arhaan Khan",
			},
			1: {
				0: "Born",
				1: "Mazhar Shaikh",
			},
			2: {
				0: "Nationality",
				1: "Indian",
			},
			3: {
				0: "Occupation",
				1: "modelActor",
			},
			4: {
				0: "Years active",
				1: "2016-Present",
			},
			5: {
				0: "Known for",
				1: "Badho Bahu  Bigg Boss 13",
			},
			6: {
				0: "Partner(s)",
				1: "Rashami Desai (2019-2020)[1]",
			},
			7: {
				0: "Children",
				1: "1",
			},
		},
		{
			0: {
				0: "Year",
				1: "Title",
				2: "Role",
				3: "Notes",
			},
			1: {
				0: "2017",
				1: "SriValli",
				2: "Majnu",
				3: "Telugu film",
			},
		},
		{
			0: {
				0: "Year",
				1: "Title",
				2: "Role",
				3: "Notes",
			},
			1: {
				0: "2016-2018",
				1: "Badho Bahu",
				2: "Rana Ahlawat",
				3: "Television debut",
			},
			2: {
				0: "2019",
				1: "Bigg Boss 13",
				2: "Contestant",
				3: "Entered on day 36 evicted on day 50, re-entered on day 65  re-evicted on day 92",
			},
		},
	}

	ArhaanKhanMatrix = [][][]string{
		{
			[]string{"Arhaan Khan", "Arhaan Khan"},
			[]string{"Born", "Mazhar Shaikh"},
			[]string{"Nationality", "Indian"},
			[]string{"Occupation", "modelActor"},
			[]string{strings.Replace("Years active", " ", string(nbsp), -1), "2016-Present"},
			[]string{strings.Replace("Known for", " ", string(nbsp), -1), "Badho Bahu  Bigg Boss 13"},
			[]string{"Partner(s)", "Rashami Desai (2019-2020)[1]"},
			[]string{"Children", "1"},
		},
		{
			[]string{"Year", "Title", "Role", "Notes"},
			[]string{"2017", "SriValli", "Majnu", "Telugu film"},
		},
		{
			[]string{"Year", "Title", "Role", "Notes"},
			[]string{"2016-2018", "Badho Bahu", "Rana Ahlawat", "Television debut"},
			[]string{"2019", "Bigg Boss 13", "Contestant", "Entered on day 36 evicted on day 50, re-entered on day 65  re-evicted on day 92"},
		},
	}

	GoldenVerbose = []server.Verbose{
		{
			0: {
				0: "Column 1",
				1: "Column 2",
				2: "Column 3",
			},
			1: {
				0: "A",
				1: "B",
				2: "B",
			},
			2: {
				0: "A",
				1: "C",
				2: "D",
			},
			3: {
				0: "E",
				1: "F",
				2: "F",
			},
			4: {
				0: "G",
				1: "F",
				2: "F",
			},
			5: {
				0: "H",
				1: "H",
				2: "H",
			},
		},
	}*/

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

	/*IssueOneVerbose = []server.Verbose{
		{
			0: {
				0: "Jeju",
				1: "South Korea",
				2: "official, in Jeju Island",
			},
			1: {
				0: "Jeju",
			},
		},
	}*/

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

	/*DataSortValueVerbose = []server.Verbose{
		{
			0: {
				0: "Abu Dhabi, United Arab Emirates",
				1: "N/A",
			},
		},
	}*/
)
