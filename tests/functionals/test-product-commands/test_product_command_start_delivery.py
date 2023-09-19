from commands.product import *
from queries.product import ProductStatus, ProductDetails, ProductLocation, ProductTrack
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product, Empty, Carrier
import time
import utils
import pytest

def arrange_start_delivery_cmd_fields(product_id):
    return Product(product_id, 'Poste Delivery Business', 'BOX', []), \
        Carrier('motorbike', 'Postman12345', 'FE764OU')


@pytest.mark.functional
@pytest.mark.skip(reason="bug LOI-501, LOI-502, LOI-504")
def test_product_command_start_delivery_OK(config,product_id):
    ## Arrange
    product, carrier = arrange_start_delivery_cmd_fields(product_id=product_id)
    cmds = [
        StartProductDelivery(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[]
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
    
    expStatus = 'LastMileTransit'

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
    assert results[2]['location'] == {}

    # track check
    assert results[3]['product']['id'] == product_id
    assert len(results[3]['trackList']) == 1
    assert results[3]['trackList'][0]['description'] == expStatus

@pytest.mark.functional
def test_product_command_start_delivery_missing_product_KO(config,product_id):
    ## Arrange
    _, carrier = arrange_start_delivery_cmd_fields(product_id=product_id)
    cmds = [
        StartProductDelivery(
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
        cmds[0].payload['product']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_start_delivery_missing_carrier_KO(config,product_id):
    ## Arrange
    product, _ = arrange_start_delivery_cmd_fields(product_id=product_id)
    cmds = [ 
        StartProductDelivery(
                product=product,
                timestamp=utils.get_timestamp(),
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
        cmds[0].payload['carrier']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_start_delivery_missing_timestamp_KO(config,product_id):
    ## Arrange
    product, carrier = arrange_start_delivery_cmd_fields(product_id=product_id)
    cmds = [
        StartProductDelivery(
                product=product,
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
        cmds[0].payload['timestamp']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_start_delivery_bad_timestamp_KO(config,product_id):
    ## Arrange
    product, carrier = arrange_start_delivery_cmd_fields(product_id=product_id)
    cmds = [
        StartProductDelivery(
                product=product,
                timestamp=123456789,
                carrier=carrier,
                attributes=[]
            ),
        StartProductDelivery(
                product=product,
                timestamp='A bad string',
                carrier=carrier,
                attributes=[]
            ),
        StartProductDelivery(
                product=product,
                timestamp='2022-12-07T17:48:33',
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
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400

@pytest.mark.functional
def test_product_command_start_delivery_missing_attributes_KO(config,product_id):
    ## Arrange
    product, carrier = arrange_start_delivery_cmd_fields(product_id=product_id)
    cmds = [
        StartProductDelivery(
                product=product,
                carrier=carrier,
                timestamp=utils.get_timestamp(),
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
def test_product_command_start_delivery_bad_attributes_KO(config,product_id):
    ## Arrange
    product, carrier = arrange_start_delivery_cmd_fields(product_id=product_id)
    cmds = [
        StartProductDelivery(
                product=product,
                carrier=carrier,
                timestamp=utils.get_timestamp(),
                attributes='a string'
            ),
        StartProductDelivery(
                product=product,
                carrier=carrier,
                timestamp=utils.get_timestamp(),
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
def test_product_command_start_delivery_missing_api_key_OK(config,product_id):
    ## Arrange
    product, carrier = arrange_start_delivery_cmd_fields(product_id=product_id)
    cmds = [
        StartProductDelivery(
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
        # missing api key
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    for response in responses:
        assert response.status_code == 403

@pytest.mark.functional
def test_product_command_start_delivery_wrong_api_key_OK(config,product_id):
    ## Arrange
    product, carrier = arrange_start_delivery_cmd_fields(product_id=product_id)
    cmds = [
        StartProductDelivery(
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
        api_key='fakeAPIKey'
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    for response in responses:
        assert response.status_code == 403
    