# Setting up JetStream Terraform Provider


```bash
$ brew install terraform
$ terraform version
Terraform v0.15.4
on darwin_amd64
```

Go to the release page as mentioned [here](https://docs.nats.io/jetstream/configuration_mgmt/terraform
) to download the release. However, the step for setting the plugins is wrong.

Copy the file `terraform-provider-jetstream_v0.0.2` to `/.terraform.d/plugins/terraform.example.com/mycorp/jetstream/1.0.0/darwin_amd64/<any_file_name_template>` instead.
`

Add this to your terraform template at the topmost, e.g. `local_terraform.tf`:
```terraform
terraform {
  required_providers {
    jetstream = {
      source = "terraform.example.com/mycorp/jetstream"
    }
  }
}
```

Run
```bash
$ terraform init
```
