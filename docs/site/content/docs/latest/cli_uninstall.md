# Uninstall the Tanzu CLI
Complete the following steps to remove the Tanzu CLI and any associated configurations

## Linux

1. Run one of the following commands to delete the tanzu binary:

    ```sh
    sudo rm -f ~/.local/share/tanzu
    ```
    or

    ```
    ~/.local/share/tanzu-cli/uninstall.sh
    ```

2. Run the following commands to delete configurations:
    ```
    rm -rf ~/.tanzu
    rm -rf ~/.config/tanzu
    rm -rf ~/tanzu-cli
    ```

##  MacOS

1. Run the following command to delete the tanzu binary:

    ```sh
    sudo rm -f "~/Library/Application Support/tanzu"
    ```

2. Run the following commands to delete configurations:

    ```sh
    shrm -rf ~/.tanzu
    rm -rf ~/.config/tanzu
    rm -rf ~/tanzu-cli
    ```

## Windows
1. Open a command prompt as an administrator and change to the tanzu-cli directory:

    ```sh
    cd %LocalAppData%\tanzu-cli
    ```
2. Run the following command:

    ```sh
    run uninstall.bat
    ```