data "podio_organization" "my_org" {
  url_label = var.org_slug
}

resource "podio_space" "test_space" {
  name   = "Testing Podio Terraform Provider"
  org_id = data.podio_organization.my_org.org_id
}

resource "podio_app" "cars" {
  space_id  = podio_space.test_space.space_id
  name      = "Cars"
  item_name = "Car"
  icon      = "251.png"
}

resource "podio_app_field" "cars_make" {
  app_id      = podio_app.cars.app_id
  type        = "text"
  label       = "Make"
  description = "The make of the vehicle (example: Ford)"
  required    = true
}

resource "podio_app_field" "cars_model" {
  app_id      = podio_app.cars.app_id
  type        = "text"
  label       = "Model"
  description = "The model of the vehicle (example: F-150)"
  required    = true
}

resource "podio_app_field" "cars_qty" {
  app_id      = podio_app.cars.app_id
  type        = "number"
  label       = "Qty"
  description = "Number of cars of this make and model in stock"
  required    = true
}
