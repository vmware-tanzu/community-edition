#!/bin/bash
#
# Convenience script to set up Tanzu Community Edition control host
#
# NOTE: Currently only supports Debian based installs
#
# Basic idea is to be able to spin up an Ubuntu VM, copy this script
# over, then run through to get a working TCE deployment going as quickly
# as possible.

IFS=$'\n'

# Helper functions
function error {
    printf '\E[31m'; echo "$@"; printf '\E[0m'
}

# We will need to access a private repo, make sure token is set
if [[ -z "$GH_USER" ]]; then
    echo "Access to GitHub private repo requires your GitHub username."

    read -r -p "Please enter your GitHub username: " GH_USER
    echo
fi

if [[ -z "$GH_TANZU_TOKEN" ]]; then
    echo "Access to GitHub private repo requires a token."
    echo "Please create a token (Settings > Developer Settings > Personal Access Tokens)"

    read -r -p "Please enter your GitHub token: " GH_TANZU_TOKEN
    echo
fi

# We need to make system changes, make sure we are running as root
if [[ $(id -u) -ne 0 && $(sudo -n true) -ne 0 ]]; then
    error "Please run this script as root/sudo"
    exit 1
fi

sudo apt update > /dev/null 2>&1
sudo apt install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common wget jq > /dev/null 2>&1

# Make sure we have Docker installed
if [[ -z "$(which docker)" ]]; then

    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository \
        "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
        $(lsb_release -cs) \
        stable"

    sudo apt update > /dev/null
    sudo apt install -y docker-ce docker-ce-cli containerd.io > /dev/null 2>&1

    if [[ $(id -u) -ne 0 ]]; then
        sudo usermod -aG docker "$(whoami)"
    fi
fi

if ! sudo docker run hello-world > /dev/null; then
    error "Unable to verify docker functionality, make sure docker is installed correctly"
    exit 1
fi

# Make sure we have kubectl
if [[ -z "$(which kubectl)" ]]; then

    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
    echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list

    sudo apt update > /dev/null
    sudo apt install -y kubectl > /dev/null 2>&1
fi

# Temporary: Get the tkg CLI
if [[ "$(which tkg)" == "" ]]; then
    echo
    echo "==============================="
    echo " IMPORT: MANUAL STEPS REQUIRED"
    echo "==============================="
    echo
    echo "The tkg CLI is required for now. This is only temporary."
    echo "Please download and install the CLI from here:"
    echo
    echo "https://www.vmware.com/go/get-tkg"
    echo
    read -r -s -p $"Press ENTER to continue once installation is complete..."
    echo
fi

if [[ -z $(which tkg) ]]; then
    error "Unable to find the tkg CLI"
    error "Please recheck manual installation and try again"
    exit 1
fi

RELINFO=$(curl -s -u "$GH_USER:$GH_TANZU_TOKEN" https://api.github.com/repos/vmware-tanzu/tce/releases/latest)
CURRENT_RELEASE=$(echo "$RELINFO" | jq -r '.tag_name')
if [[ -z "$CURRENT_RELEASE" ]]; then
    error "Unable to determine current release"
    exit 1
fi
ARTIFACT=$(echo "$RELINFO" | jq -r '.assets | map(select(.name|contains("linux")))[0].url')

echo
echo "Installing TCE $CURRENT_RELEASE"
echo

echo "Getting $ARTIFACT"
ARTIFACT=${ARTIFACT/https:\/\//}
wget -q --auth-no-challenge --header='Accept:application/octet-stream' "https://$GH_TANZU_TOKEN:@$ARTIFACT" -O tce.tgz
tar zxvf tce.tgz
rm tce.tgz
cd dist-linux || echo "Unexpected artifact structure, exiting..."; return
./install.sh
cd ..
rm -fr dist-linux

echo
echo "TKG CLI version:"
tkg version

echo
echo "Tanzu CLI version:"
tanzu version

echo
PS3="Select deployment environment: "
select TYPE in vsphere aws azure; do  break; done
echo

if [[ "$TYPE" == "vsphere" ]]; then
    # Only tested with non-TKG enabled vSphere deployment so far

    echo "Installing govc helper..."
    RELINFO=$(curl -s -u "$GH_USER:$GH_TANZU_TOKEN" https://api.github.com/repos/vmware/govmomi/releases/latest)
    ARTIFACT=$(echo "$RELINFO" | jq -r '.assets | map(select(.name|contains("linux")))[0].url')

    ARTIFACT=${ARTIFACT/https:\/\//}
    wget -q --auth-no-challenge --header='Accept:application/octet-stream' "https://$GH_TANZU_TOKEN:@$ARTIFACT" -O govc.gz
    gunzip govc.gz
    chmod +x govc*
    sudo mv govc* /usr/local/bin/govc
    echo

    read -r -p "Enter vSphere management IP/host: " VSPHERE_SERVER
    read -r -p "Enter vSphere username: " VSPHERE_USERNAME
    read -r -s -p "Enter vSphere password: " VSPHERE_PASSWORD
    echo
    echo

    export GOVC_INSECURE=1
    export GOVC_URL=$VSPHERE_SERVER
    export GOVC_USERNAME=$VSPHERE_USERNAME
    export GOVC_PASSWORD=$VSPHERE_PASSWORD

    # govc ls -t=Datacenter '*' appears to only return something if there
    # are more than one datacenters present.
    # DCS=$(govc ls -t=Datacenter '*')
    mapfile -t < <(govc ls '*' | cut -d '/' -f 2 | uniq | sort) DCS
    PS3="Select datacenter: "
    select VSPHERE_DATACENTER in "${DCS[@]}"; do break; done
    VSPHERE_DATACENTER="/$VSPHERE_DATACENTER"
    echo

    # Get the datastore to use
    mapfile -t < <(GOVC_DATACENTER=$VSPHERE_DATACENTER govc ls -t=Datastore "$VSPHERE_DATACENTER/*") DSS
    PS3="Select datastore: "
    select VSPHERE_DATASTORE in "${DSS[@]}"; do break; done
    echo

    # Add in the default VM folder, then append any others since govc will not return the default.
    FOLDERS=("$VSPHERE_DATACENTER/vm")
    mapfile -t -O 1 < <(GOVC_DATACENTER=$VSPHERE_DATACENTER govc ls -t=Folder "$VSPHERE_DATACENTER/vm") FOLDERS
    PS3="Select VM folder: "
    select VSPHERE_FOLDER in "${FOLDERS[@]}"; do break; done
    echo

    # Collect the individual hosts and clusters
    mapfile -t < <(GOVC_DATACENTER=$VSPHERE_DATACENTER govc ls -t=ClusterComputeResource "$VSPHERE_DATACENTER/host") RESOURCES
    mapfile -t -O "${#RESOURCES[@]}" < <(GOVC_DATACENTER=$VSPHERE_DATACENTER govc ls -t=ComputeResource "$VSPHERE_DATACENTER/host") RESOURCES
    PS3="Select host or cluster: "
    select VSPHERE_RESOURCE_POOL in "${RESOURCES[@]}"; do break; done
    VSPHERE_RESOURCE_POOL="$VSPHERE_RESOURCE_POOL/Resources"

    echo
    read -r -p "Paste SSH authorized public key: " VSPHERE_SSH_AUTHORIZED_KEY
    echo

    echo
    PS3="Select deployment size: "
    select SIZE in extra-large large medium small; do  break; done
    echo

    read -r -p "Enter static IP address for the management cluster: " MGMT_IP

    echo
    echo "================================================"
    echo "Using DC:              $VSPHERE_DATACENTER"
    echo "Using datastore:       $VSPHERE_DATASTORE"
    echo "Using VM folder:       $VSPHERE_FOLDER"
    echo "Using Resource pool:   $VSPHERE_RESOURCE_POOL"
    echo "Deployment sizing:     $SIZE"
    echo "Management cluster IP: $MGMT_IP"
    echo "================================================"

    # Make sure everything available in the environment
    export VSPHERE_DATACENTER=$VSPHERE_DATACENTER
    export VSPHERE_DATASTORE=$VSPHERE_DATASTORE
    export VSPHERE_FOLDER=$VSPHERE_FOLDER
    export VSPHERE_SSH_AUTHORIZED_KEY=$VSPHERE_SSH_AUTHORIZED_KEY
    export VSPHERE_RESOURCE_POOL=$VSPHERE_RESOURCE_POOL
    export VSPHERE_SERVER=$VSPHERE_SERVER
    export VSPHERE_USERNAME=$VSPHERE_USERNAME
    export VSPHERE_PASSWORD=$VSPHERE_PASSWORD

    echo tkg init --infrastructure=vsphere --plan dev --size $SIZE --vsphere-controlplane-endpoint "$MGMT_IP"
    echo
    time tkg init --infrastructure=vsphere --plan dev --size $SIZE --vsphere-controlplane-endpoint "$MGMT_IP"

elif [[ "$TYPE" == "aws" ]]; then
    echo
    echo "Not implemented yet!!!!"
    echo
    echo "Run: tkg init --ui"
elif [[ "$TYPE" == "azure" ]]; then
    echo
    echo "Not implemented yet!!!!"
    echo
    echo "Run: tkg init --ui"
fi

