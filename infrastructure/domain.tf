resource "digitalocean_domain" "default" {
   name = "rjames.me"
   ip_address = digitalocean_droplet.client.ipv4_address
}

resource "digitalocean_record" "A-gateway" {
  domain = digitalocean_domain.default.name
  type = "A"
  name = "api"
  value = digitalocean_droplet.gateway.ipv4_address
}

