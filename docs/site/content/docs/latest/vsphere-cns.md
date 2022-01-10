# Set Up vSphere CNS and Create a Storage Policy in vSphere

vSphere administrators can set up vSphere CNS and create storage policies for virtual disk (VMDK) storage, based on the needs of workload cluster users.

You can use either vSAN or local VMFS (Virtual Machine File System) for persistent storage in a Kubernetes cluster, as follows:

**vSAN Storage**:

To create a storage policy for vSAN storage in the vSphere Client, browse to **Home** > **Policies and Profiles** > **VM Storage Policies** and click **Create** to launch the **Create VM Storage Policy** wizard.

Follow the instructions in [Create a Storage Policy](https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.storage.doc/GUID-8D51CECC-ED3B-424E-BFE2-43379729A653.html) in the vSphere documentation. Make sure to:

- In the **Policy structure** pane, under **Datastore specific rules**, select **Enable rules for "vSAN" storage**.
- Configure other panes or accept defaults as needed.
- Record the storage policy name for reference as the `storagePolicyName` value in `StorageClass` objects.

**Local VMFS Storage**:

To create a storage policy for local storage, apply a tag to the storage and create a storage policy based on the tag as follows:

1. From the top-level vSphere menu, select **Tags &amp; Custom Attributes**

1. In the **Tags** pane, select **Categories** and click **New**.

1. Enter a category name, such as `tkg-storage`.
Use the checkboxes to associate it with **Datacenter** and the storage objects, **Folder** and **Datastore**.
Click **Create**.

1. From the top-level **Storage** view, select your VMFS volume, and in its **Summary** pane, click **Tags** > **Assign...**.

1. From the **Assign Tag** popup, click **Add Tag**.

1. From the **Create Tag** popup, give the tag a name, such as `tkg-storage-ds1` and assign it the **Category** you created.  Click **OK**.

1. From **Assign Tag**, select the tag and click **Assign**.

1. From top-level vSphere, select **VM Storage Policies** > **Create a Storage Policy**.  A configuration wizard starts.

1. In the **Name and description** pane, enter a name for your storage policy.
Record the storage policy name for reference as the `storagePolicyName` value in `StorageClass` objects.

1. In the **Policy structure** pane, under **Datastore specific rules**, select **Enable tag-based placement rules**.

1. In the **Tag based placement** pane, click **Add Tag Rule** and configure:

    - **Tag category**: Select your category name
    - **Usage option**: `Use storage tagged with`
    - **Tags**: Browse and select your tag name

1. Confirm and configure other panes or accept defaults as needed, then click **Review and finish**. **Finish** to create the storage policy.
