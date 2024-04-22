# Redfish

see <https://github.com/opiproject/ansible-opi-dpu> and <https://github.com/opiproject/pydpu/tree/main/pydpu/dpuredfish>

Example:

```bash
docker run --rm --pull=always ghcr.io/opiproject/ansible-opi-dpu:main all \
    --module-name include_role \
    --args name=bmc_fw_update \
    -vvv -i "172.22.4.2," \
    -e dpu_bmc_username='root' \
    -e dpu_bmc_password='NvidiaBf2#Pass' \
    -e bmc_fw_update_inventory_name='a5e46d02_running' \
    -e bmc_fw_update_image_file='/tmp/bf2-bmc-ota-24.01-5-opn.tar'
```
