# wikitable-api
An API to get Wikipedia table data.

[www.wikitable2json.com](https://www.wikitable2json.com)

## Example ([https://en.wikipedia.org/wiki/Arhaan_Khan](https://en.wikipedia.org/wiki/Arhaan_Khan))
### All tables
[https://www.wikitable2json.com/api/v1/page/Arhaan_Khan](https://www.wikitable2json.com/api/v1/page/Arhaan_Khan)
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
### Nth table
[https://www.wikitable2json.com/api/v1/page/Arhaan_Khan?n=0](https://www.wikitable2json.com/api/v1/page/Arhaan_Khan?n=0)
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
    }
  ]
}
```
For multiple individual tables, you can use a query like: ```?n=0&n=3&n=7```.

### Other query params
?lang=
The service will look for a page on the English subdomain (en.wikipedia.com) by default. To specify a page on a different language domain, like cs.wikipedia.com, use ?lang=cs.