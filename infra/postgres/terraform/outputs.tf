output "primary_public_ip" {
  value = oci_core_instance.pg_primary.public_ip
}

output "replica_public_ip" {
  value = oci_core_instance.pg_replica.public_ip
}