---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "webdock_server Resource - terraform-provider-webdock"
subcategory: ""
description: |-
  
---

# webdock_server (Resource)



## Example Usage

```terraform
resource "webdock_server" "this" {
  slug           = "example"
  name           = "example"
  location_id    = "fi"
  profile_slug   = "webdockbit-2022"
  virtualization = "container"
  image_slug     = "krellide:webdock-jammy-lemp"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `image_slug` (String) Slug of the server image. Get this from the /images endpoint. You must pass either this parameter or snapshotId
- `location_id` (String) ID of the location. Get this from the /locations endpoint.
- `name` (String)
- `profile_slug` (String) Slug of the server profile. Get this from the /profiles endpoint.

### Optional

- `slug` (String) Must be unique
- `virtualization` (String)

### Read-Only

- `date` (String)
- `image` (String)
- `ipv4` (String)
- `ipv6` (String)
- `last_updated` (String)
- `location` (String)
- `profile` (String)
- `snapshot_run_time` (Number)
- `ssh_password_auth_enabled` (Boolean)
- `status` (String)
- `web_server` (String)
- `word_press_lock_down` (Boolean)