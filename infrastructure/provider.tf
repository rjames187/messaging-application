terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

# Digital Ocean personal access token
variable "do_token" {}
# location of private key for ssh
variable "pvt_key" {}

provider "digitalocean" {
  token = var.do_token
}

# here, terraform refers to the name of the ssh key saved in Digital Ocean
data "digitalocean_ssh_key" "terraform" {
  name = "terraform"
}

