from commands.cluster import *
from tests.functionals.plan import Plan
from commands.models import Transport, Location, Cluster
import time
import utils
import pytest


def arrange_start_transport_cmd_fields(id):
    return Cluster(id), \
           Transport("520"),\
           Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])


@pytest.mark.functional
def test_cluster_command_start_transport_OK(config,cluster_id):
    ## Arrange
    cluster, transport, location = arrange_start_transport_cmd_fields(id=cluster_id)
    ap = StartTransport(
            cluster=cluster,
            transport=transport,
            timestamp=utils.get_timestamp(),
            location=location
        )
    cmds = []
    cmds.append(ap) 
    plan = Plan(
        id=cluster_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()
    assert responses[0].status_code == 202
    

@pytest.mark.functional
def test_cluster_command_start_transport_missing_cluster_KO(config,cluster_id):
    ## Arrange
    _ , transport, location = arrange_start_transport_cmd_fields(id=cluster_id)
    cmds = [
        StartTransport(
                timestamp=utils.get_timestamp(),
                location=location,
                transport=transport
            )
    ]

    plan = Plan(
        id=cluster_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['cluster']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400


@pytest.mark.functional
def test_cluster_command_start_transport_missing_location_KO(config,cluster_id):
    ## Arrange
    cluster, transport, _ = arrange_start_transport_cmd_fields(id=cluster_id)
    cmds = [
        StartTransport(
                timestamp=utils.get_timestamp(),
                cluster=cluster,
                transport=transport
            )
    ]
    plan = Plan(
        id=cluster_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['location']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

    
@pytest.mark.functional
def test_cluster_command_start_transport_missing_transport_KO(config,cluster_id):
    ## Arrange
    cluster, _, location = arrange_start_transport_cmd_fields(id=cluster_id)
    cmds = [
        StartTransport(
                timestamp=utils.get_timestamp(),
                cluster=cluster,
                location=location   
            )
    ]
    plan = Plan(
        id=cluster_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['transport']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400


@pytest.mark.functional
def test_cluster_command_start_transport_missing_timestamp_KO(config,cluster_id):
    ## Arrange
    cluster, transport,location = arrange_start_transport_cmd_fields(id=cluster_id)
    cmds = [
        StartTransport(
                cluster=cluster,
                location=location,
                transport=transport
            )
    ]

    plan = Plan(
        id=cluster_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['timestamp']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400


@pytest.mark.functional
def test_cluster_command_start_transport_bad_timestamp_KO(config,cluster_id):

    ## Arrange
    cluster, transport,location = arrange_start_transport_cmd_fields(id=cluster_id)
    cmds = [
        StartTransport(
                timestamp=123456789,   
                cluster=cluster,
                location=location,
                transport=transport
            ),
        StartTransport(
                timestamp='A bad string',
                cluster=cluster,
                location=location,
                transport=transport
            ),
        StartTransport(
                timestamp='2022-12-07T17:48:33',
                cluster=cluster,
                location=location,
                transport=transport
            )
    ]
    
    plan = Plan(
        id=cluster_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400
