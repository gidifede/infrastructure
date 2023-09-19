from commands.cluster import *
from tests.functionals.plan import Plan
from commands.models import Product, Location, Cluster
import time
import utils
import pytest

def arrange_add_product_cmd_fields(id, product_id):
    return Cluster(id), \
           Product(product_id, 'Poste Delivery Business', 'BOX', []), \
           Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])


@pytest.mark.functional
def test_cluster_command_add_product_OK(config,cluster_id, product_id):
    ## Arrange
    cluster, product, location = arrange_add_product_cmd_fields(id=cluster_id, product_id=product_id)
    ap = AddProduct(
            cluster=cluster,
            product=product,
            mode="",
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
def test_cluster_command_add_product_missing_cluster_KO(config,cluster_id, product_id):
    ## Arrange
    _ , product, location = arrange_add_product_cmd_fields(id=cluster_id, product_id=product_id)
    cmds = [
        AddProduct(
                timestamp=utils.get_timestamp(),
                location=location,
                product=product
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
def test_cluster_command_add_product_missing_location_KO(config,cluster_id, product_id):
    ## Arrange
    cluster, product, _ = arrange_add_product_cmd_fields(id=cluster_id, product_id=product_id)
    cmds = [
        AddProduct(
                timestamp=utils.get_timestamp(),
                cluster=cluster,
                product=product,
                mode=""
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
def test_cluster_command_add_product_missing_product_KO(config,cluster_id, product_id):
    ## Arrange
    cluster, _, location = arrange_add_product_cmd_fields(id=cluster_id, product_id=product_id)
    cmds = [
        AddProduct(
                timestamp=utils.get_timestamp(),
                cluster=cluster,
                location=location,
                mode=""   
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
        cmds[0].payload['product']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400


    
@pytest.mark.functional
def test_cluster_command_add_product_missing_mode_KO(config,cluster_id, product_id):
    ## Arrange
    cluster, product, location = arrange_add_product_cmd_fields(id=cluster_id, product_id=product_id)
    cmds = [
        AddProduct(
                timestamp=utils.get_timestamp(),
                cluster=cluster,
                product=product,
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
        cmds[0].payload['mode']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400


@pytest.mark.functional
def test_cluster_command_add_product_missing_timestamp_KO(config,cluster_id,product_id):
    ## Arrange
    cluster, product, location = arrange_add_product_cmd_fields(id=cluster_id, product_id=product_id)
    cmds = [
        AddProduct(
                cluster=cluster,
                product=product,
                location=location,
                mode=""   
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
def test_cluster_command_add_product_bad_timestamp_KO(config,cluster_id,product_id):
    ## Arrange
    cluster, product, location = arrange_add_product_cmd_fields(id=cluster_id, product_id=product_id)
    cmds = [
        AddProduct(
                timestamp=123456789,   
                cluster=cluster,
                product=product,
                location=location,
                mode=""  
            ),
        AddProduct(
                timestamp='A bad string',
                cluster=cluster,
                product=product,
                location=location,
                mode=""    
            ),
        AddProduct(
                timestamp='2022-12-07T17:48:33',
                cluster=cluster,
                product=product,
                location=location,
                mode="" 
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