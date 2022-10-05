terraform {
  required_version = ">= 1"

  required_providers {

    wiz = {
        source = "terraform.shell.com/deepblue/wiz"
    }
  }

  backend "local" {
    path = "./terraform.tfstate"
  }

}

provider "wiz" {
  endpoint = "https://api.eu1.app.wiz.io/graphql"
  client_id = "here"
  client_secret = "here"
}