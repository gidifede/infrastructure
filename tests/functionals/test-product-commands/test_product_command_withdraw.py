from commands.product import *
from queries.product import ProductStatus, ProductDetails, ProductLocation, ProductTrack
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product, Empty
import time
import utils
import pytest

def arrange_withdraw_cmd_fields(product_id):
    return Product(product_id, 'Poste Delivery Business', 'BOX', [])

@pytest.mark.functional
@pytest.mark.skip(reason="story LOI-504")
def test_product_command_withdraw_OK(config,product_id):
    ## Arrange
    product = arrange_withdraw_cmd_fields(product_id=product_id)
    cmds = [ 
        WithdrawProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            notes='Receiver not found'
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
    
    expStatus = 'Withdrawn'

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
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_withdraw_missing_product_KO(config,product_id):
    ## Arrange
    _ = arrange_withdraw_cmd_fields(product_id=product_id)
    cmds = [
        WithdrawProduct(
                timestamp=utils.get_timestamp(),
                notes='Receiver not found'
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
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_withdraw_missing_timestamp_KO(config,product_id):
    ## Arrange
    product = arrange_withdraw_cmd_fields(product_id=product_id)
    cmds = [
        WithdrawProduct(
                product=product,
                notes='Receiver not found'
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
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_withdraw_bad_timestamp_KO(config,product_id):
    ## Arrange
    product = arrange_withdraw_cmd_fields(product_id=product_id)
    cmds = [
        WithdrawProduct(
                product=product,
                timestamp=123456789,
                notes='Receiver not found'
            ),
        WithdrawProduct(
                product=product,
                timestamp='A bad string',
                notes='Receiver not found'
            ),
        WithdrawProduct(
                product=product,
                timestamp='2022-12-07T17:48:33',
                notes='Receiver not found'
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
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_withdraw_missing_notes_KO(config,product_id):
    ## Arrange
    product = arrange_withdraw_cmd_fields(product_id=product_id)
    cmds = [
        WithdrawProduct(
                product=product,
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
        cmds[0].payload['notes']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_withdraw_bad_notes_KO(config,product_id):
    ## Arrange
    product = arrange_withdraw_cmd_fields(product_id=product_id)
    cmds = [
        WithdrawProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                notes=12345
            ),
        WithdrawProduct(
                product=product,
                timestamp=utils.get_timestamp(),
                notes=[]
            ),
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
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_withdraw_missing_api_key_OK(config,product_id):
    ## Arrange
    product = arrange_withdraw_cmd_fields(product_id=product_id)
    cmds = [ 
        WithdrawProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            notes='Receiver not found'
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
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_withdraw_wrong_api_key_OK(config,product_id):
    ## Arrange
    product = arrange_withdraw_cmd_fields(product_id=product_id)
    cmds = [ 
        WithdrawProduct(
            product=product,
            timestamp=utils.get_timestamp(),
            notes='Receiver not found'
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