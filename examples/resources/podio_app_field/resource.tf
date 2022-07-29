resource "podio_app_field" "type_some_name_here" {
  app_id      = 1234
  type        = "text"
  label       = "Some label for your field"
  description = "Describe your field here"
  required    = true # set this to true if you want to mandate entering something into this field in order to add an item into the app this field belongs to
}
