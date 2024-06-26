# Dynamic DNS (DDNS) Service in Golang
## Overview

This project implements a Dynamic DNS (DDNS) service in Golang, designed to automate the process of updating DNS records for your domains. It supports integration with Cloudflare and Google Cloud DNS providers, allowing you to keep your DNS records synchronized with the public IP address of your server.

## Features

- **Dynamic Updates:** Automatically updates DNS records based on the current public IP address.
- **Provider Support:** Works with Cloudflare and Google Cloud DNS providers.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/OsmanTunahan/ddns.git
   cd ddns
   ```

2. **Configuration:**

   Update the `config.json` file with your DNS provider credentials and domain details:

   ```json
   {
       "dns_provider": "cloudflare", // or "google"
       "api_key": "cloudflare_api_key",
       "email": "osmantunahan@icloud.com", // Only needed for Cloudflare
       "project_id": "google_project_id", // Only needed for Google Cloud DNS
       "domain": "awoken.dev",
       "log_level": "info"
   }
   ```

## Usage

- **Run the service:**

  Start the DDNS service by running the main executable:

  ```bash
  go run main.go
  ```

  Alternatively, build the binary and run it:

  ```bash
  go build -o ddns-service main.go
  ./ddns-service
  ```

- **Logs:**

  Logs are printed to the console based on the specified `log_level` in `config.json`. Adjust `log_level` to control the verbosity of logs.