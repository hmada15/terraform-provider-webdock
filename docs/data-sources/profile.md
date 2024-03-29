---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "webdock_profile Data Source - terraform-provider-webdock"
subcategory: ""
description: |-
  
---

# webdock_profile (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `location_id` (String) Location of the profile

### Read-Only

- `profiles` (Attributes List) (see [below for nested schema](#nestedatt--profiles))

<a id="nestedatt--profiles"></a>
### Nested Schema for `profiles`

Read-Only:

- `cpu` (Attributes) CPU model (see [below for nested schema](#nestedatt--profiles--cpu))
- `disk` (Number) Disk size (in MiB)
- `name` (String) Profile name
- `price` (Attributes) Price model (see [below for nested schema](#nestedatt--profiles--price))
- `ram` (Number) RAM memory (in MiB)
- `slug` (String) Profile slug

<a id="nestedatt--profiles--cpu"></a>
### Nested Schema for `profiles.cpu`

Read-Only:

- `cores` (Number) cpu cores
- `threads` (Number) cpu threads


<a id="nestedatt--profiles--price"></a>
### Nested Schema for `profiles.price`

Read-Only:

- `amount` (Number) Price amount
- `currency` (String) Price currency
