## API Endpoints

### Create Quote

- **URL:** `POST /quote`

- **Body:**

```json
{
  "shipper": {
    "registered_number": "<your_frete_rapido_cnpj>",
    "token": "<your_frete_rapido_token>",
    "platform_code": "<your_frete_rapido_platform_code>"
  },
  "recipient": {
    "type": 0,
    "country": "BRA",
    "zipcode": "<your_dispatcher_zipcode>"
  },
  "dispatchers": [
    {
      "registered_number": "<your_frete_rapido_cnpj>",
      "zipcode": "<your_dispatcher_zipcode>",
      "volumes": [
        {
          "category": "7",
          "amount": 1,
          "unitary_weight": 5,
          "unitary_price": 349,
          "sku": "abc-teste-123",
          "height": 0.2,
          "width": 0.2,
          "length": 0.2
        },
        {
          "category": "7",
          "amount": 2,
          "unitary_weight": 4,
          "unitary_price": 556,
          "sku": "abc-teste-527",
          "height": 0.4,
          "width": 0.6,
          "length": 0.15
        }
      ]
    }
  ],
  "simulation_type": [
    0
  ]
}
```

- **Response:**

```json
{
  "carrier": [
    {
      "name": "EXPRESSO FR",
      "service": "Rodoviário",
      "deadline": 3,
      "price": 17
    },
    {
      "name": "Correios",
      "service": "SEDEX",
      "deadline": 1,
      "price": 20.99
    }
  ]
}
```

- **Error Response:** In case of an error, an error code will be returned, as established in
  the [list of codes of this API](https://dev.freterapido.com/common/codigos_de_resposta/).

### Get Metrics

- **URL:** `GET /metrics?last_quotes={?}`

- **Response:**

```json
{
  "carriers": {
    "EXPRESSO FR": {
      "count": 2,
      "total_price": 34,
      "average_price": 17
    },
    "Correios": {
      "count": 1,
      "total_price": 20.99,
      "average_price": 20.99
    }
  },
  "cheapest_quote": {
    "name": "EXPRESSO FR",
    "service": "Rodoviário",
    "deadline": 3,
    "price": 17
  },
  "most_expensive_quote": {
    "name": "Correios",
    "service": "SEDEX",
    "deadline": 1,
    "price": 20.99
  }
}
```

- **Error Response:** In case of an error, an error code will be returned, as established in
  the [list of codes of this API](https://dev.freterapido.com/common/codigos_de_resposta/).
