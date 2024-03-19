variable "gcp_project_id" {
  type        = string
  description = "The ID of the GCP project"
  default     = "cloud-csye6225-dev"
}

variable "gcp_zone" {
  type        = string
  description = "The GCP zone"
  default     = "us-east1-c"
}

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

variable "postgres_user" {
  type        = string
  description = "The PostgreSQL username"
  sensitive   = true
  default     = env("POSTGRES_USER")
}

variable "postgres_password" {
  type        = string
  description = "The PostgreSQL password"
  sensitive   = true
  default     = env("POSTGRES_PASSWORD")
}

variable "postgres_conn_str" {
  type        = string
  description = "The PostgreSQL connection string"
  sensitive   = true
  default     = "postgresql://girish:test1234@localhost:5432/postgres"
}