# wikitable-api
An API to get Wikipedia table data.

[www.wikitable2json.com](https://www.wikitable2json.com)

## Examples
### All tables
[https://www.wikitable2json.com/api/Arhaan_Khan](https://www.wikitable2json.com/api/Arhaan_Khan)
```
{
  "tables": [
    {
      "rows": {
        "0": {
          "columns": {
            "0": "Year",
            "1": "Title",
            "2": "Role",
            "3": "Notes"
          }
        },
        "1": {
          "columns": {
            "0": "2017",
            "1": "SriValli",
            "2": "Majnu",
            "3": "Telugu film"
          }
        }
      }
    },
    {
      "rows": {
        "0": {
          "columns": {
            "0": "Year",
            "1": "Title",
            "2": "Role",
            "3": "Channel",
            "4": "Notes"
          }
        },
        "1": {
          "columns": {
            "0": "2016-2018",
            "1": "Badho Bahu",
            "2": "Rana Ahlawat",
            "3": "&TV",
            "4": "Lead"
          }
        },
        "2": {
          "columns": {
            "0": "2017-2018",
            "1": "Glamx Mr and Miss India Youth Icon",
            "2": "Himself",
            "3": "GlamX Entertainment",
            "4": "Judge"
          }
        },
        "3": {
          "columns": {
            "0": "2019",
            "1": "Bigg Boss 13",
            "2": "Contestant",
            "3": "Colors TV",
            "4": "Evicted on Day 92"
          }
        }
      }
    }
  ]
}
```

### Query parameters
#### Specific tables: ?table=
[https://www.wikitable2json.com/api/Arhaan_Khan?table=0&table=1](https://www.wikitable2json.com/api/Arhaan_Khan?table=0&table=1)
#### Non-English page: ?lang=
[https://www.wikitable2json.com/api/Liste_der_Kulturdenkmale_in_Dedeleben?lang=de](https://www.wikitable2json.com/api/Liste_der_Kulturdenkmale_in_Dedeleben?lang=de)