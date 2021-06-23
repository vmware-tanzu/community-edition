## Clean-up

Standalone cluster deletion is currently unimplemented. Until it is, do the
following steps to remove resources created in Amazon EC2. Following the below steps in order will do a full clean-up of the resources.

1. Terminate all Amazon EC2 instances.

1. Delete all ELB instances.

1. Delete all NAT Gateway instances.

1. Delete the VPC.

    > This will do a final clean-up, but requires the above steps have finished
    their termination.
