from commands.product import *
from queries.product import ProductStatus, ProductDetails, ProductLocation, ProductTrack
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product, Empty, Carrier
import time
import utils
import pytest

def arrange_end_transport_cmd_fields(product_id):
    return Product(product_id, 'Poste Delivery Business', 'BOX', []), \
        Carrier('motorbike', 'Postman12345', 'FE764OU'), \
        Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])


@pytest.mark.functional
@pytest.mark.skip(reason="story LOI-504")
def test_product_command_end_transport_OK(config,product_id):
    ## Arrange
    product, carrier, location = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [
        EndProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            location=location
        )
    ]
    
    queries = [
        ProductStatus(
            id = product_id
        ), 
        ProductDetails(
            id = product_id
        ), 
        ProductLocation(
            id = product_id
        ), 
        ProductTrack(
            id = product_id)
    ]
    
    expStatus = 'WIP'

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()
    assert responses[0].status_code == 202
    time.sleep(5)
    results =  plan.start_query()

    ## Assert
    assert len(results) == 4
    # status check
    assert results[0]['product']['id'] == product_id
    assert results[0]['status'] == expStatus
    
    # details check
    assert results[1]['product'] == cmds[0].payload['product']
    assert results[1]['sender'] == {}
    assert results[1]['receiver'] == {}

    # location check
    assert results[2]['product'] == cmds[0].payload['product']
    assert results[2]['location'] == cmds[0].payload['location']

    # track check
    assert results[3]['product']['id'] == product_id
    assert len(results[3]['trackList']) == 1
    assert results[3]['trackList'][0]['description'] == expStatus

@pytest.mark.functional
def test_product_command_end_transport_missing_product_KO(config,product_id):
    ## Arrange
    _, carrier, location = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [
        EndProductTransport(
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                location=location
            )
    ]

    plan = Plan(
        id=product_id,
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
def test_product_command_end_transport_missing_carrier_KO(config,product_id):
    ## Arrange
    product, _, location = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [ 
        EndProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                attributes=[],
                location=location
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['carrier']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_end_transport_missing_location_KO(config,product_id):
    ## Arrange
    product, carrier, _ = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [ 
        EndProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[]
            )
    ]

    plan = Plan(
        id=product_id,
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
def test_product_command_end_transport_missing_timestamp_KO(config,product_id):
    ## Arrange
    product, carrier, location = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [
        EndProductTransport(
                product=product,
                carrier=carrier,
                attributes=[],
                location=location
            )
    ]

    plan = Plan(
        id=product_id,
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
def test_product_command_end_transport_bad_timestamp_KO(config,product_id):
    ## Arrange
    product, carrier, location = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [
        EndProductTransport(
                product=product,
                carrier=carrier,
                attributes=[],
                location=location,
                timestamp=123456789,
            ),
        EndProductTransport(
                timestamp='A bad string',
                product=product,
                carrier=carrier,
                attributes=[],
                location=location
            ),
        EndProductTransport(
                timestamp='2022-12-07T17:48:33',
                product=product,
                carrier=carrier,
                attributes=[],
                location=location
            )
    ]

    plan = Plan(
        id=product_id,
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

@pytest.mark.functional
def test_product_command_end_transport_missing_attributes_KO(config,product_id):
    ## Arrange
    product, carrier, location = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [
        EndProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                location=location
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['attributes']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_end_transport_bad_attributes_KO(config,product_id):
    ## Arrange
    product, carrier, location = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [
        EndProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                location=location,
                attributes='a string'
            ),
        EndProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                location=location,
                attributes=5
            )
    ]

    plan = Plan(
        id=product_id,
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

@pytest.mark.functional
def test_product_command_end_transport_missing_api_key_OK(config,product_id):
    ## Arrange
    product, carrier, location = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [
        EndProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            location=location
        )
    ]
    
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        # missing api key
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    for response in responses:
        assert response.status_code == 403

@pytest.mark.functional
def test_product_command_end_transport_wrong_api_key_OK(config,product_id):
    ## Arrange
    product, carrier, location = arrange_end_transport_cmd_fields(product_id=product_id)
    cmds = [
        EndProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            location=location
        )
    ]
    
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key='fakeAPIKey'
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    for response in responses:
        assert response.status_code == 403
    