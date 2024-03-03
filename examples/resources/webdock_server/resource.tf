resource "webdock_server" "this" {
  slug           = "example"
  name           = "example"
  location_id    = "fi"
  profile_slug   = "webdockbit-2022"
  virtualization = "container"
  image_slug     = "krellide:webdock-jammy-lemp"
}
