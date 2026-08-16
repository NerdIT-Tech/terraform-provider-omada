terraform {
  required_providers {
    omada = {
      source = "NerdIT-Tech/omada"
    }
  }
}

provider "omada" {
  # host, client_id, client_secret, and omadac_id may also be supplied via
  # the OMADA_HOST, OMADA_CLIENT_ID, OMADA_CLIENT_SECRET, and OMADA_OMADAC_ID
  # environment variables, which keeps credentials out of configuration.
  host      = "https://omada.example.com:8043"
  omadac_id = "example-omadac-id"
}
