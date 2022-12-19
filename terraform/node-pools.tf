# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_service_account
resource "google_service_account" "kubernetes" {
  account_id = "${var.project_id}-sa"
}

# https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_node_pool
resource "google_container_node_pool" "general" {
  name       = "${var.project_id}-general"
  cluster    = google_container_cluster.primary.id
  node_count = 1
  #  initial_node_count       = 1

  autoscaling {
    min_node_count = 1
    max_node_count = 5
  }
  management {
    auto_repair  = true
    auto_upgrade = true
  }

  node_config {
    preemptible  = false
    machine_type = "n1-standard-1"

    labels = {
      role = "general"
    }

    disk_size_gb    = "30"
    service_account = google_service_account.kubernetes.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

