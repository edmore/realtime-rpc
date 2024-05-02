variable "az" {
    type = list
    default = ["a", "b", "c", "d", "e", "f"]
}

variable "region" {
    type = string
    default = "us-east-1"
}

variable "service_endpoint" {
    type = string
}