from commands.cluster import *
from tests.functionals.plan import Plan
from commands.models import  Location, Cluster
import time
import utils
import pytest


def arrange_open_cmd_fields(id):
    return Cluster(id), \
        Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])


@pytest.mark.functional
def test_cluster_command_open_OK(config,cluster_id):
    ## Arrange
    cluster, location = arrange_open_cmd_fields(id=cluster_id)
    ap = Open(
            cluster=cluster,
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
def test_cluster_command_open_missing_cluster_KO(config,cluster_id):
    ## Arrange
    _ , location = arrange_open_cmd_fields(id=cluster_id)
    cmds = [
        Open(
                timestamp=utils.get_timestamp(),
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
        cmds[0].payload['cluster']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400


@pytest.mark.functional
def test_cluster_command_open_missing_location_KO(config,cluster_id):
    ## Arrange
    cluster, _ = arrange_open_cmd_fields(id=cluster_id)
    cmds = [
        Open(
                timestamp=utils.get_timestamp(),
                cluster=cluster
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
def test_cluster_command_open_missing_timestamp_KO(config,cluster_id):
    ## Arrange
    cluster, location = arrange_open_cmd_fields(id=cluster_id)
    cmds = [
        Open(
                cluster=cluster,
                location=location,
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
def test_cluster_command_open_bad_timestamp_KO(config,cluster_id):
    ## Arrange
    cluster, location = arrange_open_cmd_fields(id=cluster_id)
    cmds = [
        Open(
                timestamp=123456789,   
                cluster=cluster,
                location=location
            ),
        Open(
                timestamp='A bad string',
                cluster=cluster,
                location=location
            ),
        Open(
                timestamp='2022-12-07T17:48:33',
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
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400
