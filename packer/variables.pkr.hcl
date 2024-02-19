variable "gcp_project_id" {
  type        = string
  description = "The ID of the GCP project"
  default     = "cloud-csye6225-dev"
}

variable "gcp_zone" {
  type        = string
  description = "The GCP zone"
  default     = "us-east1-b"
}

#variable "gcp_machine_type" {
#  type        = string
#  description = "The GCP machine type"
#  default     = "f1-micro"
#}

variable "gcp_disk_size" {
  type        = string
  description = "The GCP disk size"
  default     = "50"
}

variable "gcp_disk_type" {
  type        = string
  description = "The GCP disk type"
  default     = "pd-standard"
}

variable "gcp_network" {
  type        = string
  description = "The GCP network"
  default     = "default"
}

variable "gcp_subnetwork" {
  type        = string
  description = "The GCP subnetwork"
  default     = "default"
}


variable "postgres_user" {
  type        = string
  description = "The PostgreSQL username"
#  sensitive   = true
  default     = "girish"
}

variable "postgres_password" {
  type        = string
  description = "The PostgreSQL password"
#  sensitive   = true
  default     = "test1234"
}

variable "postgres_conn_str" {
  type        = string
  description = "The PostgreSQL connection string"
  default     = "postgresql://girish:test1234@localhost:5432/postgres"
}