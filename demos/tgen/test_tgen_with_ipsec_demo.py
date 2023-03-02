import grpc
import netmiko
import pydpu
import pytest
import snappi
from passwords import *
from testbed import *

# replace with pydpu code
from pydpu.proto.v1 import ipsec_pb2
from pydpu.proto.v1 import ipsec_pb2_grpc


@pytest.fixture(scope='session')
def dpu():
    dpu_info = {
        'device_type': 'linux',
        'host': DPU_IP,
        'username': 'root',
        'password': DPU_PASSWORD
    }
    dpu_connect = netmiko.ConnectHandler(**dpu_info)
    #prompt = dpu_connect.find_prompt()
    dpu_connect.send_command('ip addr add 200.0.0.1/24 dev enp3s0f0s0')
    dpu_connect.send_command('ip addr add 201.0.0.1/24 dev enp3s0f1s0')
    dpu_connect.send_command('sysctl net.ipv4.ip_forward=1')
    dpu_connect.send_command('iptables -P FORWARD ACCEPT')
    #dpu_connect.send_command('arp -s 200.0.0.2 80:09:02:02:00:01')
    dpu_connect.send_command('arp -s 201.0.0.2 80:09:02:02:00:02')
    #dpu_connect.send_command('arp -s 40.0.0.1 80:09:02:02:00:40')
    dpu_connect.send_command('systemctl start docker')

    command = 'docker run --rm --network host --mount src=/var/run,target=/var/run,type=bind --name opi-strongswan-bridge -d ghcr.io/opiproject/opi-strongswan-bridge:main /opi-vici-bridge -port=50151'
    print(command)
    output = dpu_connect.send_command(command)
    print(output)

    # idealy replace as many of those ssh commands with OPI API as they become available

    yield dpu_connect

    dpu_connect.send_command('ip addr del 200.0.0.1/24 dev enp3s0f0s0')
    dpu_connect.send_command('ip addr del 201.0.0.1/24 dev enp3s0f1s0')
    dpu_connect.send_command('docker stop opi-strongswan-bridge')
    dpu_connect.send_command('ip neigh flush 200.0.0.0/8')
    dpu_connect.send_command('ip neigh flush 201.0.0.0/8')


@pytest.fixture(scope='session')
def server():
    server_info = {
        'device_type': 'linux',
        'host': SERVER_IP,
        'username': 'root',
        'password': SERVER_PASSWORD
    }
    server_connect = netmiko.ConnectHandler(**server_info)

    server_connect.send_command('ip addr add 200.0.0.2/24 dev enp5s0f0np0')
    

    command = 'docker compose -f /root/opi/tgen/deployment/ipsec2.yml up -d'
    print(command)
    output = server_connect.send_command(command, read_timeout=30)
    print(output)

    command = 'docker exec -i vpn-client-root bash -c "swanctl --load-all"'
    print(command)
    output = server_connect.send_command(command)
    print(output)

    yield server_connect

    command = 'docker compose -f /root/opi/tgen/deployment/ipsec2.yml down'
    print(command)
    output = server_connect.send_command(command)
    print(output)

    server_connect.send_command('ip addr del 200.0.0.2/24 dev enp5s0f0np0')


def test_server_to_server_via_ipsec_and_dpu(dpu, server):

    print('connecting to opi')
    channel = grpc.insecure_channel('10.36.78.168:50151')
    stub = ipsec_pb2_grpc.IPsecStub(channel)

    stub.IPsecVersion(ipsec_pb2.IPsecVersionReq())

    stub.IPsecStats(ipsec_pb2.IPsecStatsReq())

    print('configuring the tunnel')
    tun1_0_0 = ipsec_pb2.IPsecLoadConnReq(
        connection=ipsec_pb2.Connection(
            name='tun1_0_0',
            version='2',
            local_addrs=[ipsec_pb2.Addrs(addr='200.0.0.1')],
            remote_addrs=[ipsec_pb2.Addrs(addr='200.0.0.2')],
            local_auth=ipsec_pb2.LocalAuth(
                auth=ipsec_pb2.PSK,
                id='200.0.0.1'
            ),
            remote_auth=ipsec_pb2.RemoteAuth(
                auth=ipsec_pb2.PSK,
                id='200.0.0.2'
            ),
            children=[ipsec_pb2.Child(
                name='tun1_0_0',
                esp_proposals=ipsec_pb2.Proposals(
                    crypto_alg=[ipsec_pb2.AES128],
                    integ_alg=[ipsec_pb2.SHA1]
                ),
                remote_ts=ipsec_pb2.TrafficSelectors(
                    ts=[ipsec_pb2.TrafficSelectors.TrafficSelector(
                        cidr='40.0.0.0/24'
                    )]
                ),
                local_ts=ipsec_pb2.TrafficSelectors(
                    ts=[ipsec_pb2.TrafficSelectors.TrafficSelector(
                        cidr='201.0.0.0/24'
                    )]
                ),
            )],
            proposals=ipsec_pb2.Proposals(
                crypto_alg=[ipsec_pb2.AES128],
                integ_alg=[ipsec_pb2.SHA384],
                dhgroups=[ipsec_pb2.MODP1536]
            )
        )
    )

    connection_1 = stub.IPsecLoadConn(tun1_0_0)

    command = 'docker exec -i vpn-client-root bash -c "swanctl --initiate --child tun1_0_0"'
    print(command)
    output = server.send_command(command)
    print(output)

    tgen = snappi.api(location=f'https://{SERVER_IP}', verify=False)
    cfg = tgen.config()
    p1 = cfg.ports.port(name='server:p1', location=f'{SERVER_IP}:5555')[-1]
    p2 = cfg.ports.port(name='server:p2', location=f'{SERVER_IP}:5556')[-1]

    # add layer 1 property to configure same speed on both ports
    ly = cfg.layer1.layer1(name='ly')[-1]
    ly.port_names = [p1.name, p2.name]
    ly.speed = ly.SPEED_1_GBPS

    # enable packet capture on both ports
    cp = cfg.captures.capture(name='cp')[-1]
    cp.port_names = [p1.name, p2.name]

    # flow1 = cfg.flows.flow(name='flow server:p1 -> server:p2')[-1]
    # flow1.tx_rx.port.tx_name = p1.name
    # flow1.tx_rx.port.rx_name = p2.name
    # flow1.size.fixed = 256
    # flow1.duration.fixed_packets.packets = 1000
    # flow1.rate.pps = 1000
    # eth2, ip2, udp2 = flow1.packet.ethernet().ipv4().udp()
    # eth2.src.value = '00:1B:6E:00:00:01'
    # eth2.dst.value = '02:6a:82:5d:da:2c'
    # ip2.src.value = '40.0.0.1'
    # ip2.dst.value = '201.0.0.2'
    # inc2 = udp2.src_port.increment
    # inc2.start, inc2.step, inc2.count = 6000, 4, 10
    # udp2.dst_port.values = [8000, 8044, 8060, 8074, 8082, 8084]

    flow2 = cfg.flows.flow(name='flow server:p2 -> server:p1')[-1]
    flow2.tx_rx.port.tx_name = p2.name
    flow2.tx_rx.port.rx_name = p1.name
    flow2.size.fixed = 256
    flow2.duration.fixed_packets.packets = 1000
    flow2.rate.pps = 1000
    eth2, ip2, udp2 = flow2.packet.ethernet().ipv4().udp()
    eth2.src.value = '00:1B:6E:00:00:02'
    eth2.dst.value = '02:6a:82:5d:da:2c'
    ip2.src.value = '201.0.0.2'
    ip2.dst.value = '40.0.0.1'
    inc2 = udp2.src_port.increment
    inc2.start, inc2.step, inc2.count = 6000, 4, 10
    udp2.dst_port.values = [8000, 8044, 8060, 8074, 8082, 8084]

    print('Pushing traffic configuration ...')
    tgen.set_config(cfg)

    print('Starting transmit on all configured flows ...')
    ts = tgen.transmit_state()
    ts.state = ts.START
    tgen.set_transmit_state(ts)

    # import pdb
    # pdb.set_trace()

    print('Checking metrics on all configured ports ...')
    print('Expected\tTotal Tx\tTotal Rx')
    assert wait_for(lambda: metrics_ok(tgen)), 'Metrics validation failed!'


def metrics_ok(api):
    # create a port metrics request and filter based on port names
    cfg = api.get_config()

    req = api.metrics_request()
    req.port.port_names = [p.name for p in cfg.ports]
    #import pdb;pdb.set_trace()
    # include only sent and received packet counts
    req.port.column_names = [req.port.FRAMES_TX, req.port.FRAMES_RX]
    # fetch port metrics
    res = api.get_metrics(req)
    # calculate total frames sent and received across all configured ports
    total_tx = sum([m.frames_tx for m in res.port_metrics])
    total_rx = sum([m.frames_rx for m in res.port_metrics])
    expected = sum([f.duration.fixed_packets.packets for f in cfg.flows])

    print('%d\t\t%d\t\t%d' % (expected, total_tx, total_rx))

    return expected == total_tx and total_rx >= expected


def wait_for(func, timeout=10, interval=0.2):
    '''
    Keeps calling the `func` until it returns true or `timeout` occurs
    every `interval` seconds.
    '''
    import time

    start = time.time()

    while time.time() - start <= timeout:
        if func():
            return True
        time.sleep(interval)

    print('Timeout occurred !')
    return False
