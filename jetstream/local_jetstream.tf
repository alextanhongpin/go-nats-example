# Reference:
# https://github.com/nats-io/terraform-provider-jetstream


#https://www.terraform.io/docs/language/providers/requirements.html#source-addresses

# Place the plugin under the folder /.terraform.d/plugins/terraform.example.com/mycorp/jetstream/1.0.0/darwin_amd64/<any_file_name_template>.
terraform {
  required_providers {
    jetstream = {
      source = "terraform.example.com/mycorp/jetstream"
    }
  }
}

provider "jetstream" {
  servers = "localhost"
}

resource "jetstream_stream" "ORDERS" {
  name     = "ORDERS"
  subjects = ["ORDERS.*"]
  storage  = "file"
  max_age  = 60 * 60 * 24 * 365
}

resource "jetstream_consumer" "ORDERS_NEW" {
  stream_id      = jetstream_stream.ORDERS.id
  durable_name   = "NEW"
  deliver_all    = true
  filter_subject = "ORDERS.scratch"
  sample_freq    = 100
}
