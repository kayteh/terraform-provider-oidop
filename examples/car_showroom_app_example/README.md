# Car Showroom App Example

This shows you how to create a simple app with 3 fields in it - Make, Model, and Qty - to manage the list of car models that a showroom has in stock.

To try out this example:
1. Install the podio-terraform-provider following the [README.md](../../README.md) file in the root of this repo
1. Update [sample.tfvars]('./sampel.tfvars) with your Podio credentials and org slug. Details on how to get each of the required variable values are provided in the comments on the [sample.tfvars]('./sampel.tfvars) file.
1. Create a space in your org named "Testing Podio Terraform Provider", note down the URL at the top after you enter this space on the Podio web user interface
1. Get the space ID of your space, it will be required in the next step:
   * Go to https://developers.podio.com/doc/spaces/get-space-by-url-22481, scroll down and click the `Login here` button under the `Sandbox` section. Log into your Podio account
   * You should now see three text fields under the `Sandbox` section. Fill the `url` field with the URL you see after opening the space you created in the previous step and click the `Submit` button
   * Scroll to the bottom of the response to note your space ID
1. Run the following command to initialize your terraform state with the created Podio Space:

    ```
    terraform import -var-file=sample.tfvars podio_space.podio_space <space ID>
    ```
1. Finally, run this command to apply your terraform plan:
   ```
   terraform apply -var-file=sample.tfvars
   ```