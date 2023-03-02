import grpc
import netmiko
import pydpu
import pytest
import snappi
from passwords import *
from testbed import *


@pytest.fixture
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
    dpu_connect.send_command('ip addr add 202.0.0.1/24 dev pf0hpf')
    dpu_connect.send_command('ip addr add 203.0.0.1/24 dev pf1hpf')
    dpu_connect.send_command('sysctl net.ipv4.ip_forward=1')
    dpu_connect.send_command('iptables -P FORWARD ACCEPT')
    dpu_connect.send_command('arp -s 200.0.0.2 80:09:02:02:00:01')
    dpu_connect.send_command('arp -s 201.0.0.2 80:09:02:02:00:02')
    dpu_connect.send_command('arp -s 202.0.0.2 80:09:02:02:00:03')
    dpu_connect.send_command('arp -s 203.0.0.2 80:09:02:02:00:04')
    # idealy replace as many of those ssh commands with OPI API as they become available

    # iptables -P INPUT ACCEPT
    # iptables -P FORWARD ACCEPT
    # iptables -P OUTPUT ACCEPT
    # iptables -t nat -F
    # iptables -t mangle -F
    # iptables -F
    # iptables -X

    yield dpu_connect

    dpu_connect.send_command('ip addr del 200.0.0.1/24 dev enp3s0f0s0')
    dpu_connect.send_command('ip addr del 201.0.0.1/24 dev enp3s0f1s0')
    dpu_connect.send_command('ip addr del 202.0.0.1/24 dev pf0hpf')
    dpu_connect.send_command('ip addr del 203.0.0.1/24 dev pf1hpf')
    dpu_connect.send_command('ip neigh flush 200.0.0.0/8')
    dpu_connect.send_command('ip neigh flush 201.0.0.0/8')
    dpu_connect.send_command('ip neigh flush 202.0.0.0/8')
    dpu_connect.send_command('ip neigh flush 203.0.0.0/8')


@pytest.fixture
def server():
    server_info = {
        'device_type': 'linux',
        'host': SERVER_IP,
        'username': 'root',
        'password': SERVER_PASSWORD
    }
    server_connect = netmiko.ConnectHandler(**server_info)
    output = server_connect.send_command('docker compose -f /root/opi/tgen/deployment/tgen.yml up -d', read_timeout=30)
    print(output)

    yield server_connect

    output = server_connect.send_command('docker compose -f /root/opi/tgen/deployment/tgen.yml down')
    print(output)


@pytest.fixture
def host():
    # needs reboot 
    # mst start
    # mlxconfig -d /dev/mst/mt41686_pciconf0 s INTERNAL_CPU_MODEL=1
    # mlxconfig -d /dev/mst/mt41686_pciconf0.1 s INTERNAL_CPU_MODEL=1
    # mlxconfig -d /dev/mst/mt41686_pciconf0 s INTERNAL_CPU_MODEL=0
    # mlxconfig -d /dev/mst/mt41686_pciconf0.1 s INTERNAL_CPU_MODEL=0
    host_info = {
        'device_type': 'linux',
        'host': HOST_IP,
        'username': 'root',
        'password': HOST_PASSWORD
    }
    host_connect = netmiko.ConnectHandler(**host_info)

    output = host_connect.send_command('docker compose -f /root/opi/tgen/deployment/tgen.yml up -d', read_timeout=30)
    print(output)

    yield host_connect

    output = host_connect.send_command('docker compose -f /root/opi/tgen/deployment/tgen.yml down')
    print(output)


def test_server_to_server_via_dpu(dpu, server):
    tgen = snappi.api(location=f'https://{SERVER_IP}', verify=False)
    cfg = tgen.config()
    p1 = cfg.ports.port(name='server:p1', location=f'{SERVER_IP}:5555')[-1]
    p2 = cfg.ports.port(name='server:p2', location=f'{SERVER_IP}:5556')[-1]

    # add layer 1 property to configure same speed on both ports
    ly = cfg.layer1.layer1(name='ly')[-1]
    ly.port_names = [p1.name, p2.name]
    ly.speed = ly.SPEED_100_GBPS

    # enable packet capture on both ports
    cp = cfg.captures.capture(name='cp')[-1]
    cp.port_names = [p1.name, p2.name]

    flow = cfg.flows.flow(name='flow server:p2 -> server:p1')[-1]
    flow.tx_rx.port.tx_name = p2.name
    flow.tx_rx.port.rx_name = p1.name
    flow.size.fixed = 256
    flow.duration.fixed_packets.packets = 1000
    flow.rate.pps = 1000
    eth2, ip2, udp2 = flow.packet.ethernet().ipv4().udp()
    eth2.src.value = '00:1B:6E:00:00:01'
    eth2.dst.value = '02:6a:82:5d:da:2c'
    ip2.src.value = '201.0.0.2'
    ip2.dst.value = '200.0.0.2'
    inc2 = udp2.src_port.increment
    inc2.start, inc2.step, inc2.count = 6000, 4, 10
    udp2.dst_port.values = [8000, 8044, 8060, 8074, 8082, 8084]

    print('Pushing traffic configuration ...')
    tgen.set_config(cfg)

    print('Starting transmit on all configured flows ...')
    ts = tgen.transmit_state()
    ts.state = ts.START
    tgen.set_transmit_state(ts)

    #import pdb
    #pdb.set_trace()

    print('Checking metrics on all configured ports ...')
    print('Expected\tTotal Tx\tTotal Rx')
    assert wait_for(lambda: metrics_ok(tgen)), 'Metrics validation failed!'


def test_server_to_dpuhost_via_dpu(dpu, server, host):
    tgen = snappi.api(location=f'https://{SERVER_IP}', verify=False)
    cfg = tgen.config()
    p1 = cfg.ports.port(name='host:p1', location=f'{HOST_IP}:5555')[-1]
    p2 = cfg.ports.port(name='server:p1', location=f'{SERVER_IP}:5556')[-1]

    # add layer 1 property to configure same speed on both ports
    ly = cfg.layer1.layer1(name='ly')[-1]
    ly.port_names = [p1.name, p2.name]
    ly.speed = ly.SPEED_1_GBPS

    # enable packet capture on both ports
    cp = cfg.captures.capture(name='cp')[-1]
    cp.port_names = [p1.name, p2.name]

    flow = cfg.flows.flow(name='flow server:p2 -> host:p1')[-1]
    flow.tx_rx.port.tx_name = p2.name
    flow.tx_rx.port.rx_name = p1.name
    flow.size.fixed = 256
    flow.duration.fixed_packets.packets = 1000
    flow.rate.pps = 1000
    eth2, ip2, udp2 = flow.packet.ethernet().ipv4().udp()
    eth2.src.value = '00:1B:6E:00:00:01'
    eth2.dst.value = '02:6a:82:5d:da:2c'
    ip2.src.value = '200.0.0.2'
    ip2.dst.value = '202.0.0.2'
    inc2 = udp2.src_port.increment
    inc2.start, inc2.step, inc2.count = 6000, 4, 10
    udp2.dst_port.values = [8000, 8044, 8060, 8074, 8082, 8084]

    print('Pushing traffic configuration ...')
    tgen.set_config(cfg)

    print('Starting transmit on all configured flows ...')
    ts = tgen.transmit_state()
    ts.state = ts.START
    tgen.set_transmit_state(ts)

    print('Checking metrics on all configured ports ...')
    print('Expected\tTotal Tx\tTotal Rx')
    assert wait_for(lambda: metrics_ok(tgen)), 'Metrics validation failed!'


def test_dpuhost_to_dpuhost_via_dpu(dpu, host):
    tgen = snappi.api(location=f'https://{HOST_IP}', verify=False)
    cfg = tgen.config()
    p1 = cfg.ports.port(name='host:p1', location=f'{HOST_IP}:5555')[-1]
    p2 = cfg.ports.port(name='host:p2', location=f'{HOST_IP}:5556')[-1]

    # add layer 1 property to configure same speed on both ports
    ly = cfg.layer1.layer1(name='ly')[-1]
    ly.port_names = [p1.name, p2.name]
    ly.speed = ly.SPEED_1_GBPS

    # enable packet capture on both ports
    cp = cfg.captures.capture(name='cp')[-1]
    cp.port_names = [p1.name, p2.name]

    flow = cfg.flows.flow(name='flow host:p2 -> host:p1')[-1]
    flow.tx_rx.port.tx_name = p2.name
    flow.tx_rx.port.rx_name = p1.name
    flow.size.fixed = 256
    flow.duration.fixed_packets.packets = 1000
    flow.rate.pps = 1000
    eth2, ip2, udp2 = flow.packet.ethernet().ipv4().udp()
    eth2.src.value = '00:1B:6E:00:00:01'
    eth2.dst.value = '02:6a:82:5d:da:2c'
    ip2.src.value = '203.0.0.2'
    ip2.dst.value = '202.0.0.2'
    inc2 = udp2.src_port.increment
    inc2.start, inc2.step, inc2.count = 6000, 4, 10
    udp2.dst_port.values = [8000, 8044, 8060, 8074, 8082, 8084]

    print('Pushing traffic configuration ...')
    tgen.set_config(cfg)

    print('Starting transmit on all configured flows ...')
    ts = tgen.transmit_state()
    ts.state = ts.START
    tgen.set_transmit_state(ts)

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
