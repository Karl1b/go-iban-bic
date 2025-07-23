# GO-IBAN-BIC

Fast microservice that validates German IBANs and returns BIC codes and bank information. Built to demonstrate Go's power and simplicity.

## Overview

Created when a colleague couldn't find a good IBAN-to-BIC conversion source for German banks. Showcases Go's efficiency with practical banking data lookup.

**Key Features**: Lightning fast (data in RAM), IBAN validation, complete German bank database, simple JSON API.

## Data Source

Official Bundesbank dataset: https://www.bundesbank.de/de/aufgaben/unbarer-zahlungsverkehr/serviceangebot/bankleitzahlen/download-bankleitzahlen-602592

**Important**: Update monthly for production use.

## Usage

### API Example

```bash
curl --json '{"iban":"DE14 6001 0070 0032 3357 09"}' localhost:8080/iban
```

**Response**:
```json
{
  "iban": "DE14 6001 0070 0032 3357 09",
  "is_valid": true,
  "bic": "PBNKDEFFXXX",
  "bezeichnung": "Postbank Ndl der Deutsche Bank",
  "ort": "Stuttgart",
  "blz": "60010070"
}
```

### As Go Package

```go
import "github.com/Karl1b/go-iban-bic/pkg/ibanbic"

// Validate and get BIC info
isValid := ibanbic.ValidateIBAN(iban)
bicInfo := ibanbic.GetBic(iban)

```

## Configuration

Service reads PORT from environment variables or `.env` file. Health endpoint available at `/health`.

## Performance

O(1) lookup time, ~10MB RAM usage, thousands of requests per second. Complete German bank database loaded at startup for maximum speed.