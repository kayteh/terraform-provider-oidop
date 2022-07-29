data "podio_organization" "my_org" {
  url_label = "my-org"
}

resource "podio_space" "my_space" {
  org_id  = data.podio_organization.my_org.org_id
  name    = "My Space"
  privacy = "closed"
}

resource "podio_app" "type_some_name_here" {
  space_id    = podio_space.my_space.space_id
  name        = "My App"
  type        = "standard"
  item_name   = "Some appropriate name for your app's item"
  description = "Some description of the app" #optional
  usage       = "How the app should be used"  #optional
  # There are more attributes you can specify, please refer the docs for a full list
}
