# change the email and domain name in line 25 as needed

resource "digitalocean_droplet" "client" {
  image = "docker-20-04"
  name = "client"
  region = "nyc1"
  size = "s-1vcpu-1gb"
  ssh_keys = [
    data.digitalocean_ssh_key.terraform.id
  ]
  connection {
    host = self.ipv4_address
    user = "root"
    type = "ssh"
    private_key = file(var.pvt_key)
    timeout = "2m"
  }
  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      "sudo ufw allow 80",
      "sudo ufw allow 443",
      "sudo apt update",
      "sudo apt install -y letsencrypt",
      "sudo letsencrypt certonly --standalone -n --agree-tos --email rory.james2021@gmail.com -d rjames.me",
    ]
  }
}