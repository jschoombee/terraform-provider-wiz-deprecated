resource "wiz_project" "this" {
    name = "tf_test"
   
    cloud_account_links = [
    { 
      cloud_account_guid = "3225def3-0e0e-5cb8-955a-3583f696f778",
      environment = "DEVELOPMENT" , 
      shared = false
    },

  ]
}

