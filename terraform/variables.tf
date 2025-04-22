variable "subscription_id" {
  description = "Azure subscription ID"
  type        = string
}

variable "resource_group_name" {
  description = "リソースグループの名前"
  type        = string
  default     = "example_rg"
} 