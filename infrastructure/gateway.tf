# change the email and domain name in line 30 as needed

resource "digitalocean_droplet" "gateway" {
  image = "ubuntu-20-04-x64"
  name = "gateway"
  region = "nyc1"
  size = "s-1vcpu-512mb-10gb"
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
      "sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common",
      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
      "sudo add-apt-repository -y \"deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable\"",
      "apt-cache policy docker-ce",
      "sudo apt-get install -y docker-ce",
      "sudo apt-install -y letsencrypt",
      "sudo letsencrypt certonly --standalone -n --agree-tos --email rory.james2021@gmail.com -d api.rjames.me",
    ]
  }
}