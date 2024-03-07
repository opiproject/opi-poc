# DPU Host 3

an aditional service was created in order to set the nic at boot in 1G mode
otherwise it starts in 10G/25G mode and link will be down
`ethtool -s enP2p1s0f0np0 speed 1000`
