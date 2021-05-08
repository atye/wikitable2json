# wikitable-api
An API to get Wikipedia tables in JSON.

[www.wikitable2json.com](https://www.wikitable2json.com)

## Examples
### All tables
[https://www.wikitable2json.com/api/Arhaan_Khan](https://www.wikitable2json.com/api/Arhaan_Khan)
```
{
  "tables":[
      {
        "caption":"",
        "data":[
            [
              "Year",
              "Title",
              "Role",
              "Notes"
            ],
            [
              "2017",
              "SriValli",
              "Majnu",
              "Telugu film"
            ]
        ]
      },
      {
        "caption":"",
        "data":[
            [
              "Year",
              "Title",
              "Role",
              "Notes"
            ],
            [
              "2016-2018",
              "Badho Bahu",
              "Rana Ahlawat",
              "Television debut"
            ],
            [
              "2019",
              "Bigg Boss 13",
              "Contestant",
              "Entered on day 36 evicted on day 50, re-entered on day 65  re-evicted on day 92"
            ]
        ]
      }
  ]
}
```

### Query parameters
#### Specific tables: ?table=
[https://www.wikitable2json.com/api/Arhaan_Khan?table=0&table=1](https://www.wikitable2json.com/api/Arhaan_Khan?table=0&table=1)
#### Non-English page: ?lang=
[https://www.wikitable2json.com/api/Liste_der_Kulturdenkmale_in_Dedeleben?lang=de](https://www.wikitable2json.com/api/Liste_der_Kulturdenkmale_in_Dedeleben?lang=de)