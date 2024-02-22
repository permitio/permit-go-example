terraform {
  required_providers {
    permitio = {
      source  = "permitio/permit-io"
      version = "~> 0.0.1"
    }
  }
}

provider "permitio" {
    api_url = "https://api.permit.io"
    api_key = "<API_KEY>" # Please insert your API KEY
}

resource "permitio_resource" "blog" {
  key         = "blog"
  name        = "blog"
  description = "a new blog"
  actions = {
    "read" = {
      "name" = "read"
    }
    "write" = {
      "name" = "write"
    }
    "delete" = {
      "name"        = "delete"
      "description" = "delete a blog"
    }
  }
  attributes = {}
}

resource "permitio_role" "owner" {
  key         = "owner"
  name        = "owner"
  description = "can do everything"
  permissions = ["blog:read", "blog:write", "blog:delete"]
  depends_on = [
    permitio_resource.blog
  ]
}
resource "permitio_role" "editor" {
  key         = "editor"
  name        = "editor"
  description = "list and updates blogs"
  permissions = ["blog:read", "blog:write"]
  extends     = []
  depends_on = [
    permitio_resource.blog,
    permitio_role.owner
  ]
}

resource "permitio_role" "viewer" {
  key         = "viewer"
  name        = "viewer"
  description = "list blogs"
  permissions = ["blog:read"]
  extends     = []
  depends_on = [
    permitio_resource.blog,
    permitio_role.owner,
    permitio_role.editor
  ]
}

resource "permitio_user_set" "privileged_users" {
  key  = "privileged_users"
  name = "Privileged Users"
  conditions = jsonencode({
    "allOf" : [
      {
        "allOf" : [
          {
            "subject.is_superuser" = {
              equals = "true"
            },
          }
        ]
      }
    ]
  })
}